package auth

import "context"

type Principal struct {
	Subject string
	Issuer  string
	Scopes  map[string]bool
}
type Verifier interface {
	Verify(ctx context.Context, raw string) (Principal, error)
}
