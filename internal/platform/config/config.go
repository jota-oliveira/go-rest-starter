package config

import "os"

type Config struct{ Addr, StoreDriver, DatabaseURL, MongoURI, MongoDatabase, JWTIssuer, JWTAudience, JWTPublicKeyFile string }

func env(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
func Load() Config {
	return Config{
		Addr: env("HTTP_ADDR", ":8080"), StoreDriver: env("STORE_DRIVER", "memory"), DatabaseURL: os.Getenv("DATABASE_URL"), MongoURI: env("MONGO_URI", "mongodb://localhost:27017"), MongoDatabase: env("MONGO_DATABASE", "go_rest_starter"), JWTIssuer: env("JWT_ISSUER", "nextjs-app"), JWTAudience: env("JWT_AUDIENCE", "go-api"), JWTPublicKeyFile: os.Getenv("JWT_PUBLIC_KEY_FILE"),
	}
}
