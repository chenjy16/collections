package multiset

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// ImmutableMultiset is an immutable multiset implementation
// Once created, it cannot be modified. All modification operations return new instances
type ImmutableMultiset[E comparable] struct {
	counts map[E]int
	size   int
}

// NewImmutableMultiset creates a new empty ImmutableMultiset
func NewImmutableMultiset[E comparable]() *ImmutableMultiset[E] {
	return &ImmutableMultiset[E]{
		counts: make(map[E]int),
		size:   0,
	}
}

// NewImmutableMultisetFromSlice creates a new ImmutableMultiset from a slice
func NewImmutableMultisetFromSlice[E comparable](elements []E) *ImmutableMultiset[E] {
	counts := make(map[E]int)
	size := 0
	
	for _, element := range elements {
		counts[element]++
		size++
	}
	
	return &ImmutableMultiset[E]{
		counts: counts,
		size:   size,
	}
}

// NewImmutableMultisetFromEntries creates a new ImmutableMultiset from entries
func NewImmutableMultisetFromEntries[E comparable](entries []Entry[E]) *ImmutableMultiset[E] {
	counts := make(map[E]int)
	size := 0
	
	for _, entry := range entries {
		if entry.Count > 0 {
			counts[entry.Element] = entry.Count
			size += entry.Count
		}
	}
	
	return &ImmutableMultiset[E]{
		counts: counts,
		size:   size,
	}
}

// copyMap creates a deep copy of the counts map
func (ms *ImmutableMultiset[E]) copyMap() map[E]int {
	newCounts := make(map[E]int, len(ms.counts))
	for k, v := range ms.counts {
		newCounts[k] = v
	}
	return newCounts
}

// Add returns a new ImmutableMultiset with one occurrence of the element added
func (ms *ImmutableMultiset[E]) Add(element E) int {
	newCounts := ms.copyMap()
	prevCount := newCounts[element]
	newCounts[element] = prevCount + 1
	
	// Note: This violates the interface contract but is needed for immutability
	// The returned int is the previous count, but we can't modify this instance
	return prevCount
}

// AddCount returns a new ImmutableMultiset with the specified count of elements added
func (ms *ImmutableMultiset[E]) AddCount(element E, count int) int {
	if count < 0 {
		panic("count cannot be negative")
	}
	if count == 0 {
		return ms.Count(element)
	}
	
	newCounts := ms.copyMap()
	prevCount := newCounts[element]
	newCounts[element] = prevCount + count
	
	return prevCount
}

// WithAdd returns a new ImmutableMultiset with one occurrence of the element added
func (ms *ImmutableMultiset[E]) WithAdd(element E) *ImmutableMultiset[E] {
	newCounts := ms.copyMap()
	newCounts[element]++
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   ms.size + 1,
	}
}

// WithAddCount returns a new ImmutableMultiset with the specified count of elements added
func (ms *ImmutableMultiset[E]) WithAddCount(element E, count int) *ImmutableMultiset[E] {
	if count < 0 {
		panic("count cannot be negative")
	}
	if count == 0 {
		return ms
	}
	
	newCounts := ms.copyMap()
	newCounts[element] += count
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   ms.size + count,
	}
}

// Remove returns the previous count (but doesn't actually modify this immutable instance)
func (ms *ImmutableMultiset[E]) Remove(element E) int {
	return ms.counts[element]
}

// RemoveCount returns the previous count (but doesn't actually modify this immutable instance)
func (ms *ImmutableMultiset[E]) RemoveCount(element E, count int) int {
	if count < 0 {
		panic("count cannot be negative")
	}
	return ms.counts[element]
}

// RemoveAll returns the previous count (but doesn't actually modify this immutable instance)
func (ms *ImmutableMultiset[E]) RemoveAll(element E) int {
	return ms.counts[element]
}

