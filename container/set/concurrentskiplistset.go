package set

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/chenjianyu/collections/container/common"
)

const (
	maxLevels = 32
	p         = 0.25
)

// ConcurrentSkipListSet is a thread-safe set implementation based on skip list
// Elements are stored in sorted order
type ConcurrentSkipListSet[E comparable] struct {
	head       *skipListNode[E]
	comparator func(a, b E) int
	size       int
	mutex      sync.RWMutex
	rand       *rand.Rand
}

// skipListNode represents a node in the skip list
type skipListNode[E comparable] struct {
	value E
	next  []*skipListNode[E]
}

// NewConcurrentSkipListSet creates a new ConcurrentSkipListSet using default comparator
// Default comparator prefers Comparable.CompareTo, otherwise falls back to natural generic ordering
func NewConcurrentSkipListSet[E comparable]() *ConcurrentSkipListSet[E] {
    head := &skipListNode[E]{
        next: make([]*skipListNode[E], maxLevels),
    }

    return &ConcurrentSkipListSet[E]{
        head:       head,
        comparator: common.CompareNatural[E],
        size:       0,
        rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// NewConcurrentSkipListSetWithComparator creates a ConcurrentSkipListSet with specified comparator
func NewConcurrentSkipListSetWithComparator[E comparable](comparator func(a, b E) int) *ConcurrentSkipListSet[E] {
	head := &skipListNode[E]{
		next: make([]*skipListNode[E], maxLevels),
	}

	return &ConcurrentSkipListSet[E]{
		head:       head,
		comparator: comparator,
		size:       0,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewConcurrentSkipListSetWithComparatorStrategy creates a ConcurrentSkipListSet with a ComparatorStrategy
func NewConcurrentSkipListSetWithComparatorStrategy[E comparable](strategy common.ComparatorStrategy[E]) *ConcurrentSkipListSet[E] {
    head := &skipListNode[E]{
        next: make([]*skipListNode[E], maxLevels),
    }

    return &ConcurrentSkipListSet[E]{
        head: head,
        comparator: func(a, b E) int {
            return strategy.Compare(a, b)
        },
        size: 0,
        rand: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// randomLevel generates a random level for new nodes
func (s *ConcurrentSkipListSet[E]) randomLevel() int {
	level := 1
	for s.rand.Float64() < p && level < maxLevels {
		level++
	}
	return level
}

// findPredecessors finds the predecessors of the target value at each level
func (s *ConcurrentSkipListSet[E]) findPredecessors(target E) []*skipListNode[E] {
	predecessors := make([]*skipListNode[E], maxLevels)
	current := s.head

	for level := maxLevels - 1; level >= 0; level-- {
		for current.next[level] != nil && s.comparator(current.next[level].value, target) < 0 {
			current = current.next[level]
		}
		predecessors[level] = current
	}

	return predecessors
}

// Add adds an element to the set
// Returns false if the set already contains the element, otherwise returns true
func (s *ConcurrentSkipListSet[E]) Add(element E) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	predecessors := s.findPredecessors(element)

	// Check if element already exists
	if predecessors[0].next[0] != nil && s.comparator(predecessors[0].next[0].value, element) == 0 {
		return false
	}

	// Create new node
	level := s.randomLevel()
	newNode := &skipListNode[E]{
		value: element,
		next:  make([]*skipListNode[E], level),
	}

	// Insert the new node
	for i := 0; i < level; i++ {
		newNode.next[i] = predecessors[i].next[i]
		predecessors[i].next[i] = newNode
	}

	s.size++
	return true
}

// Remove removes the specified element from the set
// Returns true if the set contained the element, otherwise returns false
func (s *ConcurrentSkipListSet[E]) Remove(element E) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	predecessors := s.findPredecessors(element)

	// Check if element exists
	nodeToRemove := predecessors[0].next[0]
	if nodeToRemove == nil || s.comparator(nodeToRemove.value, element) != 0 {
		return false
	}

	// Remove the node from all levels
	for level := 0; level < len(nodeToRemove.next); level++ {
		predecessors[level].next[level] = nodeToRemove.next[level]
	}

	s.size--
	return true
}

// Contains checks if the set contains the specified element
func (s *ConcurrentSkipListSet[E]) Contains(element E) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	current := s.head
	for level := maxLevels - 1; level >= 0; level-- {
		for current.next[level] != nil && s.comparator(current.next[level].value, element) < 0 {
			current = current.next[level]
		}
	}

	next := current.next[0]
	return next != nil && s.comparator(next.value, element) == 0
}

// Size returns the number of elements in the set
func (s *ConcurrentSkipListSet[E]) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.size
}

// IsEmpty checks if the set is empty
func (s *ConcurrentSkipListSet[E]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear clears all elements from the set
func (s *ConcurrentSkipListSet[E]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Reset head node
	for i := 0; i < maxLevels; i++ {
		s.head.next[i] = nil
	}
	s.size = 0
}

// ToSlice returns a slice containing all elements in the set
func (s *ConcurrentSkipListSet[E]) ToSlice() []E {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]E, 0, s.size)
	current := s.head.next[0]

	for current != nil {
		result = append(result, current.value)
		current = current.next[0]
	}

	return result
}

// ForEach executes the given operation for each element in the set
func (s *ConcurrentSkipListSet[E]) ForEach(fn func(E)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	current := s.head.next[0]
	for current != nil {
		fn(current.value)
		current = current.next[0]
	}
}

// Union returns a new set containing all elements from this set and the other set
func (s *ConcurrentSkipListSet[E]) Union(other Set[E]) Set[E] {
	result := NewConcurrentSkipListSetWithComparator(s.comparator)
	s.ForEach(func(element E) {
		result.Add(element)
	})
	other.ForEach(func(element E) {
		result.Add(element)
	})
	return result
}

// Intersection returns a new set containing elements that exist in both sets
func (s *ConcurrentSkipListSet[E]) Intersection(other Set[E]) Set[E] {
	result := NewConcurrentSkipListSetWithComparator(s.comparator)
	s.ForEach(func(element E) {
		if other.Contains(element) {
			result.Add(element)
		}
	})
	return result
}

// Difference returns a new set containing elements that exist in this set but not in the other set
func (s *ConcurrentSkipListSet[E]) Difference(other Set[E]) Set[E] {
	result := NewConcurrentSkipListSetWithComparator(s.comparator)
	s.ForEach(func(element E) {
		if !other.Contains(element) {
			result.Add(element)
		}
	})
	return result
}

// IsSubsetOf checks if this set is a subset of the other set
func (s *ConcurrentSkipListSet[E]) IsSubsetOf(other Set[E]) bool {
	if s.Size() > other.Size() {
		return false
	}

	isSubset := true
	s.ForEach(func(element E) {
		if !other.Contains(element) {
			isSubset = false
		}
	})
	return isSubset
}

// IsSupersetOf checks if this set is a superset of the other set
func (s *ConcurrentSkipListSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(s)
}

// Iterator returns an iterator for the set
func (s *ConcurrentSkipListSet[E]) Iterator() common.Iterator[E] {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	elements := s.ToSlice()
	return &concurrentSkipListSetIterator[E]{
		elements: elements,
		index:    0,
	}
}

// concurrentSkipListSetIterator implements Iterator for ConcurrentSkipListSet
type concurrentSkipListSetIterator[E comparable] struct {
	elements []E
	index    int
}

// HasNext returns true if there are more elements to iterate
func (it *concurrentSkipListSetIterator[E]) HasNext() bool {
	return it.index < len(it.elements)
}

// Next returns the next element
func (it *concurrentSkipListSetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		var zero E
		return zero, false
	}
	element := it.elements[it.index]
	it.index++
	return element, true
}

// Remove removes the current element (not supported for concurrent collections)
func (it *concurrentSkipListSetIterator[E]) Remove() bool {
	return false // Not supported for concurrent collections
}

// String returns the string representation of the set
func (s *ConcurrentSkipListSet[E]) String() string {
	if s.IsEmpty() {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{")

	first := true
	s.ForEach(func(element E) {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", element))
		first = false
	})

	sb.WriteString("}")
	return sb.String()
}
