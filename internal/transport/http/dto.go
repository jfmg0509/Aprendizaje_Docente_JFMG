package http // Paquete de transporte HTTP (DTOs para API)

// Importaciones
import (
	"time" // time.Time para fechas de creación/actualización

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain" // Entidades del dominio (User, Book, Role)
)

// -------------------- USERS DTO --------------------

// UserDTO es el objeto que se envía/recibe por la API (JSON).
// Se usa para NO exponer directamente la entidad domain.User y mantener encapsulación.
// Nota: aquí los campos son públicos (Mayúscula) para que encoding/json pueda serializarlos.
type UserDTO struct {
	ID        uint64      `json:"id"`         // Identificador del usuario
	Name      string      `json:"name"`       // Nombre
	Email     string      `json:"email"`      // Correo
	Role      domain.Role `json:"role"`       // Rol (tipo del dominio)
	Active    bool        `json:"active"`     // Estado lógico (activo/inactivo)
	CreatedAt time.Time   `json:"created_at"` // Fecha de creación
	UpdatedAt time.Time   `json:"updated_at"` // Fecha de última actualización
}

// userToDTO convierte una entidad domain.User (dominio) a UserDTO (capa HTTP).
// Observa que NO accede a campos directamente; usa getters: u.ID(), u.Name()...
// Eso es encapsulación en acción.
func userToDTO(u *domain.User) UserDTO {
	return UserDTO{
		ID:        u.ID(),        // Getter del dominio
		Name:      u.Name(),      // Getter del dominio
		Email:     u.Email(),     // Getter del dominio
		Role:      u.Role(),      // Getter del dominio
		Active:    u.Active(),    // Getter del dominio
		CreatedAt: u.CreatedAt(), // Getter del dominio
		UpdatedAt: u.UpdatedAt(), // Getter del dominio
	}
}

// usersToDTO convierte un slice de usuarios del dominio a slice de DTO.
// Usa make con capacidad = len(list) para ser más eficiente al hacer append.
func usersToDTO(list []*domain.User) []UserDTO {
	out := make([]UserDTO, 0, len(list)) // slice vacío con capacidad reservada
	for _, u := range list {             // recorre cada usuario del dominio
		out = append(out, userToDTO(u)) // agrega el DTO convertido
	}
	return out // devuelve slice de DTOs listo para JSON
}

// -------------------- BOOKS DTO --------------------

// BookDTO es el objeto que se envía/recibe por la API (JSON) para libros.
// Igual que UserDTO: separa dominio de transporte.
type BookDTO struct {
	ID          uint64    `json:"id"`          // ID del libro
	Title       string    `json:"title"`       // Título
	Author      string    `json:"author"`      // Autor
	Year        int       `json:"year"`        // Año de publicación
	ISBN        string    `json:"isbn"`        // Código ISBN
	Category    string    `json:"category"`    // Categoría
	Tags        []string  `json:"tags"`        // Etiquetas (slice)
	Description string    `json:"description"` // Descripción
	Active      bool      `json:"active"`      // Estado lógico
	CreatedAt   time.Time `json:"created_at"`  // Fecha de creación
	UpdatedAt   time.Time `json:"updated_at"`  // Fecha de actualización
}

// bookToDTO convierte la entidad domain.Book a BookDTO.
// Mantiene encapsulación usando getters del dominio.
func bookToDTO(b *domain.Book) BookDTO {
	return BookDTO{
		ID:          b.ID(),          // Getter
		Title:       b.Title(),       // Getter
		Author:      b.Author(),      // Getter
		Year:        b.Year(),        // Getter
		ISBN:        b.ISBN(),        // Getter
		Category:    b.Category(),    // Getter
		Tags:        b.Tags(),        // Getter (slice)
		Description: b.Description(), // Getter
		Active:      b.Active(),      // Getter
		CreatedAt:   b.CreatedAt(),   // Getter
		UpdatedAt:   b.UpdatedAt(),   // Getter
	}
}

// booksToDTO convierte un slice de libros del dominio a un slice de BookDTO.
// Mismo patrón eficiente: make con capacidad + append.
func booksToDTO(list []*domain.Book) []BookDTO {
	out := make([]BookDTO, 0, len(list)) // reserva capacidad
	for _, b := range list {             // itera libros del dominio
		out = append(out, bookToDTO(b)) // convierte y agrega
	}
	return out // devuelve DTOs listos para writeJSON
}
