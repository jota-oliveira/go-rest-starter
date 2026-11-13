package httpapi

import (
	"context"
	"example.com/go-rest-starter/internal/application/auth"
	"net/http"
	"strings"
)

type principalKey struct{}

func PrincipalFrom(ctx context.Context) (auth.Principal, bool) {
	v, ok := ctx.Value(principalKey{}).(auth.Principal)
	return v, ok
}
func RequireAuth(verifier auth.Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeJSON(w, 401, map[string]string{"error": "unauthorized"})
				return
			}
			principal, err := verifier.Verify(r.Context(), parts[1])
			if err != nil || principal.Subject == "" {
				writeJSON(w, 401, map[string]string{"error": "unauthorized"})
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, principal)))
		})
	}
}
func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := PrincipalFrom(r.Context())
			if !ok || !principal.Scopes[scope] {
				writeJSON(w, 403, map[string]string{"error": "forbidden"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
