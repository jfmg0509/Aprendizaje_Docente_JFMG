package domain // Dominio: definición de errores del negocio

import "errors"

// Errores del dominio.
// Se definen una sola vez y se reutilizan en todo el sistema.
//
// Estos errores permiten:
// - Comparaciones seguras con errors.Is
// - Traducción correcta a códigos HTTP
// - Evitar strings hardcodeados
var (

	// ErrNotFound se usa cuando una entidad no existe en el sistema
	// Ejemplo: buscar un usuario por ID inexistente
	ErrNotFound = errors.New("not found")

	// ErrDuplicate se usa cuando se viola una restricción de unicidad
	// Ejemplo: email de usuario duplicado, ISBN duplicado
	ErrDuplicate = errors.New("duplicate")

	// ErrValidation se usa cuando los datos no cumplen reglas del dominio
	// Ejemplo: nombre corto, email inválido, año fuera de rango
	ErrValidation = errors.New("validation error")

	// ErrInvalidRole se usa cuando el rol de un usuario no es permitido
	// Ejemplo: rol distinto de ADMIN, USER, etc.
	ErrInvalidRole = errors.New("invalid role")

	// ErrInvalidAccess se usa cuando el tipo de acceso no es válido
	// Ejemplo: accessType distinto de VIEW, DOWNLOAD, etc.
	ErrInvalidAccess = errors.New("invalid access type")

	// ErrInactiveEntity se usa cuando una entidad existe pero está inactiva
	// Ejemplo: usuario o libro desactivado
	ErrInactiveEntity = errors.New("entity is inactive")
)
