package ranges

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRangeCreation(t *testing.T) {
	// Test closed range
	r1 := ClosedRange(1, 10)
	if r1.IsEmpty() {
		t.Error("Closed range should not be empty")
	}
	
	// Test open range
	r2 := OpenRange(1, 10)
	if r2.IsEmpty() {
		t.Error("Open range should not be empty")
	}
	
	// Test singleton range
	r3 := Singleton(5)
	if r3.IsEmpty() {
		t.Error("Singleton range should not be empty")
	}
	if !r3.Contains(5) {
		t.Error("Singleton range should contain its value")
	}
}

func TestTreeRangeSet(t *testing.T) {
	rs := NewTreeRangeSet[int]()
	
	// Test empty set
	if !rs.IsEmpty() {
		t.Error("New range set should be empty")
	}
	if rs.Size() != 0 {
		t.Error("Empty range set should have size 0")
	}
	
	// Test adding ranges
	rs.Add(ClosedRange(1, 5))
	if rs.IsEmpty() {
		t.Error("Range set should not be empty after adding range")
	}
	if rs.Size() != 1 {
		t.Error("Range set should have size 1 after adding one range")
	}
	
	// Test contains
	if !rs.ContainsValue(3) {
		t.Error("Range set should contain value 3")
	}
	if rs.ContainsValue(10) {
		t.Error("Range set should not contain value 10")
	}
}

func TestImmutableRangeSet(t *testing.T) {
	rs := NewImmutableRangeSet[int]()
	
	// Test empty set
	if !rs.IsEmpty() {
		t.Error("New immutable range set should be empty")
	}
	
	// Test that modification methods log warnings instead of panicking
	rs.Add(ClosedRange(1, 5)) // Should log warning but not panic
	
	// Verify the set remains unchanged
	if !rs.IsEmpty() {
		t.Error("Immutable range set should remain empty after Add operation")
	}
}

func TestTreeRangeMap(t *testing.T) {
	rm := NewTreeRangeMap[int, string]()
	
	// Test empty map
	if !rm.IsEmpty() {
		t.Error("New range map should be empty")
	}
	if rm.Size() != 0 {
		t.Error("Empty range map should have size 0")
	}
	
	// Test putting values
	rm.Put(ClosedRange(1, 5), "first")
	if rm.IsEmpty() {
		t.Error("Range map should not be empty after putting value")
	}
	if rm.Size() != 1 {
		t.Error("Range map should have size 1 after putting one value")
	}
	
	// Test getting values
	if val, ok := rm.Get(3); !ok || val != "first" {
		t.Error("Range map should return 'first' for key 3")
	}
	if _, ok := rm.Get(10); ok {
		t.Error("Range map should not have value for key 10")
	}
}

func TestImmutableRangeMap(t *testing.T) {
	rm := NewImmutableRangeMap[int, string]()
	
	// Test empty map
	if !rm.IsEmpty() {
		t.Error("New immutable range map should be empty")
	}
	
	// Test that modification methods log warnings instead of panicking
	rm.Put(ClosedRange(1, 5), "test") // Should log warning but not panic
	
	// Verify the map remains unchanged
	if !rm.IsEmpty() {
		t.Error("Immutable range map should remain empty after Put operation")
	}
}

func TestRangeOperations(t *testing.T) {
	r1 := ClosedRange(1, 5)
	r2 := ClosedRange(3, 7)
	
	// Test intersection
		intersection := r1.Intersection(r2)
		if intersection.IsEmpty() {
			t.Error("Intersection should not be empty")
		}
		
		// Test connection
		if !r1.IsConnected(r2) {
			t.Error("Ranges should be connected")
		}
		
		// Test contains range
		r3 := ClosedRange(2, 4)
		if !r1.ContainsRange(r3) {
			t.Error("Range [1,5] should contain range [2,4]")
		}
}

