// Package set provides set data structure implementations
package set

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// linkedHashSetNode represents a node in the LinkedHashSet
type linkedHashSetNode[E comparable] struct {
	data E
	prev *linkedHashSetNode[E]
	next *linkedHashSetNode[E]
}

// LinkedHashSet is a Set implementation that maintains insertion order
// It combines the fast lookup of a hash table with the ordering of a linked list
type LinkedHashSet[E comparable] struct {
	buckets      [][]E                       // Hash table for fast lookup
	nodeMap      map[E]*linkedHashSetNode[E] // Map from element to node for O(1) access
	head         *linkedHashSetNode[E]       // Head of the doubly linked list
	tail         *linkedHashSetNode[E]       // Tail of the doubly linked list
	size         int
	hashStrategy common.HashStrategy[E] // Custom hash strategy
}

// NewLinkedHashSet creates a new LinkedHashSet with default hash strategy
func NewLinkedHashSet[E comparable]() *LinkedHashSet[E] {
	return NewLinkedHashSetWithHashStrategy(common.NewComparableHashStrategy[E]())
}

// NewLinkedHashSetWithHashStrategy creates a new LinkedHashSet with custom hash strategy
func NewLinkedHashSetWithHashStrategy[E comparable](hashStrategy common.HashStrategy[E]) *LinkedHashSet[E] {
	return &LinkedHashSet[E]{
		buckets:      make([][]E, 16),
		nodeMap:      make(map[E]*linkedHashSetNode[E]),
		hashStrategy: hashStrategy,
	}
}

// LinkedHashSetFromSlice creates a new LinkedHashSet from a slice with default hash strategy
func LinkedHashSetFromSlice[E comparable](slice []E) *LinkedHashSet[E] {
	return LinkedHashSetFromSliceWithHashStrategy(slice, common.NewComparableHashStrategy[E]())
}

// LinkedHashSetFromSliceWithHashStrategy creates a new LinkedHashSet from a slice with custom hash strategy
func LinkedHashSetFromSliceWithHashStrategy[E comparable](slice []E, hashStrategy common.HashStrategy[E]) *LinkedHashSet[E] {
	set := NewLinkedHashSetWithHashStrategy(hashStrategy)
	for _, element := range slice {
		set.Add(element)
	}
	return set
}

// hash computes the hash value for an element using the hash strategy
func (s *LinkedHashSet[E]) hash(element E) uint32 {
	return uint32(s.hashStrategy.Hash(element))
}

// Add adds an element to the set
func (s *LinkedHashSet[E]) Add(element E) bool {
	// Check if element already exists using hash strategy
	for existingElement := range s.nodeMap {
		if s.hashStrategy.Equals(existingElement, element) {
			return false
		}
	}

	// Add to hash table
	index := s.hash(element) % uint32(len(s.buckets))
	s.buckets[index] = append(s.buckets[index], element)

	// Create new node and add to linked list
	newNode := &linkedHashSetNode[E]{data: element}
	s.nodeMap[element] = newNode

	if s.head == nil {
		// First element
		s.head = newNode
		s.tail = newNode
	} else {
		// Add to tail
		s.tail.next = newNode
		newNode.prev = s.tail
		s.tail = newNode
	}

	s.size++

	// Resize if load factor is too high
	if s.size > len(s.buckets)*2 {
		s.resize()
	}

	return true
}

// Remove removes the specified element from the set
func (s *LinkedHashSet[E]) Remove(element E) bool {
	// Find the element using hash strategy
	var nodeToRemove *linkedHashSetNode[E]
	var keyToRemove E
	found := false

	for existingElement, node := range s.nodeMap {
		if s.hashStrategy.Equals(existingElement, element) {
			nodeToRemove = node
			keyToRemove = existingElement
			found = true
			break
		}
	}

	if !found {
		return false
	}

	// Remove from hash table
	index := s.hash(keyToRemove) % uint32(len(s.buckets))
	bucket := s.buckets[index]
	for i, existing := range bucket {
		if s.hashStrategy.Equals(existing, element) {
			s.buckets[index] = append(bucket[:i], bucket[i+1:]...)
			break
		}
	}

	// Remove from linked list
	if nodeToRemove.prev != nil {
		nodeToRemove.prev.next = nodeToRemove.next
	} else {
		s.head = nodeToRemove.next
	}

	if nodeToRemove.next != nil {
		nodeToRemove.next.prev = nodeToRemove.prev
	} else {
		s.tail = nodeToRemove.prev
	}

	// Remove from node map
	delete(s.nodeMap, keyToRemove)
	s.size--

	return true
}

