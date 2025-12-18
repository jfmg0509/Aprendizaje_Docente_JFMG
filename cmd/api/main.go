package main // Paquete principal, punto de entrada del programa

import (
	"html/template" // Manejo de templates HTML para el frontend
	"log"           // Registro de logs
	"net/http"      // Servidor HTTP
	"os"            // Acceso a variables del sistema operativo
	"os/signal"     // Captura señales del sistema (Ctrl+C)
	"syscall"       // Señales del sistema (SIGINT, SIGTERM)
	"time"          // Manejo de tiempos y timeouts

	// Handlers y servidor HTTP
	apphttp "github.com/jfmg0509/sistema_libros_funcional_go/internal/transport/http"

	// Configuración (variables de entorno)
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/config"

	// Base de datos y repositorios MySQL
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/db"

	// Casos de uso y lógica de negocio
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

func main() {

	// Carga la configuración desde .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err) // Detiene el programa si falla
	}

	// Abre la conexión a MySQL
	database, err := db.Open(cfg.DSN())
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.SQL.Close() // Cierra la BD al finalizar el programa

	// Creación de repositorios (implementan interfaces)
	userRepo := db.NewMySQLUserRepo(database.SQL)
	bookRepo := db.NewMySQLBookRepo(database.SQL)
	accessRepo := db.NewMySQLAccessRepo(database.SQL)

	// Cola concurrente para registrar accesos (canales + goroutines)
	accessQueue := usecase.NewAccessQueue(
		accessRepo,
		cfg.AccessQueueSize,
		cfg.AccessWorkers,
	)
	defer accessQueue.Close() // Cierre ordenado de workers

	// Servicios de negocio
	userSvc := usecase.NewUserService(userRepo)
	bookSvc := usecase.NewBookService(
		bookRepo,
		userRepo,
		accessRepo,
		accessQueue,
	)

	// Carga de templates HTML
	tpl := template.Must(
		template.ParseGlob("web/templates/*.html"),
	)

	// Inicialización de handlers HTTP
	handlers := apphttp.NewHandlers(userSvc, bookSvc, tpl)

	// Creación del servidor con router (Gorilla Mux)
	server := apphttp.NewServer(handlers)

	// Configuración del servidor HTTP
	httpSrv := &http.Server{
		Addr:              cfg.Addr,        // Dirección y puerto
		Handler:           server.Router,   // Router HTTP
		ReadHeaderTimeout: 5 * time.Second, // Timeout de seguridad
	}

	// Arranque del servidor en una goroutine
	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := httpSrv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// Canal para capturar señales del sistema
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Bloquea el main hasta recibir Ctrl+C
	<-sig

	// Apagado controlado del servidor
	log.Println("shutting down...")
	_ = httpSrv.Close()
}
