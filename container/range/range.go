// Package range provides range-based collection implementations similar to Guava's Range collections
package ranges

import (
    "fmt"
    "log"
    "github.com/chenjianyu/collections/container/common"
)

// BoundType represents the type of bound (open or closed)
type BoundType int

const (
	Open   BoundType = iota // (value) - exclusive
	Closed                  // [value] - inclusive
)

// Range represents a contiguous span of values of some Comparable type
type Range[T comparable] interface {
	// LowerBound returns the lower bound of this range
	LowerBound() (T, BoundType, bool)
	
	// UpperBound returns the upper bound of this range
	UpperBound() (T, BoundType, bool)
	
	// Contains returns true if the value is within this range
	Contains(value T) bool
	
	// ContainsRange returns true if the other range is entirely contained within this range
	ContainsRange(other Range[T]) bool
	
	// IsConnected returns true if there exists a (possibly empty) range which is
	// enclosed by both this range and other
	IsConnected(other Range[T]) bool
	
	// Intersection returns the maximal range enclosed by both this range and other
	// Returns nil if no such range exists
	Intersection(other Range[T]) Range[T]
	
	// Span returns the minimal range that encloses both this range and other
	Span(other Range[T]) Range[T]
	
	// IsEmpty returns true if this range is empty
	IsEmpty() bool
	
	// String returns the string representation of this range
	String() string
}

// RangeSet represents a set of non-empty, disconnected ranges
type RangeSet[T comparable] interface {
	// Size returns the number of ranges in this set
	Size() int
	
	// IsEmpty returns true if this range set is empty
	IsEmpty() bool
	
	// Clear removes all ranges from this set
	Clear()
	
	// String returns the string representation of this range set
	String() string
	
	// Add adds a range to this range set
	Add(rangeToAdd Range[T])
	
	// AddRange adds a range defined by bounds to this range set
	AddRange(lower T, lowerType BoundType, upper T, upperType BoundType)
	
	// Remove removes a range from this range set
	Remove(rangeToRemove Range[T])
	
	// RemoveRange removes a range defined by bounds from this range set
	RemoveRange(lower T, lowerType BoundType, upper T, upperType BoundType)
	
	// ContainsValue returns true if the value is contained in any range in this set
	ContainsValue(value T) bool
	
	// ContainsRange returns true if the range is entirely contained in this set
	ContainsRange(rangeToCheck Range[T]) bool
	
	// Encloses returns true if this range set encloses the other range set
	Encloses(other RangeSet[T]) bool
	
	// AsRanges returns a view of the disconnected ranges that make up this range set
	AsRanges() []Range[T]
	
	// Complement returns the complement of this range set
	Complement() RangeSet[T]
	
	// Union returns the union of this range set with another
	Union(other RangeSet[T]) RangeSet[T]
	
	// Intersection returns the intersection of this range set with another
	Intersection(other RangeSet[T]) RangeSet[T]
	
	// Difference returns the difference of this range set with another
	Difference(other RangeSet[T]) RangeSet[T]
}

// RangeMap represents a mapping from disjoint nonempty ranges to non-null values
type RangeMap[K comparable, V any] interface {
	// Size returns the number of range-value mappings in this map
	Size() int
	
	// IsEmpty returns true if this range map is empty
	IsEmpty() bool
	
	// Clear removes all mappings from this range map
	Clear()
	
	// String returns the string representation of this range map
	String() string
	
	// Get returns the value associated with the specified key, or nil if no such value exists
	Get(key K) (V, bool)
	
	// GetEntry returns the range-value entry containing the specified key, or nil if no such entry exists
	GetEntry(key K) (Range[K], V, bool)
	
	// Put associates the specified value with the specified range
	Put(rangeKey Range[K], value V)
	
	// PutRange associates the specified value with the specified range defined by bounds
	PutRange(lower K, lowerType BoundType, upper K, upperType BoundType, value V)
	
	// Remove removes all associations from this range map in the specified range
	Remove(rangeToRemove Range[K])
	
	// RemoveRange removes all associations from this range map in the specified range defined by bounds
	RemoveRange(lower K, lowerType BoundType, upper K, upperType BoundType)
	
	// AsMapOfRanges returns a view of this range map as a Map<Range<K>, V>
	AsMapOfRanges() map[Range[K]]V
	
	// AsDescendingMapOfRanges returns a view of this range map as a descending Map<Range<K>, V>
	AsDescendingMapOfRanges() map[Range[K]]V
	
	// Span returns the minimal range enclosing the ranges in this RangeMap
	Span() (Range[K], bool)
	
    // SubRangeMap returns a view of the portion of this range map that intersects with range
    SubRangeMap(range_ Range[K]) RangeMap[K, V]

    // Entries returns all range-value pairs as common entries with Range[K] as key
    Entries() []common.Entry[Range[K], V]
}

// Entry represents a range-value pair in a RangeMap
type Entry[K comparable, V any] struct {
	Range Range[K]
	Value V
}

// String returns the string representation of this entry
func (e Entry[K, V]) String() string {
	return fmt.Sprintf("%s=%v", e.Range.String(), e.Value)
}

// Comparator function type for comparing values
type Comparator[T any] func(a, b T) int

// DefaultComparator provides a default comparison function for comparable types
func DefaultComparator[T comparable](a, b T) int {
    // Prefer Comparable.CompareTo if available
    if comparableA, ok := any(a).(interface{ CompareTo(interface{}) int }); ok {
        return comparableA.CompareTo(any(b))
    }
    switch any(a).(type) {
        case int:
		aInt, bInt := any(a).(int), any(b).(int)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case int8:
		aInt, bInt := any(a).(int8), any(b).(int8)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case int16:
		aInt, bInt := any(a).(int16), any(b).(int16)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case int32:
		aInt, bInt := any(a).(int32), any(b).(int32)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case int64:
		aInt, bInt := any(a).(int64), any(b).(int64)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case uint:
		aInt, bInt := any(a).(uint), any(b).(uint)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case uint8:
		aInt, bInt := any(a).(uint8), any(b).(uint8)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case uint16:
		aInt, bInt := any(a).(uint16), any(b).(uint16)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case uint32:
		aInt, bInt := any(a).(uint32), any(b).(uint32)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case uint64:
		aInt, bInt := any(a).(uint64), any(b).(uint64)
		if aInt < bInt {
			return -1
		} else if aInt > bInt {
			return 1
		}
		return 0
	case float32:
		aFloat, bFloat := any(a).(float32), any(b).(float32)
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0
	case float64:
		aFloat, bFloat := any(a).(float64), any(b).(float64)
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0
	case string:
		aStr, bStr := any(a).(string), any(b).(string)
		if aStr < bStr {
			return -1
		} else if aStr > bStr {
			return 1
		}
		return 0
	default:
		// For other types, we can't provide a meaningful comparison
		// This should be handled by providing a custom comparator
		log.Printf("Warning: no default comparator available for type %T", a)
		return 0 // Return 0 to indicate equality as fallback
	}
}

// ComparatorFromStrategy adapts a common.ComparatorStrategy to a ranges Comparator function.
// This allows ordered range structures to accept a unified comparator strategy.
func ComparatorFromStrategy[T comparable](strategy common.ComparatorStrategy[T]) Comparator[T] {
    if strategy == nil {
        // Fallback to default comparator when no strategy provided
        return DefaultComparator[T]
    }
    return func(a, b T) int {
        return strategy.Compare(a, b)
    }
}