package usecase

import (
    "context"
    "testing"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

func TestUserServiceCreateAndGet(t *testing.T) {
    repo := newMemUserRepo()
    svc := NewUserService(repo)

    u, err := svc.Create(context.Background(), "Ana", "ana@example.com", domain.RoleAdmin)
    if err != nil { t.Fatalf("create: %v", err) }
    if u.ID() == 0 { t.Fatalf("expected id") }

    got, err := svc.Get(context.Background(), u.ID())
    if err != nil { t.Fatalf("get: %v", err) }
    if got.Email() != "ana@example.com" { t.Fatalf("email mismatch") }
}

func TestUserValidation(t *testing.T) {
    repo := newMemUserRepo()
    svc := NewUserService(repo)

    if _, err := svc.Create(context.Background(), "A", "bad", domain.RoleAdmin); err == nil {
        t.Fatalf("expected validation error")
    }
}
