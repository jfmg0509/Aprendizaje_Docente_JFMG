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
				_, _ = q.repo.Create(context.Background(), e)
			}
		}()
	}

	return q
}

func (q *AccessQueue) Enqueue(e *domain.AccessEvent) {
	q.closeMu.Lock()
	defer q.closeMu.Unlock()

	if q.closed {
		return
	}
	q.ch <- e
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
