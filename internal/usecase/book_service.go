package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

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

func (s *BookService) Create(ctx context.Context, title, author string, year int, isbn, category string, tags []string, description string) (*domain.Book, error) {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	isbn = strings.TrimSpace(isbn)

	if title == "" || author == "" || isbn == "" {
		return nil, errors.New("title, author e isbn son obligatorios")
	}

	// evita ISBN duplicado si tu repo lo soporta
	if existing, _ := s.books.GetByISBN(ctx, isbn); existing != nil {
		return nil, errors.New("ya existe un libro con ese ISBN")
	}

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

func (s *BookService) Delete(ctx context.Context, id uint64) error {
	return s.books.Delete(ctx, id)
}

// Update aplica cambios usando setters reales del dominio.
func (s *BookService) Update(ctx context.Context, id uint64, in UpdateBookInput) (*domain.Book, error) {
	b, err := s.books.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Solo aplica lo que venga en PATCH
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
		// opcional: validar duplicado si cambias ISBN
		if existing, _ := s.books.GetByISBN(ctx, *in.ISBN); existing != nil && existing.ID() != id {
			return nil, errors.New("ya existe un libro con ese ISBN")
		}
		if err := b.SetISBN(*in.ISBN); err != nil {
			return nil, err
		}
	}
	if in.Category != nil {
		if err := b.SetCategory(*in.Category); err != nil {
			return nil, err
		}
	}
	if in.Tags != nil {
		b.SetTags(*in.Tags)
	}
	if in.Description != nil {
		b.SetDescription(*in.Description)
	}
	if in.Active != nil {
		if *in.Active {
			b.Activate()
		} else {
			b.Deactivate()
		}
	}

	// Persistir cambios
	if err := s.books.Update(ctx, b); err != nil {
		return nil, err
	}

	return s.books.GetByID(ctx, id)
}

func (s *BookService) RecordAccess(ctx context.Context, userID, bookID uint64, t domain.AccessType) error {
	if userID == 0 || bookID == 0 {
		return errors.New("user_id y book_id son obligatorios")
	}

	// valida existencia
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return err
	}
	if _, err := s.books.GetByID(ctx, bookID); err != nil {
		return err
	}

	e, err := domain.NewAccessEvent(userID, bookID, t)
	if err != nil {
		return err
	}

	// Si hay cola -> async
	if s.queue != nil {
		s.queue.Enqueue(e)
		return nil
	}

	// Directo a repo
	_, err = s.access.Create(ctx, e)
	return err
}

func (s *BookService) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
	return s.access.StatsByBook(ctx, bookID)
}

// ===== helpers opcionales (si te sirven en algún punto) =====

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func now() time.Time { return time.Now() }
