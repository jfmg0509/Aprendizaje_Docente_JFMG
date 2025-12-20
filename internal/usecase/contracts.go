package usecase

import (
	"context"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// ====== Repos contracts ======

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id uint64) error
}

type BookRepo interface {
	Create(ctx context.Context, b *domain.Book) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*domain.Book, error)
	GetByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	List(ctx context.Context) ([]*domain.Book, error)
	Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error)
	Update(ctx context.Context, b *domain.Book) error
	Delete(ctx context.Context, id uint64) error
}

type AccessRepo interface {
	Create(ctx context.Context, e *domain.AccessEvent) (uint64, error)
	StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error)
}

// ====== DTO for Update ======
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
