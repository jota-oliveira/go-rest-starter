package httpapi

import (
	"net/http"
	"time"

	"example.com/go-rest-starter/internal/application/auth"
	"example.com/go-rest-starter/internal/application/items"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(service *items.Service, verifier auth.Verifier) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, middleware.Timeout(30*time.Second))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, map[string]string{"status": "ok"}) })
	h := handler{service: service}
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(RequireAuth(verifier))
		r.With(RequireScope("items:read")).Get("/items", h.list)
		r.With(RequireScope("items:write")).Post("/items", h.create)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			p, _ := PrincipalFrom(r.Context())
			writeJSON(w, 200, map[string]any{"sub": p.Subject, "iss": p.Issuer})
		})
	})
	return r
}
