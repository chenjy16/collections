// Package set provides set data structure implementations
package set

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// HashSet is a Set implementation based on hash table
type HashSet[E comparable] struct {
	buckets      [][]E
	size         int
	hashStrategy common.HashStrategy[E]
}

// New creates a new HashSet with default hash strategy
func New[E comparable]() *HashSet[E] {
	return NewWithHashStrategy(common.NewComparableHashStrategy[E]())
}

// NewWithHashStrategy creates a new HashSet with custom hash strategy
func NewWithHashStrategy[E comparable](hashStrategy common.HashStrategy[E]) *HashSet[E] {
	return &HashSet[E]{
		buckets:      make([][]E, 16),
		hashStrategy: hashStrategy,
	}
}

// FromSlice creates a new HashSet from a slice with default hash strategy
func FromSlice[E comparable](slice []E) *HashSet[E] {
	return FromSliceWithHashStrategy(slice, common.NewComparableHashStrategy[E]())
}

// FromSliceWithHashStrategy creates a new HashSet from a slice with custom hash strategy
func FromSliceWithHashStrategy[E comparable](slice []E, hashStrategy common.HashStrategy[E]) *HashSet[E] {
	set := NewWithHashStrategy(hashStrategy)
	for _, element := range slice {
		set.Add(element)
	}
	return set
}

// Add adds an element to the set
func (s *HashSet[E]) Add(element E) bool {
	index := s.hash(element) % len(s.buckets)
	bucket := s.buckets[index]

	// Check if the element is really the same (handle hash collisions)
	for _, existing := range bucket {
		if s.hashStrategy.Equals(existing, element) {
			return false
		}
	}

	s.buckets[index] = append(bucket, element)
	s.size++
	return true
}

// Remove removes the specified element from the set
func (s *HashSet[E]) Remove(element E) bool {
	index := s.hash(element) % len(s.buckets)
	bucket := s.buckets[index]

	for i, existing := range bucket {
		// Check if the element is really the same (handle hash collisions)
		if s.hashStrategy.Equals(existing, element) {
			s.buckets[index] = append(bucket[:i], bucket[i+1:]...)
			s.size--
			return true
		}
	}
	return false
}

// Contains checks if the set contains the specified element
func (s *HashSet[E]) Contains(element E) bool {
	index := s.hash(element) % len(s.buckets)
	bucket := s.buckets[index]

	for _, existing := range bucket {
		// Check if the element is really the same (handle hash collisions)
		if s.hashStrategy.Equals(existing, element) {
			return true
		}
	}
	return false
}

// Size returns the number of elements in the set
func (s *HashSet[E]) Size() int {
	return s.size
}

// IsEmpty checks if the set is empty
func (s *HashSet[E]) IsEmpty() bool {
	return s.size == 0
}

// Clear empties the set
func (s *HashSet[E]) Clear() {
	s.buckets = make([][]E, 16)
	s.size = 0
}

// ToSlice returns a slice containing all elements in the set
func (s *HashSet[E]) ToSlice() []E {
	result := make([]E, 0, s.size)
	for _, bucket := range s.buckets {
		result = append(result, bucket...)
	}
	return result
}

// ForEach executes the given operation on each element in the set
func (s *HashSet[E]) ForEach(fn func(E)) {
	for _, bucket := range s.buckets {
		for _, element := range bucket {
			fn(element)
		}
	}
}

// Union returns the union of this set and another set
func (s *HashSet[E]) Union(other Set[E]) Set[E] {
	result := NewWithHashStrategy(s.hashStrategy)

	// Add all elements from this set
	s.ForEach(func(element E) {
		result.Add(element)
	})

	// Add all elements from the other set
	other.ForEach(func(element E) {
		result.Add(element)
	})

	return result
}

// Intersection returns the intersection of this set and another set
func (s *HashSet[E]) Intersection(other Set[E]) Set[E] {
	result := NewWithHashStrategy(s.hashStrategy)

	// Add elements that exist in both sets
	s.ForEach(func(element E) {
		if other.Contains(element) {
			result.Add(element)
		}
	})

	return result
}

// Difference returns the difference of this set and another set
func (s *HashSet[E]) Difference(other Set[E]) Set[E] {
	result := NewWithHashStrategy(s.hashStrategy)

	// Add elements that are in this set but not in the other set
	s.ForEach(func(element E) {
		if !other.Contains(element) {
			result.Add(element)
		}
	})

	return result
}

// IsSubsetOf checks if this set is a subset of another set
func (s *HashSet[E]) IsSubsetOf(other Set[E]) bool {
	// Empty set is a subset of any set
	if s.IsEmpty() {
		return true
	}

	// If this set is larger than the other set, it cannot be a subset
	if s.Size() > other.Size() {
		return false
	}

	// Check if every element in this set is in the other set
	for _, bucket := range s.buckets {
		for _, element := range bucket {
			if !other.Contains(element) {
				return false
			}
		}
	}
	return true
}

// IsSupersetOf checks if this set is a superset of another set
func (s *HashSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(s)
}

// String returns the string representation of the set
func (s *HashSet[E]) String() string {
	if s.IsEmpty() {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	first := true
	for _, bucket := range s.buckets {
		for _, element := range bucket {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%v", element))
			first = false
		}
	}
	builder.WriteString("]")
	return builder.String()
}

// Iterator returns an iterator for traversing elements in the set
func (s *HashSet[E]) Iterator() common.Iterator[E] {
	return &hashSetIterator[E]{
		set:      s,
		elements: s.ToSlice(),
		cursor:   0,
		lastRet:  -1,
	}
}

// hashSetIterator is the iterator implementation for HashSet
type hashSetIterator[E comparable] struct {
	set      *HashSet[E]
	elements []E
	cursor   int // Index of the next element
	lastRet  int // Index of the last returned element, -1 if none
}

// HasNext checks if the iterator has a next element
func (it *hashSetIterator[E]) HasNext() bool {
	return it.cursor < len(it.elements)
}

// Next returns the next element in the iterator
func (it *hashSetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		var zero E
		return zero, false
	}
	element := it.elements[it.cursor]
	it.lastRet = it.cursor
	it.cursor++
	return element, true
}

// Remove removes the last element returned by the iterator
func (it *hashSetIterator[E]) Remove() bool {
	if it.lastRet == -1 {
		return false
	}
	removed := it.set.Remove(it.elements[it.lastRet])
	if removed {
		// Update the iterator's element list
		it.elements = it.set.ToSlice()
		it.cursor = 0 // Reset cursor since element list has changed
		it.lastRet = -1
	}
	return removed
}

// hash computes hash value for the element using the hash strategy
func (s *HashSet[E]) hash(element E) int {
	hashValue := int(s.hashStrategy.Hash(element))
	if hashValue < 0 {
		hashValue = -hashValue
	}
	return hashValue
}
