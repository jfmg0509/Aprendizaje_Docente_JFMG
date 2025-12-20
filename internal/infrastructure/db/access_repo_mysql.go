package db

import (
	"context"
	"database/sql"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

type MySQLAccessRepo struct{ db *sql.DB }

func NewMySQLAccessRepo(db *sql.DB) *MySQLAccessRepo { return &MySQLAccessRepo{db: db} }

func (r *MySQLAccessRepo) Create(ctx context.Context, e *domain.AccessEvent) (uint64, error) {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO access_events (user_id, book_id, access_type) VALUES (?,?,?)`,
		e.UserID(), e.BookID(), string(e.AccessType()),
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *MySQLAccessRepo) StatsByBook(ctx context.Context, bookID uint64) (map[domain.AccessType]int, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT access_type, COUNT(*) 
         FROM access_events 
         WHERE book_id=? 
         GROUP BY access_type`,
		bookID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Siempre devolvemos las 3 llaves, aunque no existan filas aún
	stats := map[domain.AccessType]int{
		domain.AccessApertura: 0,
		domain.AccessLectura:  0,
		domain.AccessDescarga: 0,
	}

	for rows.Next() {
		var t string
		var c int
		if err := rows.Scan(&t, &c); err != nil {
			return nil, err
		}

		// Normaliza: solo contamos los tipos válidos
		switch domain.AccessType(t) {
		case domain.AccessApertura, domain.AccessLectura, domain.AccessDescarga:
			stats[domain.AccessType(t)] = c
		default:
			// si en BD hay un tipo raro, lo ignoramos para no romper la UI
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
