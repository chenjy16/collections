package multiset

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chenjianyu/collections/container/common"
)

// LinkedHashMultiset is a multiset implementation that maintains insertion order
// It combines a hash map for O(1) operations with a doubly-linked list for order preservation
type LinkedHashMultiset[E comparable] struct {
	counts map[E]*linkedEntry[E]
	head   *linkedEntry[E]
	tail   *linkedEntry[E]
	size   int
	mu     sync.RWMutex
}

type linkedEntry[E comparable] struct {
	element E
	count   int
	prev    *linkedEntry[E]
	next    *linkedEntry[E]
}

// NewLinkedHashMultiset creates a new empty LinkedHashMultiset
func NewLinkedHashMultiset[E comparable]() *LinkedHashMultiset[E] {
	ms := &LinkedHashMultiset[E]{
		counts: make(map[E]*linkedEntry[E]),
	}
	// Create sentinel nodes
	ms.head = &linkedEntry[E]{}
	ms.tail = &linkedEntry[E]{}
	ms.head.next = ms.tail
	ms.tail.prev = ms.head
	return ms
}

// NewLinkedHashMultisetFromSlice creates a new LinkedHashMultiset from a slice
func NewLinkedHashMultisetFromSlice[E comparable](elements []E) *LinkedHashMultiset[E] {
	ms := NewLinkedHashMultiset[E]()
	for _, element := range elements {
		ms.Add(element)
	}
	return ms
}

// Add adds one occurrence of the specified element
func (ms *LinkedHashMultiset[E]) Add(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	if entry, exists := ms.counts[element]; exists {
		prevCount := entry.count
		entry.count++
		ms.size++
		return prevCount
	}
	
	// Create new entry and add to end of list
	entry := &linkedEntry[E]{
		element: element,
		count:   1,
	}
	ms.counts[element] = entry
	ms.addToTail(entry)
	ms.size++
	return 0
}

// AddCount adds the specified number of occurrences of the element
func (ms *LinkedHashMultiset[E]) AddCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	if count == 0 {
		return ms.Count(element), nil
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	if entry, exists := ms.counts[element]; exists {
		prevCount := entry.count
		entry.count += count
		ms.size += count
		return prevCount, nil
	}
	
	// Create new entry and add to end of list
	entry := &linkedEntry[E]{
		element: element,
		count:   count,
	}
	ms.counts[element] = entry
	ms.addToTail(entry)
	ms.size += count
	return 0, nil
}

// Remove removes one occurrence of the specified element
func (ms *LinkedHashMultiset[E]) Remove(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	entry, exists := ms.counts[element]
	if !exists {
		return 0
	}
	
	prevCount := entry.count
	entry.count--
	ms.size--
	
	if entry.count == 0 {
		ms.removeFromList(entry)
		delete(ms.counts, element)
	}
	
	return prevCount
}

// RemoveCount removes the specified number of occurrences of the element
func (ms *LinkedHashMultiset[E]) RemoveCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	if count == 0 {
		return ms.Count(element), nil
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	entry, exists := ms.counts[element]
	if !exists {
		return 0, nil
	}
	
	prevCount := entry.count
	removeCount := count
	if removeCount > prevCount {
		removeCount = prevCount
	}
	
	entry.count -= removeCount
	ms.size -= removeCount
	
	if entry.count == 0 {
		ms.removeFromList(entry)
		delete(ms.counts, element)
	}
	
	return prevCount, nil
}

// RemoveAll removes all occurrences of the specified element
func (ms *LinkedHashMultiset[E]) RemoveAll(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	entry, exists := ms.counts[element]
	if !exists {
		return 0
	}
	
	prevCount := entry.count
	ms.size -= prevCount
	ms.removeFromList(entry)
	delete(ms.counts, element)
	
	return prevCount
}

// Count returns the number of occurrences of the specified element
func (ms *LinkedHashMultiset[E]) Count(element E) int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	if entry, exists := ms.counts[element]; exists {
		return entry.count
	}
	return 0
}

// SetCount sets the count of the specified element to the given value
func (ms *LinkedHashMultiset[E]) SetCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	entry, exists := ms.counts[element]
	var prevCount int
	
	if exists {
		prevCount = entry.count
		if count == 0 {
			ms.size -= prevCount
			ms.removeFromList(entry)
			delete(ms.counts, element)
		} else {
			ms.size += count - prevCount
			entry.count = count
		}
	} else {
		prevCount = 0
		if count > 0 {
			newEntry := &linkedEntry[E]{
				element: element,
				count:   count,
			}
			ms.counts[element] = newEntry
			ms.addToTail(newEntry)
			ms.size += count
		}
	}
	
	return prevCount, nil
}

// Contains checks if the multiset contains the specified element
func (ms *LinkedHashMultiset[E]) Contains(element E) bool {
	return ms.Count(element) > 0
}

// IsEmpty returns true if the multiset contains no elements
func (ms *LinkedHashMultiset[E]) IsEmpty() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.size == 0
}

// Size returns the number of distinct elements in the multiset
func (ms *LinkedHashMultiset[E]) Size() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.counts)
}

// TotalSize returns the total number of elements (including duplicates)
func (ms *LinkedHashMultiset[E]) TotalSize() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.size
}

