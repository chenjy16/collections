package multiset

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chenjianyu/collections/container/common"
)

// HashMultiset is a multiset implementation based on a hash map
// It provides O(1) average time complexity for basic operations
type HashMultiset[E comparable] struct {
	counts map[E]int
	size   int
	mu     sync.RWMutex
}

// NewHashMultiset creates a new empty HashMultiset
func NewHashMultiset[E comparable]() *HashMultiset[E] {
	return &HashMultiset[E]{
		counts: make(map[E]int),
		size:   0,
	}
}

// NewHashMultisetFromSlice creates a new HashMultiset from a slice
func NewHashMultisetFromSlice[E comparable](elements []E) *HashMultiset[E] {
	ms := NewHashMultiset[E]()
	for _, element := range elements {
		ms.Add(element)
	}
	return ms
}

// Add adds one occurrence of the specified element
func (ms *HashMultiset[E]) Add(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := ms.counts[element]
	ms.counts[element] = prevCount + 1
	ms.size++
	return prevCount
}

// AddCount adds the specified number of occurrences of the element
func (ms *HashMultiset[E]) AddCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	if count == 0 {
		return ms.Count(element), nil
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := ms.counts[element]
	ms.counts[element] = prevCount + count
	ms.size += count
	return prevCount, nil
}

// Remove removes one occurrence of the specified element
func (ms *HashMultiset[E]) Remove(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := ms.counts[element]
	if prevCount > 0 {
		if prevCount == 1 {
			delete(ms.counts, element)
		} else {
			ms.counts[element] = prevCount - 1
		}
		ms.size--
	}
	return prevCount
}

// RemoveCount removes the specified number of occurrences of the element
func (ms *HashMultiset[E]) RemoveCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	if count == 0 {
		return ms.Count(element), nil
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := ms.counts[element]
	if prevCount > 0 {
		removeCount := count
		if removeCount > prevCount {
			removeCount = prevCount
		}
		
		newCount := prevCount - removeCount
		if newCount == 0 {
			delete(ms.counts, element)
		} else {
			ms.counts[element] = newCount
		}
		ms.size -= removeCount
	}
	return prevCount, nil
}

// RemoveAll removes all occurrences of the specified element
func (ms *HashMultiset[E]) RemoveAll(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := ms.counts[element]
	if prevCount > 0 {
		delete(ms.counts, element)
		ms.size -= prevCount
	}
	return prevCount
}

// Count returns the number of occurrences of the specified element
func (ms *HashMultiset[E]) Count(element E) int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.counts[element]
}

// SetCount sets the count of the specified element to the given value
func (ms *HashMultiset[E]) SetCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := ms.counts[element]
	
	if count == 0 {
		if prevCount > 0 {
			delete(ms.counts, element)
			ms.size -= prevCount
		}
	} else {
		ms.counts[element] = count
		ms.size += count - prevCount
	}
	
	return prevCount, nil
}

// Contains checks if the multiset contains the specified element
func (ms *HashMultiset[E]) Contains(element E) bool {
	return ms.Count(element) > 0
}

// IsEmpty returns true if the multiset contains no elements
func (ms *HashMultiset[E]) IsEmpty() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.size == 0
}

// Size returns the number of distinct elements in the multiset
func (ms *HashMultiset[E]) Size() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.counts)
}

// TotalSize returns the total number of elements (including duplicates)
func (ms *HashMultiset[E]) TotalSize() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.size
}

// DistinctElements returns the number of distinct elements
func (ms *HashMultiset[E]) DistinctElements() int {
	return ms.Size()
}

// Clear removes all elements from the multiset
func (ms *HashMultiset[E]) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counts = make(map[E]int)
	ms.size = 0
}

// ElementSet returns a slice of distinct elements
func (ms *HashMultiset[E]) ElementSet() []E {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	elements := make([]E, 0, len(ms.counts))
	for element := range ms.counts {
		elements = append(elements, element)
	}
	return elements
}

// EntrySet returns a slice of entries (element-count pairs)
func (ms *HashMultiset[E]) EntrySet() []Entry[E] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	entries := make([]Entry[E], 0, len(ms.counts))
	for element, count := range ms.counts {
		entries = append(entries, Entry[E]{Element: element, Count: count})
	}
	return entries
}

// ToSlice returns a slice containing all elements (including duplicates)
func (ms *HashMultiset[E]) ToSlice() []E {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	result := make([]E, 0, ms.size)
	for element, count := range ms.counts {
		for i := 0; i < count; i++ {
			result = append(result, element)
		}
	}
	return result
}

