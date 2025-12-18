package db

import (
    "context"
    "database/sql"
    "errors"
    "strings"
    "time"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type MySQLBookRepo struct{ db *sql.DB }

func NewMySQLBookRepo(db *sql.DB) *MySQLBookRepo { return &MySQLBookRepo{db: db} }

func (r *MySQLBookRepo) Create(ctx context.Context, b *domain.Book) (uint64, error) {
    res, err := r.db.ExecContext(ctx,
        `INSERT INTO books (title,author,year,isbn,category,tags,description,active) VALUES (?,?,?,?,?,?,?,?)`,
        b.Title(), b.Author(), b.Year(), b.ISBN(), b.Category(), domain.JoinTags(b.Tags()), b.Description(), boolToTiny(b.Active()),
    )
    if err != nil {
        if isMySQLDuplicate(err) { return 0, domain.ErrDuplicate }
        return 0, err
    }
    id, _ := res.LastInsertId()
    return uint64(id), nil
}

func (r *MySQLBookRepo) GetByID(ctx context.Context, id uint64) (*domain.Book, error) {
    row := r.db.QueryRowContext(ctx, `SELECT id,title,author,year,isbn,category,tags,COALESCE(description,''),active,created_at,COALESCE(updated_at,created_at) FROM books WHERE id=?`, id)
    var (
        rid uint64
        title, author, isbn, category, tags, desc string
        year int
        active int
        createdAt, updatedAt time.Time
    )
    if err := row.Scan(&rid, &title, &author, &year, &isbn, &category, &tags, &desc, &active, &createdAt, &updatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
        return nil, err
    }
    return domain.HydrateBook(rid, title, author, year, isbn, category, tags, desc, active==1, createdAt, updatedAt)
}

func (r *MySQLBookRepo) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
    isbn = strings.TrimSpace(isbn)
    row := r.db.QueryRowContext(ctx, `SELECT id,title,author,year,isbn,category,tags,COALESCE(description,''),active,created_at,COALESCE(updated_at,created_at) FROM books WHERE isbn=?`, isbn)
    var (
        rid uint64
        title, author, is, category, tags, desc string
        year int
        active int
        createdAt, updatedAt time.Time
    )
    if err := row.Scan(&rid, &title, &author, &year, &is, &category, &tags, &desc, &active, &createdAt, &updatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
        return nil, err
    }
    return domain.HydrateBook(rid, title, author, year, is, category, tags, desc, active==1, createdAt, updatedAt)
}

func (r *MySQLBookRepo) List(ctx context.Context) ([]*domain.Book, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id,title,author,year,isbn,category,tags,COALESCE(description,''),active,created_at,COALESCE(updated_at,created_at) FROM books ORDER BY id DESC`)
    if err != nil { return nil, err }
    defer rows.Close()

    out := []*domain.Book{}
    for rows.Next() {
        var (
            rid uint64
            title, author, isbn, category, tags, desc string
            year int
            active int
            createdAt, updatedAt time.Time
        )
        if err := rows.Scan(&rid, &title, &author, &year, &isbn, &category, &tags, &desc, &active, &createdAt, &updatedAt); err != nil {
            return nil, err
        }
        b, err := domain.HydrateBook(rid, title, author, year, isbn, category, tags, desc, active==1, createdAt, updatedAt)
        if err != nil { return nil, err }
        out = append(out, b)
    }
    return out, nil
}

func (r *MySQLBookRepo) Search(ctx context.Context, f domain.BookFilter) ([]*domain.Book, error) {
    // Construcción simple con filtros (slices para params)
    where := []string{"1=1"}
    args := []any{}

    if q := strings.TrimSpace(f.Q); q != "" {
        where = append(where, "(LOWER(title) LIKE ? OR LOWER(author) LIKE ? OR LOWER(tags) LIKE ?)")
        like := "%" + strings.ToLower(q) + "%"
        args = append(args, like, like, like)
    }
    if a := strings.TrimSpace(f.Author); a != "" {
        where = append(where, "LOWER(author) LIKE ?")
        args = append(args, "%"+strings.ToLower(a)+"%")
    }
    if c := strings.TrimSpace(f.Category); c != "" {
        where = append(where, "LOWER(category) LIKE ?")
        args = append(args, "%"+strings.ToLower(c)+"%")
    }

    query := `SELECT id,title,author,year,isbn,category,tags,COALESCE(description,''),active,created_at,COALESCE(updated_at,created_at)
              FROM books WHERE ` + strings.Join(where, " AND ") + ` ORDER BY id DESC LIMIT 200`
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil { return nil, err }
    defer rows.Close()

    out := []*domain.Book{}
    for rows.Next() {
        var (
            rid uint64
            title, author, isbn, category, tags, desc string
            year int
            active int
            createdAt, updatedAt time.Time
        )
        if err := rows.Scan(&rid, &title, &author, &year, &isbn, &category, &tags, &desc, &active, &createdAt, &updatedAt); err != nil {
            return nil, err
        }
        b, err := domain.HydrateBook(rid, title, author, year, isbn, category, tags, desc, active==1, createdAt, updatedAt)
        if err != nil { return nil, err }
        out = append(out, b)
    }
    return out, nil
}

func (r *MySQLBookRepo) Update(ctx context.Context, b *domain.Book) error {
    _, err := r.db.ExecContext(ctx,
        `UPDATE books SET title=?,author=?,year=?,isbn=?,category=?,tags=?,description=?,active=? WHERE id=?`,
        b.Title(), b.Author(), b.Year(), b.ISBN(), b.Category(), domain.JoinTags(b.Tags()), b.Description(), boolToTiny(b.Active()), b.ID(),
    )
    if err != nil {
        if isMySQLDuplicate(err) { return domain.ErrDuplicate }
        return err
    }
    return nil
}

func (r *MySQLBookRepo) Delete(ctx context.Context, id uint64) error {
    res, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id=?`, id)
    if err != nil { return err }
    n, _ := res.RowsAffected()
    if n == 0 { return domain.ErrNotFound }
    return nil
}
