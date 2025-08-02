package ranges

import (
	"fmt"
)

// rangeImpl is the concrete implementation of Range
type rangeImpl[T comparable] struct {
	hasLowerBound bool
	lowerBound    T
	lowerType     BoundType
	hasUpperBound bool
	upperBound    T
	upperType     BoundType
	comparator    Comparator[T]
}

// NewRange creates a new range with the specified bounds
func NewRange[T comparable](lower T, lowerType BoundType, upper T, upperType BoundType) Range[T] {
	return NewRangeWithComparator(lower, lowerType, upper, upperType, DefaultComparator[T])
}

// NewRangeWithComparator creates a new range with the specified bounds and comparator
func NewRangeWithComparator[T comparable](lower T, lowerType BoundType, upper T, upperType BoundType, cmp Comparator[T]) Range[T] {
	return &rangeImpl[T]{
		hasLowerBound: true,
		lowerBound:    lower,
		lowerType:     lowerType,
		hasUpperBound: true,
		upperBound:    upper,
		upperType:     upperType,
		comparator:    cmp,
	}
}

// OpenRange creates an open range (lower, upper)
func OpenRange[T comparable](lower, upper T) Range[T] {
	return NewRange(lower, Open, upper, Open)
}

// ClosedRange creates a closed range [lower, upper]
func ClosedRange[T comparable](lower, upper T) Range[T] {
	return NewRange(lower, Closed, upper, Closed)
}

// OpenClosed creates a range (lower, upper]
func OpenClosed[T comparable](lower, upper T) Range[T] {
	return NewRange(lower, Open, upper, Closed)
}

// ClosedOpen creates a range [lower, upper)
func ClosedOpen[T comparable](lower, upper T) Range[T] {
	return NewRange(lower, Closed, upper, Open)
}

// GreaterThan creates a range (lower, +∞)
func GreaterThan[T comparable](lower T) Range[T] {
	return &rangeImpl[T]{
		hasLowerBound: true,
		lowerBound:    lower,
		lowerType:     Open,
		hasUpperBound: false,
		comparator:    DefaultComparator[T],
	}
}

// AtLeast creates a range [lower, +∞)
func AtLeast[T comparable](lower T) Range[T] {
	return &rangeImpl[T]{
		hasLowerBound: true,
		lowerBound:    lower,
		lowerType:     Closed,
		hasUpperBound: false,
		comparator:    DefaultComparator[T],
	}
}

// LessThan creates a range (-∞, upper)
func LessThan[T comparable](upper T) Range[T] {
	return &rangeImpl[T]{
		hasLowerBound: false,
		hasUpperBound: true,
		upperBound:    upper,
		upperType:     Open,
		comparator:    DefaultComparator[T],
	}
}

// AtMost creates a range (-∞, upper]
func AtMost[T comparable](upper T) Range[T] {
	return &rangeImpl[T]{
		hasLowerBound: false,
		hasUpperBound: true,
		upperBound:    upper,
		upperType:     Closed,
		comparator:    DefaultComparator[T],
	}
}

// All creates a range (-∞, +∞)
func All[T comparable]() Range[T] {
	return &rangeImpl[T]{
		hasLowerBound: false,
		hasUpperBound: false,
		comparator:    DefaultComparator[T],
	}
}

// Singleton creates a range [value, value]
func Singleton[T comparable](value T) Range[T] {
	return ClosedRange(value, value)
}

// LowerBound returns the lower bound of this range
func (r *rangeImpl[T]) LowerBound() (T, BoundType, bool) {
	return r.lowerBound, r.lowerType, r.hasLowerBound
}

// UpperBound returns the upper bound of this range
func (r *rangeImpl[T]) UpperBound() (T, BoundType, bool) {
	return r.upperBound, r.upperType, r.hasUpperBound
}

// Contains returns true if the value is within this range
func (r *rangeImpl[T]) Contains(value T) bool {
	if r.hasLowerBound {
		cmp := r.comparator(value, r.lowerBound)
		if cmp < 0 || (cmp == 0 && r.lowerType == Open) {
			return false
		}
	}
	
	if r.hasUpperBound {
		cmp := r.comparator(value, r.upperBound)
		if cmp > 0 || (cmp == 0 && r.upperType == Open) {
			return false
		}
	}
	
	return true
}

