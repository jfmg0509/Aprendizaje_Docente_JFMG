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
    // manejo de errores (requerimiento)
    status := http.StatusInternalServerError
    msg := err.Error()

    switch {
    case errors.Is(err, domain.ErrNotFound):
        status = http.StatusNotFound
        msg = "not found"
    case errors.Is(err, domain.ErrDuplicate):
        status = http.StatusConflict
        msg = "duplicate"
    case errors.Is(err, domain.ErrValidation), errors.Is(err, domain.ErrInvalidRole), errors.Is(err, domain.ErrInvalidAccess):
        status = http.StatusBadRequest
        msg = err.Error()
    case errors.Is(err, domain.ErrInactiveEntity):
        status = http.StatusForbidden
        msg = "inactive entity"
    }

    writeJSON(w, status, apiError{Error: msg})
}
