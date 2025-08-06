package set

import (
	"testing"
)

func TestTreeSet_New(t *testing.T) {
	ts := NewTreeSet[int]()
	if ts == nil {
		t.Error("NewTreeSet() should not return nil")
	}
	if !ts.IsEmpty() {
		t.Error("New TreeSet should be empty")
	}
	if ts.Size() != 0 {
		t.Errorf("New TreeSet size should be 0, got %d", ts.Size())
	}
}

func TestTreeSet_NewWithComparator(t *testing.T) {
	comparator := func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	
	ts := NewTreeSetWithComparator(comparator)
	if ts == nil {
		t.Error("NewTreeSetWithComparator() should not return nil")
	}
	if !ts.IsEmpty() {
		t.Error("New TreeSet should be empty")
	}
	if ts.Size() != 0 {
		t.Errorf("New TreeSet size should be 0, got %d", ts.Size())
	}
}

func TestTreeSet_Add(t *testing.T) {
	ts := NewTreeSet[int]()
	
	// Test adding to empty set
	added := ts.Add(5)
	if !added {
		t.Error("Add should return true when adding new element")
	}
	if ts.Size() != 1 {
		t.Errorf("TreeSet size should be 1, got %d", ts.Size())
	}
	if ts.IsEmpty() {
		t.Error("TreeSet should not be empty")
	}
	
	// Test adding duplicate element
	added = ts.Add(5)
	if added {
		t.Error("Add should return false when adding duplicate element")
	}
	if ts.Size() != 1 {
		t.Errorf("TreeSet size should remain 1, got %d", ts.Size())
	}
	
	// Test adding multiple elements
	ts.Add(3)
	ts.Add(7)
	ts.Add(1)
	ts.Add(9)
	
	if ts.Size() != 5 {
		t.Errorf("TreeSet size should be 5, got %d", ts.Size())
	}
}

func TestTreeSet_Contains(t *testing.T) {
	ts := NewTreeSet[int]()
	
	// Test contains on empty set
	if ts.Contains(5) {
		t.Error("Empty TreeSet should not contain any element")
	}
	
	// Test contains after adding elements
	ts.Add(5)
	ts.Add(3)
	ts.Add(7)
	
	if !ts.Contains(5) {
		t.Error("TreeSet should contain 5")
	}
	if !ts.Contains(3) {
		t.Error("TreeSet should contain 3")
	}
	if !ts.Contains(7) {
		t.Error("TreeSet should contain 7")
	}
	if ts.Contains(10) {
		t.Error("TreeSet should not contain 10")
	}
}

func TestTreeSet_Remove(t *testing.T) {
	ts := NewTreeSet[int]()
	
	// Test removing from empty set
	removed := ts.Remove(5)
	if removed {
		t.Error("Remove should return false when removing from empty set")
	}
	
	// Test removing existing element
	ts.Add(5)
	ts.Add(3)
	ts.Add(7)
	
	removed = ts.Remove(5)
	if !removed {
		t.Error("Remove should return true when removing existing element")
	}
	if ts.Size() != 2 {
		t.Errorf("TreeSet size should be 2, got %d", ts.Size())
	}
	if ts.Contains(5) {
		t.Error("TreeSet should not contain 5 after removal")
	}
	
	// Test removing non-existing element
	removed = ts.Remove(10)
	if removed {
		t.Error("Remove should return false when removing non-existing element")
	}
	if ts.Size() != 2 {
		t.Errorf("TreeSet size should remain 2, got %d", ts.Size())
	}
}

func TestTreeSet_Clear(t *testing.T) {
	ts := NewTreeSet[int]()
	ts.Add(5)
	ts.Add(3)
	ts.Add(7)
	
	ts.Clear()
	if !ts.IsEmpty() {
		t.Error("TreeSet should be empty after clear")
	}
	if ts.Size() != 0 {
		t.Errorf("TreeSet size should be 0 after clear, got %d", ts.Size())
	}
	if ts.Contains(5) {
		t.Error("TreeSet should not contain any elements after clear")
	}
}

func TestTreeSet_ToSlice(t *testing.T) {
	ts := NewTreeSet[int]()
	
	// Test empty set
	slice := ts.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Empty TreeSet slice should have length 0, got %d", len(slice))
	}
	
	// Test non-empty set - elements should be in sorted order
	ts.Add(5)
	ts.Add(3)
	ts.Add(7)
	ts.Add(1)
	ts.Add(9)
	
	slice = ts.ToSlice()
	if len(slice) != 5 {
		t.Errorf("TreeSet slice should have length 5, got %d", len(slice))
	}
	
	// Check if elements are in sorted order
	expected := []int{1, 3, 5, 7, 9}
	for i, val := range slice {
		if val != expected[i] {
			t.Errorf("Element at index %d should be %d, got %d", i, expected[i], val)
		}
	}
}

func TestTreeSet_Union(t *testing.T) {
	ts1 := NewTreeSet[int]()
	ts1.Add(1)
	ts1.Add(2)
	ts1.Add(3)
	
	ts2 := NewTreeSet[int]()
	ts2.Add(3)
	ts2.Add(4)
	ts2.Add(5)
	
	union := ts1.Union(ts2)
	if union.Size() != 5 {
		t.Errorf("Union size should be 5, got %d", union.Size())
	}
	
	expected := []int{1, 2, 3, 4, 5}
	for _, val := range expected {
		if !union.Contains(val) {
			t.Errorf("Union should contain %d", val)
		}
	}
}

