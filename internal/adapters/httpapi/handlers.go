package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"example.com/go-rest-starter/internal/application/items"
	"example.com/go-rest-starter/internal/domain/item"
)

type handler struct{ service *items.Service }

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
func (h handler) list(w http.ResponseWriter, r *http.Request) {
	principal, _ := PrincipalFrom(r.Context())
	result, err := h.service.List(r.Context(), principal.Subject)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal_error"})
		return
	}
	writeJSON(w, 200, map[string]any{"data": result})
}
func (h handler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	if decoder.Decode(new(any)) != io.EOF {
		writeJSON(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	principal, _ := PrincipalFrom(r.Context())
	result, err := h.service.Create(r.Context(), principal.Subject, strings.TrimSpace(req.Name))
	if errors.Is(err, item.ErrInvalidName) {
		writeJSON(w, 422, map[string]string{"error": "invalid_name"})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal_error"})
		return
	}
	writeJSON(w, 201, result)
}