// TestRangeImplMethods tests additional methods in rangeImpl
func TestRangeImplMethods(t *testing.T) {
	t.Run("OpenRange", func(t *testing.T) {
		r := OpenRange(1, 5)
		assert.False(t, r.Contains(1))
		assert.True(t, r.Contains(3))
		assert.False(t, r.Contains(5))
		assert.Equal(t, "(1..5)", r.String())
	})

	t.Run("OpenClosed", func(t *testing.T) {
		r := OpenClosed(1, 5)
		assert.False(t, r.Contains(1))
		assert.True(t, r.Contains(3))
		assert.True(t, r.Contains(5))
		assert.Equal(t, "(1..5]", r.String())
	})

	t.Run("ClosedOpen", func(t *testing.T) {
		r := ClosedOpen(1, 5)
		assert.True(t, r.Contains(1))
		assert.True(t, r.Contains(3))
		assert.False(t, r.Contains(5))
		assert.Equal(t, "[1..5)", r.String())
	})

	t.Run("GreaterThan", func(t *testing.T) {
		r := GreaterThan(5)
		assert.False(t, r.Contains(5))
		assert.True(t, r.Contains(6))
		assert.True(t, r.Contains(100))
		assert.Equal(t, "(5..+∞)", r.String())
	})

	t.Run("AtLeast", func(t *testing.T) {
		r := AtLeast(5)
		assert.True(t, r.Contains(5))
		assert.True(t, r.Contains(6))
		assert.True(t, r.Contains(100))
		assert.Equal(t, "[5..+∞)", r.String())
	})

	t.Run("LessThan", func(t *testing.T) {
		r := LessThan(5)
		assert.True(t, r.Contains(4))
		assert.False(t, r.Contains(5))
		assert.False(t, r.Contains(6))
		assert.Equal(t, "(-∞..5)", r.String())
	})

	t.Run("AtMost", func(t *testing.T) {
		r := AtMost(5)
		assert.True(t, r.Contains(4))
		assert.True(t, r.Contains(5))
		assert.False(t, r.Contains(6))
		assert.Equal(t, "(-∞..5]", r.String())
	})

	t.Run("All", func(t *testing.T) {
		r := All[int]()
		assert.True(t, r.Contains(-1000))
		assert.True(t, r.Contains(0))
		assert.True(t, r.Contains(1000))
		assert.Equal(t, "(-∞..+∞)", r.String())
	})
}

// TestRangeSpan tests the Span method
func TestRangeSpan(t *testing.T) {
	t.Run("OverlappingRanges", func(t *testing.T) {
		r1 := ClosedRange(1, 5)
		r2 := ClosedRange(3, 8)
		span := r1.Span(r2)
		assert.True(t, span.Contains(1))
		assert.True(t, span.Contains(8))
		assert.Equal(t, "[1..8]", span.String())
	})

	t.Run("DisjointRanges", func(t *testing.T) {
		r1 := ClosedRange(1, 3)
		r2 := ClosedRange(6, 8)
		span := r1.Span(r2)
		assert.True(t, span.Contains(1))
		assert.True(t, span.Contains(5)) // Should include gap
		assert.True(t, span.Contains(8))
		assert.Equal(t, "[1..8]", span.String())
	})

	t.Run("UnboundedRanges", func(t *testing.T) {
		r1 := AtLeast(5)
		r2 := AtMost(10)
		span := r1.Span(r2)
		assert.True(t, span.Contains(-100))
		assert.True(t, span.Contains(100))
		assert.Equal(t, "(-∞..+∞)", span.String())
	})
}

// TestRangeIsEmpty tests the IsEmpty method
func TestRangeIsEmpty(t *testing.T) {
	t.Run("ValidRange", func(t *testing.T) {
		r := ClosedRange(1, 5)
		assert.False(t, r.IsEmpty())
	})

	t.Run("SingletonRange", func(t *testing.T) {
		r := Singleton(5)
		assert.False(t, r.IsEmpty())
	})

	t.Run("OpenRangeWithSameBounds", func(t *testing.T) {
		r := OpenRange(5, 5)
		assert.True(t, r.IsEmpty())
	})

	t.Run("InvalidRange", func(t *testing.T) {
		r := ClosedRange(5, 3) // upper < lower
		assert.True(t, r.IsEmpty())
	})
}

