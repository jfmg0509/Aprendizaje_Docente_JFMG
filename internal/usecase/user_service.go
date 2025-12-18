package usecase // Capa de casos de uso (lógica de negocio)

import (
	"context" // Permite cancelación y timeouts desde handlers
	"fmt"     // Para envolver errores con contexto
	"time"    // Para simular trabajo / mostrar timeouts

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain" // Dominio: User, errores, interfaces
)

// UserService contiene la lógica de negocio relacionada con usuarios.
// No depende de MySQL directo: depende de la interfaz domain.UserRepository (POO + SOLID).
type UserService struct {
	repo domain.UserRepository // Repositorio (interfaz) para persistencia de usuarios
}

// NewUserService es el "constructor" del servicio.
// Inyecta el repositorio para desacoplar la lógica de negocio de la BD.
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo} // Retorna servicio listo para usar
}

// Create crea un usuario aplicando reglas de negocio:
// - Validación mediante constructor del dominio
// - Regla: email único
// - Persistencia mediante repo
func (s *UserService) Create(ctx context.Context, name, email string, role domain.Role) (*domain.User, error) {

	// Crea entidad User usando un constructor del dominio (valida campos).
	u, err := domain.NewUser(name, email, role)
	if err != nil {
		return nil, err // Si falla validación, se retorna error del dominio
	}

	// Regla del negocio: el email debe ser único.
	// Si encontramos un usuario con el mismo email, retornamos ErrDuplicate.
	if _, err := s.repo.GetByEmail(ctx, u.Email()); err == nil {
		return nil, domain.ErrDuplicate
	}

	// Guarda el usuario en la BD (por medio del repositorio).
	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err // Error de persistencia
	}

	// Recarga el usuario para obtener timestamps reales de BD (created_at/updated_at).
	return s.repo.GetByID(ctx, id)
}

// List devuelve todos los usuarios.
// No hay reglas extra aquí: delega al repositorio.
func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

// Get devuelve un usuario por ID.
// Delegación directa al repositorio.
func (s *UserService) Get(ctx context.Context, id uint64) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update actualiza un usuario por ID.
// Aplica cambios usando setters del dominio (encapsulación) y luego persiste.
func (s *UserService) Update(
	ctx context.Context,
	id uint64,
	name, email string,
	role domain.Role,
	active *bool, // puntero para permitir "opcional": nil => no cambiar
) (*domain.User, error) {

	// Obtiene el usuario actual desde BD.
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Si name viene vacío, no se actualiza.
	// Si viene con valor, se usa setter (valida internamente).
	if name != "" {
		if err := u.SetName(name); err != nil {
			return nil, err
		}
	}

	// Si email viene vacío, no se actualiza.
	// Si viene con valor, se valida con SetEmail.
	if email != "" {
		if err := u.SetEmail(email); err != nil {
			return nil, err
		}
	}

	// Si role viene vacío, no se actualiza.
	// Si viene con valor, se valida con SetRole.
	if role != "" {
		if err := u.SetRole(role); err != nil {
			return nil, err
		}
	}

	// active es puntero: si es nil, no se cambia.
	// Si tiene valor, se activa o desactiva el usuario.
	if active != nil {
		if *active {
			u.Activate()
		} else {
			u.Deactivate()
		}
	}

	// Simulación de trabajo (sirve para demostrar uso de context/timeouts).
	// No cambia estado del usuario; solo es para “consumir tiempo” si se usara.
	_ = time.Now()

	// Persistimos cambios.
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	// Recargamos de BD para devolver estado final actualizado.
	return s.repo.GetByID(ctx, id)
}

// Delete elimina un usuario por ID.
// Envuelve el error con contexto (mejor trazabilidad).
func (s *UserService) Delete(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		// fmt.Errorf con %w permite "wrap" del error (errors.Is seguirá funcionando).
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
