package domain // Paquete dominio: entidades y reglas del negocio

import (
	"fmt"     // Construcción de errores con contexto
	"strings" // Limpieza y normalización de strings
	"time"    // Manejo de fechas
)

// Book representa la entidad Libro del dominio.
// Todos los campos son privados → encapsulación real.
type Book struct {
	id          uint64    // ID único (asignado por BD)
	title       string    // Título del libro
	author      string    // Autor
	year        int       // Año de publicación
	isbn        string    // ISBN (regla: único)
	category    string    // Categoría
	tags        []string  // Etiquetas (slice)
	description string    // Descripción
	active      bool      // Estado lógico
	createdAt   time.Time // Fecha de creación
	updatedAt   time.Time // Fecha de última actualización
}

// NewBook es el constructor del dominio.
// Crea un libro válido aplicando reglas de negocio.
func NewBook(
	title, author string,
	year int,
	isbn, category string,
	tags []string,
	description string,
) (*Book, error) {

	// Libro activo por defecto
	b := &Book{active: true}

	// Validaciones mediante setters (encapsulación)
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

	// Tags y descripción no retornan error
	b.SetTags(tags)
	b.SetDescription(description)

	// Fecha de creación
	b.createdAt = time.Now()

	return b, nil
}

// HydrateBook reconstruye un Book desde la base de datos.
// Se usa al leer desde MySQL.
func HydrateBook(
	id uint64,
	title, author string,
	year int,
	isbn, category, tagsCSV, description string,
	active bool,
	createdAt, updatedAt time.Time,
) (*Book, error) {

	// Convierte CSV a slice
	tags := splitTags(tagsCSV)

	// Reutiliza constructor para validar datos
	b, err := NewBook(title, author, year, isbn, category, tags, description)
	if err != nil {
		return nil, err
	}

	// Sobrescribe datos de persistencia
	b.id = id
	b.active = active
	b.createdAt = createdAt
	b.updatedAt = updatedAt

	return b, nil
}

// -------------------- Getters --------------------

// ID devuelve el ID del libro
func (b *Book) ID() uint64 { return b.id }

// Title devuelve el título
func (b *Book) Title() string { return b.title }

// Author devuelve el autor
func (b *Book) Author() string { return b.author }

// Year devuelve el año
func (b *Book) Year() int { return b.year }

// ISBN devuelve el ISBN
func (b *Book) ISBN() string { return b.isbn }

// Category devuelve la categoría
func (b *Book) Category() string { return b.category }

// Tags devuelve una COPIA del slice (protege encapsulación)
func (b *Book) Tags() []string {
	return append([]string{}, b.tags...)
}

// Description devuelve la descripción
func (b *Book) Description() string { return b.description }

// Active indica si el libro está activo
func (b *Book) Active() bool { return b.active }

// CreatedAt devuelve fecha creación
func (b *Book) CreatedAt() time.Time { return b.createdAt }

// UpdatedAt devuelve fecha actualización
func (b *Book) UpdatedAt() time.Time { return b.updatedAt }

// -------------------- Setters --------------------

// SetTitle valida y asigna el título
func (b *Book) SetTitle(title string) error {
	title = strings.TrimSpace(title)
	if len(title) < 2 {
		return fmt.Errorf("%w: title must have at least 2 characters", ErrValidation)
	}
	b.title = title
	b.updatedAt = time.Now()
	return nil
}

// SetAuthor valida y asigna el autor
func (b *Book) SetAuthor(author string) error {
	author = strings.TrimSpace(author)
	if len(author) < 2 {
		return fmt.Errorf("%w: author must have at least 2 characters", ErrValidation)
	}
	b.author = author
	b.updatedAt = time.Now()
	return nil
}

// SetYear valida el año de publicación
func (b *Book) SetYear(year int) error {
	if year < 1400 || year > time.Now().Year()+1 {
		return fmt.Errorf("%w: invalid year", ErrValidation)
	}
	b.year = year
	b.updatedAt = time.Now()
	return nil
}

// SetISBN valida y asigna ISBN
func (b *Book) SetISBN(isbn string) error {
	isbn = strings.TrimSpace(isbn)
	if len(isbn) < 5 {
		return fmt.Errorf("%w: isbn too short", ErrValidation)
	}
	b.isbn = isbn
	b.updatedAt = time.Now()
	return nil
}

// SetCategory valida y asigna categoría
func (b *Book) SetCategory(category string) error {
	category = strings.TrimSpace(category)
	if len(category) < 2 {
		return fmt.Errorf("%w: category too short", ErrValidation)
	}
	b.category = category
	b.updatedAt = time.Now()
	return nil
}

// SetTags limpia y asigna etiquetas (slice)
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

// SetDescription asigna descripción
func (b *Book) SetDescription(desc string) {
	b.description = strings.TrimSpace(desc)
	b.updatedAt = time.Now()
}

// Deactivate desactiva el libro
func (b *Book) Deactivate() {
	b.active = false
	b.updatedAt = time.Now()
}

// Activate activa el libro
func (b *Book) Activate() {
	b.active = true
	b.updatedAt = time.Now()
}

// -------------------- Helpers --------------------

// splitTags convierte un CSV en slice
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

// JoinTags convierte un slice en CSV (útil para guardar en BD)
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
