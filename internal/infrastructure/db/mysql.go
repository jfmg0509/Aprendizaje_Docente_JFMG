package db

import (
    "context"
    "database/sql"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type DB struct {
    SQL *sql.DB
}

func Open(dsn string) (*DB, error) {
    sqlDB, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := sqlDB.PingContext(ctx); err != nil {
        _ = sqlDB.Close()
        return nil, err
    }
    return &DB{SQL: sqlDB}, nil
}
