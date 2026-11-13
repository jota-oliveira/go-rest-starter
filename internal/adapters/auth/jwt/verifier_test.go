package jwt_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	verifier "example.com/go-rest-starter/internal/adapters/auth/jwt"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

func TestVerifier(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "public.pem")
	if err := os.WriteFile(file, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	v, err := verifier.New(file, "nextjs-app", "go-api")
	if err != nil {
		t.Fatal(err)
	}
	claims := jwtlib.MapClaims{"iss": "nextjs-app", "aud": "go-api", "sub": "user-1", "exp": time.Now().Add(time.Minute).Unix(), "iat": time.Now().Unix(), "scope": "items:read"}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodEdDSA, claims)
	signed, err := token.SignedString(priv)
	if err != nil {
		t.Fatal(err)
	}
	p, err := v.Verify(context.Background(), signed)
	if err != nil {
		t.Fatal(err)
	}
	if p.Subject != "user-1" || !p.Scopes["items:read"] {
		t.Fatal("claims incorrectly parsed")
	}
	claims["aud"] = "another-api"
	bad, _ := jwtlib.NewWithClaims(jwtlib.SigningMethodEdDSA, claims).SignedString(priv)
	if _, err := v.Verify(context.Background(), bad); err == nil {
		t.Fatal("wrong audience accepted")
	}
}
