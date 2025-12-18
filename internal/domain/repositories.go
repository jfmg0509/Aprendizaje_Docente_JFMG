package domain

import "context"

// Interfaces (requerimiento)
type UserRepository interface {
    Create(ctx context.Context, u *User) (uint64, error)
    GetByID(ctx context.Context, id uint64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    List(ctx context.Context) ([]*User, error)
    Update(ctx context.Context, u *User) error
    Delete(ctx context.Context, id uint64) error
}

type BookFilter struct {
    Q        string
    Author   string
    Category string
}

type BookRepository interface {
    Create(ctx context.Context, b *Book) (uint64, error)
    GetByID(ctx context.Context, id uint64) (*Book, error)
    GetByISBN(ctx context.Context, isbn string) (*Book, error)
    List(ctx context.Context) ([]*Book, error)
    Search(ctx context.Context, f BookFilter) ([]*Book, error)
    Update(ctx context.Context, b *Book) error
    Delete(ctx context.Context, id uint64) error
}

type AccessLogRepository interface {
    Create(ctx context.Context, e *AccessEvent) (uint64, error)
    StatsByBook(ctx context.Context, bookID uint64) (map[AccessType]int, error)
}
