package domain

import (
    "fmt"
    "strings"
    "time"
)

type Book struct {
    id          uint64
    title       string
    author      string
    year        int
    isbn        string
    category    string
    tags        []string
    description string
    active      bool
    createdAt   time.Time
    updatedAt   time.Time
}

func NewBook(title, author string, year int, isbn, category string, tags []string, description string) (*Book, error) {
    b := &Book{active: true}
    if err := b.SetTitle(title); err != nil {
        return nil, err
    }
    if err := b.SetAuthor(author); err != nil {
        return nil, err
    }
    if err := b.SetYear(year); err != nil {
        return nil, err
    }
    if err := b.SetISBN(isbn); err != nil {
        return nil, err
    }
    if err := b.SetCategory(category); err != nil {
        return nil, err
    }
    b.SetTags(tags)
    b.SetDescription(description)
    b.createdAt = time.Now()
    return b, nil
}

func HydrateBook(id uint64, title, author string, year int, isbn, category, tagsCSV, description string, active bool, createdAt, updatedAt time.Time) (*Book, error) {
    tags := splitTags(tagsCSV)
    b, err := NewBook(title, author, year, isbn, category, tags, description)
    if err != nil {
        return nil, err
    }
    b.id = id
    b.active = active
    b.createdAt = createdAt
    b.updatedAt = updatedAt
    return b, nil
}

// Getters
func (b *Book) ID() uint64 { return b.id }
func (b *Book) Title() string { return b.title }
func (b *Book) Author() string { return b.author }
func (b *Book) Year() int { return b.year }
func (b *Book) ISBN() string { return b.isbn }
func (b *Book) Category() string { return b.category }
func (b *Book) Tags() []string { return append([]string{}, b.tags...) } // slice copy
func (b *Book) Description() string { return b.description }
func (b *Book) Active() bool { return b.active }
func (b *Book) CreatedAt() time.Time { return b.createdAt }
func (b *Book) UpdatedAt() time.Time { return b.updatedAt }

// Setters
func (b *Book) SetTitle(title string) error {
    title = strings.TrimSpace(title)
    if len(title) < 2 {
        return fmt.Errorf("%w: title must have at least 2 characters", ErrValidation)
    }
    b.title = title
    b.updatedAt = time.Now()
    return nil
}

func (b *Book) SetAuthor(author string) error {
    author = strings.TrimSpace(author)
    if len(author) < 2 {
        return fmt.Errorf("%w: author must have at least 2 characters", ErrValidation)
    }
    b.author = author
    b.updatedAt = time.Now()
    return nil
}

func (b *Book) SetYear(year int) error {
    if year < 1400 || year > time.Now().Year()+1 {
        return fmt.Errorf("%w: invalid year", ErrValidation)
    }
    b.year = year
    b.updatedAt = time.Now()
    return nil
}

func (b *Book) SetISBN(isbn string) error {
    isbn = strings.TrimSpace(isbn)
    if len(isbn) < 5 {
        return fmt.Errorf("%w: isbn too short", ErrValidation)
    }
    b.isbn = isbn
    b.updatedAt = time.Now()
    return nil
}

func (b *Book) SetCategory(category string) error {
    category = strings.TrimSpace(category)
    if len(category) < 2 {
        return fmt.Errorf("%w: category too short", ErrValidation)
    }
    b.category = category
    b.updatedAt = time.Now()
    return nil
}

func (b *Book) SetTags(tags []string) {
    cleaned := make([]string, 0, len(tags))
    for _, t := range tags {
        t = strings.TrimSpace(t)
        if t != "" {
            cleaned = append(cleaned, t)
        }
    }
    b.tags = cleaned
    b.updatedAt = time.Now()
}

func (b *Book) SetDescription(desc string) {
    b.description = strings.TrimSpace(desc)
    b.updatedAt = time.Now()
}

func (b *Book) Deactivate() { b.active = false; b.updatedAt = time.Now() }
func (b *Book) Activate() { b.active = true; b.updatedAt = time.Now() }

func splitTags(csv string) []string {
    csv = strings.TrimSpace(csv)
    if csv == "" {
        return nil
    }
    parts := strings.Split(csv, ",")
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p != "" {
            out = append(out, p)
        }
    }
    return out
}

func JoinTags(tags []string) string {
    if len(tags) == 0 {
        return ""
    }
    cleaned := make([]string, 0, len(tags))
    for _, t := range tags {
        t = strings.TrimSpace(t)
        if t != "" {
            cleaned = append(cleaned, t)
        }
    }
    return strings.Join(cleaned, ",")
}
