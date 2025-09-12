package ranges

import (
	"strings"
)

// ImmutableRangeSet is an immutable implementation of RangeSet
// All modification operations return new instances
type ImmutableRangeSet[T comparable] struct {
	ranges     []Range[T]
	comparator Comparator[T]
}

// NewImmutableRangeSet creates a new empty ImmutableRangeSet
func NewImmutableRangeSet[T comparable]() RangeSet[T] {
	return &ImmutableRangeSet[T]{
		ranges:     make([]Range[T], 0),
		comparator: DefaultComparator[T],
	}
}

// NewImmutableRangeSetWithComparator creates a new ImmutableRangeSet with custom comparator
func NewImmutableRangeSetWithComparator[T comparable](cmp Comparator[T]) RangeSet[T] {
	return &ImmutableRangeSet[T]{
		ranges:     make([]Range[T], 0),
		comparator: cmp,
	}
}

// NewImmutableRangeSetFromRanges creates a new ImmutableRangeSet from existing ranges
func NewImmutableRangeSetFromRanges[T comparable](ranges []Range[T]) RangeSet[T] {
	// Create a mutable set to handle merging and sorting
	mutableSet := NewTreeRangeSet[T]()
	for _, r := range ranges {
		mutableSet.Add(r)
	}

	return &ImmutableRangeSet[T]{
		ranges:     mutableSet.AsRanges(),
		comparator: DefaultComparator[T],
	}
}

// Size returns the number of ranges in this set
func (irs *ImmutableRangeSet[T]) Size() int {
	return len(irs.ranges)
}

// IsEmpty returns true if this range set is empty
func (irs *ImmutableRangeSet[T]) IsEmpty() bool {
	return len(irs.ranges) == 0
}

// Clear is not supported for ImmutableRangeSet - this is a no-op
func (irs *ImmutableRangeSet[T]) Clear() {
	// No-op for immutable collections
}

// WithClear returns a new empty ImmutableRangeSet
func (irs *ImmutableRangeSet[T]) WithClear() RangeSet[T] {
	return NewImmutableRangeSetWithComparator(irs.comparator)
}

