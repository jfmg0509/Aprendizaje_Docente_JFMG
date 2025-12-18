package usecase

import (
    "context"
    "fmt"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type BookService struct {
    books  domain.BookRepository
    users  domain.UserRepository
    access domain.AccessLogRepository
    queue  *AccessQueue
}

func NewBookService(books domain.BookRepository, users domain.UserRepository, access domain.AccessLogRepository, queue *AccessQueue) *BookService {
    return &BookService{books: books, users: users, access: access, queue: queue}
}

func (s *BookService) Create(ctx context.Context, title, author string, year int, isbn, category string, tags []string, description string) (*domain.Book, error) {
    b, err := domain.NewBook(title, author, year, isbn, category, tags, description)
    if err != nil { return nil, err }

    // regla: ISBN único
    if _, err := s.books.GetByISBN(ctx, b.ISBN()); err == nil {
        return nil, domain.ErrDuplicate
    }
    id, err := s.books.Create(ctx, b)
    if err != nil { return nil, err }
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

func (s *BookService) Update(ctx context.Context, id uint64, in UpdateBookInput) (*domain.Book, error) {
    b, err := s.books.GetByID(ctx, id)
    if err != nil { return nil, err }

    if in.Title != nil {
        if err := b.SetTitle(*in.Title); err != nil { return nil, err }
    }
    if in.Author != nil {
        if err := b.SetAuthor(*in.Author); err != nil { return nil, err }
    }
    if in.Year != nil {
        if err := b.SetYear(*in.Year); err != nil { return nil, err }
    }
    if in.ISBN != nil {
        if err := b.SetISBN(*in.ISBN); err != nil { return nil, err }
    }
    if in.Category != nil {
        if err := b.SetCategory(*in.Category); err != nil { return nil, err }
    }
    if in.Tags != nil {
        b.SetTags(*in.Tags)
    }
    if in.Description != nil {
        b.SetDescription(*in.Description)
    }
    if in.Active != nil {
        if *in.Active { b.Activate() } else { b.Deactivate() }
    }

    if err := s.books.Update(ctx, b); err != nil { return nil, err }
    return s.books.GetByID(ctx, id)
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

func (s *BookService) Delete(ctx context.Context, id uint64) error {
    return s.books.Delete(ctx, id)
}

func (s *BookService) RecordAccess(ctx context.Context, userID, bookID uint64, t domain.AccessType) error {
    // validar existencia + activo
    u, err := s.users.GetByID(ctx, userID)
    if err != nil { return fmt.Errorf("user: %w", err) }
    if !u.Active() { return domain.ErrInactiveEntity }

    b, err := s.books.GetByID(ctx, bookID)
    if err != nil { return fmt.Errorf("book: %w", err) }
    if !b.Active() { return domain.ErrInactiveEntity }

    e, err := domain.NewAccessEvent(userID, bookID, t)
    if err != nil { return err }

    // enqueue async
    if s.queue != nil {
        ok := s.queue.Enqueue(e)
        if !ok {
            // fallback sync
            _, err := s.access.Create(ctx, e)
            return err
        }
        return nil
    }
    _, err = s.access.Create(ctx, e)
    return err
}

func (s *BookService) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
    // map requerido en la consigna
    return s.access.StatsByBook(ctx, bookID)
}
