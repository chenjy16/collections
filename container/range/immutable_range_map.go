package ranges

import (
	"strings"
	"log"
	"github.com/chenjianyu/collections/container/common"
)

// ImmutableRangeMap is an immutable implementation of RangeMap
// All modification operations return new instances
type ImmutableRangeMap[K comparable, V any] struct {
	entries    []Entry[K, V]
	comparator Comparator[K]
}

// NewImmutableRangeMap creates a new empty ImmutableRangeMap
func NewImmutableRangeMap[K comparable, V any]() RangeMap[K, V] {
	return &ImmutableRangeMap[K, V]{
		entries:    make([]Entry[K, V], 0),
		comparator: DefaultComparator[K],
	}
}

// NewImmutableRangeMapWithComparator creates a new ImmutableRangeMap with custom comparator
func NewImmutableRangeMapWithComparator[K comparable, V any](cmp Comparator[K]) RangeMap[K, V] {
	return &ImmutableRangeMap[K, V]{
		entries:    make([]Entry[K, V], 0),
		comparator: cmp,
	}
}

// NewImmutableRangeMapFromEntries creates a new ImmutableRangeMap from existing entries
func NewImmutableRangeMapFromEntries[K comparable, V any](entries []Entry[K, V]) RangeMap[K, V] {
	// Create a mutable map to handle overlapping ranges
	mutableMap := NewTreeRangeMap[K, V]()
	for _, entry := range entries {
		mutableMap.Put(entry.Range, entry.Value)
	}
	
	return &ImmutableRangeMap[K, V]{
		entries:    convertToEntries(mutableMap.AsMapOfRanges()),
		comparator: DefaultComparator[K],
	}
}

// Helper function to convert map to entries slice
func convertToEntries[K comparable, V any](rangeMap map[Range[K]]V) []Entry[K, V] {
	entries := make([]Entry[K, V], 0, len(rangeMap))
	for r, v := range rangeMap {
		entries = append(entries, Entry[K, V]{Range: r, Value: v})
	}
	return entries
}

// Size returns the number of range-value mappings in this map
func (irm *ImmutableRangeMap[K, V]) Size() int {
	return len(irm.entries)
}

// IsEmpty returns true if this range map is empty
func (irm *ImmutableRangeMap[K, V]) IsEmpty() bool {
	return len(irm.entries) == 0
}

// Clear logs an error and returns as ImmutableRangeMap is immutable
func (irm *ImmutableRangeMap[K, V]) Clear() {
	err := common.ImmutableOperationError("Clear", "WithClear()")
	log.Printf("Warning: %v", err)
}

// WithClear returns a new empty ImmutableRangeMap
func (irm *ImmutableRangeMap[K, V]) WithClear() RangeMap[K, V] {
	return NewImmutableRangeMapWithComparator[K, V](irm.comparator)
}

