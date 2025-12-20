package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

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

// ==============================
// Libros
// ==============================

func (s *BookService) Create(ctx context.Context, title, author string, year int, isbn, category string, tags []string, description string) (*domain.Book, error) {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	isbn = strings.TrimSpace(isbn)

	if title == "" || author == "" || isbn == "" {
		return nil, errors.New("title, author e isbn son obligatorios")
	}

	b, err := domain.NewBook(title, author, year, isbn, category, tags, description)
	if err != nil {
		return nil, err
	}

	id, err := s.books.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	// Para no depender de setters, devolvemos el libro consultándolo
	return s.books.Get(ctx, id)
}

func (s *BookService) List(ctx context.Context) ([]*domain.Book, error) {
	return s.books.List(ctx)
}

func (s *BookService) Get(ctx context.Context, id uint64) (*domain.Book, error) {
	return s.books.Get(ctx, id)
}

func (s *BookService) Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error) {
	return s.books.Search(ctx, f)
}

func (s *BookService) Delete(ctx context.Context, id uint64) error {
	return s.books.Delete(ctx, id)
}

// UpdateBookInput ya existe en tu proyecto (lo usa handlers.go).
// Aquí lo usamos tal cual.
func (s *BookService) Update(ctx context.Context, id uint64, in UpdateBookInput) (*domain.Book, error) {
	// Obtengo el actual
	b, err := s.books.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Aplicar cambios con métodos que ya tengas en domain.Book.
	// Si tu domain.Book NO tiene estos setters, dime y lo ajusto a tu modelo real.
	if in.Title != nil {
		b.SetTitle(*in.Title)
	}
	if in.Author != nil {
		b.SetAuthor(*in.Author)
	}
	if in.Year != nil {
		b.SetYear(*in.Year)
	}
	if in.ISBN != nil {
		b.SetISBN(*in.ISBN)
	}
	if in.Category != nil {
		b.SetCategory(*in.Category)
	}
	if in.Tags != nil {
		b.SetTags(*in.Tags)
	}
	if in.Description != nil {
		b.SetDescription(*in.Description)
	}
	if in.Active != nil {
		b.SetActive(*in.Active)
	}

	if err := s.books.Update(ctx, b); err != nil {
		return nil, err
	}
	return s.books.Get(ctx, id)
}

// ==============================
// Accesos (ESTE ERA EL BUG)
// ==============================

func (s *BookService) RecordAccess(ctx context.Context, userID, bookID uint64, t domain.AccessType) error {
	if userID == 0 || bookID == 0 {
		return errors.New("user_id y book_id son obligatorios")
	}

	// (Opcional) Validar existencia
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return err
	}
	if _, err := s.books.Get(ctx, bookID); err != nil {
		return err
	}

	// ✅ IMPORTANTE: tu domain.NewAccessEvent (según tu error) acepta SOLO 3 args
	e, err := domain.NewAccessEvent(userID, bookID, t)
	if err != nil {
		return err
	}

	// ✅ Si hay cola, SOLO ENCOLA (NO insertar directo)
	if s.queue != nil {
		s.queue.Enqueue(e)
		return nil
	}

	// Sin cola => directo
	_, err = s.access.Create(ctx, e)
	return err
}

func (s *BookService) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
	return s.access.StatsByBook(ctx, bookID)
}
