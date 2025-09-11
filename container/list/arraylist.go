// Package list provides implementations of list data structures
package list

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// ArrayList is a List implementation based on dynamic arrays
type ArrayList[E any] struct {
	elements []E
}

// New creates a new ArrayList
func New[E any]() *ArrayList[E] {
	return &ArrayList[E]{elements: make([]E, 0)}
}

// WithCapacity creates an ArrayList with the specified initial capacity
func WithCapacity[E any](capacity int) *ArrayList[E] {
	if capacity < 0 {
		capacity = 0
	}
	return &ArrayList[E]{elements: make([]E, 0, capacity)}
}

// FromSlice creates a new ArrayList from a slice
func FromSlice[E any](slice []E) *ArrayList[E] {
	elements := make([]E, len(slice))
	copy(elements, slice)
	return &ArrayList[E]{elements: elements}
}

// Add adds an element to the end of the list
func (list *ArrayList[E]) Add(element E) bool {
	list.elements = append(list.elements, element)
	return true
}

// Insert inserts an element at the specified position
func (list *ArrayList[E]) Insert(index int, element E) error {
	if index < 0 || index > len(list.elements) {
		return common.IndexOutOfBoundsError(index, len(list.elements))
	}

	// Add to the end
	if index == len(list.elements) {
		list.elements = append(list.elements, element)
		return nil
	}

	// Insert in the middle
	list.elements = append(list.elements, common.ZeroValue[E]()) // Extend the slice
	copy(list.elements[index+1:], list.elements[index:])
	list.elements[index] = element
	return nil
}

// Get retrieves the element at the specified index
func (list *ArrayList[E]) Get(index int) (E, error) {
	if index < 0 || index >= len(list.elements) {
		return common.ZeroValue[E](), common.IndexOutOfBoundsError(index, len(list.elements))
	}
	return list.elements[index], nil
}

// Set replaces the element at the specified index
func (list *ArrayList[E]) Set(index int, element E) (E, bool) {
	if index < 0 || index >= len(list.elements) {
		return common.ZeroValue[E](), false
	}
	oldElement := list.elements[index]
	list.elements[index] = element
	return oldElement, true
}

// RemoveAt removes the element at the specified index
func (list *ArrayList[E]) RemoveAt(index int) (E, bool) {
	if index < 0 || index >= len(list.elements) {
		return common.ZeroValue[E](), false
	}
	element := list.elements[index]
	// Move elements to fill the gap
	copy(list.elements[index:], list.elements[index+1:])
	// Clear the last element to avoid memory leak
	list.elements[len(list.elements)-1] = common.ZeroValue[E]()
	// Shrink the slice
	list.elements = list.elements[:len(list.elements)-1]
	return element, true
}

// Remove removes the first occurrence of the specified element
func (list *ArrayList[E]) Remove(element E) bool {
	index := list.IndexOf(element)
	if index >= 0 {
		_, removed := list.RemoveAt(index)
		return removed
	}
	return false
}

// Contains checks if the list contains the specified element
func (list *ArrayList[E]) Contains(element E) bool {
	return list.IndexOf(element) >= 0
}

// IndexOf returns the index of the first occurrence of the specified element in the list
func (list *ArrayList[E]) IndexOf(element E) int {
	for i, e := range list.elements {
		if common.Equal(e, element) {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of the specified element in the list
func (list *ArrayList[E]) LastIndexOf(element E) int {
	for i := len(list.elements) - 1; i >= 0; i-- {
		if common.Equal(list.elements[i], element) {
			return i
		}
	}
	return -1
}

// Size returns the number of elements in the list
func (list *ArrayList[E]) Size() int {
	return len(list.elements)
}

// IsEmpty checks if the list is empty
func (list *ArrayList[E]) IsEmpty() bool {
	return len(list.elements) == 0
}

// Clear empties the list and releases memory by setting elements to zero values
func (list *ArrayList[E]) Clear() {
	// Clear all elements to prevent memory leaks
	for i := range list.elements {
		list.elements[i] = common.ZeroValue[E]()
	}
	// Reset the slice to zero length but keep some capacity for reuse
	if cap(list.elements) > 100 { // Only shrink if capacity is large
		list.elements = make([]E, 0, 16) // Reset to smaller capacity
	} else {
		list.elements = list.elements[:0]
	}
}

// ToSlice returns a slice containing all elements in the list
func (list *ArrayList[E]) ToSlice() []E {
	result := make([]E, len(list.elements))
	copy(result, list.elements)
	return result
}

// SubList returns a view of the specified range in the list
func (list *ArrayList[E]) SubList(fromIndex, toIndex int) (List[E], error) {
	if fromIndex < 0 || toIndex > len(list.elements) || fromIndex > toIndex {
		return nil, common.InvalidRangeError(fromIndex, toIndex)
	}

	subElements := make([]E, toIndex-fromIndex)
	copy(subElements, list.elements[fromIndex:toIndex])
	return &ArrayList[E]{elements: subElements}, nil
}

// ForEach executes the given operation on each element in the list
func (list *ArrayList[E]) ForEach(f func(E)) {
	for _, element := range list.elements {
		f(element)
	}
}

// String returns the string representation of the list
func (list *ArrayList[E]) String() string {
	if len(list.elements) == 0 {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	for i, element := range list.elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", element))
	}
	builder.WriteString("]")
	return builder.String()
}

// Iterator returns an iterator for traversing the elements in the list
func (list *ArrayList[E]) Iterator() common.Iterator[E] {
	return &arrayListIterator[E]{list: list, cursor: 0, lastRet: -1}
}

// arrayListIterator is the iterator implementation for ArrayList
type arrayListIterator[E any] struct {
	list    *ArrayList[E]
	cursor  int // Index of the next element
	lastRet int // Index of the last returned element, -1 if none
}

// HasNext checks if the iterator has a next element
func (it *arrayListIterator[E]) HasNext() bool {
	return it.cursor < len(it.list.elements)
}

// Next returns the next element in the iterator
func (it *arrayListIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		return common.ZeroValue[E](), false
	}

	element := it.list.elements[it.cursor]
	it.lastRet = it.cursor
	it.cursor++
	return element, true
}

// Remove removes the last element returned by the iterator
func (it *arrayListIterator[E]) Remove() bool {
	if it.lastRet < 0 {
		return false
	}

	_, removed := it.list.RemoveAt(it.lastRet)
	if removed {
		it.cursor = it.lastRet
		it.lastRet = -1
	}
	return removed
}
