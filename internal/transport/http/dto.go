package http

import (
    "time"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// DTOs para no exponer directamente campos privados (encapsulación)
type UserDTO struct {
    ID uint64 `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
    Role domain.Role `json:"role"`
    Active bool `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func userToDTO(u *domain.User) UserDTO {
    return UserDTO{
        ID: u.ID(),
        Name: u.Name(),
        Email: u.Email(),
        Role: u.Role(),
        Active: u.Active(),
        CreatedAt: u.CreatedAt(),
        UpdatedAt: u.UpdatedAt(),
    }
}

func usersToDTO(list []*domain.User) []UserDTO {
    out := make([]UserDTO, 0, len(list))
    for _, u := range list {
        out = append(out, userToDTO(u))
    }
    return out
}

type BookDTO struct {
    ID uint64 `json:"id"`
    Title string `json:"title"`
    Author string `json:"author"`
    Year int `json:"year"`
    ISBN string `json:"isbn"`
    Category string `json:"category"`
    Tags []string `json:"tags"`
    Description string `json:"description"`
    Active bool `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func bookToDTO(b *domain.Book) BookDTO {
    return BookDTO{
        ID: b.ID(),
        Title: b.Title(),
        Author: b.Author(),
        Year: b.Year(),
        ISBN: b.ISBN(),
        Category: b.Category(),
        Tags: b.Tags(),
        Description: b.Description(),
        Active: b.Active(),
        CreatedAt: b.CreatedAt(),
        UpdatedAt: b.UpdatedAt(),
    }
}

func booksToDTO(list []*domain.Book) []BookDTO {
    out := make([]BookDTO, 0, len(list))
    for _, b := range list {
        out = append(out, bookToDTO(b))
    }
    return out
}
