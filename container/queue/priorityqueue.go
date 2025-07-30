// Package queue provides queue data structure implementations
package queue

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// ErrEmptyQueue indicates the queue is empty error
var ErrEmptyQueue = errors.New("queue is empty")

// ErrFullQueue indicates the queue is full error
var ErrFullQueue = errors.New("queue is full")

// PriorityQueue is a priority queue implementation based on binary heap
// Default is min heap, can create max heap by providing custom comparator
type PriorityQueue[E any] struct {
	heap       []E
	comparator func(a, b E) int
	maxCap     int // maximum capacity, 0 means unbounded
}

// NewPriorityQueue creates a new unbounded priority queue using default comparator
// Default comparator requires element type to implement common.Comparable interface
func NewPriorityQueue[E common.Comparable]() *PriorityQueue[E] {
	return &PriorityQueue[E]{
		heap: make([]E, 0),
		comparator: func(a, b E) int {
			return a.CompareTo(b)
		},
		maxCap: 0,
	}
}

// NewPriorityQueueWithComparator creates a new unbounded priority queue using custom comparator
func NewPriorityQueueWithComparator[E any](comparator func(a, b E) int) *PriorityQueue[E] {
	return &PriorityQueue[E]{
		heap:       make([]E, 0),
		comparator: comparator,
		maxCap:     0,
	}
}

// WithCapacity creates a priority queue with specified maximum capacity using default comparator
func WithCapacityComparable[E common.Comparable](capacity int) *PriorityQueue[E] {
	return &PriorityQueue[E]{
		heap: make([]E, 0, capacity),
		comparator: func(a, b E) int {
			return a.CompareTo(b)
		},
		maxCap: capacity,
	}
}

// WithCapacityAndComparator creates a priority queue with specified maximum capacity using custom comparator
func WithCapacityAndComparator[E any](capacity int, comparator func(a, b E) int) *PriorityQueue[E] {
	return &PriorityQueue[E]{
		heap:       make([]E, 0, capacity),
		comparator: comparator,
		maxCap:     capacity,
	}
}

// NewFromSlice creates a new priority queue from a slice using default comparator
func NewFromSliceComparable[E common.Comparable](slice []E) *PriorityQueue[E] {
	pq := NewPriorityQueue[E]()
	for _, item := range slice {
		pq.Add(item)
	}
	return pq
}

// NewFromSliceWithComparator creates a new priority queue from a slice using custom comparator
func NewFromSliceWithComparator[E any](slice []E, comparator func(a, b E) int) *PriorityQueue[E] {
	pq := NewPriorityQueueWithComparator(comparator)
	for _, item := range slice {
		pq.Add(item)
	}
	return pq
}

// Size returns the number of elements in the queue
func (pq *PriorityQueue[E]) Size() int {
	return len(pq.heap)
}

// IsEmpty checks if the queue is empty
func (pq *PriorityQueue[E]) IsEmpty() bool {
	return len(pq.heap) == 0
}

// isFull checks if the queue is full
func (pq *PriorityQueue[E]) isFull() bool {
	return pq.maxCap > 0 && len(pq.heap) >= pq.maxCap
}

// Clear clears the queue
func (pq *PriorityQueue[E]) Clear() {
	pq.heap = pq.heap[:0]
}

// Contains checks if the queue contains the specified element
func (pq *PriorityQueue[E]) Contains(element E) bool {
	for _, item := range pq.heap {
		if common.Equal(item, element) {
			return true
		}
	}
	return false
}

// ForEach executes the given operation for each element in the queue
// Note: traversal order is not guaranteed to be in priority order
func (pq *PriorityQueue[E]) ForEach(fn func(E)) {
	for _, item := range pq.heap {
		fn(item)
	}
}

// String returns the string representation of the queue
func (pq *PriorityQueue[E]) String() string {
	if pq.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	for i, item := range pq.heap {
		sb.WriteString(fmt.Sprintf("%v", item))
		if i < len(pq.heap)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]")
	return sb.String()
}

// Add adds an element to the queue
func (pq *PriorityQueue[E]) Add(element E) error {
	if pq.isFull() {
		return ErrFullQueue
	}

	pq.heap = append(pq.heap, element)
	pq.heapifyUp(len(pq.heap) - 1)
	return nil
}

// Offer adds an element to the queue
func (pq *PriorityQueue[E]) Offer(element E) bool {
	if pq.isFull() {
		return false
	}

	pq.heap = append(pq.heap, element)
	pq.heapifyUp(len(pq.heap) - 1)
	return true
}

// Remove removes and returns the highest priority element from the queue
func (pq *PriorityQueue[E]) Remove() (E, error) {
	if pq.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}

	root := pq.heap[0]
	lastIndex := len(pq.heap) - 1
	pq.heap[0] = pq.heap[lastIndex]
	pq.heap = pq.heap[:lastIndex]

	if len(pq.heap) > 0 {
		pq.heapifyDown(0)
	}

	return root, nil
}

// Poll removes and returns the highest priority element from the queue
func (pq *PriorityQueue[E]) Poll() (E, bool) {
	element, err := pq.Remove()
	return element, err == nil
}

// Element returns the highest priority element from the queue without removing it
func (pq *PriorityQueue[E]) Element() (E, error) {
	if pq.IsEmpty() {
		var zero E
		return zero, ErrEmptyQueue
	}
	return pq.heap[0], nil
}

// Peek returns the highest priority element from the queue without removing it
func (pq *PriorityQueue[E]) Peek() (E, bool) {
	element, err := pq.Element()
	return element, err == nil
}

// ToSlice returns a slice containing all elements in the queue
// Note: returned slice is not guaranteed to be sorted by priority
func (pq *PriorityQueue[E]) ToSlice() []E {
	result := make([]E, len(pq.heap))
	copy(result, pq.heap)
	return result
}

// ToSortedSlice returns a slice of queue elements sorted by priority
func (pq *PriorityQueue[E]) ToSortedSlice() []E {
	if pq.IsEmpty() {
		return []E{}
	}

	// Create a copy of the priority queue
	tempPQ := &PriorityQueue[E]{
		heap:       make([]E, len(pq.heap)),
		comparator: pq.comparator,
		maxCap:     0,
	}
	copy(tempPQ.heap, pq.heap)

	result := make([]E, 0, len(pq.heap))
	for !tempPQ.IsEmpty() {
		element, _ := tempPQ.Remove()
		result = append(result, element)
	}

	return result
}

// Internal method: heapify up operation to maintain heap property
func (pq *PriorityQueue[E]) heapifyUp(index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if pq.comparator(pq.heap[index], pq.heap[parentIndex]) >= 0 {
			break
		}
		pq.heap[index], pq.heap[parentIndex] = pq.heap[parentIndex], pq.heap[index]
		index = parentIndex
	}
}

// Internal method: heapify down operation to maintain heap property
func (pq *PriorityQueue[E]) heapifyDown(index int) {
	for {
		leftChild := 2*index + 1
		rightChild := 2*index + 2
		smallest := index

		if leftChild < len(pq.heap) && pq.comparator(pq.heap[leftChild], pq.heap[smallest]) < 0 {
			smallest = leftChild
		}

		if rightChild < len(pq.heap) && pq.comparator(pq.heap[rightChild], pq.heap[smallest]) < 0 {
			smallest = rightChild
		}

		if smallest == index {
			break
		}

		pq.heap[index], pq.heap[smallest] = pq.heap[smallest], pq.heap[index]
		index = smallest
	}
}