func TestTreeSet_Intersection(t *testing.T) {
	ts1 := NewTreeSet[int]()
	ts1.Add(1)
	ts1.Add(2)
	ts1.Add(3)
	ts1.Add(4)
	
	ts2 := NewTreeSet[int]()
	ts2.Add(3)
	ts2.Add(4)
	ts2.Add(5)
	ts2.Add(6)
	
	intersection := ts1.Intersection(ts2)
	if intersection.Size() != 2 {
		t.Errorf("Intersection size should be 2, got %d", intersection.Size())
	}
	
	if !intersection.Contains(3) {
		t.Error("Intersection should contain 3")
	}
	if !intersection.Contains(4) {
		t.Error("Intersection should contain 4")
	}
	if intersection.Contains(1) {
		t.Error("Intersection should not contain 1")
	}
	if intersection.Contains(5) {
		t.Error("Intersection should not contain 5")
	}
}

func TestTreeSet_Difference(t *testing.T) {
	ts1 := NewTreeSet[int]()
	ts1.Add(1)
	ts1.Add(2)
	ts1.Add(3)
	ts1.Add(4)
	
	ts2 := NewTreeSet[int]()
	ts2.Add(3)
	ts2.Add(4)
	ts2.Add(5)
	ts2.Add(6)
	
	difference := ts1.Difference(ts2)
	if difference.Size() != 2 {
		t.Errorf("Difference size should be 2, got %d", difference.Size())
	}
	
	if !difference.Contains(1) {
		t.Error("Difference should contain 1")
	}
	if !difference.Contains(2) {
		t.Error("Difference should contain 2")
	}
	if difference.Contains(3) {
		t.Error("Difference should not contain 3")
	}
	if difference.Contains(4) {
		t.Error("Difference should not contain 4")
	}
}

func TestTreeSet_IsSubsetOf(t *testing.T) {
	ts1 := NewTreeSet[int]()
	ts1.Add(1)
	ts1.Add(2)
	
	ts2 := NewTreeSet[int]()
	ts2.Add(1)
	ts2.Add(2)
	ts2.Add(3)
	ts2.Add(4)
	
	// Test subset
	if !ts1.IsSubsetOf(ts2) {
		t.Error("ts1 should be a subset of ts2")
	}
	
	// Test not subset
	if ts2.IsSubsetOf(ts1) {
		t.Error("ts2 should not be a subset of ts1")
	}
	
	// Test equal sets
	ts3 := NewTreeSet[int]()
	ts3.Add(1)
	ts3.Add(2)
	
	if !ts1.IsSubsetOf(ts3) {
		t.Error("ts1 should be a subset of ts3 (equal sets)")
	}
	
	// Test empty set
	ts4 := NewTreeSet[int]()
	if !ts4.IsSubsetOf(ts1) {
		t.Error("Empty set should be a subset of any set")
	}
}

func TestTreeSet_ForEach(t *testing.T) {
	ts := NewTreeSet[int]()
	ts.Add(1)
	ts.Add(2)
	ts.Add(3)
	
	sum := 0
	ts.ForEach(func(val int) {
		sum += val
	})
	
	if sum != 6 {
		t.Errorf("Sum should be 6, got %d", sum)
	}
}

func TestTreeSet_String(t *testing.T) {
	ts := NewTreeSet[int]()
	
	// Test empty set
	str := ts.String()
	if str != "{}" {
		t.Errorf("Empty TreeSet string should be '{}', got '%s'", str)
	}
	
	// Test non-empty set
	ts.Add(1)
	ts.Add(2)
	ts.Add(3)
	
	str = ts.String()
	if str == "" {
		t.Error("Non-empty TreeSet string should not be empty")
	}
	// The exact format may vary, but it should contain the elements
	if len(str) < 5 { // At least "{1,2,3}" or similar
		t.Errorf("TreeSet string seems too short: '%s'", str)
	}
}

func TestTreeSet_CustomComparator(t *testing.T) {
	// Test with reverse comparator
	reverseComparator := func(a, b int) int {
		if a > b {
			return -1
		} else if a < b {
			return 1
		}
		return 0
	}
	
	ts := NewTreeSetWithComparator(reverseComparator)
	ts.Add(5)
	ts.Add(3)
	ts.Add(7)
	ts.Add(1)
	ts.Add(9)
	
	slice := ts.ToSlice()
	// With reverse comparator, elements should be in descending order
	expected := []int{9, 7, 5, 3, 1}
	for i, val := range slice {
		if val != expected[i] {
			t.Errorf("Element at index %d should be %d, got %d", i, expected[i], val)
		}
	}
}

func TestTreeSet_EdgeCases(t *testing.T) {
	ts := NewTreeSet[int]()
	
	// Test operations on empty set
	if ts.Remove(1) {
		t.Error("Remove on empty set should return false")
	}
	
	if ts.Contains(1) {
		t.Error("Contains on empty set should return false")
	}
	
	// Test single element operations
	ts.Add(42)
	if !ts.Contains(42) {
		t.Error("TreeSet should contain the single added element")
	}
	
	if ts.Remove(42) {
		if !ts.IsEmpty() {
			t.Error("TreeSet should be empty after removing the only element")
		}
	}
	
	// Test adding same element multiple times
	ts.Clear()
	for i := 0; i < 5; i++ {
		added := ts.Add(10)
		if i == 0 && !added {
			t.Error("First add should return true")
		}
		if i > 0 && added {
			t.Errorf("Add %d should return false for duplicate", i)
		}
	}
	
	if ts.Size() != 1 {
		t.Errorf("TreeSet size should be 1 after adding same element multiple times, got %d", ts.Size())
	}
}