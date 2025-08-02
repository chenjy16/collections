package ranges

import (
	"strings"
	"sync"
)

// TreeRangeMap is a mutable implementation of RangeMap backed by a tree structure
type TreeRangeMap[K comparable, V any] struct {
	entries    []Entry[K, V]
	comparator Comparator[K]
	mutex      sync.RWMutex
}

// NewTreeRangeMap creates a new empty TreeRangeMap
func NewTreeRangeMap[K comparable, V any]() RangeMap[K, V] {
	return &TreeRangeMap[K, V]{
		entries:    make([]Entry[K, V], 0),
		comparator: DefaultComparator[K],
	}
}

// NewTreeRangeMapWithComparator creates a new TreeRangeMap with custom comparator
func NewTreeRangeMapWithComparator[K comparable, V any](cmp Comparator[K]) RangeMap[K, V] {
	return &TreeRangeMap[K, V]{
		entries:    make([]Entry[K, V], 0),
		comparator: cmp,
	}
}

// Size returns the number of range-value mappings in this map
func (trm *TreeRangeMap[K, V]) Size() int {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	return len(trm.entries)
}

// IsEmpty returns true if this range map is empty
func (trm *TreeRangeMap[K, V]) IsEmpty() bool {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	return len(trm.entries) == 0
}

// Clear removes all mappings from this range map
func (trm *TreeRangeMap[K, V]) Clear() {
	trm.mutex.Lock()
	defer trm.mutex.Unlock()
	trm.entries = make([]Entry[K, V], 0)
}

// String returns the string representation of this range map
func (trm *TreeRangeMap[K, V]) String() string {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	if len(trm.entries) == 0 {
		return "{}"
	}
	
	var parts []string
	for _, entry := range trm.entries {
		parts = append(parts, entry.String())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// Put associates the specified value with the specified range
func (trm *TreeRangeMap[K, V]) Put(rangeKey Range[K], value V) {
	if rangeKey == nil || rangeKey.IsEmpty() {
		return
	}
	
	trm.mutex.Lock()
	defer trm.mutex.Unlock()
	
	// Remove any overlapping ranges first
	trm.removeOverlapping(rangeKey)
	
	// Add the new entry
	entry := Entry[K, V]{
		Range: rangeKey,
		Value: value,
	}
	
	trm.entries = append(trm.entries, entry)
	trm.sortEntries()
}

// PutRange associates the specified value with the specified range
func (trm *TreeRangeMap[K, V]) PutRange(lower K, lowerType BoundType, upper K, upperType BoundType, value V) {
	rangeKey := NewRangeWithComparator(lower, lowerType, upper, upperType, trm.comparator)
	trm.Put(rangeKey, value)
}

// Get returns the value associated with the specified key, or nil if no such value exists
func (trm *TreeRangeMap[K, V]) Get(key K) (V, bool) {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	for _, entry := range trm.entries {
		if entry.Range.Contains(key) {
			return entry.Value, true
		}
	}
	
	var zero V
	return zero, false
}

// GetEntry returns the range-value entry that contains the specified key
func (trm *TreeRangeMap[K, V]) GetEntry(key K) (Range[K], V, bool) {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	for _, entry := range trm.entries {
		if entry.Range.Contains(key) {
			return entry.Range, entry.Value, true
		}
	}
	
	var zeroV V
	var zeroR Range[K]
	return zeroR, zeroV, false
}

// Remove removes all mappings from the specified range
func (trm *TreeRangeMap[K, V]) Remove(rangeToRemove Range[K]) {
	if rangeToRemove == nil || rangeToRemove.IsEmpty() {
		return
	}
	
	trm.mutex.Lock()
	defer trm.mutex.Unlock()
	
	var newEntries []Entry[K, V]
	
	for _, entry := range trm.entries {
		if !entry.Range.IsConnected(rangeToRemove) {
			// No overlap, keep the entry
			newEntries = append(newEntries, entry)
		} else if rangeToRemove.ContainsRange(entry.Range) {
			// Range to remove completely contains the entry, remove it
			continue
		} else {
			// For partial overlaps, we'll keep the entry for now
			// TODO: Implement proper range splitting when Range interface is complete
			newEntries = append(newEntries, entry)
		}
	}
	
	trm.entries = newEntries
	trm.sortEntries()
}

// RemoveRange removes all mappings from the specified range
func (trm *TreeRangeMap[K, V]) RemoveRange(lower K, lowerType BoundType, upper K, upperType BoundType) {
	rangeToRemove := NewRangeWithComparator(lower, lowerType, upper, upperType, trm.comparator)
	trm.Remove(rangeToRemove)
}

// Span returns the minimal range that contains all ranges in this map
func (trm *TreeRangeMap[K, V]) Span() (Range[K], bool) {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	if len(trm.entries) == 0 {
		var zeroR Range[K]
		return zeroR, false
	}
	
	if len(trm.entries) == 1 {
		return trm.entries[0].Range, true
	}
	
	// For now, return the first range as a placeholder
	// TODO: Implement proper span calculation when Range interface is complete
	return trm.entries[0].Range, true
}

// AsMapOfRanges returns a view of this range map as a map from ranges to values
func (trm *TreeRangeMap[K, V]) AsMapOfRanges() map[Range[K]]V {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	result := make(map[Range[K]]V)
	for _, entry := range trm.entries {
		result[entry.Range] = entry.Value
	}
	return result
}

// AsDescendingMapOfRanges returns a view of this range map as a map from ranges to values in descending order
func (trm *TreeRangeMap[K, V]) AsDescendingMapOfRanges() map[Range[K]]V {
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	result := make(map[Range[K]]V)
	for i := len(trm.entries) - 1; i >= 0; i-- {
		entry := trm.entries[i]
		result[entry.Range] = entry.Value
	}
	return result
}

// SubRangeMap returns a view of the portion of this range map that intersects with the given range
func (trm *TreeRangeMap[K, V]) SubRangeMap(subRange Range[K]) RangeMap[K, V] {
	if subRange == nil || subRange.IsEmpty() {
		return NewTreeRangeMapWithComparator[K, V](trm.comparator)
	}
	
	trm.mutex.RLock()
	defer trm.mutex.RUnlock()
	
	result := NewTreeRangeMapWithComparator[K, V](trm.comparator)
	
	for _, entry := range trm.entries {
		if entry.Range.IsConnected(subRange) {
			intersection := entry.Range.Intersection(subRange)
			if !intersection.IsEmpty() {
				result.Put(intersection, entry.Value)
			}
		}
	}
	
	return result
}

// Helper method to remove overlapping ranges
func (trm *TreeRangeMap[K, V]) removeOverlapping(rangeKey Range[K]) {
	var newEntries []Entry[K, V]
	
	for _, entry := range trm.entries {
		if !entry.Range.IsConnected(rangeKey) {
			newEntries = append(newEntries, entry)
		}
	}
	
	trm.entries = newEntries
}

// Helper method to sort entries by range
func (trm *TreeRangeMap[K, V]) sortEntries() {
	// TODO: Implement proper sorting when Range interface has LowerEndpoint method
	// For now, entries are kept in insertion order
}