// Package stack provides stack data structure implementations
package stack

import (
	"github.com/chenjianyu/collections/container/common"
)

// Stack represents a last-in-first-out (LIFO) stack
type Stack[E any] interface {
	common.Container[E]

	// Push pushes an element onto the top of the stack
	// Returns an error if the stack is full (for bounded stacks)
	Push(element E) error

	// Pop removes and returns the element at the top of the stack
	// Returns an error if the stack is empty
	Pop() (E, error)

	// Peek returns the element at the top of the stack without removing it
	// Returns an error if the stack is empty
	Peek() (E, error)

	// Search searches for an element in the stack
	// If found, returns the 1-based position index (1 means top of stack)
	// If not found, returns -1
	Search(element E) int

	// ToSlice returns a slice containing all elements in the stack
	// The first element in the slice is the bottom element, the last element is the top element
	ToSlice() []E
}
