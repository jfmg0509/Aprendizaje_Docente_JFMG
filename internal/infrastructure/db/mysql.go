package db // Capa de infraestructura: acceso a base de datos

import (
	"context"      // Contextos para timeout y cancelación
	"database/sql" // API estándar de SQL en Go
	"time"         // Manejo de tiempos

	_ "github.com/go-sql-driver/mysql" // Driver MySQL (import anónimo)
)

// DB es un wrapper alrededor de *sql.DB.
// Permite centralizar la conexión y extenderla en el futuro.
type DB struct {
	SQL *sql.DB // Conexión principal a la base de datos
}

// Open abre una conexión a MySQL usando el DSN recibido.
// Además:
// - Configura el pool de conexiones
// - Verifica conectividad con Ping
func Open(dsn string) (*DB, error) {

	// Inicializa la conexión con el driver mysql
	// NOTA: sql.Open NO abre realmente la conexión todavía
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// ---------------- POOL DE CONEXIONES ----------------

	// Tiempo máximo que una conexión puede vivir antes de reciclarse
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Máximo de conexiones abiertas simultáneamente
	sqlDB.SetMaxOpenConns(25)

	// Máximo de conexiones inactivas en el pool
	sqlDB.SetMaxIdleConns(10)

	// ---------------- HEALTH CHECK ----------------

	// Contexto con timeout para evitar bloqueos
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ping fuerza una conexión real a la base
	// Si falla aquí, la app no debe continuar
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close() // cierre defensivo
		return nil, err
	}

	// Si todo salió bien, retornamos nuestro wrapper DB
	return &DB{SQL: sqlDB}, nil
}
