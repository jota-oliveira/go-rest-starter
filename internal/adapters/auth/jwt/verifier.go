package jwt

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"os"
	"strings"

	"example.com/go-rest-starter/internal/application/auth"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid bearer token")

type Verifier struct {
	key              ed25519.PublicKey
	issuer, audience string
}

func New(publicKeyFile, issuer, audience string) (*Verifier, error) {
	pemBytes, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("read JWT public key: %w", err)
	}
	key, err := jwtlib.ParseEdPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("parse JWT public key: %w", err)
	}
	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("expected Ed25519 public key")
	}
	if issuer == "" || audience == "" {
		return nil, errors.New("JWT issuer and audience required")
	}
	return &Verifier{key: pub, issuer: issuer, audience: audience}, nil
}

type Claims struct {
	Scope string `json:"scope,omitempty"`
	jwtlib.RegisteredClaims
}

func (v *Verifier) Verify(_ context.Context, raw string) (auth.Principal, error) {
	claims := &Claims{}
	token, err := jwtlib.ParseWithClaims(raw, claims, func(t *jwtlib.Token) (any, error) { return v.key, nil }, jwtlib.WithValidMethods([]string{jwtlib.SigningMethodEdDSA.Alg()}), jwtlib.WithIssuer(v.issuer), jwtlib.WithAudience(v.audience), jwtlib.WithExpirationRequired(), jwtlib.WithIssuedAt(), jwtlib.WithLeeway(15_000_000_000))
	if err != nil || !token.Valid || claims.Subject == "" {
		return auth.Principal{}, ErrInvalidToken
	}
	scopes := make(map[string]bool)
	for _, s := range strings.Fields(claims.Scope) {
		scopes[s] = true
	}
	return auth.Principal{Subject: claims.Subject, Issuer: claims.Issuer, Scopes: scopes}, nil
}
