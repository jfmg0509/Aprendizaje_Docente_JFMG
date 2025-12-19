package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type apiError struct {
	Error string `json:"error"`
}

// writeJSON responde JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeErr traduce errores del dominio a HTTP status.
func writeErr(w http.ResponseWriter, err error) {
	if err == nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "unknown error"})
		return
	}

	// Dominio -> status code
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, apiError{Error: err.Error()})
	case errors.Is(err, domain.ErrDuplicate):
		writeJSON(w, http.StatusConflict, apiError{Error: err.Error()})
	case errors.Is(err, domain.ErrValidation),
		errors.Is(err, domain.ErrInvalidRole),
		errors.Is(err, domain.ErrInvalidAccess):
		writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()})
	}
}
