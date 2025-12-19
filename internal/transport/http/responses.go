package http

import (
	"encoding/json"
	"net/http"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// writeJSON responde con JSON estándar
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeErr convierte errores de dominio a HTTP
func writeErr(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNotFound:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case domain.ErrDuplicate:
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}
