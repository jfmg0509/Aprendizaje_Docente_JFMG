package usecase // Casos de uso: lógica del negocio (no HTTP, no SQL directo)

import (
	"context" // Propaga cancelación/timeouts desde handlers
	"fmt"     // Para envolver errores con contexto (error wrapping)

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain" // Dominio: Book, User, errores, interfaces
)

// BookService contiene lógica de negocio de libros.
// Depende de interfaces (repositorios) y opcionalmente de una cola concurrente de accesos.
type BookService struct {
	books  domain.BookRepository      // Repositorio de libros (interfaz)
	users  domain.UserRepository      // Repositorio de usuarios (interfaz) para validar existencia/estado
	access domain.AccessLogRepository // Repositorio de accesos (interfaz) para registrar eventos
	queue  *AccessQueue               // Cola concurrente para registrar accesos sin bloquear
}

// NewBookService construye el servicio e inyecta sus dependencias.
// Esto desacopla el servicio de implementaciones concretas (MySQL, etc.).
func NewBookService(
	books domain.BookRepository,
	users domain.UserRepository,
	access domain.AccessLogRepository,
	queue *AccessQueue,
) *BookService {
	return &BookService{
		books:  books,
		users:  users,
		access: access,
		queue:  queue,
	}
}

// Create crea un libro aplicando reglas del dominio y del negocio.
// Regla de negocio: ISBN único.
func (s *BookService) Create(
	ctx context.Context,
	title, author string,
	year int,
	isbn, category string,
	tags []string,
	description string,
) (*domain.Book, error) {

	// Crea la entidad Book desde el dominio (valida campos internamente).
	b, err := domain.NewBook(title, author, year, isbn, category, tags, description)
	if err != nil {
		return nil, err // Error de validación del dominio
	}

	// Regla: ISBN único. Si existe un libro con ese ISBN => ErrDuplicate.
	if _, err := s.books.GetByISBN(ctx, b.ISBN()); err == nil {
		return nil, domain.ErrDuplicate
	}

	// Persiste el libro.
	id, err := s.books.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	// Recarga el libro para devolver timestamps reales desde BD.
	return s.books.GetByID(ctx, id)
}

// List devuelve todos los libros.
func (s *BookService) List(ctx context.Context) ([]*domain.Book, error) {
	return s.books.List(ctx)
}

// Get devuelve un libro por ID.
func (s *BookService) Get(ctx context.Context, id uint64) (*domain.Book, error) {
	return s.books.GetByID(ctx, id)
}

// Search busca libros usando un filtro del dominio (BookFilter).
func (s *BookService) Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error) {
	return s.books.Search(ctx, f)
}

// Update actualiza un libro por ID.
// Usa un input con punteros para permitir “campos opcionales”.
func (s *BookService) Update(ctx context.Context, id uint64, in UpdateBookInput) (*domain.Book, error) {

	// Trae el libro actual desde repo.
	b, err := s.books.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cada campo opcional se actualiza solo si viene distinto de nil.
	// Encapsulación: se usan setters del dominio, no campos directos.

	if in.Title != nil {
		if err := b.SetTitle(*in.Title); err != nil {
			return nil, err
		}
	}

	if in.Author != nil {
		if err := b.SetAuthor(*in.Author); err != nil {
			return nil, err
		}
	}

	if in.Year != nil {
		if err := b.SetYear(*in.Year); err != nil {
			return nil, err
		}
	}

	if in.ISBN != nil {
		if err := b.SetISBN(*in.ISBN); err != nil {
			return nil, err
		}
	}

	if in.Category != nil {
		if err := b.SetCategory(*in.Category); err != nil {
			return nil, err
		}
	}

	// Tags y Description podrían no requerir validación, por eso se setean directo.
	// Aun así, se mantiene encapsulación llamando métodos del objeto.
	if in.Tags != nil {
		b.SetTags(*in.Tags)
	}

	if in.Description != nil {
		b.SetDescription(*in.Description)
	}

	// Manejo de estado activo/inactivo.
	if in.Active != nil {
		if *in.Active {
			b.Activate()
		} else {
			b.Deactivate()
		}
	}

	// Persiste cambios.
	if err := s.books.Update(ctx, b); err != nil {
		return nil, err
	}

	// Recarga desde BD para devolver estado final.
	return s.books.GetByID(ctx, id)
}

// UpdateBookInput representa datos opcionales para actualizar un libro.
// Punteros => campo opcional: nil significa “no actualizar”.
type UpdateBookInput struct {
	Title       *string   // nil => no actualizar
	Author      *string   // nil => no actualizar
	Year        *int      // nil => no actualizar
	ISBN        *string   // nil => no actualizar
	Category    *string   // nil => no actualizar
	Tags        *[]string // nil => no actualizar (slice de tags)
	Description *string   // nil => no actualizar
	Active      *bool     // nil => no actualizar
}

// Delete elimina un libro por ID.
func (s *BookService) Delete(ctx context.Context, id uint64) error {
	return s.books.Delete(ctx, id)
}

// RecordAccess registra un evento de acceso (view/download/etc.).
// Aquí se aplican reglas:
// - Usuario debe existir y estar activo
// - Libro debe existir y estar activo
// - AccessType debe ser válido (dominio)
// Luego se registra el evento:
// - Preferentemente de forma asíncrona (cola con goroutines/canales)
// - Si la cola está llena, se hace fallback síncrono (directo a repo)
func (s *BookService) RecordAccess(ctx context.Context, userID, bookID uint64, t domain.AccessType) error {

	// 1) Validar que el usuario existe.
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user: %w", err) // wrap para mantener causa original
	}

	// 2) Validar que el usuario esté activo.
	if !u.Active() {
		return domain.ErrInactiveEntity
	}

	// 3) Validar que el libro existe.
	b, err := s.books.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("book: %w", err)
	}

	// 4) Validar que el libro esté activo.
	if !b.Active() {
		return domain.ErrInactiveEntity
	}

	// 5) Crear el evento desde el dominio (valida accessType, etc.).
	e, err := domain.NewAccessEvent(userID, bookID, t)
	if err != nil {
		return err
	}

	// 6) Encolar asíncronamente para no bloquear el request.
	// Si queue existe, se intenta Enqueue.
	if s.queue != nil {
		ok := s.queue.Enqueue(e)

		// Si la cola está llena o cerrada (ok=false), hacemos fallback síncrono.
		if !ok {
			_, err := s.access.Create(ctx, e) // persistencia directa
			return err
		}

		// Si encoló bien, devolvemos nil.
		return nil
	}

	// Si no hay cola configurada, registra directo en repo.
	_, err = s.access.Create(ctx, e)
	return err
}

// StatsByBook devuelve estadísticas por libro.
// Retorna map[AccessType]int para cumplir el requisito de maps.
func (s *BookService) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
	// Delegación al repositorio de accesos: hace agregación/consulta.
	return s.access.StatsByBook(ctx, bookID)
}