// ContainsRange returns true if the other range is entirely contained within this range
func (r *rangeImpl[T]) ContainsRange(other Range[T]) bool {
	if other.IsEmpty() {
		return true
	}
	
	// Check lower bound
	if lowerBound, lowerType, hasLower := other.LowerBound(); hasLower {
		if !r.Contains(lowerBound) {
			return false
		}
		// If other has a closed lower bound and this range has the same bound but open, it's not contained
		if r.hasLowerBound && r.comparator(lowerBound, r.lowerBound) == 0 {
			if lowerType == Closed && r.lowerType == Open {
				return false
			}
		}
	}
	
	// Check upper bound
	if upperBound, upperType, hasUpper := other.UpperBound(); hasUpper {
		if !r.Contains(upperBound) {
			return false
		}
		// If other has a closed upper bound and this range has the same bound but open, it's not contained
		if r.hasUpperBound && r.comparator(upperBound, r.upperBound) == 0 {
			if upperType == Closed && r.upperType == Open {
				return false
			}
		}
	}
	
	return true
}

// IsConnected returns true if there exists a (possibly empty) range which is
// enclosed by both this range and other
func (r *rangeImpl[T]) IsConnected(other Range[T]) bool {
	// Two ranges are connected if they overlap or are adjacent
	return !r.Intersection(other).IsEmpty() || r.isAdjacent(other)
}

// isAdjacent checks if two ranges are adjacent (touching but not overlapping)
func (r *rangeImpl[T]) isAdjacent(other Range[T]) bool {
	otherImpl, ok := other.(*rangeImpl[T])
	if !ok {
		return false
	}
	
	// Check if r's upper bound touches other's lower bound
	if r.hasUpperBound && otherImpl.hasLowerBound {
		cmp := r.comparator(r.upperBound, otherImpl.lowerBound)
		if cmp == 0 && (r.upperType == Closed || otherImpl.lowerType == Closed) {
			return true
		}
	}
	
	// Check if other's upper bound touches r's lower bound
	if otherImpl.hasUpperBound && r.hasLowerBound {
		cmp := r.comparator(otherImpl.upperBound, r.lowerBound)
		if cmp == 0 && (otherImpl.upperType == Closed || r.lowerType == Closed) {
			return true
		}
	}
	
	return false
}

// Intersection returns the maximal range enclosed by both this range and other
func (r *rangeImpl[T]) Intersection(other Range[T]) Range[T] {
	otherImpl, ok := other.(*rangeImpl[T])
	if !ok {
		return nil
	}
	
	// Calculate intersection bounds
	var lowerBound T
	var lowerType BoundType
	var hasLower bool
	
	var upperBound T
	var upperType BoundType
	var hasUpper bool
	
	// Lower bound of intersection
	if !r.hasLowerBound && !otherImpl.hasLowerBound {
		hasLower = false
	} else if !r.hasLowerBound {
		lowerBound = otherImpl.lowerBound
		lowerType = otherImpl.lowerType
		hasLower = true
	} else if !otherImpl.hasLowerBound {
		lowerBound = r.lowerBound
		lowerType = r.lowerType
		hasLower = true
	} else {
		cmp := r.comparator(r.lowerBound, otherImpl.lowerBound)
		if cmp > 0 {
			lowerBound = r.lowerBound
			lowerType = r.lowerType
		} else if cmp < 0 {
			lowerBound = otherImpl.lowerBound
			lowerType = otherImpl.lowerType
		} else {
			// Same bound, use the more restrictive type
			lowerBound = r.lowerBound
			if r.lowerType == Open || otherImpl.lowerType == Open {
				lowerType = Open
			} else {
				lowerType = Closed
			}
		}
		hasLower = true
	}
	
	// Upper bound of intersection
	if !r.hasUpperBound && !otherImpl.hasUpperBound {
		hasUpper = false
	} else if !r.hasUpperBound {
		upperBound = otherImpl.upperBound
		upperType = otherImpl.upperType
		hasUpper = true
	} else if !otherImpl.hasUpperBound {
		upperBound = r.upperBound
		upperType = r.upperType
		hasUpper = true
	} else {
		cmp := r.comparator(r.upperBound, otherImpl.upperBound)
		if cmp < 0 {
			upperBound = r.upperBound
			upperType = r.upperType
		} else if cmp > 0 {
			upperBound = otherImpl.upperBound
			upperType = otherImpl.upperType
		} else {
			// Same bound, use the more restrictive type
			upperBound = r.upperBound
			if r.upperType == Open || otherImpl.upperType == Open {
				upperType = Open
			} else {
				upperType = Closed
			}
		}
		hasUpper = true
	}
	
	// Check if the intersection is valid
	if hasLower && hasUpper {
		cmp := r.comparator(lowerBound, upperBound)
		if cmp > 0 || (cmp == 0 && (lowerType == Open || upperType == Open)) {
			// Empty intersection
			return &rangeImpl[T]{
				hasLowerBound: true,
				lowerBound:    lowerBound,
				lowerType:     Open,
				hasUpperBound: true,
				upperBound:    lowerBound,
				upperType:     Open,
				comparator:    r.comparator,
			}
		}
	}
	
	result := &rangeImpl[T]{
		hasLowerBound: hasLower,
		hasUpperBound: hasUpper,
		comparator:    r.comparator,
	}
	
	if hasLower {
		result.lowerBound = lowerBound
		result.lowerType = lowerType
	}
	
	if hasUpper {
		result.upperBound = upperBound
		result.upperType = upperType
	}
	
	return result
}

