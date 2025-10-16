package ranges

import (
    "sort"
    "strings"
    "sync"
    "github.com/chenjianyu/collections/container/common"
)

// TreeRangeSet is a RangeSet implementation based on a TreeMap
// It maintains ranges in sorted order and provides efficient range operations
type TreeRangeSet[T comparable] struct {
	ranges     []Range[T]
	comparator Comparator[T]
	mutex      sync.RWMutex
}

// NewTreeRangeSet creates a new TreeRangeSet with default comparator
func NewTreeRangeSet[T comparable]() RangeSet[T] {
	return &TreeRangeSet[T]{
		ranges:     make([]Range[T], 0),
		comparator: DefaultComparator[T],
	}
}

// NewTreeRangeSetWithComparator creates a new TreeRangeSet with custom comparator
func NewTreeRangeSetWithComparator[T comparable](cmp Comparator[T]) RangeSet[T] {
    return &TreeRangeSet[T]{
        ranges:     make([]Range[T], 0),
        comparator: cmp,
    }
}

// NewTreeRangeSetWithComparatorStrategy creates a new TreeRangeSet with a ComparatorStrategy
func NewTreeRangeSetWithComparatorStrategy[T comparable](strategy common.ComparatorStrategy[T]) RangeSet[T] {
    return &TreeRangeSet[T]{
        ranges:     make([]Range[T], 0),
        comparator: ComparatorFromStrategy[T](strategy),
    }
}

// Size returns the number of ranges in this set
func (ts *TreeRangeSet[T]) Size() int {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return len(ts.ranges)
}

// IsEmpty returns true if this range set is empty
func (ts *TreeRangeSet[T]) IsEmpty() bool {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return len(ts.ranges) == 0
}

// Clear removes all ranges from this set
func (ts *TreeRangeSet[T]) Clear() {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.ranges = ts.ranges[:0]
}

