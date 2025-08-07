package list

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// ImmutableList is an immutable implementation of List
// Once created, it cannot be modified. All modification operations return new instances
type ImmutableList[E any] struct {
	elements []E
}

// NewImmutableList creates a new empty ImmutableList
func NewImmutableList[E any]() *ImmutableList[E] {
	return &ImmutableList[E]{
		elements: make([]E, 0),
	}
}

// NewImmutableListFromSlice creates a new ImmutableList from a slice
func NewImmutableListFromSlice[E any](elements []E) *ImmutableList[E] {
	// Create a copy to ensure immutability
	elementsCopy := make([]E, len(elements))
	copy(elementsCopy, elements)
	
	return &ImmutableList[E]{
		elements: elementsCopy,
	}
}

// Of creates a new ImmutableList from the given elements
func Of[E any](elements ...E) *ImmutableList[E] {
	return NewImmutableListFromSlice(elements)
}

// copyElements creates a copy of the elements slice
func (il *ImmutableList[E]) copyElements() []E {
	elementsCopy := make([]E, len(il.elements))
	copy(elementsCopy, il.elements)
	return elementsCopy
}

// Size returns the number of elements in the list
func (il *ImmutableList[E]) Size() int {
	return len(il.elements)
}

// IsEmpty returns true if the list is empty
func (il *ImmutableList[E]) IsEmpty() bool {
	return len(il.elements) == 0
}

// Clear does nothing as ImmutableList is immutable
// Use WithClear() to get a new empty list
func (il *ImmutableList[E]) Clear() {
	// No-op for immutable list
}

// Contains returns true if the list contains the specified element
func (il *ImmutableList[E]) Contains(element E) bool {
	for _, e := range il.elements {
		if any(e) == any(element) {
			return true
		}
	}
	return false
}

// Add returns false as ImmutableList is immutable
// Use WithAdd() to get a new list with the element added
func (il *ImmutableList[E]) Add(element E) bool {
	return false
}

// Insert returns an error as ImmutableList is immutable
// Use WithInsert() to get a new list with the element inserted
func (il *ImmutableList[E]) Insert(index int, element E) error {
	return common.ImmutableOperationError("Insert", "use WithInsert method")
}

// Get retrieves the element at the specified index
func (il *ImmutableList[E]) Get(index int) (E, error) {
	var zero E
	if index < 0 || index >= len(il.elements) {
		return zero, common.IndexOutOfBoundsError(index, len(il.elements))
	}
	return il.elements[index], nil
}

// Set returns zero value and false as ImmutableList is immutable
// Use WithSet() to get a new list with the element set
func (il *ImmutableList[E]) Set(index int, element E) (E, bool) {
	var zero E
	return zero, false
}

// RemoveAt returns zero value and false as ImmutableList is immutable
// Use WithRemoveAt() to get a new list with the element removed
func (il *ImmutableList[E]) RemoveAt(index int) (E, bool) {
	var zero E
	return zero, false
}

// Remove returns false as ImmutableList is immutable
// Use WithRemove() to get a new list with the element removed
func (il *ImmutableList[E]) Remove(element E) bool {
	return false
}

