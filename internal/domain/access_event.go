package domain // Dominio: eventos y reglas del negocio

import (
	"fmt"     // Para construir errores con contexto
	"strings" // Para normalizar strings
	"time"    // Para timestamps
)

// AccessEvent representa un evento de acceso de un usuario a un libro.
// Es una entidad del dominio basada en eventos (event-driven).
type AccessEvent struct {
	id         uint64     // ID único del evento (asignado por BD)
	userID     uint64     // ID del usuario que accede
	bookID     uint64     // ID del libro accedido
	accessType AccessType // Tipo de acceso (VIEW, DOWNLOAD, etc.)
	createdAt  time.Time  // Fecha y hora del evento
}

// NewAccessEvent es el constructor del evento.
// Aplica validaciones del dominio antes de crear el objeto.
func NewAccessEvent(userID, bookID uint64, accessType AccessType) (*AccessEvent, error) {

	// Normaliza el tipo de acceso a mayúsculas
	accessType = AccessType(strings.ToUpper(string(accessType)))

	// Valida que el tipo de acceso sea permitido
	if !accessType.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidAccess, accessType)
	}

	// Valida que los IDs sean obligatorios
	if userID == 0 || bookID == 0 {
		return nil, fmt.Errorf("%w: user_id and book_id are required", ErrValidation)
	}

	// Crea y retorna el evento
	return &AccessEvent{
		userID:     userID,
		bookID:     bookID,
		accessType: accessType,
		createdAt:  time.Now(), // Timestamp del evento
	}, nil
}

// HydrateAccessEvent reconstruye un evento desde la base de datos.
// Se usa cuando se leen registros ya persistidos.
func HydrateAccessEvent(
	id, userID, bookID uint64,
	accessType AccessType,
	createdAt time.Time,
) (*AccessEvent, error) {

	// Reutiliza constructor para validar consistencia
	e, err := NewAccessEvent(userID, bookID, accessType)
	if err != nil {
		return nil, err
	}

	// Sobrescribe datos provenientes de BD
	e.id = id
	e.createdAt = createdAt

	return e, nil
}

// -------------------- Getters --------------------

// ID devuelve el ID del evento
func (e *AccessEvent) ID() uint64 { return e.id }

// UserID devuelve el ID del usuario
func (e *AccessEvent) UserID() uint64 { return e.userID }

// BookID devuelve el ID del libro
func (e *AccessEvent) BookID() uint64 { return e.bookID }

// AccessType devuelve el tipo de acceso
func (e *AccessEvent) AccessType() AccessType { return e.accessType }

// CreatedAt devuelve la fecha del evento
func (e *AccessEvent) CreatedAt() time.Time { return e.createdAt }
