package multiset

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chenjianyu/collections/container/common"
)

// ConcurrentHashMultiset is a thread-safe multiset implementation
// It uses fine-grained locking with segment-based approach for better concurrency
type ConcurrentHashMultiset[E comparable] struct {
	segments []*segment[E]
	segMask  uint32
	size     int64
}

type segment[E comparable] struct {
	counts map[E]int
	mu     sync.RWMutex
}

const (
	defaultSegmentCount = 16
	defaultSegmentSize  = 16
)

// NewConcurrentHashMultiset creates a new empty ConcurrentHashMultiset
func NewConcurrentHashMultiset[E comparable]() *ConcurrentHashMultiset[E] {
	return NewConcurrentHashMultisetWithSegments[E](defaultSegmentCount)
}

// NewConcurrentHashMultisetWithSegments creates a new ConcurrentHashMultiset with specified segment count
func NewConcurrentHashMultisetWithSegments[E comparable](segmentCount int) *ConcurrentHashMultiset[E] {
	// Ensure segment count is a power of 2
	if segmentCount <= 0 {
		segmentCount = defaultSegmentCount
	}
	
	// Round up to next power of 2
	actualSegmentCount := 1
	for actualSegmentCount < segmentCount {
		actualSegmentCount <<= 1
	}
	
	segments := make([]*segment[E], actualSegmentCount)
	for i := range segments {
		segments[i] = &segment[E]{
			counts: make(map[E]int, defaultSegmentSize),
		}
	}
	
	return &ConcurrentHashMultiset[E]{
		segments: segments,
		segMask:  uint32(actualSegmentCount - 1),
	}
}

// NewConcurrentHashMultisetFromSlice creates a new ConcurrentHashMultiset from a slice
func NewConcurrentHashMultisetFromSlice[E comparable](elements []E) *ConcurrentHashMultiset[E] {
	ms := NewConcurrentHashMultiset[E]()
	for _, element := range elements {
		ms.Add(element)
	}
	return ms
}

// hash computes a hash value for the given element
func (ms *ConcurrentHashMultiset[E]) hash(element E) uint32 {
	// Simple hash function - in production, use a better hash function
	str := fmt.Sprintf("%v", element)
	var hash uint32 = 5381
	for _, c := range str {
		hash = ((hash << 5) + hash) + uint32(c)
	}
	return hash
}

// getSegment returns the segment for the given element
func (ms *ConcurrentHashMultiset[E]) getSegment(element E) *segment[E] {
	hash := ms.hash(element)
	return ms.segments[hash&ms.segMask]
}

// Add adds one occurrence of the specified element
func (ms *ConcurrentHashMultiset[E]) Add(element E) int {
	seg := ms.getSegment(element)
	seg.mu.Lock()
	defer seg.mu.Unlock()
	
	prevCount := seg.counts[element]
	seg.counts[element] = prevCount + 1
	atomic.AddInt64(&ms.size, 1)
	return prevCount
}

// AddCount adds the specified number of occurrences of the element
func (ms *ConcurrentHashMultiset[E]) AddCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	if count == 0 {
		return ms.Count(element), nil
	}
	
	seg := ms.getSegment(element)
	seg.mu.Lock()
	defer seg.mu.Unlock()
	
	prevCount := seg.counts[element]
	seg.counts[element] = prevCount + count
	atomic.AddInt64(&ms.size, int64(count))
	return prevCount, nil
}

// Remove removes one occurrence of the specified element
func (ms *ConcurrentHashMultiset[E]) Remove(element E) int {
	seg := ms.getSegment(element)
	seg.mu.Lock()
	defer seg.mu.Unlock()
	
	prevCount := seg.counts[element]
	if prevCount > 0 {
		if prevCount == 1 {
			delete(seg.counts, element)
		} else {
			seg.counts[element] = prevCount - 1
		}
		atomic.AddInt64(&ms.size, -1)
	}
	return prevCount
}