// IndexOf returns the index of the first occurrence of the specified element
func (il *ImmutableList[E]) IndexOf(element E) int {
	for i, e := range il.elements {
		if any(e) == any(element) {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of the specified element
func (il *ImmutableList[E]) LastIndexOf(element E) int {
	for i := len(il.elements) - 1; i >= 0; i-- {
		if any(il.elements[i]) == any(element) {
			return i
		}
	}
	return -1
}

// SubList returns a view of the specified range in the list
func (il *ImmutableList[E]) SubList(fromIndex, toIndex int) (List[E], error) {
	if fromIndex < 0 || toIndex > len(il.elements) || fromIndex > toIndex {
		return nil, common.InvalidRangeError(fromIndex, toIndex)
	}
	
	subElements := make([]E, toIndex-fromIndex)
	copy(subElements, il.elements[fromIndex:toIndex])
	
	return &ImmutableList[E]{elements: subElements}, nil
}

// ToSlice returns a slice containing all elements in the list
func (il *ImmutableList[E]) ToSlice() []E {
	return il.copyElements()
}

// WithAdd returns a new ImmutableList with the element added to the end
func (il *ImmutableList[E]) WithAdd(element E) *ImmutableList[E] {
	newElements := make([]E, len(il.elements)+1)
	copy(newElements, il.elements)
	newElements[len(il.elements)] = element
	
	return &ImmutableList[E]{elements: newElements}
}

// WithInsert returns a new ImmutableList with the element inserted at the specified position
func (il *ImmutableList[E]) WithInsert(index int, element E) (*ImmutableList[E], error) {
	if index < 0 || index > len(il.elements) {
		return nil, common.IndexOutOfBoundsError(index, len(il.elements))
	}
	
	newElements := make([]E, len(il.elements)+1)
	copy(newElements[:index], il.elements[:index])
	newElements[index] = element
	copy(newElements[index+1:], il.elements[index:])
	
	return &ImmutableList[E]{elements: newElements}, nil
}

// WithSet returns a new ImmutableList with the element at the specified index replaced
func (il *ImmutableList[E]) WithSet(index int, element E) (*ImmutableList[E], error) {
	if index < 0 || index >= len(il.elements) {
		return nil, common.IndexOutOfBoundsError(index, len(il.elements))
	}
	
	newElements := il.copyElements()
	newElements[index] = element
	
	return &ImmutableList[E]{elements: newElements}, nil
}

// WithRemoveAt returns a new ImmutableList with the element at the specified index removed
func (il *ImmutableList[E]) WithRemoveAt(index int) (*ImmutableList[E], error) {
	if index < 0 || index >= len(il.elements) {
		return nil, common.IndexOutOfBoundsError(index, len(il.elements))
	}
	
	newElements := make([]E, len(il.elements)-1)
	copy(newElements[:index], il.elements[:index])
	copy(newElements[index:], il.elements[index+1:])
	
	return &ImmutableList[E]{elements: newElements}, nil
}

// WithRemove returns a new ImmutableList with the first occurrence of the element removed
func (il *ImmutableList[E]) WithRemove(element E) *ImmutableList[E] {
	index := il.IndexOf(element)
	if index == -1 {
		return il // Element not found, return same instance
	}
	
	newList, _ := il.WithRemoveAt(index)
	return newList
}

// WithClear returns a new empty ImmutableList
func (il *ImmutableList[E]) WithClear() *ImmutableList[E] {
	return NewImmutableList[E]()
}

// Iterator returns an iterator for the list
func (il *ImmutableList[E]) Iterator() common.Iterator[E] {
	return &immutableListIterator[E]{
		list:  il,
		index: 0,
	}
}

// ForEach executes the given function for each element in the list
func (il *ImmutableList[E]) ForEach(fn func(E)) {
	for _, element := range il.elements {
		fn(element)
	}
}

// String returns a string representation of the list
func (il *ImmutableList[E]) String() string {
	if il.IsEmpty() {
		return "[]"
	}
	
	var builder strings.Builder
	builder.WriteString("[")
	
	for i, element := range il.elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", element))
	}
	
	builder.WriteString("]")
	return builder.String()
}

// immutableListIterator implements Iterator for ImmutableList
type immutableListIterator[E any] struct {
	list  *ImmutableList[E]
	index int
}

func (it *immutableListIterator[E]) HasNext() bool {
	return it.index < len(it.list.elements)
}

func (it *immutableListIterator[E]) Next() (E, bool) {
	var zero E
	if !it.HasNext() {
		return zero, false
	}
	
	element := it.list.elements[it.index]
	it.index++
	return element, true
}

func (it *immutableListIterator[E]) Remove() bool {
	// Remove operation is not supported for immutable collections
	return false
}