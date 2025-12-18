package db // Infraestructura DB: repositorios MySQL

import (
	"context"      // Para timeouts/cancelación desde handlers -> DB
	"database/sql" // Driver SQL estándar
	"errors"       // Para comparar errores (errors.Is)
	"strings"      // Para normalizar email
	"time"         // Para timestamps leídos de BD

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain" // Dominio (entidad + errores)
)

// MySQLUserRepo es la implementación concreta de domain.UserRepository usando MySQL.
// Guarda internamente un *sql.DB para ejecutar consultas.
type MySQLUserRepo struct{ db *sql.DB }

// NewMySQLUserRepo es el "constructor" del repositorio.
// Inyecta la dependencia db.
func NewMySQLUserRepo(db *sql.DB) *MySQLUserRepo { return &MySQLUserRepo{db: db} }

// Create inserta un usuario en la tabla users.
// Retorna el ID generado por MySQL.
func (r *MySQLUserRepo) Create(ctx context.Context, u *domain.User) (uint64, error) {

	// Inserta usuario. Nota: usamos getters (encapsulación).
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (name,email,role,active) VALUES (?,?,?,?)`,
		u.Name(),               // nombre validado por dominio
		u.Email(),              // email validado por dominio
		string(u.Role()),       // rol como string
		boolToTiny(u.Active()), // bool -> 0/1 para MySQL
	)

	// Manejo de error: si es duplicado, lo traducimos al error del dominio
	if err != nil {
		if isMySQLDuplicate(err) { // helper que detecta error "duplicate key"
			return 0, domain.ErrDuplicate
		}
		return 0, err
	}

	// Obtiene ID autoincremental generado
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

// GetByID trae un usuario por ID.
// Si no existe, retorna domain.ErrNotFound.
func (r *MySQLUserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {

	// QueryRowContext retorna 1 fila máximo (o error)
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id,name,email,role,active,created_at,COALESCE(updated_at,created_at)
		 FROM users WHERE id=?`,
		id,
	)

	// Variables donde se escanean columnas
	var (
		rid                  uint64
		name, email, role    string
		active               int
		createdAt, updatedAt time.Time
	)

	// Scan copia los valores del row a las variables
	if err := row.Scan(&rid, &name, &email, &role, &active, &createdAt, &updatedAt); err != nil {

		// Si no hay filas, se traduce a ErrNotFound
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		// Cualquier otro error se retorna tal cual
		return nil, err
	}

	// HydrateUser reconstruye la entidad respetando validaciones y encapsulación
	return domain.HydrateUser(
		rid,
		name,
		email,
		domain.Role(role),
		active == 1,
		createdAt,
		updatedAt,
	)
}

// GetByEmail trae usuario por email.
func (r *MySQLUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	// Normaliza email para evitar diferencias por mayúsculas/espacios
	email = strings.ToLower(strings.TrimSpace(email))

	row := r.db.QueryRowContext(
		ctx,
		`SELECT id,name,email,role,active,created_at,COALESCE(updated_at,created_at)
		 FROM users WHERE email=?`,
		email,
	)

	var (
		rid                  uint64
		name, em, role       string
		active               int
		createdAt, updatedAt time.Time
	)

	if err := row.Scan(&rid, &name, &em, &role, &active, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return domain.HydrateUser(
		rid,
		name,
		em,
		domain.Role(role),
		active == 1,
		createdAt,
		updatedAt,
	)
}

// List retorna todos los usuarios ordenados descendentemente por ID.
func (r *MySQLUserRepo) List(ctx context.Context) ([]*domain.User, error) {

	// QueryContext devuelve múltiples filas
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id,name,email,role,active,created_at,COALESCE(updated_at,created_at)
		 FROM users ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // cierre de recursos

	// Slice de salida (lista de usuarios)
	out := []*domain.User{}

	// Itera fila por fila
	for rows.Next() {
		var (
			rid                  uint64
			name, em, role       string
			active               int
			createdAt, updatedAt time.Time
		)

		// Scan por cada fila
		if err := rows.Scan(&rid, &name, &em, &role, &active, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		// Reconstruye entidad dominio
		u, err := domain.HydrateUser(rid, name, em, domain.Role(role), active == 1, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		// Agrega al slice
		out = append(out, u)
	}

	return out, nil
}

// Update actualiza un usuario existente.
func (r *MySQLUserRepo) Update(ctx context.Context, u *domain.User) error {

	_, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET name=?, email=?, role=?, active=? WHERE id=?`,
		u.Name(),
		u.Email(),
		string(u.Role()),
		boolToTiny(u.Active()),
		u.ID(),
	)

	// Si falla por duplicado (email unique), traduce a ErrDuplicate
	if err != nil {
		if isMySQLDuplicate(err) {
			return domain.ErrDuplicate
		}
		return err
	}

	return nil
}

// Delete elimina un usuario por ID.
// Si no afectó filas, significa que no existía.
func (r *MySQLUserRepo) Delete(ctx context.Context, id uint64) error {

	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, id)
	if err != nil {
		return err
	}

	// RowsAffected dice cuántas filas fueron eliminadas
	n, _ := res.RowsAffected()

	// Si no eliminó nada, no existía
	if n == 0 {
		return domain.ErrNotFound
	}

	return nil
}
