// Package queue provides queue data structure implementations
package queue

import (
	"github.com/chenjianyu/collections/container/list"
)

// LinkedList is a Deque implementation based on doubly linked list, reusing LinkedList from list package
type LinkedList[E any] struct {
	list   *list.LinkedList[E]
	maxCap int // maximum capacity, 0 means unbounded
}

// New creates a new unbounded LinkedList
func New[E any]() *LinkedList[E] {
	return &LinkedList[E]{
		list:   list.NewLinkedList[E](),
		maxCap: 0,
	}
}

// WithCapacity creates a LinkedList with specified maximum capacity
func WithCapacity[E any](capacity int) *LinkedList[E] {
	return &LinkedList[E]{
		list:   list.NewLinkedList[E](),
		maxCap: capacity,
	}
}

// FromSlice creates a new LinkedList from a slice
func FromSlice[E any](slice []E) *LinkedList[E] {
	ll := New[E]()
	for _, item := range slice {
		ll.list.Add(item)
	}
	return ll
}

// Size returns the number of elements in the queue
func (ll *LinkedList[E]) Size() int {
	return ll.list.Size()
}

// IsEmpty checks if the queue is empty
func (ll *LinkedList[E]) IsEmpty() bool {
	return ll.list.IsEmpty()
}

// isFull checks if the queue is full
func (ll *LinkedList[E]) isFull() bool {
	return ll.maxCap > 0 && ll.list.Size() >= ll.maxCap
}

// Clear clears the queue
func (ll *LinkedList[E]) Clear() {
	ll.list.Clear()
}

// Contains checks if the queue contains the specified element
func (ll *LinkedList[E]) Contains(element E) bool {
	return ll.list.Contains(element)
}

// ForEach executes the given operation for each element in the queue
func (ll *LinkedList[E]) ForEach(fn func(E)) {
	ll.list.ForEach(fn)
}

// String returns the string representation of the queue
func (ll *LinkedList[E]) String() string {
	return ll.list.String()
}

// Add adds an element to the tail of the queue
func (ll *LinkedList[E]) Add(element E) error {
	if ll.isFull() {
		return ErrFullQueue
	}
	ll.list.Add(element)
	return nil
}

// Offer adds an element to the tail of the queue
func (ll *LinkedList[E]) Offer(element E) bool {
	if ll.isFull() {
		return false
	}
	ll.list.Add(element)
	return true
}

// Remove removes and returns the element at the head of the queue
func (ll *LinkedList[E]) Remove() (E, error) {
	if ll.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}
	val, success := ll.list.RemoveFirst()
	if !success {
		var zero E
		return zero, ErrEmptyQueue
	}
	return val, nil
}

// Poll removes and returns the element at the head of the queue
func (ll *LinkedList[E]) Poll() (E, bool) {
	val, success := ll.list.RemoveFirst()
	return val, success
}

// Element returns the element at the head of the queue without removing it
func (ll *LinkedList[E]) Element() (E, error) {
	return ll.list.GetFirst()
}

// Peek returns the element at the head of the queue without removing it
func (ll *LinkedList[E]) Peek() (E, bool) {
	val, err := ll.list.GetFirst()
	return val, err == nil
}

// AddFirst adds an element to the head of the queue
func (ll *LinkedList[E]) AddFirst(element E) error {
	if ll.isFull() {
		return ErrFullQueue
	}
	ll.list.AddFirst(element)
	return nil
}

// AddLast adds an element to the tail of the queue
func (ll *LinkedList[E]) AddLast(element E) error {
	if ll.isFull() {
		return ErrFullQueue
	}
	ll.list.AddLast(element)
	return nil
}

// OfferFirst adds an element to the head of the queue
func (ll *LinkedList[E]) OfferFirst(element E) bool {
	if ll.isFull() {
		return false
	}
	ll.list.AddFirst(element)
	return true
}

// OfferLast adds an element to the tail of the queue
func (ll *LinkedList[E]) OfferLast(element E) bool {
	if ll.isFull() {
		return false
	}
	ll.list.AddLast(element)
	return true
}

// RemoveFirst removes and returns the element at the head of the queue
func (ll *LinkedList[E]) RemoveFirst() (E, error) {
	if ll.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}
	val, success := ll.list.RemoveFirst()
	if !success {
		var zero E
		return zero, ErrEmptyQueue
	}
	return val, nil
}

// RemoveLast removes and returns the element at the tail of the queue
func (ll *LinkedList[E]) RemoveLast() (E, error) {
	if ll.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}
	val, success := ll.list.RemoveLast()
	if !success {
		var zero E
		return zero, ErrEmptyQueue
	}
	return val, nil
}

// PollFirst removes and returns the element at the head of the queue
func (ll *LinkedList[E]) PollFirst() (E, bool) {
	val, err := ll.RemoveFirst()
	return val, err == nil
}

// PollLast removes and returns the element at the tail of the queue
func (ll *LinkedList[E]) PollLast() (E, bool) {
	val, err := ll.RemoveLast()
	return val, err == nil
}

// GetFirst returns the element at the head of the queue without removing it
func (ll *LinkedList[E]) GetFirst() (E, error) {
	if ll.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}
	return ll.list.GetFirst()
}

// GetLast returns the element at the tail of the queue without removing it
func (ll *LinkedList[E]) GetLast() (E, error) {
	if ll.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}
	return ll.list.GetLast()
}

// PeekFirst returns the element at the head of the queue without removing it
func (ll *LinkedList[E]) PeekFirst() (E, bool) {
	val, err := ll.GetFirst()
	return val, err == nil
}

// PeekLast returns the element at the tail of the queue without removing it
func (ll *LinkedList[E]) PeekLast() (E, bool) {
	val, err := ll.GetLast()
	return val, err == nil
}

// ToSlice returns a slice containing all elements in the queue
func (ll *LinkedList[E]) ToSlice() []E {
	return ll.list.ToSlice()
}