// String returns the string representation of this range set
func (ts *TreeRangeSet[T]) String() string {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	if len(ts.ranges) == 0 {
		return "{}"
	}
	
	var parts []string
	for _, r := range ts.ranges {
		parts = append(parts, r.String())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// Add adds a range to this range set
func (ts *TreeRangeSet[T]) Add(rangeToAdd Range[T]) {
	if rangeToAdd == nil || rangeToAdd.IsEmpty() {
		return
	}
	
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	ts.addInternal(rangeToAdd)
}

// AddRange adds a range defined by bounds to this range set
func (ts *TreeRangeSet[T]) AddRange(lower T, lowerType BoundType, upper T, upperType BoundType) {
	rangeToAdd := NewRangeWithComparator(lower, lowerType, upper, upperType, ts.comparator)
	ts.Add(rangeToAdd)
}

// addInternal adds a range without locking (internal use)
func (ts *TreeRangeSet[T]) addInternal(rangeToAdd Range[T]) {
	if len(ts.ranges) == 0 {
		ts.ranges = append(ts.ranges, rangeToAdd)
		return
	}
	
	// Find overlapping and adjacent ranges
	var toMerge []Range[T]
	var newRanges []Range[T]
	
	for _, existing := range ts.ranges {
		if existing.IsConnected(rangeToAdd) {
			toMerge = append(toMerge, existing)
		} else {
			newRanges = append(newRanges, existing)
		}
	}
	
	// Merge all connected ranges
	merged := rangeToAdd
	for _, r := range toMerge {
		merged = merged.Span(r)
	}
	
	// Insert the merged range in the correct position
	newRanges = append(newRanges, merged)
	ts.ranges = ts.sortRanges(newRanges)
}

// Remove removes a range from this range set
func (ts *TreeRangeSet[T]) Remove(rangeToRemove Range[T]) {
	if rangeToRemove == nil || rangeToRemove.IsEmpty() {
		return
	}
	
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	ts.removeInternal(rangeToRemove)
}

// RemoveRange removes a range defined by bounds from this range set
func (ts *TreeRangeSet[T]) RemoveRange(lower T, lowerType BoundType, upper T, upperType BoundType) {
	rangeToRemove := NewRangeWithComparator(lower, lowerType, upper, upperType, ts.comparator)
	ts.Remove(rangeToRemove)
}

// removeInternal removes a range without locking (internal use)
func (ts *TreeRangeSet[T]) removeInternal(rangeToRemove Range[T]) {
	var newRanges []Range[T]
	
	for _, existing := range ts.ranges {
		intersection := existing.Intersection(rangeToRemove)
		if intersection.IsEmpty() {
			// No intersection, keep the existing range
			newRanges = append(newRanges, existing)
		} else {
			// There's an intersection, we need to split the existing range
			remainders := ts.subtractRange(existing, rangeToRemove)
			newRanges = append(newRanges, remainders...)
		}
	}
	
	ts.ranges = newRanges
}

// subtractRange subtracts rangeToRemove from existing and returns the remaining ranges
func (ts *TreeRangeSet[T]) subtractRange(existing, rangeToRemove Range[T]) []Range[T] {
	var result []Range[T]
	
	existingLower, existingLowerType, hasExistingLower := existing.LowerBound()
	existingUpper, existingUpperType, hasExistingUpper := existing.UpperBound()
	removeLower, removeLowerType, hasRemoveLower := rangeToRemove.LowerBound()
	removeUpper, removeUpperType, hasRemoveUpper := rangeToRemove.UpperBound()
	
	// Left remainder: from existing lower to remove lower
	if hasExistingLower && hasRemoveLower {
		cmp := ts.comparator(existingLower, removeLower)
		if cmp < 0 || (cmp == 0 && existingLowerType == Closed && removeLowerType == Open) {
			// Create left remainder
			var upperType BoundType
			if removeLowerType == Closed {
				upperType = Open
			} else {
				upperType = Closed
			}
			leftRange := &rangeImpl[T]{
				hasLowerBound: true,
				lowerBound:    existingLower,
				lowerType:     existingLowerType,
				hasUpperBound: true,
				upperBound:    removeLower,
				upperType:     upperType,
				comparator:    ts.comparator,
			}
			if !leftRange.IsEmpty() {
				result = append(result, leftRange)
			}
		}
	} else if hasExistingLower && !hasRemoveLower {
		// Remove range has no lower bound, so no left remainder
	}
	
	// Right remainder: from remove upper to existing upper
	if hasExistingUpper && hasRemoveUpper {
		cmp := ts.comparator(removeUpper, existingUpper)
		if cmp < 0 || (cmp == 0 && removeUpperType == Open && existingUpperType == Closed) {
			// Create right remainder
			var lowerType BoundType
			if removeUpperType == Closed {
				lowerType = Open
			} else {
				lowerType = Closed
			}
			rightRange := &rangeImpl[T]{
				hasLowerBound: true,
				lowerBound:    removeUpper,
				lowerType:     lowerType,
				hasUpperBound: true,
				upperBound:    existingUpper,
				upperType:     existingUpperType,
				comparator:    ts.comparator,
			}
			if !rightRange.IsEmpty() {
				result = append(result, rightRange)
			}
		}
	} else if hasExistingUpper && !hasRemoveUpper {
		// Remove range has no upper bound, so no right remainder
	}
	
	return result
}

// ContainsValue returns true if the value is contained in any range in this set
func (ts *TreeRangeSet[T]) ContainsValue(value T) bool {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	for _, r := range ts.ranges {
		if r.Contains(value) {
			return true
		}
	}
	return false
}

// ContainsRange returns true if the range is entirely contained in this set
func (ts *TreeRangeSet[T]) ContainsRange(rangeToCheck Range[T]) bool {
	if rangeToCheck == nil || rangeToCheck.IsEmpty() {
		return true
	}
	
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	for _, r := range ts.ranges {
		if r.ContainsRange(rangeToCheck) {
			return true
		}
	}
	return false
}

// Encloses returns true if this range set encloses the other range set
func (ts *TreeRangeSet[T]) Encloses(other RangeSet[T]) bool {
	if other == nil || other.IsEmpty() {
		return true
	}
	
	otherRanges := other.AsRanges()
	for _, r := range otherRanges {
		if !ts.ContainsRange(r) {
			return false
		}
	}
	return true
}

// AsRanges returns a view of the disconnected ranges that make up this range set
func (ts *TreeRangeSet[T]) AsRanges() []Range[T] {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	result := make([]Range[T], len(ts.ranges))
	copy(result, ts.ranges)
	return result
}

// Complement returns the complement of this range set
func (ts *TreeRangeSet[T]) Complement() RangeSet[T] {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	
	complement := NewTreeRangeSetWithComparator(ts.comparator)
	
	if len(ts.ranges) == 0 {
		// If this set is empty, complement is all values
		complement.Add(All[T]())
		return complement
	}
	
	// Add range from -∞ to first range's lower bound
	firstRange := ts.ranges[0]
	if lower, lowerType, hasLower := firstRange.LowerBound(); hasLower {
		var upperType BoundType
		if lowerType == Closed {
			upperType = Open
		} else {
			upperType = Closed
		}
		beforeFirst := &rangeImpl[T]{
			hasLowerBound: false,
			hasUpperBound: true,
			upperBound:    lower,
			upperType:     upperType,
			comparator:    ts.comparator,
		}
		if !beforeFirst.IsEmpty() {
			complement.Add(beforeFirst)
		}
	}
	
	// Add ranges between consecutive ranges
	for i := 0; i < len(ts.ranges)-1; i++ {
		current := ts.ranges[i]
		next := ts.ranges[i+1]
		
		currentUpper, currentUpperType, hasCurrentUpper := current.UpperBound()
		nextLower, nextLowerType, hasNextLower := next.LowerBound()
		
		if hasCurrentUpper && hasNextLower {
			var lowerType, upperType BoundType
			if currentUpperType == Closed {
				lowerType = Open
			} else {
				lowerType = Closed
			}
			if nextLowerType == Closed {
				upperType = Open
			} else {
				upperType = Closed
			}
			
			between := &rangeImpl[T]{
				hasLowerBound: true,
				lowerBound:    currentUpper,
				lowerType:     lowerType,
				hasUpperBound: true,
				upperBound:    nextLower,
				upperType:     upperType,
				comparator:    ts.comparator,
			}
			if !between.IsEmpty() {
				complement.Add(between)
			}
		}
	}
	
	// Add range from last range's upper bound to +∞
	lastRange := ts.ranges[len(ts.ranges)-1]
	if upper, upperType, hasUpper := lastRange.UpperBound(); hasUpper {
		var lowerType BoundType
		if upperType == Closed {
			lowerType = Open
		} else {
			lowerType = Closed
		}
		afterLast := &rangeImpl[T]{
			hasLowerBound: true,
			lowerBound:    upper,
			lowerType:     lowerType,
			hasUpperBound: false,
			comparator:    ts.comparator,
		}
		if !afterLast.IsEmpty() {
			complement.Add(afterLast)
		}
	}
	
	return complement
}

// Union returns the union of this range set with another
func (ts *TreeRangeSet[T]) Union(other RangeSet[T]) RangeSet[T] {
	result := NewTreeRangeSetWithComparator(ts.comparator)
	
	// Add all ranges from this set
	for _, r := range ts.AsRanges() {
		result.Add(r)
	}
	
	// Add all ranges from other set
	for _, r := range other.AsRanges() {
		result.Add(r)
	}
	
	return result
}

// Intersection returns the intersection of this range set with another
func (ts *TreeRangeSet[T]) Intersection(other RangeSet[T]) RangeSet[T] {
	result := NewTreeRangeSetWithComparator(ts.comparator)
	
	thisRanges := ts.AsRanges()
	otherRanges := other.AsRanges()
	
	for _, thisRange := range thisRanges {
		for _, otherRange := range otherRanges {
			intersection := thisRange.Intersection(otherRange)
			if !intersection.IsEmpty() {
				result.Add(intersection)
			}
		}
	}
	
	return result
}

// Difference returns the difference of this range set with another
func (ts *TreeRangeSet[T]) Difference(other RangeSet[T]) RangeSet[T] {
	result := NewTreeRangeSetWithComparator(ts.comparator)
	
	// Start with all ranges from this set
	for _, r := range ts.AsRanges() {
		result.Add(r)
	}
	
	// Remove all ranges from other set
	for _, r := range other.AsRanges() {
		result.Remove(r)
	}
	
	return result
}

// sortRanges sorts ranges by their lower bounds
func (ts *TreeRangeSet[T]) sortRanges(ranges []Range[T]) []Range[T] {
	sort.Slice(ranges, func(i, j int) bool {
		rangeI := ranges[i]
		rangeJ := ranges[j]
		
		lowerI, _, hasLowerI := rangeI.LowerBound()
		lowerJ, _, hasLowerJ := rangeJ.LowerBound()
		
		if !hasLowerI && !hasLowerJ {
			return false // Both unbounded, consider equal
		}
		if !hasLowerI {
			return true // Unbounded comes first
		}
		if !hasLowerJ {
			return false // Bounded comes after unbounded
		}
		
		return ts.comparator(lowerI, lowerJ) < 0
	})
	
	return ranges
}