// Iterator returns an iterator over the multiset elements
func (ms *HashMultiset[E]) Iterator() common.Iterator[E] {
	return &hashMultisetIterator[E]{
		multiset: ms,
		entries:  ms.EntrySet(),
		index:    0,
		current:  0,
	}
}

// ForEach executes the given function for each element in the multiset
func (ms *HashMultiset[E]) ForEach(fn func(E)) {
	for _, entry := range ms.EntrySet() {
		for i := 0; i < entry.Count; i++ {
			fn(entry.Element)
		}
	}
}

// Union returns a new multiset containing the union of this and another multiset
func (ms *HashMultiset[E]) Union(other Multiset[E]) Multiset[E] {
	result := NewHashMultiset[E]()
	
	// Add all elements from this multiset
	for _, entry := range ms.EntrySet() {
		result.AddCount(entry.Element, entry.Count)
	}
	
	// Add elements from other multiset, taking maximum count
	for _, entry := range other.EntrySet() {
		currentCount := result.Count(entry.Element)
		if entry.Count > currentCount {
			result.SetCount(entry.Element, entry.Count)
		}
	}
	
	return result
}

// Intersection returns a new multiset containing the intersection
func (ms *HashMultiset[E]) Intersection(other Multiset[E]) Multiset[E] {
	result := NewHashMultiset[E]()
	
	for _, entry := range ms.EntrySet() {
		otherCount := other.Count(entry.Element)
		if otherCount > 0 {
			minCount := entry.Count
			if otherCount < minCount {
				minCount = otherCount
			}
			result.AddCount(entry.Element, minCount)
		}
	}
	
	return result
}

// Difference returns a new multiset containing elements in this but not in other
func (ms *HashMultiset[E]) Difference(other Multiset[E]) Multiset[E] {
	result := NewHashMultiset[E]()
	
	for _, entry := range ms.EntrySet() {
		otherCount := other.Count(entry.Element)
		if entry.Count > otherCount {
			result.AddCount(entry.Element, entry.Count-otherCount)
		}
	}
	
	return result
}

// IsSubsetOf checks if this multiset is a subset of another
func (ms *HashMultiset[E]) IsSubsetOf(other Multiset[E]) bool {
	for _, entry := range ms.EntrySet() {
		if other.Count(entry.Element) < entry.Count {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this multiset is a superset of another
func (ms *HashMultiset[E]) IsSupersetOf(other Multiset[E]) bool {
	return other.IsSubsetOf(ms)
}

// String returns a string representation of the multiset
func (ms *HashMultiset[E]) String() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	if ms.size == 0 {
		return "HashMultiset[]"
	}
	
	var builder strings.Builder
	builder.WriteString("HashMultiset[")
	
	first := true
	for element, count := range ms.counts {
		if !first {
			builder.WriteString(", ")
		}
		if count == 1 {
			builder.WriteString(fmt.Sprintf("%v", element))
		} else {
			builder.WriteString(fmt.Sprintf("%v x %d", element, count))
		}
		first = false
	}
	
	builder.WriteString("]")
	return builder.String()
}

// hashMultisetIterator implements Iterator for HashMultiset
type hashMultisetIterator[E comparable] struct {
	multiset *HashMultiset[E]
	entries  []Entry[E]
	index    int
	current  int
}

func (it *hashMultisetIterator[E]) HasNext() bool {
	return it.index < len(it.entries) && (it.current < it.entries[it.index].Count || it.index+1 < len(it.entries))
}

func (it *hashMultisetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		var zero E
		return zero, false
	}
	
	if it.current >= it.entries[it.index].Count {
		it.index++
		it.current = 0
	}
	
	element := it.entries[it.index].Element
	it.current++
	return element, true
}

func (it *hashMultisetIterator[E]) Reset() {
	it.entries = it.multiset.EntrySet()
	it.index = 0
	it.current = 0
}

func (it *hashMultisetIterator[E]) Remove() bool {
	if it.index >= len(it.entries) || it.current == 0 {
		return false
	}
	
	element := it.entries[it.index].Element
	it.multiset.Remove(element)
	
	// Refresh entries after removal
	it.entries = it.multiset.EntrySet()
	if it.index >= len(it.entries) {
		it.index = len(it.entries)
		it.current = 0
	} else if it.current > it.entries[it.index].Count {
		it.current = it.entries[it.index].Count
	}
	
	return true
}