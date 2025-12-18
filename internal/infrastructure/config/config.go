package config

import (
    "fmt"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Addr         string
    BaseURL      string
    DBUser       string
    DBPass       string
    DBHost       string
    DBPort       string
    DBName       string
    DBParams     string
    AccessQueueSize int
    AccessWorkers   int
}

func Load() (Config, error) {
    // Carga .env si existe (no falla si no existe)
    _ = godotenv.Load()

    cfg := Config{
        Addr:    getenv("APP_ADDR", ":8081"),
        BaseURL: getenv("APP_BASE_URL", "http://localhost:8081"),
        DBUser:  getenv("DB_USER", "root"),
        DBPass:  os.Getenv("DB_PASS"),
        DBHost:  getenv("DB_HOST", "127.0.0.1"),
        DBPort:  getenv("DB_PORT", "3306"),
        DBName:  getenv("DB_NAME", "libros_poo"),
        DBParams: getenv("DB_PARAMS", "parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"),
        AccessQueueSize: atoi(getenv("ACCESS_QUEUE_SIZE", "200"), 200),
        AccessWorkers:   atoi(getenv("ACCESS_WORKERS", "4"), 4),
    }
    if cfg.DBName == "" {
        return Config{}, fmt.Errorf("DB_NAME is required")
    }
    return cfg, nil
}

func (c Config) DSN() string {
    // user:pass@tcp(host:port)/dbname?params
    if c.DBPass == "" {
        return fmt.Sprintf("%s@tcp(%s:%s)/%s?%s", c.DBUser, c.DBHost, c.DBPort, c.DBName, c.DBParams)
    }
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, c.DBParams)
}

func getenv(k, def string) string {
    v := os.Getenv(k)
    if v == "" {
        return def
    }
    return v
}

func atoi(s string, def int) int {
    v, err := strconv.Atoi(s)
    if err != nil {
        return def
    }
    return v
}
