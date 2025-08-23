package utils

import "sync"

type Queue struct {
	mu       sync.Mutex
	elements []interface{}
}

func (q *Queue) Enqueue(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.elements = append(q.elements, item)
}

func (q *Queue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.elements) == 0 {
		return nil
	}
	item := q.elements[0]
	q.elements = q.elements[1:]
	return item
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.elements) == 0
}

func (q *Queue) Peek() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.elements) == 0 {
		return nil
	}
	return q.elements[0]
}
