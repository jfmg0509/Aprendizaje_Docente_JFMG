package usecase

import (
    "context"
    "fmt"
    "time"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type UserService struct {
    repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, name, email string, role domain.Role) (*domain.User, error) {
    u, err := domain.NewUser(name, email, role)
    if err != nil { return nil, err }

    // regla: email único
    if _, err := s.repo.GetByEmail(ctx, u.Email()); err == nil {
        return nil, domain.ErrDuplicate
    }
    id, err := s.repo.Create(ctx, u)
    if err != nil { return nil, err }

    // recargar para tener timestamps reales
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
    return s.repo.List(ctx)
}

func (s *UserService) Get(ctx context.Context, id uint64) (*domain.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id uint64, name, email string, role domain.Role, active *bool) (*domain.User, error) {
    u, err := s.repo.GetByID(ctx, id)
    if err != nil { return nil, err }

    if name != "" {
        if err := u.SetName(name); err != nil { return nil, err }
    }
    if email != "" {
        if err := u.SetEmail(email); err != nil { return nil, err }
    }
    if role != "" {
        if err := u.SetRole(role); err != nil { return nil, err }
    }
    if active != nil {
        if *active { u.Activate() } else { u.Deactivate() }
    }

    // simular algo de trabajo (para mostrar context/timeouts en llamadas)
    _ = time.Now()

    if err := s.repo.Update(ctx, u); err != nil { return nil, err }
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id uint64) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        return fmt.Errorf("delete user: %w", err)
    }
    return nil
}