// WithRemove returns a new ImmutableMultiset with one occurrence of the element removed
func (ms *ImmutableMultiset[E]) WithRemove(element E) *ImmutableMultiset[E] {
	currentCount := ms.counts[element]
	if currentCount == 0 {
		return ms
	}
	
	newCounts := ms.copyMap()
	if currentCount == 1 {
		delete(newCounts, element)
	} else {
		newCounts[element] = currentCount - 1
	}
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   ms.size - 1,
	}
}

// WithRemoveCount returns a new ImmutableMultiset with the specified count of elements removed
func (ms *ImmutableMultiset[E]) WithRemoveCount(element E, count int) *ImmutableMultiset[E] {
	if count < 0 {
		panic("count cannot be negative")
	}
	if count == 0 {
		return ms
	}
	
	currentCount := ms.counts[element]
	if currentCount == 0 {
		return ms
	}
	
	removeCount := count
	if removeCount > currentCount {
		removeCount = currentCount
	}
	
	newCounts := ms.copyMap()
	newCount := currentCount - removeCount
	if newCount == 0 {
		delete(newCounts, element)
	} else {
		newCounts[element] = newCount
	}
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   ms.size - removeCount,
	}
}

// WithRemoveAll returns a new ImmutableMultiset with all occurrences of the element removed
func (ms *ImmutableMultiset[E]) WithRemoveAll(element E) *ImmutableMultiset[E] {
	currentCount := ms.counts[element]
	if currentCount == 0 {
		return ms
	}
	
	newCounts := ms.copyMap()
	delete(newCounts, element)
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   ms.size - currentCount,
	}
}

// Count returns the number of occurrences of the specified element
func (ms *ImmutableMultiset[E]) Count(element E) int {
	return ms.counts[element]
}

// SetCount returns the previous count (but doesn't actually modify this immutable instance)
func (ms *ImmutableMultiset[E]) SetCount(element E, count int) int {
	if count < 0 {
		panic("count cannot be negative")
	}
	return ms.counts[element]
}

// WithSetCount returns a new ImmutableMultiset with the element count set to the specified value
func (ms *ImmutableMultiset[E]) WithSetCount(element E, count int) *ImmutableMultiset[E] {
	if count < 0 {
		panic("count cannot be negative")
	}
	
	currentCount := ms.counts[element]
	if currentCount == count {
		return ms
	}
	
	newCounts := ms.copyMap()
	sizeDiff := count - currentCount
	
	if count == 0 {
		delete(newCounts, element)
	} else {
		newCounts[element] = count
	}
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   ms.size + sizeDiff,
	}
}

// Contains checks if the multiset contains the specified element
func (ms *ImmutableMultiset[E]) Contains(element E) bool {
	return ms.counts[element] > 0
}

// IsEmpty returns true if the multiset contains no elements
func (ms *ImmutableMultiset[E]) IsEmpty() bool {
	return ms.size == 0
}

// Size returns the number of distinct elements in the multiset
func (ms *ImmutableMultiset[E]) Size() int {
	return len(ms.counts)
}

// TotalSize returns the total number of elements (including duplicates)
func (ms *ImmutableMultiset[E]) TotalSize() int {
	return ms.size
}

// DistinctElements returns the number of distinct elements
func (ms *ImmutableMultiset[E]) DistinctElements() int {
	return len(ms.counts)
}

// Clear returns an empty ImmutableMultiset
func (ms *ImmutableMultiset[E]) Clear() {
	// This violates the interface but is needed for immutability
	// Use WithClear() instead for proper immutable behavior
}

// WithClear returns a new empty ImmutableMultiset
func (ms *ImmutableMultiset[E]) WithClear() *ImmutableMultiset[E] {
	return NewImmutableMultiset[E]()
}

// ElementSet returns a slice of distinct elements
func (ms *ImmutableMultiset[E]) ElementSet() []E {
	elements := make([]E, 0, len(ms.counts))
	for element := range ms.counts {
		elements = append(elements, element)
	}
	return elements
}

// EntrySet returns a slice of entries (element-count pairs)
func (ms *ImmutableMultiset[E]) EntrySet() []Entry[E] {
	entries := make([]Entry[E], 0, len(ms.counts))
	for element, count := range ms.counts {
		entries = append(entries, Entry[E]{Element: element, Count: count})
	}
	return entries
}

