package domain // Dominio: tipos base (enumeraciones) y reglas de validación

import (
	"strings" // Para normalizar strings (ToUpper)
)

// Role representa el rol de un usuario en el sistema.
// Se implementa como "enum" en Go: type + const.
type Role string

// Constantes de roles permitidos (valores cerrados).
const (
	RoleAdmin     Role = "ADMIN"     // Administrador
	RoleReader    Role = "READER"    // Lector
	RoleConsultor Role = "CONSULTOR" // Consultor
)

// AllowedRoles es un ARRAY (requerimiento) con el conjunto fijo de roles permitidos.
// Array => tamaño fijo (3) porque este catálogo no debería cambiar en runtime.
var AllowedRoles = [3]Role{RoleAdmin, RoleReader, RoleConsultor}

// IsValid valida si un Role está dentro del catálogo permitido.
// Normaliza a mayúsculas para admitir entradas como "admin", "Admin", etc.
func (r Role) IsValid() bool {

	// Normalización: compara siempre en mayúsculas
	rr := strings.ToUpper(string(r))

	// Recorre el array de roles permitidos
	for _, allowed := range AllowedRoles {
		if rr == string(allowed) {
			return true // es válido
		}
	}

	// No coincide con ningún rol permitido
	return false
}

// -------------------------------------------------------------

// AccessType representa el tipo de acceso/acción que un usuario realiza sobre un libro.
type AccessType string

// Constantes del catálogo de accesos permitidos.
const (
	AccessApertura AccessType = "APERTURA" // Cuando se abre/ingresa al detalle
	AccessLectura  AccessType = "LECTURA"  // Cuando se lee/consulta
	AccessDescarga AccessType = "DESCARGA" // Cuando se descarga
)

// AllowedAccessTypes es un ARRAY (requerimiento) con el conjunto fijo de accesos permitidos.
var AllowedAccessTypes = [3]AccessType{AccessApertura, AccessLectura, AccessDescarga}

// IsValid valida si un AccessType pertenece al catálogo permitido.
func (a AccessType) IsValid() bool {

	// Normalización: mayúsculas para comparar
	aa := strings.ToUpper(string(a))

	// Recorre el array permitido
	for _, allowed := range AllowedAccessTypes {
		if aa == string(allowed) {
			return true // válido
		}
	}

	// No pertenece al catálogo
	return false
}
