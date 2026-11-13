package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-rest-starter/internal/adapters/auth/jwt"
	"example.com/go-rest-starter/internal/adapters/httpapi"
	"example.com/go-rest-starter/internal/adapters/persistence/memory"
	"example.com/go-rest-starter/internal/adapters/persistence/mongodb"
	"example.com/go-rest-starter/internal/adapters/persistence/postgres"
	"example.com/go-rest-starter/internal/application/items"
	"example.com/go-rest-starter/internal/platform/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
func run() error {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if cfg.JWTPublicKeyFile == "" {
		return errors.New("JWT_PUBLIC_KEY_FILE must be configured (make keygen)")
	}
	verifier, err := jwt.New(cfg.JWTPublicKeyFile, cfg.JWTIssuer, cfg.JWTAudience)
	if err != nil {
		return err
	}
	var repo items.Repository
	switch cfg.StoreDriver {
	case "memory":
		repo = memory.New()
	case "postgres":
		if cfg.DatabaseURL == "" {
			return errors.New("DATABASE_URL required")
		}
		pool, e := pgxpool.New(ctx, cfg.DatabaseURL)
		if e != nil {
			return e
		}
		defer pool.Close()
		if e = pool.Ping(ctx); e != nil {
			return e
		}
		repo = postgres.New(pool)
	case "mongodb":
		client, e := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
		if e != nil {
			return e
		}
		defer func() { _ = client.Disconnect(context.Background()) }()
		if e = client.Ping(ctx, nil); e != nil {
			return e
		}
		repo = mongodb.New(client.Database(cfg.MongoDatabase))
	default:
		return errors.New("STORE_DRIVER must be memory, postgres or mongodb")
	}
	srv := &http.Server{Addr: cfg.Addr, Handler: httpapi.NewRouter(items.NewService(repo), verifier), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 35 * time.Second, IdleTimeout: 60 * time.Second}
	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", cfg.Addr, "store", cfg.StoreDriver)
		errCh <- srv.ListenAndServe()
	}()
	select {
	case err = <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
