package usecase

import (
    "context"
    "sync"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type memUserRepo struct {
    mu sync.Mutex
    next uint64
    byID map[uint64]*domain.User
    byEmail map[string]uint64
}

func newMemUserRepo() *memUserRepo {
    return &memUserRepo{
        next: 1,
        byID: map[uint64]*domain.User{},
        byEmail: map[string]uint64{},
    }
}

func (r *memUserRepo) Create(ctx context.Context, u *domain.User) (uint64, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    if _, ok := r.byEmail[u.Email()]; ok { return 0, domain.ErrDuplicate }
    id := r.next; r.next++
    hu, _ := domain.HydrateUser(id, u.Name(), u.Email(), u.Role(), u.Active(), u.CreatedAt(), u.UpdatedAt())
    r.byID[id] = hu
    r.byEmail[u.Email()] = id
    return id, nil
}

func (r *memUserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    u, ok := r.byID[id]
    if !ok { return nil, domain.ErrNotFound }
    return u, nil
}

func (r *memUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    id, ok := r.byEmail[email]
    if !ok { return nil, domain.ErrNotFound }
    return r.byID[id], nil
}

func (r *memUserRepo) List(ctx context.Context) ([]*domain.User, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    out := make([]*domain.User, 0, len(r.byID))
    for _, u := range r.byID { out = append(out, u) }
    return out, nil
}

func (r *memUserRepo) Update(ctx context.Context, u *domain.User) error {
    r.mu.Lock(); defer r.mu.Unlock()
    if _, ok := r.byID[u.ID()]; !ok { return domain.ErrNotFound }
    r.byID[u.ID()] = u
    r.byEmail[u.Email()] = u.ID()
    return nil
}

func (r *memUserRepo) Delete(ctx context.Context, id uint64) error {
    r.mu.Lock(); defer r.mu.Unlock()
    u, ok := r.byID[id]
    if !ok { return domain.ErrNotFound }
    delete(r.byEmail, u.Email())
    delete(r.byID, id)
    return nil
}

type memBookRepo struct {
    mu sync.Mutex
    next uint64
    byID map[uint64]*domain.Book
    byISBN map[string]uint64
}

func newMemBookRepo() *memBookRepo {
    return &memBookRepo{
        next: 1,
        byID: map[uint64]*domain.Book{},
        byISBN: map[string]uint64{},
    }
}

func (r *memBookRepo) Create(ctx context.Context, b *domain.Book) (uint64, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    if _, ok := r.byISBN[b.ISBN()]; ok { return 0, domain.ErrDuplicate }
    id := r.next; r.next++
    hb, _ := domain.HydrateBook(id, b.Title(), b.Author(), b.Year(), b.ISBN(), b.Category(), domain.JoinTags(b.Tags()), b.Description(), b.Active(), b.CreatedAt(), b.UpdatedAt())
    r.byID[id] = hb
    r.byISBN[b.ISBN()] = id
    return id, nil
}

func (r *memBookRepo) GetByID(ctx context.Context, id uint64) (*domain.Book, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    b, ok := r.byID[id]
    if !ok { return nil, domain.ErrNotFound }
    return b, nil
}

func (r *memBookRepo) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    id, ok := r.byISBN[isbn]
    if !ok { return nil, domain.ErrNotFound }
    return r.byID[id], nil
}

func (r *memBookRepo) List(ctx context.Context) ([]*domain.Book, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    out := make([]*domain.Book, 0, len(r.byID))
    for _, b := range r.byID { out = append(out, b) }
    return out, nil
}

func (r *memBookRepo) Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error) {
    // implement naive search demonstrating maps/slices
    r.mu.Lock(); defer r.mu.Unlock()
    out := []*domain.Book{}
    q := stringsLower(f.Q)
    a := stringsLower(f.Author)
    c := stringsLower(f.Category)
    for _, b := range r.byID {
        if q != "" && !(contains(stringsLower(b.Title()), q) || contains(stringsLower(b.Author()), q)) {
            continue
        }
        if a != "" && !contains(stringsLower(b.Author()), a) { continue }
        if c != "" && !contains(stringsLower(b.Category()), c) { continue }
        out = append(out, b)
    }
    return out, nil
}

func (r *memBookRepo) Update(ctx context.Context, b *domain.Book) error {
    r.mu.Lock(); defer r.mu.Unlock()
    if _, ok := r.byID[b.ID()]; !ok { return domain.ErrNotFound }
    r.byID[b.ID()] = b
    r.byISBN[b.ISBN()] = b.ID()
    return nil
}

func (r *memBookRepo) Delete(ctx context.Context, id uint64) error {
    r.mu.Lock(); defer r.mu.Unlock()
    b, ok := r.byID[id]
    if !ok { return domain.ErrNotFound }
    delete(r.byISBN, b.ISBN())
    delete(r.byID, id)
    return nil
}

type memAccessRepo struct {
    mu sync.Mutex
    events []*domain.AccessEvent
}

func newMemAccessRepo() *memAccessRepo { return &memAccessRepo{} }

func (r *memAccessRepo) Create(ctx context.Context, e *domain.AccessEvent) (uint64, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    r.events = append(r.events, e)
    return uint64(len(r.events)), nil
}

func (r *memAccessRepo) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    stats := map[domain.AccessType]int{
        domain.AccessApertura: 0,
        domain.AccessLectura: 0,
        domain.AccessDescarga: 0,
    }
    for _, e := range r.events {
        if e.BookID() == bookID {
            stats[e.AccessType()]++
        }
    }
    return stats, nil
}

// tiny helpers to avoid importing strings in multiple files
func stringsLower(s string) string {
    b := []byte(s)
    for i := range b {
        if b[i] >= 'A' && b[i] <= 'Z' {
            b[i] = b[i] + ('a' - 'A')
        }
    }
    return string(b)
}
func contains(s, sub string) bool {
    if sub == "" { return true }
    // simple contains
    for i := 0; i+len(sub) <= len(s); i++ {
        if s[i:i+len(sub)] == sub { return true }
    }
    return false
}
