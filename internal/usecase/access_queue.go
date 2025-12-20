package usecase

import (
	"context"
	"sync"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

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

	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			for e := range q.ch {
				if e == nil {
					continue
				}
				// best-effort async
				_, _ = q.repo.Create(context.Background(), e)
			}
		}()
	}

	return q
}

// ✅ TryEnqueue NO bloquea el UI: si está llena o cerrada -> devuelve false
func (q *AccessQueue) TryEnqueue(ctx context.Context, e *domain.AccessEvent) bool {
	if e == nil {
		return false
	}

	// 1) ver si está cerrada (lock corto)
	q.closeMu.Lock()
	closed := q.closed
	q.closeMu.Unlock()
	if closed {
		return false
	}

	// 2) envío no bloqueante
	select {
	case q.ch <- e:
		return true
	default:
		// cola llena
		return false
	}
}

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
