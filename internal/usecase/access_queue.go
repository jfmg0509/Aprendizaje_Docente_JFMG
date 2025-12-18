package usecase // Capa de casos de uso: concurrencia y procesamiento asíncrono

import (
	"context" // Contextos para timeout en operaciones de BD
	"log"     // Logging de errores
	"sync"    // Sincronización (WaitGroup)
	"time"    // Timeouts

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain" // Dominio: AccessEvent, interfaces
)

// AccessQueue implementa una cola concurrente para registrar accesos.
// Usa canales, goroutines y WaitGroup.
type AccessQueue struct {
	ch      chan *domain.AccessEvent   // Canal bufferizado de eventos
	repo    domain.AccessLogRepository // Repositorio para guardar accesos
	wg      sync.WaitGroup             // Espera a que terminen los workers
	closing chan struct{}              // Canal para señalizar cierre
}

// NewAccessQueue crea la cola y lanza los workers.
// buffer: tamaño del canal (capacidad)
// workers: número de goroutines consumidoras
func NewAccessQueue(repo domain.AccessLogRepository, buffer int, workers int) *AccessQueue {

	// Valores por defecto si vienen mal configurados
	if buffer <= 0 {
		buffer = 100
	}
	if workers <= 0 {
		workers = 2
	}

	// Inicializa la cola
	q := &AccessQueue{
		ch:      make(chan *domain.AccessEvent, buffer), // canal bufferizado
		repo:    repo,                                   // repo inyectado
		closing: make(chan struct{}),                    // señal de cierre
	}

	// Lanza N workers como goroutines
	for i := 0; i < workers; i++ {
		q.wg.Add(1)        // incrementa contador de goroutines activas
		go q.worker(i + 1) // inicia worker en goroutine
	}

	return q
}

// Enqueue intenta insertar un evento en la cola.
// Retorna:
// - true si se encoló correctamente
// - false si la cola está llena (backpressure)
func (q *AccessQueue) Enqueue(e *domain.AccessEvent) bool {
	select {
	case q.ch <- e:
		// Evento encolado con éxito
		return true
	default:
		// Canal lleno: no bloquea la request
		return false
	}
}

// worker consume eventos del canal y los guarda en BD.
// Cada worker corre en su propia goroutine.
func (q *AccessQueue) worker(n int) {
	defer q.wg.Done() // Indica que el worker terminó al salir

	for {
		select {

		// Caso 1: llega un evento al canal
		case e := <-q.ch:
			if e == nil {
				continue
			}

			// Contexto con timeout para la operación de BD
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			// Inserta el evento en la base de datos
			if _, err := q.repo.Create(ctx, e); err != nil {
				log.Printf("[access-worker-%d] insert error: %v", n, err)
			}

			// Libera recursos del contexto
			cancel()

		// Caso 2: señal de cierre del sistema
		case <-q.closing:
			// Sale del worker de forma ordenada
			return
		}
	}
}

// Close cierra la cola de forma controlada.
// Señala a los workers que deben detenerse y espera a que terminen.
func (q *AccessQueue) Close() {
	close(q.closing) // Notifica a todos los workers que deben salir
	q.wg.Wait()      // Espera a que todas las goroutines finalicen
}