// ToSlice returns a slice containing all elements (including duplicates)
func (ms *ImmutableMultiset[E]) ToSlice() []E {
	result := make([]E, 0, ms.size)
	for element, count := range ms.counts {
		for i := 0; i < count; i++ {
			result = append(result, element)
		}
	}
	return result
}

// Iterator returns an iterator over the multiset elements
func (ms *ImmutableMultiset[E]) Iterator() common.Iterator[E] {
	return &immutableMultisetIterator[E]{
		multiset: ms,
		entries:  ms.EntrySet(),
		index:    0,
		current:  0,
	}
}

// ForEach executes the given function for each element in the multiset
func (ms *ImmutableMultiset[E]) ForEach(fn func(E)) {
	for element, count := range ms.counts {
		for i := 0; i < count; i++ {
			fn(element)
		}
	}
}

// Union returns a new multiset containing the union of this and another multiset
func (ms *ImmutableMultiset[E]) Union(other Multiset[E]) Multiset[E] {
	newCounts := ms.copyMap()
	newSize := ms.size
	
	// Add elements from other multiset, taking maximum count
	for _, entry := range other.EntrySet() {
		currentCount := newCounts[entry.Element]
		if entry.Count > currentCount {
			newSize += entry.Count - currentCount
			newCounts[entry.Element] = entry.Count
		}
	}
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   newSize,
	}
}

// Intersection returns a new multiset containing the intersection
func (ms *ImmutableMultiset[E]) Intersection(other Multiset[E]) Multiset[E] {
	newCounts := make(map[E]int)
	newSize := 0
	
	for element, count := range ms.counts {
		otherCount := other.Count(element)
		if otherCount > 0 {
			minCount := count
			if otherCount < minCount {
				minCount = otherCount
			}
			newCounts[element] = minCount
			newSize += minCount
		}
	}
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   newSize,
	}
}

// Difference returns a new multiset containing elements in this but not in other
func (ms *ImmutableMultiset[E]) Difference(other Multiset[E]) Multiset[E] {
	newCounts := make(map[E]int)
	newSize := 0
	
	for element, count := range ms.counts {
		otherCount := other.Count(element)
		if count > otherCount {
			diffCount := count - otherCount
			newCounts[element] = diffCount
			newSize += diffCount
		}
	}
	
	return &ImmutableMultiset[E]{
		counts: newCounts,
		size:   newSize,
	}
}

// IsSubsetOf checks if this multiset is a subset of another
func (ms *ImmutableMultiset[E]) IsSubsetOf(other Multiset[E]) bool {
	for element, count := range ms.counts {
		if other.Count(element) < count {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this multiset is a superset of another
func (ms *ImmutableMultiset[E]) IsSupersetOf(other Multiset[E]) bool {
	return other.IsSubsetOf(ms)
}

// String returns a string representation of the multiset
func (ms *ImmutableMultiset[E]) String() string {
	if ms.size == 0 {
		return "ImmutableMultiset[]"
	}
	
	var builder strings.Builder
	builder.WriteString("ImmutableMultiset[")
	
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

// immutableMultisetIterator implements Iterator for ImmutableMultiset
type immutableMultisetIterator[E comparable] struct {
	multiset *ImmutableMultiset[E]
	entries  []Entry[E]
	index    int
	current  int
}

func (it *immutableMultisetIterator[E]) HasNext() bool {
	return it.index < len(it.entries) && (it.current < it.entries[it.index].Count || it.index+1 < len(it.entries))
}

func (it *immutableMultisetIterator[E]) Next() (E, bool) {
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

func (it *immutableMultisetIterator[E]) Reset() {
	it.entries = it.multiset.EntrySet()
	it.index = 0
	it.current = 0
}

func (it *immutableMultisetIterator[E]) Remove() bool {
	// Remove operation is not supported for immutable collections
	return false
}