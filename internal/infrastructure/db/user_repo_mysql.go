package db

import (
    "context"
    "database/sql"
    "errors"
    "strings"
    "time"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type MySQLUserRepo struct{ db *sql.DB }

func NewMySQLUserRepo(db *sql.DB) *MySQLUserRepo { return &MySQLUserRepo{db: db} }

func (r *MySQLUserRepo) Create(ctx context.Context, u *domain.User) (uint64, error) {
    // detectar duplicado por email (unique index)
    res, err := r.db.ExecContext(ctx,
        `INSERT INTO users (name,email,role,active) VALUES (?,?,?,?)`,
        u.Name(), u.Email(), string(u.Role()), boolToTiny(u.Active()),
    )
    if err != nil {
        if isMySQLDuplicate(err) {
            return 0, domain.ErrDuplicate
        }
        return 0, err
    }
    id, _ := res.LastInsertId()
    return uint64(id), nil
}

func (r *MySQLUserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
    row := r.db.QueryRowContext(ctx, `SELECT id,name,email,role,active,created_at,COALESCE(updated_at,created_at) FROM users WHERE id=?`, id)
    var (
        rid uint64
        name, email, role string
        active int
        createdAt, updatedAt time.Time
    )
    if err := row.Scan(&rid, &name, &email, &role, &active, &createdAt, &updatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return domain.HydrateUser(rid, name, email, domain.Role(role), active == 1, createdAt, updatedAt)
}

func (r *MySQLUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    email = strings.ToLower(strings.TrimSpace(email))
    row := r.db.QueryRowContext(ctx, `SELECT id,name,email,role,active,created_at,COALESCE(updated_at,created_at) FROM users WHERE email=?`, email)
    var (
        rid uint64
        name, em, role string
        active int
        createdAt, updatedAt time.Time
    )
    if err := row.Scan(&rid, &name, &em, &role, &active, &createdAt, &updatedAt); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return domain.HydrateUser(rid, name, em, domain.Role(role), active == 1, createdAt, updatedAt)
}

func (r *MySQLUserRepo) List(ctx context.Context) ([]*domain.User, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id,name,email,role,active,created_at,COALESCE(updated_at,created_at) FROM users ORDER BY id DESC`)
    if err != nil { return nil, err }
    defer rows.Close()

    out := []*domain.User{}
    for rows.Next() {
        var (
            rid uint64
            name, em, role string
            active int
            createdAt, updatedAt time.Time
        )
        if err := rows.Scan(&rid, &name, &em, &role, &active, &createdAt, &updatedAt); err != nil {
            return nil, err
        }
        u, err := domain.HydrateUser(rid, name, em, domain.Role(role), active == 1, createdAt, updatedAt)
        if err != nil { return nil, err }
        out = append(out, u)
    }
    return out, nil
}

func (r *MySQLUserRepo) Update(ctx context.Context, u *domain.User) error {
    _, err := r.db.ExecContext(ctx, `UPDATE users SET name=?, email=?, role=?, active=? WHERE id=?`,
        u.Name(), u.Email(), string(u.Role()), boolToTiny(u.Active()), u.ID(),
    )
    if err != nil {
        if isMySQLDuplicate(err) {
            return domain.ErrDuplicate
        }
        return err
    }
    return nil
}

func (r *MySQLUserRepo) Delete(ctx context.Context, id uint64) error {
    res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, id)
    if err != nil { return err }
    n, _ := res.RowsAffected()
    if n == 0 { return domain.ErrNotFound }
    return nil
}
