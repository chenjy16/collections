// Package queue provides queue data structure implementations
package queue

import (
	"github.com/chenjianyu/collections/container/common"
)

// Queue represents a first-in-first-out (FIFO) queue
type Queue[E any] interface {
	common.Container[E]

	// Add adds an element to the tail of the queue
	// Returns an error if the queue is full (for bounded queues)
	Add(element E) error

	// Offer adds an element to the tail of the queue
	// Returns false if the queue is full (for bounded queues)
	Offer(element E) bool

	// Remove removes and returns the element at the head of the queue
	// Returns an error if the queue is empty
	Remove() (E, error)

	// Poll removes and returns the element at the head of the queue
	// Returns zero value and false if the queue is empty
	Poll() (E, bool)

	// Element returns the element at the head of the queue without removing it
	// Returns an error if the queue is empty
	Element() (E, error)

	// Peek returns the element at the head of the queue without removing it
	// Returns zero value and false if the queue is empty
	Peek() (E, bool)
}

// Deque represents a double-ended queue that supports adding and removing elements from both ends
type Deque[E any] interface {
	Queue[E]

	// AddFirst adds an element to the head of the queue
	// Returns an error if the queue is full (for bounded queues)
	AddFirst(element E) error

	// AddLast adds an element to the tail of the queue
	// Returns an error if the queue is full (for bounded queues)
	AddLast(element E) error

	// OfferFirst adds an element to the head of the queue
	// Returns false if the queue is full (for bounded queues)
	OfferFirst(element E) bool

	// OfferLast adds an element to the tail of the queue
	// Returns false if the queue is full (for bounded queues)
	OfferLast(element E) bool

	// RemoveFirst removes and returns the element at the head of the queue
	// Returns an error if the queue is empty
	RemoveFirst() (E, error)

	// RemoveLast removes and returns the element at the tail of the queue
	// Returns an error if the queue is empty
	RemoveLast() (E, error)

	// PollFirst removes and returns the element at the head of the queue
	// Returns zero value and false if the queue is empty
	PollFirst() (E, bool)

	// PollLast removes and returns the element at the tail of the queue
	// Returns zero value and false if the queue is empty
	PollLast() (E, bool)

	// GetFirst returns the element at the head of the queue without removing it
	// Returns an error if the queue is empty
	GetFirst() (E, error)

	// GetLast returns the element at the tail of the queue without removing it
	// Returns an error if the queue is empty
	GetLast() (E, error)

	// PeekFirst returns the element at the head of the queue without removing it
	// Returns zero value and false if the queue is empty
	PeekFirst() (E, bool)

	// PeekLast returns the element at the tail of the queue without removing it
	// Returns zero value and false if the queue is empty
	PeekLast() (E, bool)
}