// DistinctElements returns the number of distinct elements
func (ms *LinkedHashMultiset[E]) DistinctElements() int {
	return ms.Size()
}

// Clear removes all elements from the multiset
func (ms *LinkedHashMultiset[E]) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.counts = make(map[E]*linkedEntry[E])
	ms.head.next = ms.tail
	ms.tail.prev = ms.head
	ms.size = 0
}

// ElementSet returns a slice of distinct elements in insertion order
func (ms *LinkedHashMultiset[E]) ElementSet() []E {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	elements := make([]E, 0, len(ms.counts))
	current := ms.head.next
	for current != ms.tail {
		elements = append(elements, current.element)
		current = current.next
	}
	return elements
}

// EntrySet returns a slice of entries (element-count pairs) in insertion order
func (ms *LinkedHashMultiset[E]) EntrySet() []Entry[E] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	entries := make([]Entry[E], 0, len(ms.counts))
	current := ms.head.next
	for current != ms.tail {
		entries = append(entries, Entry[E]{Element: current.element, Count: current.count})
		current = current.next
	}
	return entries
}

// ToSlice returns a slice containing all elements (including duplicates) in insertion order
func (ms *LinkedHashMultiset[E]) ToSlice() []E {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	result := make([]E, 0, ms.size)
	current := ms.head.next
	for current != ms.tail {
		for i := 0; i < current.count; i++ {
			result = append(result, current.element)
		}
		current = current.next
	}
	return result
}

// Iterator returns an iterator over the multiset elements in insertion order
func (ms *LinkedHashMultiset[E]) Iterator() common.Iterator[E] {
    return &baseMultisetIterator[E]{
        entries: ms.EntrySet(),
        index:   0,
        current: 0,
        removeFunc: func(elem E) bool {
            ms.Remove(elem)
            return true
        },
        refresh: func() []Entry[E] { return ms.EntrySet() },
    }
}

// ForEach executes the given function for each element in the multiset in insertion order
func (ms *LinkedHashMultiset[E]) ForEach(fn func(E)) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	current := ms.head.next
	for current != ms.tail {
		for i := 0; i < current.count; i++ {
			fn(current.element)
		}
		current = current.next
	}
}

// Union returns a new multiset containing the union of this and another multiset
func (ms *LinkedHashMultiset[E]) Union(other Multiset[E]) Multiset[E] {
	result := NewLinkedHashMultiset[E]()
	
	// Add all elements from this multiset (preserving order)
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
func (ms *LinkedHashMultiset[E]) Intersection(other Multiset[E]) Multiset[E] {
	result := NewLinkedHashMultiset[E]()
	
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
func (ms *LinkedHashMultiset[E]) Difference(other Multiset[E]) Multiset[E] {
	result := NewLinkedHashMultiset[E]()
	
	for _, entry := range ms.EntrySet() {
		otherCount := other.Count(entry.Element)
		if entry.Count > otherCount {
			result.AddCount(entry.Element, entry.Count-otherCount)
		}
	}
	
	return result
}

// IsSubsetOf checks if this multiset is a subset of another
func (ms *LinkedHashMultiset[E]) IsSubsetOf(other Multiset[E]) bool {
	for _, entry := range ms.EntrySet() {
		if other.Count(entry.Element) < entry.Count {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this multiset is a superset of another
func (ms *LinkedHashMultiset[E]) IsSupersetOf(other Multiset[E]) bool {
	return other.IsSubsetOf(ms)
}

// String returns a string representation of the multiset
func (ms *LinkedHashMultiset[E]) String() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	if ms.size == 0 {
		return "LinkedHashMultiset[]"
	}
	
	var builder strings.Builder
	builder.WriteString("LinkedHashMultiset[")
	
	first := true
	current := ms.head.next
	for current != ms.tail {
		if !first {
			builder.WriteString(", ")
		}
		if current.count == 1 {
			builder.WriteString(fmt.Sprintf("%v", current.element))
		} else {
			builder.WriteString(fmt.Sprintf("%v x %d", current.element, current.count))
		}
		first = false
		current = current.next
	}
	
	builder.WriteString("]")
	return builder.String()
}

// Helper methods for doubly-linked list operations
func (ms *LinkedHashMultiset[E]) addToTail(entry *linkedEntry[E]) {
	prev := ms.tail.prev
	prev.next = entry
	entry.prev = prev
	entry.next = ms.tail
	ms.tail.prev = entry
}

func (ms *LinkedHashMultiset[E]) removeFromList(entry *linkedEntry[E]) {
	entry.prev.next = entry.next
	entry.next.prev = entry.prev
}

// linkedHashMultisetIterator implements Iterator for LinkedHashMultiset
type linkedHashMultisetIterator[E comparable] struct {
	multiset *LinkedHashMultiset[E]
	entries  []Entry[E]
	index    int
	current  int
}

func (it *linkedHashMultisetIterator[E]) HasNext() bool {
	return it.index < len(it.entries) && (it.current < it.entries[it.index].Count || it.index+1 < len(it.entries))
}

func (it *linkedHashMultisetIterator[E]) Next() (E, bool) {
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

func (it *linkedHashMultisetIterator[E]) Reset() {
	it.entries = it.multiset.EntrySet()
	it.index = 0
	it.current = 0
}

func (it *linkedHashMultisetIterator[E]) Remove() bool {
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