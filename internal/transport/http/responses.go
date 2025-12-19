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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	if err == nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "unknown error"})
		return
	}

	// IMPORTANTE:
	// Aquí NO usamos domain.ErrInvalid porque tu dominio NO lo tiene.
	// Solo mapeamos los errores que sí existen.
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, apiError{Error: err.Error()})
	case errors.Is(err, domain.ErrDuplicate):
		writeJSON(w, http.StatusConflict, apiError{Error: err.Error()})
	default:
		// cualquier otro error lo mandamos como 400 para simplificar
		writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
	}
}
