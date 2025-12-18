package http // Paquete HTTP: utilidades para respuestas y errores REST

import (
	"encoding/json" // Serialización JSON
	"errors"        // Comparación de errores con errors.Is
	"net/http"      // Códigos de estado y ResponseWriter

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain" // Errores del dominio
)

// writeJSON escribe una respuesta HTTP en formato JSON.
// Recibe:
// - w: ResponseWriter
// - status: código HTTP (200, 201, 400, etc.)
// - v: cualquier estructura a serializar (any)
func writeJSON(w http.ResponseWriter, status int, v any) {
	// Define el tipo de contenido como JSON UTF-8
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Escribe el código de estado HTTP
	w.WriteHeader(status)

	// Codifica la estructura 'v' a JSON y la envía en el body
	// El error se ignora porque, en este punto, no hay mucho que hacer si falla
	_ = json.NewEncoder(w).Encode(v)
}

// writeErr centraliza el manejo de errores del sistema.
// Convierte errores del dominio en respuestas HTTP REST apropiadas.
func writeErr(w http.ResponseWriter, err error) {

	// Por defecto, cualquier error es un error interno del servidor
	status := http.StatusInternalServerError
	msg := err.Error()

	// Switch basado en comparación de errores
	// errors.Is permite comparar incluso errores envueltos (wrapping)
	switch {

	// Recurso no encontrado → 404
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		msg = "not found"

	// Recurso duplicado → 409 Conflict
	case errors.Is(err, domain.ErrDuplicate):
		status = http.StatusConflict
		msg = "duplicate"

	// Errores de validación → 400 Bad Request
	case errors.Is(err, domain.ErrValidation),
		errors.Is(err, domain.ErrInvalidRole),
		errors.Is(err, domain.ErrInvalidAccess):
		status = http.StatusBadRequest
		msg = err.Error()

	// Entidad inactiva → 403 Forbidden
	case errors.Is(err, domain.ErrInactiveEntity):
		status = http.StatusForbidden
		msg = "inactive entity"
	}

	// Responde el error en formato JSON usando el DTO apiError
	writeJSON(w, status, apiError{Error: msg})
}