// Span returns the minimal range that encloses both this range and other
func (r *rangeImpl[T]) Span(other Range[T]) Range[T] {
	otherImpl, ok := other.(*rangeImpl[T])
	if !ok {
		return r
	}
	
	// Calculate span bounds
	var lowerBound T
	var lowerType BoundType
	var hasLower bool
	
	var upperBound T
	var upperType BoundType
	var hasUpper bool
	
	// Lower bound of span
	if !r.hasLowerBound || !otherImpl.hasLowerBound {
		hasLower = false
	} else {
		cmp := r.comparator(r.lowerBound, otherImpl.lowerBound)
		if cmp < 0 {
			lowerBound = r.lowerBound
			lowerType = r.lowerType
		} else if cmp > 0 {
			lowerBound = otherImpl.lowerBound
			lowerType = otherImpl.lowerType
		} else {
			// Same bound, use the less restrictive type
			lowerBound = r.lowerBound
			if r.lowerType == Closed || otherImpl.lowerType == Closed {
				lowerType = Closed
			} else {
				lowerType = Open
			}
		}
		hasLower = true
	}
	
	// Upper bound of span
	if !r.hasUpperBound || !otherImpl.hasUpperBound {
		hasUpper = false
	} else {
		cmp := r.comparator(r.upperBound, otherImpl.upperBound)
		if cmp > 0 {
			upperBound = r.upperBound
			upperType = r.upperType
		} else if cmp < 0 {
			upperBound = otherImpl.upperBound
			upperType = otherImpl.upperType
		} else {
			// Same bound, use the less restrictive type
			upperBound = r.upperBound
			if r.upperType == Closed || otherImpl.upperType == Closed {
				upperType = Closed
			} else {
				upperType = Open
			}
		}
		hasUpper = true
	}
	
	result := &rangeImpl[T]{
		hasLowerBound: hasLower,
		hasUpperBound: hasUpper,
		comparator:    r.comparator,
	}
	
	if hasLower {
		result.lowerBound = lowerBound
		result.lowerType = lowerType
	}
	
	if hasUpper {
		result.upperBound = upperBound
		result.upperType = upperType
	}
	
	return result
}

// IsEmpty returns true if this range is empty
func (r *rangeImpl[T]) IsEmpty() bool {
	if !r.hasLowerBound || !r.hasUpperBound {
		return false
	}
	
	cmp := r.comparator(r.lowerBound, r.upperBound)
	return cmp > 0 || (cmp == 0 && (r.lowerType == Open || r.upperType == Open))
}

// String returns the string representation of this range
func (r *rangeImpl[T]) String() string {
	if r.IsEmpty() {
		return "(empty)"
	}
	
	var result string
	
	// Lower bound
	if !r.hasLowerBound {
		result += "(-∞"
	} else {
		if r.lowerType == Closed {
			result += fmt.Sprintf("[%v", r.lowerBound)
		} else {
			result += fmt.Sprintf("(%v", r.lowerBound)
		}
	}
	
	result += ".."
	
	// Upper bound
	if !r.hasUpperBound {
		result += "+∞)"
	} else {
		if r.upperType == Closed {
			result += fmt.Sprintf("%v]", r.upperBound)
		} else {
			result += fmt.Sprintf("%v)", r.upperBound)
		}
	}
	
	return result
}