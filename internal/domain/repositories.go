package domain // Dominio: contratos (interfaces) y tipos de búsqueda

import "context" // Context permite cancelación/timeout desde HTTP hacia DB

// -------------------- UserRepository --------------------

// UserRepository es una interfaz (requerimiento) que define el contrato
// para persistir y consultar usuarios.
//
// Importante:
// - Está en el dominio para que la lógica de negocio dependa de abstracciones,
//   no de MySQL (inversión de dependencias).
// - Facilita testing: puedes implementar un repo fake/mock sin BD.
type UserRepository interface {

	// Create guarda un usuario y retorna el ID generado (por BD).
	Create(ctx context.Context, u *User) (uint64, error)

	// GetByID obtiene un usuario por ID.
	// Debe retornar ErrNotFound si no existe.
	GetByID(ctx context.Context, id uint64) (*User, error)

	// GetByEmail obtiene un usuario por email (regla de unicidad).
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List retorna todos los usuarios.
	List(ctx context.Context) ([]*User, error)

	// Update actualiza los datos del usuario.
	Update(ctx context.Context, u *User) error

	// Delete elimina el usuario por ID.
	// Puede ser borrado físico o lógico según implementación.
	Delete(ctx context.Context, id uint64) error
}

// -------------------- BookFilter --------------------

// BookFilter encapsula filtros de búsqueda de libros.
// Se mantiene simple para construir consultas dinámicas en repositorio.
type BookFilter struct {
	Q        string // Texto libre (título/ISBN/descripción, etc.)
	Author   string // Filtro por autor
	Category string // Filtro por categoría
}

// -------------------- BookRepository --------------------

// BookRepository define el contrato para persistir y consultar libros.
type BookRepository interface {

	// Create guarda un libro y retorna el ID generado.
	Create(ctx context.Context, b *Book) (uint64, error)

	// GetByID obtiene libro por ID.
	GetByID(ctx context.Context, id uint64) (*Book, error)

	// GetByISBN obtiene libro por ISBN (regla de unicidad).
	GetByISBN(ctx context.Context, isbn string) (*Book, error)

	// List retorna todos los libros.
	List(ctx context.Context) ([]*Book, error)

	// Search busca libros por filtros (Q, author, category).
	Search(ctx context.Context, f BookFilter) ([]*Book, error)

	// Update actualiza datos del libro.
	Update(ctx context.Context, b *Book) error

	// Delete elimina libro por ID.
	Delete(ctx context.Context, id uint64) error
}

// -------------------- AccessLogRepository --------------------

// AccessLogRepository define el contrato para guardar accesos y consultar estadísticas.
type AccessLogRepository interface {

	// Create registra un evento de acceso.
	Create(ctx context.Context, e *AccessEvent) (uint64, error)

	// StatsByBook retorna estadísticas agrupadas por tipo de acceso.
	// Ejemplo retorno: map[AccessType]int{"LECTURA": 10, "DESCARGA": 3}
	StatsByBook(ctx context.Context, bookID uint64) (map[AccessType]int, error)
}
