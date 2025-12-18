package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    apphttp "github.com/jfmg0509/sistema_libros_funcional_go/internal/transport/http"
    "github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/config"
    "github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/db"
    "github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config: %v", err)
    }

    database, err := db.Open(cfg.DSN())
    if err != nil {
        log.Fatalf("db: %v", err)
    }
    defer database.SQL.Close()

    userRepo := db.NewMySQLUserRepo(database.SQL)
    bookRepo := db.NewMySQLBookRepo(database.SQL)
    accessRepo := db.NewMySQLAccessRepo(database.SQL)

    // goroutines + canales para registrar accesos sin bloquear request
    accessQueue := usecase.NewAccessQueue(accessRepo, cfg.AccessQueueSize, cfg.AccessWorkers)
    defer accessQueue.Close()

    userSvc := usecase.NewUserService(userRepo)
    bookSvc := usecase.NewBookService(bookRepo, userRepo, accessRepo, accessQueue)

    tpl := template.Must(template.ParseGlob("web/templates/*.html"))
    handlers := apphttp.NewHandlers(userSvc, bookSvc, tpl)
    server := apphttp.NewServer(handlers)

    httpSrv := &http.Server{
        Addr:              cfg.Addr,
        Handler:           server.Router,
        ReadHeaderTimeout: 5 * time.Second,
    }

    go func() {
        log.Printf("listening on %s", cfg.Addr)
        if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server: %v", err)
        }
    }()

    // Shutdown graceful
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig

    log.Println("shutting down...")
    _ = httpSrv.Close()
}