// TestTreeRangeSetAdvanced tests advanced TreeRangeSet operations
func TestTreeRangeSetAdvanced(t *testing.T) {
	t.Run("AsRanges", func(t *testing.T) {
		rs := NewTreeRangeSet[int]()
		rs.Add(ClosedRange(1, 3))
		rs.Add(ClosedRange(5, 7))
		rs.Add(ClosedRange(9, 11))

		ranges := rs.AsRanges()
		assert.Equal(t, 3, len(ranges))
		assert.True(t, ranges[0].Contains(2))
		assert.True(t, ranges[1].Contains(6))
		assert.True(t, ranges[2].Contains(10))
	})

	t.Run("ContainsValue", func(t *testing.T) {
		rs := NewTreeRangeSet[int]()
		rs.Add(ClosedRange(1, 5))
		rs.Add(ClosedRange(10, 15))

		assert.True(t, rs.ContainsValue(3))
		assert.True(t, rs.ContainsValue(12))
		assert.False(t, rs.ContainsValue(7))
		assert.False(t, rs.ContainsValue(20))
	})

	t.Run("ContainsRange", func(t *testing.T) {
		rs := NewTreeRangeSet[int]()
		rs.Add(ClosedRange(1, 10))

		assert.True(t, rs.ContainsRange(ClosedRange(3, 7)))
		assert.False(t, rs.ContainsRange(ClosedRange(5, 15)))
		assert.True(t, rs.ContainsRange(Singleton(5)))
	})

	t.Run("Encloses", func(t *testing.T) {
		rs1 := NewTreeRangeSet[int]()
		rs1.Add(ClosedRange(1, 10))
		rs1.Add(ClosedRange(20, 30))

		rs2 := NewTreeRangeSet[int]()
		rs2.Add(ClosedRange(2, 5))
		rs2.Add(ClosedRange(22, 25))

		assert.True(t, rs1.Encloses(rs2))
		assert.False(t, rs2.Encloses(rs1))
	})

	t.Run("Complement", func(t *testing.T) {
		rs := NewTreeRangeSet[int]()
		rs.Add(ClosedRange(5, 10))

		complement := rs.Complement()
		assert.True(t, complement.ContainsValue(3))
		assert.False(t, complement.ContainsValue(7))
		assert.True(t, complement.ContainsValue(15))
	})

	t.Run("Union", func(t *testing.T) {
		rs1 := NewTreeRangeSet[int]()
		rs1.Add(ClosedRange(1, 5))

		rs2 := NewTreeRangeSet[int]()
		rs2.Add(ClosedRange(3, 8))

		union := rs1.Union(rs2)
		assert.True(t, union.ContainsValue(1))
		assert.True(t, union.ContainsValue(4))
		assert.True(t, union.ContainsValue(8))
	})

	t.Run("Intersection", func(t *testing.T) {
		rs1 := NewTreeRangeSet[int]()
		rs1.Add(ClosedRange(1, 8))

		rs2 := NewTreeRangeSet[int]()
		rs2.Add(ClosedRange(3, 10))

		intersection := rs1.Intersection(rs2)
		assert.False(t, intersection.ContainsValue(2))
		assert.True(t, intersection.ContainsValue(5))
		assert.False(t, intersection.ContainsValue(9))
	})

	t.Run("Difference", func(t *testing.T) {
		rs1 := NewTreeRangeSet[int]()
		rs1.Add(ClosedRange(1, 10))

		rs2 := NewTreeRangeSet[int]()
		rs2.Add(ClosedRange(3, 7))

		difference := rs1.Difference(rs2)
		assert.True(t, difference.ContainsValue(2))
		assert.False(t, difference.ContainsValue(5))
		assert.True(t, difference.ContainsValue(9))
	})
}

// TestTreeRangeMapAdvanced tests advanced TreeRangeMap operations
func TestTreeRangeMapAdvanced(t *testing.T) {
	t.Run("GetEntry", func(t *testing.T) {
		rm := NewTreeRangeMap[int, string]()
		rm.Put(ClosedRange(1, 5), "first")
		rm.Put(ClosedRange(10, 15), "second")

		rangeKey, value, found := rm.GetEntry(3)
		assert.True(t, found)
		assert.Equal(t, "first", value)
		assert.True(t, rangeKey.Contains(3))

		_, _, found = rm.GetEntry(7)
		assert.False(t, found)
	})

	t.Run("AsMapOfRanges", func(t *testing.T) {
		rm := NewTreeRangeMap[int, string]()
		rm.Put(ClosedRange(1, 5), "first")
		rm.Put(ClosedRange(10, 15), "second")

		rangeMap := rm.AsMapOfRanges()
		assert.Equal(t, 2, len(rangeMap))
	})

	t.Run("AsDescendingMapOfRanges", func(t *testing.T) {
		rm := NewTreeRangeMap[int, string]()
		rm.Put(ClosedRange(1, 5), "first")
		rm.Put(ClosedRange(10, 15), "second")

		descendingMap := rm.AsDescendingMapOfRanges()
		assert.Equal(t, 2, len(descendingMap))
	})

	t.Run("SubRangeMap", func(t *testing.T) {
		rm := NewTreeRangeMap[int, string]()
		rm.Put(ClosedRange(1, 10), "full")
		rm.Put(ClosedRange(15, 20), "other")

		subMap := rm.SubRangeMap(ClosedRange(3, 7))
		value, found := subMap.Get(5)
		assert.True(t, found)
		assert.Equal(t, "full", value)

		_, found = subMap.Get(18)
		assert.False(t, found)
	})

	t.Run("Span", func(t *testing.T) {
		rm := NewTreeRangeMap[int, string]()
		rm.Put(ClosedRange(1, 5), "first")
		rm.Put(ClosedRange(10, 15), "second")

		span, hasSpan := rm.Span()
		assert.True(t, hasSpan)
		assert.True(t, span.Contains(1))
		assert.True(t, span.Contains(15))
		assert.True(t, span.Contains(7)) // Should include gap
	})

	t.Run("EmptyMapSpan", func(t *testing.T) {
		rm := NewTreeRangeMap[int, string]()
		_, hasSpan := rm.Span()
		assert.False(t, hasSpan)
	})
}