// String returns the string representation of this range map
func (irm *ImmutableRangeMap[K, V]) String() string {
	if len(irm.entries) == 0 {
		return "{}"
	}
	
	var parts []string
	for _, entry := range irm.entries {
		parts = append(parts, entry.String())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// Get returns the value associated with the specified key, or nil if no such value exists
func (irm *ImmutableRangeMap[K, V]) Get(key K) (V, bool) {
	for _, entry := range irm.entries {
		if entry.Range.Contains(key) {
			return entry.Value, true
		}
	}
	
	var zero V
	return zero, false
}

// GetEntry returns the range-value entry that contains the specified key
func (irm *ImmutableRangeMap[K, V]) GetEntry(key K) (Range[K], V, bool) {
	for _, entry := range irm.entries {
		if entry.Range.Contains(key) {
			return entry.Range, entry.Value, true
		}
	}
	
	var zeroV V
	var zeroR Range[K]
	return zeroR, zeroV, false
}

// Put logs an error and returns as ImmutableRangeMap is immutable
func (irm *ImmutableRangeMap[K, V]) Put(rangeKey Range[K], value V) {
	err := common.ImmutableOperationError("Put", "WithPut()")
	log.Printf("Warning: %v", err)
}

// WithPut returns a new ImmutableRangeMap with the range-value mapping added
func (irm *ImmutableRangeMap[K, V]) WithPut(rangeKey Range[K], value V) RangeMap[K, V] {
	if rangeKey == nil || rangeKey.IsEmpty() {
		return irm
	}
	
	// Create a mutable copy and add the mapping
	mutableMap := NewTreeRangeMapWithComparator[K, V](irm.comparator)
	for _, entry := range irm.entries {
		mutableMap.Put(entry.Range, entry.Value)
	}
	mutableMap.Put(rangeKey, value)
	
	return &ImmutableRangeMap[K, V]{
		entries:    convertToEntries(mutableMap.AsMapOfRanges()),
		comparator: irm.comparator,
	}
}

// PutRange logs an error and returns as ImmutableRangeMap is immutable
func (irm *ImmutableRangeMap[K, V]) PutRange(lower K, lowerType BoundType, upper K, upperType BoundType, value V) {
	err := common.ImmutableOperationError("PutRange", "WithPutRange()")
	log.Printf("Warning: %v", err)
}

// WithPutRange returns a new ImmutableRangeMap with the range-value mapping added
func (irm *ImmutableRangeMap[K, V]) WithPutRange(lower K, lowerType BoundType, upper K, upperType BoundType, value V) RangeMap[K, V] {
	rangeKey := NewRangeWithComparator(lower, lowerType, upper, upperType, irm.comparator)
	return irm.WithPut(rangeKey, value)
}

// Remove logs an error and returns as ImmutableRangeMap is immutable
func (irm *ImmutableRangeMap[K, V]) Remove(rangeToRemove Range[K]) {
	err := common.ImmutableOperationError("Remove", "WithRemove()")
	log.Printf("Warning: %v", err)
}

// WithRemove returns a new ImmutableRangeMap with the range removed
func (irm *ImmutableRangeMap[K, V]) WithRemove(rangeToRemove Range[K]) RangeMap[K, V] {
	if rangeToRemove == nil || rangeToRemove.IsEmpty() {
		return irm
	}
	
	// Create a mutable copy and remove the range
	mutableMap := NewTreeRangeMapWithComparator[K, V](irm.comparator)
	for _, entry := range irm.entries {
		mutableMap.Put(entry.Range, entry.Value)
	}
	mutableMap.Remove(rangeToRemove)
	
	return &ImmutableRangeMap[K, V]{
		entries:    convertToEntries(mutableMap.AsMapOfRanges()),
		comparator: irm.comparator,
	}
}

// RemoveRange logs an error and returns as ImmutableRangeMap is immutable
func (irm *ImmutableRangeMap[K, V]) RemoveRange(lower K, lowerType BoundType, upper K, upperType BoundType) {
	err := common.ImmutableOperationError("RemoveRange", "WithRemoveRange()")
	log.Printf("Warning: %v", err)
}

// WithRemoveRange returns a new ImmutableRangeMap with the range removed
func (irm *ImmutableRangeMap[K, V]) WithRemoveRange(lower K, lowerType BoundType, upper K, upperType BoundType) RangeMap[K, V] {
	rangeToRemove := NewRangeWithComparator(lower, lowerType, upper, upperType, irm.comparator)
	return irm.WithRemove(rangeToRemove)
}

// AsMapOfRanges returns a view of this range map as a map from ranges to values
func (irm *ImmutableRangeMap[K, V]) AsMapOfRanges() map[Range[K]]V {
	result := make(map[Range[K]]V)
	for _, entry := range irm.entries {
		result[entry.Range] = entry.Value
	}
	return result
}

// AsDescendingMapOfRanges returns a view of this range map as a map from ranges to values in descending order
func (irm *ImmutableRangeMap[K, V]) AsDescendingMapOfRanges() map[Range[K]]V {
	result := make(map[Range[K]]V)
	for i := len(irm.entries) - 1; i >= 0; i-- {
		entry := irm.entries[i]
		result[entry.Range] = entry.Value
	}
	return result
}

// Span returns the minimal range that contains all ranges in this map
func (irm *ImmutableRangeMap[K, V]) Span() (Range[K], bool) {
	if len(irm.entries) == 0 {
		var zeroR Range[K]
		return zeroR, false
	}
	
	if len(irm.entries) == 1 {
		return irm.entries[0].Range, true
	}
	
	// Calculate the span of all ranges
	span := irm.entries[0].Range
	for i := 1; i < len(irm.entries); i++ {
		span = span.Span(irm.entries[i].Range)
	}
	
	return span, true
}

// SubRangeMap returns a view of the portion of this range map that intersects with the given range
func (irm *ImmutableRangeMap[K, V]) SubRangeMap(subRange Range[K]) RangeMap[K, V] {
	if subRange == nil || subRange.IsEmpty() {
		return NewImmutableRangeMapWithComparator[K, V](irm.comparator)
	}
	
	var resultEntries []Entry[K, V]
	
	for _, entry := range irm.entries {
		if entry.Range.IsConnected(subRange) {
			intersection := entry.Range.Intersection(subRange)
			if !intersection.IsEmpty() {
				resultEntries = append(resultEntries, Entry[K, V]{
					Range: intersection,
					Value: entry.Value,
				})
			}
		}
	}
	
	return &ImmutableRangeMap[K, V]{
		entries:    resultEntries,
		comparator: irm.comparator,
	}
}