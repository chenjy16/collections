// Package set provides set data structure implementations
package set

import (
	"github.com/chenjianyu/collections/container/common"
)

// Set represents a collection that contains no duplicate elements
type Set[E comparable] interface {
	common.Container[E]
	common.Iterable[E]

	// Add adds an element to the set
	// Returns false if the set already contains the element, otherwise returns true
	Add(element E) bool

	// Remove removes the specified element from the set
	// Returns true if the set contained the element, otherwise returns false
	Remove(element E) bool

	// Contains checks if the set contains the specified element
	Contains(element E) bool

	// ToSlice returns a slice containing all elements in the set
	ToSlice() []E

	// Union returns the union of this set and another set
	Union(other Set[E]) Set[E]

	// Intersection returns the intersection of this set and another set
	Intersection(other Set[E]) Set[E]

	// Difference returns the difference of this set and another set
	// i.e., elements that are in this set but not in the other set
	Difference(other Set[E]) Set[E]

	// IsSubsetOf checks if this set is a subset of another set
	IsSubsetOf(other Set[E]) bool

	// IsSupersetOf checks if this set is a superset of another set
	IsSupersetOf(other Set[E]) bool
}
