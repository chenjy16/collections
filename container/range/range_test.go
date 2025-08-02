package ranges

import (
	"testing"
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
	
	// Test that modification methods panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Add should panic on immutable range set")
		}
	}()
	rs.Add(ClosedRange(1, 5))
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
	
	// Test that modification methods panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Put should panic on immutable range map")
		}
	}()
	rm.Put(ClosedRange(1, 5), "test")
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