// RemoveCount removes the specified number of occurrences of the element
func (ms *ConcurrentHashMultiset[E]) RemoveCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	if count == 0 {
		return ms.Count(element), nil
	}
	
	seg := ms.getSegment(element)
	seg.mu.Lock()
	defer seg.mu.Unlock()
	
	prevCount := seg.counts[element]
	if prevCount > 0 {
		removeCount := count
		if removeCount > prevCount {
			removeCount = prevCount
		}
		
		newCount := prevCount - removeCount
		if newCount == 0 {
			delete(seg.counts, element)
		} else {
			seg.counts[element] = newCount
		}
		atomic.AddInt64(&ms.size, -int64(removeCount))
	}
	return prevCount, nil
}

// RemoveAll removes all occurrences of the specified element
func (ms *ConcurrentHashMultiset[E]) RemoveAll(element E) int {
	seg := ms.getSegment(element)
	seg.mu.Lock()
	defer seg.mu.Unlock()
	
	prevCount := seg.counts[element]
	if prevCount > 0 {
		delete(seg.counts, element)
		atomic.AddInt64(&ms.size, -int64(prevCount))
	}
	return prevCount
}

// Count returns the number of occurrences of the specified element
func (ms *ConcurrentHashMultiset[E]) Count(element E) int {
	seg := ms.getSegment(element)
	seg.mu.RLock()
	defer seg.mu.RUnlock()
	return seg.counts[element]
}

// SetCount sets the count of the specified element to the given value
func (ms *ConcurrentHashMultiset[E]) SetCount(element E, count int) (int, error) {
	if count < 0 {
		return 0, common.NegativeCountError(count)
	}
	
	seg := ms.getSegment(element)
	seg.mu.Lock()
	defer seg.mu.Unlock()
	
	prevCount := seg.counts[element]
	
	if count == 0 {
		if prevCount > 0 {
			delete(seg.counts, element)
			atomic.AddInt64(&ms.size, -int64(prevCount))
		}
	} else {
		seg.counts[element] = count
		atomic.AddInt64(&ms.size, int64(count-prevCount))
	}
	
	return prevCount, nil
}

// Contains checks if the multiset contains the specified element
func (ms *ConcurrentHashMultiset[E]) Contains(element E) bool {
	return ms.Count(element) > 0
}

// IsEmpty returns true if the multiset contains no elements
func (ms *ConcurrentHashMultiset[E]) IsEmpty() bool {
	return atomic.LoadInt64(&ms.size) == 0
}

// Size returns the number of distinct elements in the multiset
func (ms *ConcurrentHashMultiset[E]) Size() int {
	count := 0
	for _, seg := range ms.segments {
		seg.mu.RLock()
		count += len(seg.counts)
		seg.mu.RUnlock()
	}
	return count
}

// TotalSize returns the total number of elements (including duplicates)
func (ms *ConcurrentHashMultiset[E]) TotalSize() int {
	return int(atomic.LoadInt64(&ms.size))
}

// DistinctElements returns the number of distinct elements
func (ms *ConcurrentHashMultiset[E]) DistinctElements() int {
	return ms.Size()
}

// Clear removes all elements from the multiset
func (ms *ConcurrentHashMultiset[E]) Clear() {
	for _, seg := range ms.segments {
		seg.mu.Lock()
		seg.counts = make(map[E]int, defaultSegmentSize)
		seg.mu.Unlock()
	}
	atomic.StoreInt64(&ms.size, 0)
}

// ElementSet returns a slice of distinct elements
func (ms *ConcurrentHashMultiset[E]) ElementSet() []E {
	var elements []E
	
	for _, seg := range ms.segments {
		seg.mu.RLock()
		for element := range seg.counts {
			elements = append(elements, element)
		}
		seg.mu.RUnlock()
	}
	
	return elements
}