// Contains checks if the set contains the specified element
func (s *LinkedHashSet[E]) Contains(element E) bool {
	for existingElement := range s.nodeMap {
		if s.hashStrategy.Equals(existingElement, element) {
			return true
		}
	}
	return false
}

// Size returns the number of elements in the set
func (s *LinkedHashSet[E]) Size() int {
	return s.size
}

// IsEmpty checks if the set is empty
func (s *LinkedHashSet[E]) IsEmpty() bool {
	return s.size == 0
}

// Clear empties the set
func (s *LinkedHashSet[E]) Clear() {
	s.buckets = make([][]E, 16)
	s.nodeMap = make(map[E]*linkedHashSetNode[E])
	s.head = nil
	s.tail = nil
	s.size = 0
}

// ToSlice returns a slice containing all elements in insertion order
func (s *LinkedHashSet[E]) ToSlice() []E {
	result := make([]E, 0, s.size)
	current := s.head
	for current != nil {
		result = append(result, current.data)
		current = current.next
	}
	return result
}

// ForEach executes the given operation on each element in insertion order
func (s *LinkedHashSet[E]) ForEach(f func(E)) {
	current := s.head
	for current != nil {
		f(current.data)
		current = current.next
	}
}

// Union returns the union of this set and another set
func (s *LinkedHashSet[E]) Union(other Set[E]) Set[E] {
	result := NewLinkedHashSetWithHashStrategy(s.hashStrategy)

	// Add all elements from this set (maintains order)
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
func (s *LinkedHashSet[E]) Intersection(other Set[E]) Set[E] {
	result := NewLinkedHashSetWithHashStrategy(s.hashStrategy)

	// Add elements that exist in both sets (maintains order from this set)
	s.ForEach(func(element E) {
		if other.Contains(element) {
			result.Add(element)
		}
	})

	return result
}

// Difference returns the difference of this set and another set
func (s *LinkedHashSet[E]) Difference(other Set[E]) Set[E] {
	result := NewLinkedHashSetWithHashStrategy(s.hashStrategy)

	// Add elements that are in this set but not in the other set
	s.ForEach(func(element E) {
		if !other.Contains(element) {
			result.Add(element)
		}
	})

	return result
}

// IsSubsetOf checks if this set is a subset of another set
func (s *LinkedHashSet[E]) IsSubsetOf(other Set[E]) bool {
	// Empty set is a subset of any set
	if s.IsEmpty() {
		return true
	}

	// Check if all elements in this set are in the other set
	isSubset := true
	s.ForEach(func(element E) {
		if !other.Contains(element) {
			isSubset = false
		}
	})

	return isSubset
}

// IsSupersetOf checks if this set is a superset of another set
func (s *LinkedHashSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(s)
}

// String returns the string representation of the set in insertion order
func (s *LinkedHashSet[E]) String() string {
	if s.IsEmpty() {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	current := s.head
	first := true
	for current != nil {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", current.data))
		first = false
		current = current.next
	}
	builder.WriteString("]")
	return builder.String()
}

// Iterator returns an iterator for traversing elements in insertion order
func (s *LinkedHashSet[E]) Iterator() common.Iterator[E] {
	return &linkedHashSetIterator[E]{
		set:     s,
		current: s.head,
		lastRet: nil,
	}
}

// resize increases the capacity of the hash table
func (s *LinkedHashSet[E]) resize() {
	oldBuckets := s.buckets
	s.buckets = make([][]E, len(oldBuckets)*2)

	// Rehash all elements
	for _, bucket := range oldBuckets {
		for _, element := range bucket {
			index := s.hash(element) % uint32(len(s.buckets))
			s.buckets[index] = append(s.buckets[index], element)
		}
	}
}

// linkedHashSetIterator is the iterator implementation for LinkedHashSet
type linkedHashSetIterator[E comparable] struct {
	set     *LinkedHashSet[E]
	current *linkedHashSetNode[E]
	lastRet *linkedHashSetNode[E]
}

// HasNext checks if the iterator has a next element
func (it *linkedHashSetIterator[E]) HasNext() bool {
	return it.current != nil
}

// Next returns the next element in insertion order
func (it *linkedHashSetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		var zero E
		return zero, false
	}

	element := it.current.data
	it.lastRet = it.current
	it.current = it.current.next
	return element, true
}

// Remove removes the last element returned by the iterator
func (it *linkedHashSetIterator[E]) Remove() bool {
	if it.lastRet == nil {
		return false
	}

	removed := it.set.Remove(it.lastRet.data)
	it.lastRet = nil
	return removed
}
