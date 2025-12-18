package domain // Paquete dominio: entidades y reglas de negocio (POO)

import (
	"fmt"     // Para construir errores con contexto (fmt.Errorf)
	"strings" // Para limpiar y normalizar textos (TrimSpace, ToLower, etc.)
	"time"    // Para timestamps (createdAt/updatedAt)
)

// User representa la entidad Usuario en el dominio.
// POO + Encapsulación: los campos son privados (minúscula) y solo se accede por métodos.
type User struct {
	id        uint64    // Identificador único (normalmente lo asigna la BD)
	name      string    // Nombre del usuario
	email     string    // Email del usuario (normalizado)
	role      Role      // Rol (tipo del dominio)
	active    bool      // Estado lógico (activo/inactivo)
	createdAt time.Time // Fecha creación
	updatedAt time.Time // Fecha actualización
}

// NewUser es el constructor del dominio.
// Crea un usuario válido aplicando reglas de validación mediante setters.
// Nota: activa al usuario por defecto.
func NewUser(name, email string, role Role) (*User, error) {

	// Crea el objeto con valores por defecto
	u := &User{active: true}

	// Usa setters para validar y asignar cada campo (encapsulación).
	if err := u.SetName(name); err != nil {
		return nil, err
	}
	if err := u.SetEmail(email); err != nil {
		return nil, err
	}
	if err := u.SetRole(role); err != nil {
		return nil, err
	}

	// createdAt se marca al momento de creación en memoria.
	// En BD, puede ser reemplazado por timestamps reales al guardar.
	u.createdAt = time.Now()

	return u, nil
}

// HydrateUser es una "factory" para reconstruir un User desde la base de datos.
// Se usa porque la BD ya tiene id, createdAt, updatedAt y active.
// IMPORTANTE: aún aplica validaciones (llama NewUser y setters), garantizando consistencia del dominio.
func HydrateUser(
	id uint64,
	name, email string,
	role Role,
	active bool,
	createdAt, updatedAt time.Time,
) (*User, error) {

	// Reutiliza el constructor para validar name/email/role.
	u, err := NewUser(name, email, role)
	if err != nil {
		return nil, err
	}

	// Sobrescribe datos que vienen desde BD (persistencia).
	u.id = id
	u.active = active
	u.createdAt = createdAt
	u.updatedAt = updatedAt

	return u, nil
}

// -------------------- Getters (encapsulación) --------------------

// ID devuelve el id del usuario (campo privado).
func (u *User) ID() uint64 { return u.id }

// Name devuelve el nombre.
func (u *User) Name() string { return u.name }

// Email devuelve el email.
func (u *User) Email() string { return u.email }

// Role devuelve el rol.
func (u *User) Role() Role { return u.role }

// Active devuelve si el usuario está activo.
func (u *User) Active() bool { return u.active }

// CreatedAt devuelve fecha creación.
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt devuelve fecha actualización.
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// -------------------- Setters (encapsulación + validación) --------------------

// SetName valida y asigna el nombre.
// Regla: nombre mínimo 2 caracteres.
func (u *User) SetName(name string) error {

	// Limpia espacios
	name = strings.TrimSpace(name)

	// Validación mínima
	if len(name) < 2 {
		// Envuelve ErrValidation para que errors.Is funcione (y HTTP lo traduzca bien).
		return fmt.Errorf("%w: name must have at least 2 characters", ErrValidation)
	}

	// Asigna valor
	u.name = name

	// Actualiza timestamp de modificación
	u.updatedAt = time.Now()

	return nil
}

// SetEmail normaliza y valida email.
// Normaliza: trim + lower.
// Validación: contiene "@" y longitud mínima.
func (u *User) SetEmail(email string) error {

	// Normalización: elimina espacios y pasa a minúsculas
	email = strings.TrimSpace(strings.ToLower(email))

	// Validación simple (podría ser más estricta, pero cumple el objetivo académico)
	if !strings.Contains(email, "@") || len(email) < 5 {
		return fmt.Errorf("%w: invalid email", ErrValidation)
	}

	// Asigna y actualiza timestamp
	u.email = email
	u.updatedAt = time.Now()

	return nil
}

// SetRole valida el rol.
// Normaliza a mayúsculas para permitir entradas como "admin" o "Admin".
func (u *User) SetRole(role Role) error {

	// Normaliza el rol a mayúsculas
	role = Role(strings.ToUpper(string(role)))

	// Valida contra roles permitidos (método del tipo Role)
	if !role.IsValid() {
		// ErrInvalidRole permite que la capa HTTP responda 400 de forma consistente
		return fmt.Errorf("%w: %s", ErrInvalidRole, role)
	}

	// Asigna y actualiza timestamp
	u.role = role
	u.updatedAt = time.Now()

	return nil
}

// Deactivate desactiva el usuario (método de comportamiento, POO).
func (u *User) Deactivate() {
	u.active = false         // Cambia estado
	u.updatedAt = time.Now() // Marca actualización
}

// Activate activa el usuario (método de comportamiento, POO).
func (u *User) Activate() {
	u.active = true          // Cambia estado
	u.updatedAt = time.Now() // Marca actualización
}
