package usecase

import (
	"context"
	"sync"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// AccessRepo debe existir en tu proyecto. Si ya existe en otro archivo de usecase,
// NO lo vuelvas a declarar aquí.
// Lo usamos como dependencia en la cola.
//
// type AccessRepo interface {
//     Create(ctx context.Context, e *domain.AccessEvent) (uint64, error)
// }

type AccessQueue struct {
	repo    AccessRepo
	ch      chan *domain.AccessEvent
	wg      sync.WaitGroup
	closeMu sync.Mutex
	closed  bool
}

func NewAccessQueue(repo AccessRepo, buffer int, workers int) *AccessQueue {
	if buffer <= 0 {
		buffer = 64
	}
	if workers <= 0 {
		workers = 1
	}

	q := &AccessQueue{
		repo: repo,
		ch:   make(chan *domain.AccessEvent, buffer),
	}

	// Workers
	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			for e := range q.ch {
				if e == nil {
					continue
				}
				// Insert async. Si falla, no reventamos el worker.
				_, _ = q.repo.Create(context.Background(), e)
			}
		}()
	}

	return q
}

// Enqueue agrega el evento a la cola (si no está cerrada)
func (q *AccessQueue) Enqueue(e *domain.AccessEvent) {
	q.closeMu.Lock()
	defer q.closeMu.Unlock()

	if q.closed {
		return
	}
	q.ch <- e
}

// Close cierra la cola y espera a los workers
func (q *AccessQueue) Close() {
	q.closeMu.Lock()
	if q.closed {
		q.closeMu.Unlock()
		return
	}
	q.closed = true
	close(q.ch)
	q.closeMu.Unlock()

	q.wg.Wait()
}
