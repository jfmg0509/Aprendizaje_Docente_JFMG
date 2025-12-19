package http

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

func TestAPICreateUser(t *testing.T) {
	ur := newMemUserRepo()
	br := newMemBookRepo()
	ar := newMemAccessRepo()

	q := usecase.NewAccessQueue(ar, 10, 1)
	defer q.Close()

	userSvc := usecase.NewUserService(ur)
	bookSvc := usecase.NewBookService(br, ur, ar, q)

	// Renderer “dummy” para tests:
	// No renderizamos UI en este test, pero NewHandler lo requiere.
	r := &Renderer{t: template.New("x")}

	h := NewHandler(userSvc, bookSvc, r)
	router := NewRouter(h)

	body, _ := json.Marshal(map[string]any{
		"name":  "Ana",
		"email": "ana@example.com",
		"role":  "ADMIN",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rec.Code, rec.Body.String())
	}
}

// --- minimal in-memory repos for handler tests ---

type memUserRepo struct {
	mu      sync.Mutex
	next    uint64
	byID    map[uint64]*domain.User
	byEmail map[string]uint64
}

func newMemUserRepo() *memUserRepo {
	return &memUserRepo{next: 1, byID: map[uint64]*domain.User{}, byEmail: map[string]uint64{}}
}

func (r *memUserRepo) Create(ctx context.Context, u *domain.User) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byEmail[u.Email()]; ok {
		return 0, domain.ErrDuplicate
	}
	id := r.next
	r.next++

	hu, _ := domain.HydrateUser(id, u.Name(), u.Email(), u.Role(), u.Active(), u.CreatedAt(), u.UpdatedAt())
	r.byID[id] = hu
	r.byEmail[u.Email()] = id
	return id, nil
}

func (r *memUserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (r *memUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.byEmail[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return r.byID[id], nil
}

func (r *memUserRepo) List(ctx context.Context) ([]*domain.User, error) { return nil, nil }
func (r *memUserRepo) Update(ctx context.Context, u *domain.User) error { return nil }
func (r *memUserRepo) Delete(ctx context.Context, id uint64) error      { return nil }

type memBookRepo struct{}

func newMemBookRepo() *memBookRepo { return &memBookRepo{} }

func (r *memBookRepo) Create(ctx context.Context, b *domain.Book) (uint64, error) { return 1, nil }
func (r *memBookRepo) GetByID(ctx context.Context, id uint64) (*domain.Book, error) {
	return nil, domain.ErrNotFound
}
func (r *memBookRepo) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	return nil, domain.ErrNotFound
}
func (r *memBookRepo) List(ctx context.Context) ([]*domain.Book, error) { return nil, nil }
func (r *memBookRepo) Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error) {
	return nil, nil
}
func (r *memBookRepo) Update(ctx context.Context, b *domain.Book) error { return nil }
func (r *memBookRepo) Delete(ctx context.Context, id uint64) error      { return nil }

type memAccessRepo struct{}

func newMemAccessRepo() *memAccessRepo { return &memAccessRepo{} }

func (r *memAccessRepo) Create(ctx context.Context, e *domain.AccessEvent) (uint64, error) {
	return 1, nil
}
func (r *memAccessRepo) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
	return map[domain.AccessType]int{}, nil
}
