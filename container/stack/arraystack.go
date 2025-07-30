// Package stack provides stack data structure implementations
package stack

import (
	"errors"
	"fmt"
	"strings"
)

// ErrEmptyStack indicates that the stack is empty
var ErrEmptyStack = errors.New("stack is empty")

// ErrFullStack indicates that the stack is full
var ErrFullStack = errors.New("stack is full")

// ArrayStack is a stack implementation based on slices
type ArrayStack[E any] struct {
	elements []E
	maxCap   int // Maximum capacity, 0 means unbounded
}

// New creates a new unbounded ArrayStack
func New[E any]() *ArrayStack[E] {
	return &ArrayStack[E]{
		elements: make([]E, 0),
		maxCap:   0,
	}
}

// WithCapacity creates an ArrayStack with specified maximum capacity
func WithCapacity[E any](capacity int) *ArrayStack[E] {
	return &ArrayStack[E]{
		elements: make([]E, 0, capacity),
		maxCap:   capacity,
	}
}

// FromSlice creates a new ArrayStack from a slice
// The first element in the slice will be the bottom element, the last element will be the top element
func FromSlice[E any](slice []E) *ArrayStack[E] {
	elements := make([]E, len(slice))
	copy(elements, slice)
	return &ArrayStack[E]{
		elements: elements,
		maxCap:   0,
	}
}

// Size returns the number of elements in the stack
func (s *ArrayStack[E]) Size() int {
	return len(s.elements)
}

// IsEmpty checks if the stack is empty
func (s *ArrayStack[E]) IsEmpty() bool {
	return len(s.elements) == 0
}

// isFull checks if the stack is full
func (s *ArrayStack[E]) isFull() bool {
	return s.maxCap > 0 && len(s.elements) >= s.maxCap
}

// Clear empties the stack
func (s *ArrayStack[E]) Clear() {
	s.elements = s.elements[:0]
}

// Contains checks if the stack contains the specified element
func (s *ArrayStack[E]) Contains(element E) bool {
	for _, e := range s.elements {
		if any(e) == any(element) {
			return true
		}
	}
	return false
}

// ForEach executes the given operation on each element in the stack
// Traversal order is from bottom to top
func (s *ArrayStack[E]) ForEach(fn func(E)) {
	for _, element := range s.elements {
		fn(element)
	}
}

// String returns the string representation of the stack
func (s *ArrayStack[E]) String() string {
	if s.IsEmpty() {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	for i, element := range s.elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", element))
	}
	builder.WriteString("]")
	return builder.String()
}

// Push pushes an element onto the top of the stack
func (s *ArrayStack[E]) Push(element E) error {
	if s.isFull() {
		return ErrFullStack
	}
	s.elements = append(s.elements, element)
	return nil
}

// Pop removes and returns the element at the top of the stack
func (s *ArrayStack[E]) Pop() (E, error) {
	if s.IsEmpty() {
		var zero E
		return zero, ErrEmptyStack
	}

	index := len(s.elements) - 1
	element := s.elements[index]
	s.elements = s.elements[:index]
	return element, nil
}

// Peek returns the element at the top of the stack without removing it
func (s *ArrayStack[E]) Peek() (E, error) {
	if s.IsEmpty() {
		var zero E
		return zero, ErrEmptyStack
	}
	return s.elements[len(s.elements)-1], nil
}

// Search searches for an element in the stack
func (s *ArrayStack[E]) Search(element E) int {
	for i := len(s.elements) - 1; i >= 0; i-- {
		if any(s.elements[i]) == any(element) {
			return len(s.elements) - i
		}
	}
	return -1
}

// ToSlice returns a slice containing all elements in the stack
func (s *ArrayStack[E]) ToSlice() []E {
	result := make([]E, len(s.elements))
	copy(result, s.elements)
	return result
}
