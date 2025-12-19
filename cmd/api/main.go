package main

import (
	"log"
	"net/http"
	"time"

	apphttp "github.com/jfmg0509/sistema_libros_funcional_go/internal/transport/http"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/config"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/db"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

func main() {
	// =========================
	// 1) Cargar Config
	// =========================
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// =========================
	// 2) Abrir DB (MySQL)
	// =========================
	database, err := db.Open(cfg.DSN())
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer database.SQL.Close()

	// =========================
	// 3) Repositorios (MySQL)
	// =========================
	userRepo := db.NewMySQLUserRepo(database.SQL)
	bookRepo := db.NewMySQLBookRepo(database.SQL)
	accessRepo := db.NewMySQLAccessRepo(database.SQL)

	// =========================
	// 4) Cola de Accesos (goroutines + channel)
	// =========================
	queue := usecase.NewAccessQueue(accessRepo, cfg.AccessQueueSize, cfg.AccessWorkers)
	defer queue.Close()

	// =========================
	// 5) Servicios (Usecases)
	// =========================
	userService := usecase.NewUserService(userRepo)
	bookService := usecase.NewBookService(bookRepo, userRepo, accessRepo, queue)

	// =========================
	// 6) Renderer (HTML templates)
	// =========================
	// Tus templates están en: web/templates/*.html
	renderer, err := apphttp.NewRenderer("web/templates")
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	// =========================
	// 7) Handler único (API + UI)
	// =========================
	h := apphttp.NewHandler(userService, bookService, renderer)

	// =========================
	// 8) Router
	// =========================
	router := apphttp.NewRouter(h)

	// =========================
	// 9) HTTP Server
	// =========================
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
