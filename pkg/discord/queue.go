package discord

import (
	"sync"
)

type Queue struct {
	channel chan interface{}
	wg      sync.WaitGroup
}

// Create New Queue
func NewQueue(maxSize uint8) *Queue {
	return &Queue{
		channel: make(chan interface{}, maxSize),
	}
}

// Add item to Queue
func (q *Queue) Enqueue(value interface{}) {
	q.wg.Add(1)
	go func() {
		q.channel <- value
		q.wg.Done()
	}()
}

// Remove item from Queue
func (q *Queue) Dequeue() interface{} {
	return <-q.channel
}

// Close Queue
func (q *Queue) Close() {
	q.wg.Wait()
	close(q.channel)
}
