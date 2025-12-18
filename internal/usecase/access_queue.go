package usecase

import (
    "context"
    "log"
    "sync"
    "time"

    "github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// Concurrencia + goroutines + canales
type AccessQueue struct {
    ch      chan *domain.AccessEvent
    repo    domain.AccessLogRepository
    wg      sync.WaitGroup
    closing chan struct{}
}

func NewAccessQueue(repo domain.AccessLogRepository, buffer int, workers int) *AccessQueue {
    if buffer <= 0 { buffer = 100 }
    if workers <= 0 { workers = 2 }

    q := &AccessQueue{
        ch:      make(chan *domain.AccessEvent, buffer),
        repo:    repo,
        closing: make(chan struct{}),
    }

    for i := 0; i < workers; i++ {
        q.wg.Add(1)
        go q.worker(i+1) // goroutine
    }
    return q
}

func (q *AccessQueue) Enqueue(e *domain.AccessEvent) bool {
    select {
    case q.ch <- e:
        return true
    default:
        // Cola llena: backpressure
        return false
    }
}

func (q *AccessQueue) worker(n int) {
    defer q.wg.Done()
    for {
        select {
        case e := <-q.ch:
            if e == nil {
                continue
            }
            ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
            if _, err := q.repo.Create(ctx, e); err != nil {
                log.Printf("[access-worker-%d] insert error: %v", n, err)
            }
            cancel()
        case <-q.closing:
            return
        }
    }
}

func (q *AccessQueue) Close() {
    close(q.closing)
    q.wg.Wait()
}
