package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	// Mapeo básico de errores del dominio a HTTP codes
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
	case errors.Is(err, domain.ErrDuplicate):
		writeJSON(w, http.StatusConflict, map[string]any{"error": err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
}