// EntrySet returns a slice of entries (element-count pairs)
func (ms *ConcurrentHashMultiset[E]) EntrySet() []Entry[E] {
	var entries []Entry[E]
	
	for _, seg := range ms.segments {
		seg.mu.RLock()
		for element, count := range seg.counts {
			entries = append(entries, Entry[E]{Element: element, Count: count})
		}
		seg.mu.RUnlock()
	}
	
	return entries
}

// ToSlice returns a slice containing all elements (including duplicates)
func (ms *ConcurrentHashMultiset[E]) ToSlice() []E {
	var result []E
	
	for _, seg := range ms.segments {
		seg.mu.RLock()
		for element, count := range seg.counts {
			for i := 0; i < count; i++ {
				result = append(result, element)
			}
		}
		seg.mu.RUnlock()
	}
	
	return result
}

// Iterator returns an iterator over the multiset elements
func (ms *ConcurrentHashMultiset[E]) Iterator() common.Iterator[E] {
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

// ForEach executes the given function for each element in the multiset
func (ms *ConcurrentHashMultiset[E]) ForEach(fn func(E)) {
	for _, entry := range ms.EntrySet() {
		for i := 0; i < entry.Count; i++ {
			fn(entry.Element)
		}
	}
}

// Union returns a new multiset containing the union of this and another multiset
func (ms *ConcurrentHashMultiset[E]) Union(other Multiset[E]) Multiset[E] {
	result := NewConcurrentHashMultiset[E]()
	
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
func (ms *ConcurrentHashMultiset[E]) Intersection(other Multiset[E]) Multiset[E] {
	result := NewConcurrentHashMultiset[E]()
	
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
func (ms *ConcurrentHashMultiset[E]) Difference(other Multiset[E]) Multiset[E] {
	result := NewConcurrentHashMultiset[E]()
	
	for _, entry := range ms.EntrySet() {
		otherCount := other.Count(entry.Element)
		if entry.Count > otherCount {
			result.AddCount(entry.Element, entry.Count-otherCount)
		}
	}
	
	return result
}

// IsSubsetOf checks if this multiset is a subset of another
func (ms *ConcurrentHashMultiset[E]) IsSubsetOf(other Multiset[E]) bool {
	for _, entry := range ms.EntrySet() {
		if other.Count(entry.Element) < entry.Count {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this multiset is a superset of another
func (ms *ConcurrentHashMultiset[E]) IsSupersetOf(other Multiset[E]) bool {
	return other.IsSubsetOf(ms)
}

// String returns a string representation of the multiset
func (ms *ConcurrentHashMultiset[E]) String() string {
	if ms.TotalSize() == 0 {
		return "ConcurrentHashMultiset[]"
	}
	
	var builder strings.Builder
	builder.WriteString("ConcurrentHashMultiset[")
	
	entries := ms.EntrySet()
	for i, entry := range entries {
		if i > 0 {
			builder.WriteString(", ")
		}
		if entry.Count == 1 {
			builder.WriteString(fmt.Sprintf("%v", entry.Element))
		} else {
			builder.WriteString(fmt.Sprintf("%v x %d", entry.Element, entry.Count))
		}
	}
	
	builder.WriteString("]")
	return builder.String()
}

// concurrentHashMultisetIterator implements Iterator for ConcurrentHashMultiset
type concurrentHashMultisetIterator[E comparable] struct {
	multiset *ConcurrentHashMultiset[E]
	entries  []Entry[E]
	index    int
	current  int
}

func (it *concurrentHashMultisetIterator[E]) HasNext() bool {
	return it.index < len(it.entries) && (it.current < it.entries[it.index].Count || it.index+1 < len(it.entries))
}

func (it *concurrentHashMultisetIterator[E]) Next() (E, bool) {
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

func (it *concurrentHashMultisetIterator[E]) Reset() {
	it.entries = it.multiset.EntrySet()
	it.index = 0
	it.current = 0
}

func (it *concurrentHashMultisetIterator[E]) Remove() bool {
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