package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// apiError es el formato estándar de error JSON.
type apiError struct {
	Error string `json:"error"`
}

// writeJSON escribe respuesta JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeErr mapea errores del dominio a HTTP status codes.
func writeErr(w http.ResponseWriter, err error) {
	if err == nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "unknown error"})
		return
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, apiError{Error: err.Error()})
	case errors.Is(err, domain.ErrDuplicate):
		writeJSON(w, http.StatusConflict, apiError{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()})
	}
}
