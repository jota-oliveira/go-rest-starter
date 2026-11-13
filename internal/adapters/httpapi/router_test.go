package httpapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"example.com/go-rest-starter/internal/adapters/httpapi"
	"example.com/go-rest-starter/internal/adapters/persistence/memory"
	"example.com/go-rest-starter/internal/application/auth"
	"example.com/go-rest-starter/internal/application/items"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testVerifier struct{}

func (testVerifier) Verify(_ context.Context, token string) (auth.Principal, error) {
	if token != "valid" {
		return auth.Principal{}, errors.New("invalid")
	}
	return auth.Principal{Subject: "user-1", Issuer: "test", Scopes: map[string]bool{"items:read": true, "items:write": true}}, nil
}
func TestEndpoints(t *testing.T) {
	router := httpapi.NewRouter(items.NewService(memory.New()), testVerifier{})
	health := httptest.NewRecorder()
	router.ServeHTTP(health, httptest.NewRequest("GET", "/health", nil))
	if health.Code != 200 {
		t.Fatalf("health: %d", health.Code)
	}
	unauth := httptest.NewRecorder()
	router.ServeHTTP(unauth, httptest.NewRequest("GET", "/api/v1/items", nil))
	if unauth.Code != 401 {
		t.Fatalf("unauthorized: %d", unauth.Code)
	}
	post := httptest.NewRequest("POST", "/api/v1/items", strings.NewReader(`{"name":"first"}`))
	post.Header.Set("Authorization", "Bearer valid")
	created := httptest.NewRecorder()
	router.ServeHTTP(created, post)
	if created.Code != 201 {
		t.Fatalf("create: %d %s", created.Code, created.Body.String())
	}
	get := httptest.NewRequest("GET", "/api/v1/items", nil)
	get.Header.Set("Authorization", "Bearer valid")
	result := httptest.NewRecorder()
	router.ServeHTTP(result, get)
	var body struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(result.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if result.Code != 200 || len(body.Data) != 1 || body.Data[0].Name != "first" {
		t.Fatalf("list: %d %s", result.Code, result.Body.String())
	}
}
