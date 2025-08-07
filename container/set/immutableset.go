package set

import (
	"fmt"
	"log"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// ImmutableSet is an immutable implementation of Set
// Once created, it cannot be modified. All modification operations return new instances
type ImmutableSet[E comparable] struct {
	elements map[E]struct{}
}

// NewImmutableSet creates a new empty ImmutableSet
func NewImmutableSet[E comparable]() *ImmutableSet[E] {
	return &ImmutableSet[E]{
		elements: make(map[E]struct{}),
	}
}

// NewImmutableSetFromSlice creates a new ImmutableSet from a slice
func NewImmutableSetFromSlice[E comparable](elements []E) *ImmutableSet[E] {
	elementMap := make(map[E]struct{})
	for _, element := range elements {
		elementMap[element] = struct{}{}
	}
	
	return &ImmutableSet[E]{
		elements: elementMap,
	}
}

// SetOf creates a new ImmutableSet from the given elements
func SetOf[E comparable](elements ...E) *ImmutableSet[E] {
	return NewImmutableSetFromSlice(elements)
}

// copyElements creates a copy of the elements map
func (is *ImmutableSet[E]) copyElements() map[E]struct{} {
	elementsCopy := make(map[E]struct{}, len(is.elements))
	for k, v := range is.elements {
		elementsCopy[k] = v
	}
	return elementsCopy
}

// Size returns the number of elements in the set
func (is *ImmutableSet[E]) Size() int {
	return len(is.elements)
}

// IsEmpty returns true if the set is empty
func (is *ImmutableSet[E]) IsEmpty() bool {
	return len(is.elements) == 0
}

// Clear logs an error as ImmutableSet is immutable
func (is *ImmutableSet[E]) Clear() {
	err := common.ImmutableOperationError("Clear", "create a new empty set")
	log.Printf("Warning: %v", err)
}

// Contains returns true if the set contains the specified element
func (is *ImmutableSet[E]) Contains(element E) bool {
	_, exists := is.elements[element]
	return exists
}

// Add logs an error and returns false as ImmutableSet is immutable
func (is *ImmutableSet[E]) Add(element E) bool {
	err := common.ImmutableOperationError("Add", "use WithAdd method")
	log.Printf("Warning: %v", err)
	return false
}

// Remove logs an error and returns false as ImmutableSet is immutable
func (is *ImmutableSet[E]) Remove(element E) bool {
	err := common.ImmutableOperationError("Remove", "use WithRemove method")
	log.Printf("Warning: %v", err)
	return false
}

// ToSlice returns a slice containing all elements in the set
func (is *ImmutableSet[E]) ToSlice() []E {
	result := make([]E, 0, len(is.elements))
	for element := range is.elements {
		result = append(result, element)
	}
	return result
}

// WithAdd returns a new ImmutableSet with the element added
func (is *ImmutableSet[E]) WithAdd(element E) *ImmutableSet[E] {
	if is.Contains(element) {
		return is // Element already exists, return same instance
	}
	
	newElements := is.copyElements()
	newElements[element] = struct{}{}
	
	return &ImmutableSet[E]{elements: newElements}
}

// WithRemove returns a new ImmutableSet with the element removed
func (is *ImmutableSet[E]) WithRemove(element E) *ImmutableSet[E] {
	if !is.Contains(element) {
		return is // Element doesn't exist, return same instance
	}
	
	newElements := is.copyElements()
	delete(newElements, element)
	
	return &ImmutableSet[E]{elements: newElements}
}

// WithClear returns a new empty ImmutableSet
func (is *ImmutableSet[E]) WithClear() *ImmutableSet[E] {
	return NewImmutableSet[E]()
}

// Union returns the union of this set and another set
func (is *ImmutableSet[E]) Union(other Set[E]) Set[E] {
	result := is.copyElements()
	
	// Add all elements from the other set
	other.ForEach(func(element E) {
		result[element] = struct{}{}
	})
	
	return &ImmutableSet[E]{elements: result}
}

// Intersection returns the intersection of this set and another set
func (is *ImmutableSet[E]) Intersection(other Set[E]) Set[E] {
	result := make(map[E]struct{})
	
	// Only include elements that exist in both sets
	for element := range is.elements {
		if other.Contains(element) {
			result[element] = struct{}{}
		}
	}
	
	return &ImmutableSet[E]{elements: result}
}

// Difference returns the difference of this set and another set
func (is *ImmutableSet[E]) Difference(other Set[E]) Set[E] {
	result := make(map[E]struct{})
	
	// Only include elements that exist in this set but not in the other
	for element := range is.elements {
		if !other.Contains(element) {
			result[element] = struct{}{}
		}
	}
	
	return &ImmutableSet[E]{elements: result}
}

// IsSubsetOf checks if this set is a subset of another set
func (is *ImmutableSet[E]) IsSubsetOf(other Set[E]) bool {
	for element := range is.elements {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this set is a superset of another set
func (is *ImmutableSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(is)
}

// Iterator returns an iterator for the set
func (is *ImmutableSet[E]) Iterator() common.Iterator[E] {
	elements := is.ToSlice()
	return &immutableSetIterator[E]{
		elements: elements,
		index:    0,
	}
}

// ForEach executes the given function for each element in the set
func (is *ImmutableSet[E]) ForEach(fn func(E)) {
	for element := range is.elements {
		fn(element)
	}
}

// String returns a string representation of the set
func (is *ImmutableSet[E]) String() string {
	if is.IsEmpty() {
		return "{}"
	}
	
	var builder strings.Builder
	builder.WriteString("{")
	
	first := true
	for element := range is.elements {
		if !first {
			builder.WriteString(", ")
		}
		first = false
		builder.WriteString(fmt.Sprintf("%v", element))
	}
	
	builder.WriteString("}")
	return builder.String()
}

// immutableSetIterator implements Iterator for ImmutableSet
type immutableSetIterator[E comparable] struct {
	elements []E
	index    int
}

func (it *immutableSetIterator[E]) HasNext() bool {
	return it.index < len(it.elements)
}

func (it *immutableSetIterator[E]) Next() (E, bool) {
	var zero E
	if !it.HasNext() {
		return zero, false
	}
	
	element := it.elements[it.index]
	it.index++
	return element, true
}

func (it *immutableSetIterator[E]) Remove() bool {
	// Remove operation is not supported for immutable collections
	return false
}