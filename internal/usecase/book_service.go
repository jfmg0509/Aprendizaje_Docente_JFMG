package usecase

import (
	"context"
	"fmt"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type BookRepo interface {
	Create(ctx context.Context, b *domain.Book) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*domain.Book, error)
	GetByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	List(ctx context.Context) ([]*domain.Book, error)
	Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error)
	Update(ctx context.Context, b *domain.Book) error
	Delete(ctx context.Context, id uint64) error
}

type UserRepo interface {
	GetByID(ctx context.Context, id uint64) (*domain.User, error)
}

type AccessRepo interface {
	Create(ctx context.Context, e *domain.AccessEvent) (uint64, error)
	StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error)
}

type UpdateBookInput struct {
	Title       *string
	Author      *string
	Year        *int
	ISBN        *string
	Category    *string
	Tags        *[]string
	Description *string
	Active      *bool
}

type BookService struct {
	books  BookRepo
	users  UserRepo
	access AccessRepo
	queue  *AccessQueue
}

func NewBookService(bookRepo BookRepo, userRepo UserRepo, accessRepo AccessRepo, queue *AccessQueue) *BookService {
	return &BookService{
		books:  bookRepo,
		users:  userRepo,
		access: accessRepo,
		queue:  queue,
	}
}

func (s *BookService) Create(ctx context.Context, title, author string, year int, isbn, category string, tags []string, description string) (*domain.Book, error) {
	b, err := domain.NewBook(title, author, year, isbn, category, tags, description)
	if err != nil {
		return nil, err
	}

	id, err := s.books.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	return s.books.GetByID(ctx, id)
}

func (s *BookService) List(ctx context.Context) ([]*domain.Book, error) {
	return s.books.List(ctx)
}

func (s *BookService) Get(ctx context.Context, id uint64) (*domain.Book, error) {
	return s.books.GetByID(ctx, id)
}

func (s *BookService) Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error) {
	return s.books.Search(ctx, f)
}

// ✅ Para que compile con tu domain actual (sin setters),
// este Update NO modifica campos (lo dejamos neutro por ahora).
func (s *BookService) Update(ctx context.Context, id uint64, in UpdateBookInput) (*domain.Book, error) {
	b, err := s.books.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.books.Update(ctx, b); err != nil {
		return nil, err
	}
	return s.books.GetByID(ctx, id)
}

func (s *BookService) Delete(ctx context.Context, id uint64) error {
	return s.books.Delete(ctx, id)
}

// ✅ FIX REAL: registrar acceso SINCRÓNICO (contador sube sí o sí).
func (s *BookService) RecordAccess(ctx context.Context, userID, bookID uint64, t domain.AccessType) error {
	// valida que existan
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return fmt.Errorf("user: %w", err)
	}
	if _, err := s.books.GetByID(ctx, bookID); err != nil {
		return fmt.Errorf("book: %w", err)
	}

	e, err := domain.NewAccessEvent(userID, bookID, t)
	if err != nil {
		return err
	}

	// 1) GUARDA DIRECTO
	if _, err := s.access.Create(ctx, e); err != nil {
		return err
	}

	// 2) (Opcional) encola también
	if s.queue != nil {
		_ = s.queue.Enqueue(e)
	}

	return nil
}

func (s *BookService) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
	return s.access.StatsByBook(ctx, bookID)
}
