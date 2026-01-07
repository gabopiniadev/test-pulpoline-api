package queue

import (
	"log"
	"sync"

	"test-pulpoline-api/internal/queue/models"
	"test-pulpoline-api/pkg/errors"
)

type RequestQueue struct {
	requests chan models.Request
	workers  int
	wg       sync.WaitGroup
	closed   bool
	mu       sync.Mutex
}

func NewRequestQueue(bufferSize int) *RequestQueue {
	queue := &RequestQueue{
		requests: make(chan models.Request, bufferSize),
		workers:  5,
	}

	queue.startWorkers()

	return queue
}

func (q *RequestQueue) startWorkers() {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

func (q *RequestQueue) worker(id int) {
	defer q.wg.Done()
	log.Printf("Worker %d iniciado", id)

	for req := range q.requests {
		select {
		case <-req.Context.Done():
			req.Error <- req.Context.Err()
			continue
		default:
		}

		log.Printf("Worker %d procesando solicitud %s", id, req.ID)
	}
}

func (q *RequestQueue) Enqueue(req models.Request) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return errors.ErrQueueClosed
	}

	select {
	case q.requests <- req:
		return nil
	default:
		return errors.ErrQueueFull
	}
}

func (q *RequestQueue) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	close(q.requests)
	q.mu.Unlock()

	q.wg.Wait()
	log.Println("Cola cerrada, todos los workers han terminado")
}

func (q *RequestQueue) GetChannel() <-chan models.Request {
	return q.requests
}
