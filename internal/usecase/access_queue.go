package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

var ErrQueueClosed = errors.New("access queue closed")

// ✅ IMPORTANTE: este interface NO debe llamarse AccessRepo porque
// ya existe otro AccessRepo en book_service.go con más métodos.
type AccessWriter interface {
	Create(ctx context.Context, e *domain.AccessEvent) (uint64, error)
}

// AccessQueue procesa accesos en background.
type AccessQueue struct {
	repo    AccessWriter
	ch      chan *domain.AccessEvent
	wg      sync.WaitGroup
	closing chan struct{}
	once    sync.Once
}

func NewAccessQueue(repo AccessWriter, bufferSize int, workers int) *AccessQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	if workers <= 0 {
		workers = 1
	}

	q := &AccessQueue{
		repo:    repo,
		ch:      make(chan *domain.AccessEvent, bufferSize),
		closing: make(chan struct{}),
	}

	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go q.worker()
	}

	return q
}

func (q *AccessQueue) Enqueue(e *domain.AccessEvent) error {
	select {
	case <-q.closing:
		return ErrQueueClosed
	default:
	}

	select {
	case q.ch <- e:
		return nil
	case <-q.closing:
		return ErrQueueClosed
	}
}

func (q *AccessQueue) worker() {
	defer q.wg.Done()

	for {
		select {
		case <-q.closing:
			return
		case e, ok := <-q.ch:
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, _ = q.repo.Create(ctx, e) // si falla, no tumba el server
			cancel()
		}
	}
}

func (q *AccessQueue) Close() {
	q.once.Do(func() {
		close(q.closing)
		close(q.ch)
		q.wg.Wait()
	})
}
