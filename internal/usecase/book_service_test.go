package usecase

import (
    "context"
    "testing"
    "time"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

func TestBookServiceRecordAccessQueued(t *testing.T) {
    users := newMemUserRepo()
    books := newMemBookRepo()
    access := newMemAccessRepo()

    // seed
    u, _ := domain.NewUser("Juan", "juan@example.com", domain.RoleReader)
    uid, _ := users.Create(context.Background(), u)
    b, _ := domain.NewBook("Go POO", "Autor", 2024, "ISBN-1", "Programación", []string{"go","poo"}, "")
    bid, _ := books.Create(context.Background(), b)

    q := NewAccessQueue(access, 10, 1)
    defer q.Close()

    svc := NewBookService(books, users, access, q)

    if err := svc.RecordAccess(context.Background(), uid, bid, domain.AccessLectura); err != nil {
        t.Fatalf("record: %v", err)
    }

    // esperar a que el worker procese
    time.Sleep(60 * time.Millisecond)

    stats, _ := svc.StatsByBook(context.Background(), bid)
    if stats[domain.AccessLectura] != 1 {
        t.Fatalf("expected 1 lectura, got %d", stats[domain.AccessLectura])
    }
}
