// Package multiset provides multiset data structure implementations
// A multiset (also known as a bag) is a collection that allows duplicate elements
// and keeps track of the count of each element
package multiset

import (
	"github.com/chenjianyu/collections/container/common"
)

// Entry represents an element and its count in the multiset
type Entry[E comparable] struct {
	Element E
	Count   int
}

// Multiset represents a collection that allows duplicate elements
// and keeps track of the count of each element
type Multiset[E comparable] interface {
	common.Container[E]
	common.Iterable[E]

	// Add adds one occurrence of the specified element to this multiset
	// Returns the previous count of the element (0 if not present)
	Add(element E) int

	// AddCount adds the specified number of occurrences of the element
	// Returns the previous count of the element (0 if not present)
	AddCount(element E, count int) int

	// Remove removes one occurrence of the specified element from this multiset
	// Returns the previous count of the element (0 if not present)
	Remove(element E) int

	// RemoveCount removes the specified number of occurrences of the element
	// Returns the previous count of the element (0 if not present)
	RemoveCount(element E, count int) int

	// RemoveAll removes all occurrences of the specified element
	// Returns the previous count of the element (0 if not present)
	RemoveAll(element E) int

	// Count returns the number of occurrences of the specified element
	Count(element E) int

	// SetCount sets the count of the specified element to the given value
	// Returns the previous count of the element (0 if not present)
	SetCount(element E, count int) int

	// ElementSet returns a set view of the distinct elements in this multiset
	ElementSet() []E

	// EntrySet returns a set view of the entries (element-count pairs) in this multiset
	EntrySet() []Entry[E]

	// TotalSize returns the total number of elements in this multiset (including duplicates)
	TotalSize() int

	// DistinctElements returns the number of distinct elements in this multiset
	DistinctElements() int

	// ToSlice returns a slice containing all elements in this multiset (including duplicates)
	ToSlice() []E

	// Union returns a new multiset containing the union of this multiset and another
	// The count of each element is the maximum count from either multiset
	Union(other Multiset[E]) Multiset[E]

	// Intersection returns a new multiset containing the intersection of this multiset and another
	// The count of each element is the minimum count from either multiset
	Intersection(other Multiset[E]) Multiset[E]

	// Difference returns a new multiset containing elements in this multiset but not in the other
	// The count of each element is the count in this multiset minus the count in the other
	Difference(other Multiset[E]) Multiset[E]

	// IsSubsetOf checks if this multiset is a subset of another multiset
	// Returns true if for every element, its count in this multiset is <= its count in the other
	IsSubsetOf(other Multiset[E]) bool

	// IsSupersetOf checks if this multiset is a superset of another multiset
	IsSupersetOf(other Multiset[E]) bool
}