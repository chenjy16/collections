// Package list provides implementations of list data structures
package list

import (
	"github.com/chenjianyu/collections/container/common"
)

// List represents an ordered collection of elements that allows duplicates
type List[E any] interface {
	common.Container[E]
	common.Iterable[E]

	// Add adds an element to the end of the list
	// Returns whether the addition was successful
	Add(element E) bool

	// Insert inserts an element at the specified position
	// Returns an error if the index is invalid
	Insert(index int, element E) error

	// Get retrieves the element at the specified index
	// Returns the zero value and an error if the index is invalid
	Get(index int) (E, error)

	// Set replaces the element at the specified index
	// Returns the replaced element and whether the operation was successful
	Set(index int, element E) (E, bool)

	// RemoveAt removes the element at the specified index
	// Returns the removed element and whether the operation was successful
	RemoveAt(index int) (E, bool)

	// Remove removes the first occurrence of the specified element
	// Returns whether the removal was successful
	Remove(element E) bool

	// IndexOf returns the index of the first occurrence of the specified element in the list
	// Returns -1 if not found
	IndexOf(element E) int

	// LastIndexOf returns the index of the last occurrence of the specified element in the list
	// Returns -1 if not found
	LastIndexOf(element E) int

	// SubList returns a view of the specified range in the list
	// Returns an error if the indices are invalid
	SubList(fromIndex, toIndex int) (List[E], error)

	// ToSlice returns a slice containing all elements in the list
	ToSlice() []E
}