// String returns the string representation of this range set
func (irs *ImmutableRangeSet[T]) String() string {
	if len(irs.ranges) == 0 {
		return "{}"
	}

	var parts []string
	for _, r := range irs.ranges {
		parts = append(parts, r.String())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// Add is not supported for ImmutableRangeSet - this is a no-op
func (irs *ImmutableRangeSet[T]) Add(rangeToAdd Range[T]) {
	// No-op for immutable collections
}

// WithAdd returns a new ImmutableRangeSet with the range added
func (irs *ImmutableRangeSet[T]) WithAdd(rangeToAdd Range[T]) RangeSet[T] {
	if rangeToAdd == nil || rangeToAdd.IsEmpty() {
		return irs
	}

	// Create a mutable copy and add the range
	mutableSet := NewTreeRangeSetWithComparator(irs.comparator)
	for _, r := range irs.ranges {
		mutableSet.Add(r)
	}
	mutableSet.Add(rangeToAdd)

	return &ImmutableRangeSet[T]{
		ranges:     mutableSet.AsRanges(),
		comparator: irs.comparator,
	}
}

// AddRange is not supported for ImmutableRangeSet - this is a no-op
func (irs *ImmutableRangeSet[T]) AddRange(lower T, lowerType BoundType, upper T, upperType BoundType) {
	// No-op for immutable collections
}

// WithAddRange returns a new ImmutableRangeSet with the range added
func (irs *ImmutableRangeSet[T]) WithAddRange(lower T, lowerType BoundType, upper T, upperType BoundType) RangeSet[T] {
	rangeToAdd := NewRangeWithComparator(lower, lowerType, upper, upperType, irs.comparator)
	return irs.WithAdd(rangeToAdd)
}

// Remove is not supported for ImmutableRangeSet - this is a no-op
func (irs *ImmutableRangeSet[T]) Remove(rangeToRemove Range[T]) {
	// No-op for immutable collections
}

// WithRemove returns a new ImmutableRangeSet with the range removed
func (irs *ImmutableRangeSet[T]) WithRemove(rangeToRemove Range[T]) RangeSet[T] {
	if rangeToRemove == nil || rangeToRemove.IsEmpty() {
		return irs
	}

	// Create a mutable copy and remove the range
	mutableSet := NewTreeRangeSetWithComparator(irs.comparator)
	for _, r := range irs.ranges {
		mutableSet.Add(r)
	}
	mutableSet.Remove(rangeToRemove)

	return &ImmutableRangeSet[T]{
		ranges:     mutableSet.AsRanges(),
		comparator: irs.comparator,
	}
}

// RemoveRange is a no-op as ImmutableRangeSet is immutable
func (irs *ImmutableRangeSet[T]) RemoveRange(lower T, lowerType BoundType, upper T, upperType BoundType) {
	// No-op for immutable collections
}

// WithRemoveRange returns a new ImmutableRangeSet with the range removed
func (irs *ImmutableRangeSet[T]) WithRemoveRange(lower T, lowerType BoundType, upper T, upperType BoundType) RangeSet[T] {
	rangeToRemove := NewRangeWithComparator(lower, lowerType, upper, upperType, irs.comparator)
	return irs.WithRemove(rangeToRemove)
}

// ContainsValue returns true if the value is contained in any range in this set
func (irs *ImmutableRangeSet[T]) ContainsValue(value T) bool {
	for _, r := range irs.ranges {
		if r.Contains(value) {
			return true
		}
	}
	return false
}

// ContainsRange returns true if the range is entirely contained in this set
func (irs *ImmutableRangeSet[T]) ContainsRange(rangeToCheck Range[T]) bool {
	if rangeToCheck == nil || rangeToCheck.IsEmpty() {
		return true
	}

	for _, r := range irs.ranges {
		if r.ContainsRange(rangeToCheck) {
			return true
		}
	}
	return false
}

// Encloses returns true if this range set encloses the other range set
func (irs *ImmutableRangeSet[T]) Encloses(other RangeSet[T]) bool {
	if other == nil || other.IsEmpty() {
		return true
	}

	otherRanges := other.AsRanges()
	for _, r := range otherRanges {
		if !irs.ContainsRange(r) {
			return false
		}
	}
	return true
}

// AsRanges returns a view of the disconnected ranges that make up this range set
func (irs *ImmutableRangeSet[T]) AsRanges() []Range[T] {
	result := make([]Range[T], len(irs.ranges))
	copy(result, irs.ranges)
	return result
}

// Complement returns the complement of this range set
func (irs *ImmutableRangeSet[T]) Complement() RangeSet[T] {
	// Create a mutable set to compute complement
	mutableSet := NewTreeRangeSetWithComparator(irs.comparator)
	for _, r := range irs.ranges {
		mutableSet.Add(r)
	}

	complement := mutableSet.Complement()
	return &ImmutableRangeSet[T]{
		ranges:     complement.AsRanges(),
		comparator: irs.comparator,
	}
}

// Union returns the union of this range set with another
func (irs *ImmutableRangeSet[T]) Union(other RangeSet[T]) RangeSet[T] {
	// Create a mutable set to compute union
	mutableSet := NewTreeRangeSetWithComparator(irs.comparator)

	// Add all ranges from this set
	for _, r := range irs.ranges {
		mutableSet.Add(r)
	}

	// Add all ranges from other set
	for _, r := range other.AsRanges() {
		mutableSet.Add(r)
	}

	return &ImmutableRangeSet[T]{
		ranges:     mutableSet.AsRanges(),
		comparator: irs.comparator,
	}
}

// Intersection returns the intersection of this range set with another
func (irs *ImmutableRangeSet[T]) Intersection(other RangeSet[T]) RangeSet[T] {
	result := NewTreeRangeSetWithComparator(irs.comparator)

	otherRanges := other.AsRanges()

	for _, thisRange := range irs.ranges {
		for _, otherRange := range otherRanges {
			intersection := thisRange.Intersection(otherRange)
			if !intersection.IsEmpty() {
				result.Add(intersection)
			}
		}
	}

	return &ImmutableRangeSet[T]{
		ranges:     result.AsRanges(),
		comparator: irs.comparator,
	}
}

// Difference returns the difference of this range set with another
func (irs *ImmutableRangeSet[T]) Difference(other RangeSet[T]) RangeSet[T] {
	// Create a mutable set to compute difference
	mutableSet := NewTreeRangeSetWithComparator(irs.comparator)

	// Start with all ranges from this set
	for _, r := range irs.ranges {
		mutableSet.Add(r)
	}

	// Remove all ranges from other set
	for _, r := range other.AsRanges() {
		mutableSet.Remove(r)
	}

	return &ImmutableRangeSet[T]{
		ranges:     mutableSet.AsRanges(),
		comparator: irs.comparator,
	}
}