// TestImmutableRangeSetAdvanced tests advanced ImmutableRangeSet operations
func TestImmutableRangeSetAdvanced(t *testing.T) {
	t.Run("WithOperations", func(t *testing.T) {
		irs := NewImmutableRangeSet[int]()
		
		// Test WithAdd
		irs2 := irs.(*ImmutableRangeSet[int]).WithAdd(ClosedRange(1, 5))
		assert.True(t, irs.IsEmpty())
		assert.False(t, irs2.IsEmpty())
		assert.True(t, irs2.ContainsValue(3))

		// Test WithAddRange
		irs3 := irs2.(*ImmutableRangeSet[int]).WithAddRange(10, Closed, 15, Closed)
		assert.Equal(t, 1, irs2.Size())
		assert.Equal(t, 2, irs3.Size())

		// Test WithRemove
		irs4 := irs3.(*ImmutableRangeSet[int]).WithRemove(ClosedRange(1, 5))
		assert.Equal(t, 2, irs3.Size())
		assert.Equal(t, 1, irs4.Size())
		assert.True(t, irs4.ContainsValue(12))

		// Test WithRemoveRange
		irs5 := irs4.(*ImmutableRangeSet[int]).WithRemoveRange(10, Closed, 15, Closed)
		assert.Equal(t, 1, irs4.Size())
		assert.True(t, irs5.IsEmpty())

		// Test WithClear
		irs6 := irs3.(*ImmutableRangeSet[int]).WithClear()
		assert.Equal(t, 2, irs3.Size())
		assert.True(t, irs6.IsEmpty())
	})

	t.Run("ImmutableOperationWarnings", func(t *testing.T) {
		irs := NewImmutableRangeSet[int]()
		
		// These should log warnings but not panic
		irs.Add(ClosedRange(1, 5))
		irs.AddRange(1, Closed, 5, Closed)
		irs.Remove(ClosedRange(1, 5))
		irs.RemoveRange(1, Closed, 5, Closed)
		irs.Clear()
		
		// Original should remain unchanged
		assert.True(t, irs.IsEmpty())
	})

	t.Run("NewImmutableRangeSetFromRanges", func(t *testing.T) {
		ranges := []Range[int]{
			ClosedRange(1, 3),
			ClosedRange(2, 5), // Overlapping - should be merged
			ClosedRange(10, 15),
		}
		
		irs := NewImmutableRangeSetFromRanges(ranges)
		assert.Equal(t, 2, irs.Size()) // Should be merged to 2 ranges
		assert.True(t, irs.ContainsValue(4))
		assert.True(t, irs.ContainsValue(12))
		assert.False(t, irs.ContainsValue(7))
	})
}

// TestRangeEdgeCases tests edge cases and error conditions
func TestRangeEdgeCases(t *testing.T) {
	t.Run("NilRangeOperations", func(t *testing.T) {
		rs := NewTreeRangeSet[int]()
		rs.Add(nil) // Should not panic
		assert.True(t, rs.IsEmpty())

		rm := NewTreeRangeMap[int, string]()
		rm.Put(nil, "test") // Should not panic
		assert.True(t, rm.IsEmpty())
	})

	t.Run("EmptyRangeOperations", func(t *testing.T) {
		emptyRange := OpenRange(5, 5) // Empty range
		
		rs := NewTreeRangeSet[int]()
		rs.Add(emptyRange)
		assert.True(t, rs.IsEmpty())

		rm := NewTreeRangeMap[int, string]()
		rm.Put(emptyRange, "test")
		assert.True(t, rm.IsEmpty())
	})

	t.Run("CustomComparator", func(t *testing.T) {
		// Reverse comparator
		reverseComparator := func(a, b int) int {
			if a > b {
				return -1
			} else if a < b {
				return 1
			}
			return 0
		}

		rs := NewTreeRangeSetWithComparator(reverseComparator)
		// With reverse comparator, we need to create a valid range
		// In reverse order, 5 should come before 1, so [5,1] is valid
		rs.Add(NewRangeWithComparator(5, Closed, 1, Closed, reverseComparator))
		assert.False(t, rs.IsEmpty())
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		rs := NewTreeRangeSet[int]()
		
		// Test concurrent reads and writes
		done := make(chan bool, 2)
		
		go func() {
			for i := 0; i < 100; i++ {
				rs.Add(ClosedRange(i, i+10))
			}
			done <- true
		}()
		
		go func() {
			for i := 0; i < 100; i++ {
				rs.ContainsValue(i)
				rs.Size()
			}
			done <- true
		}()
		
		<-done
		<-done
		
		assert.False(t, rs.IsEmpty())
	})
}