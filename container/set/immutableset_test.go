package set

import (
	"testing"
)

func TestNewImmutableSet(t *testing.T) {
	set := NewImmutableSet[int]()
	
	if !set.IsEmpty() {
		t.Error("New set should be empty")
	}
	
	if set.Size() != 0 {
		t.Errorf("Expected size 0, got %d", set.Size())
	}
}

func TestNewImmutableSetFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 2, 4} // Note: duplicate 2
	set := NewImmutableSetFromSlice(slice)
	
	if set.Size() != 4 { // Should be 4 unique elements
		t.Errorf("Expected size 4, got %d", set.Size())
	}
	
	expectedElements := []int{1, 2, 3, 4}
	for _, expected := range expectedElements {
		if !set.Contains(expected) {
			t.Errorf("Set should contain %d", expected)
		}
	}
}

func TestSetOf(t *testing.T) {
	set := SetOf(1, 2, 3, 2, 4) // Note: duplicate 2
	
	if set.Size() != 4 { // Should be 4 unique elements
		t.Errorf("Expected size 4, got %d", set.Size())
	}
	
	expectedElements := []int{1, 2, 3, 4}
	for _, expected := range expectedElements {
		if !set.Contains(expected) {
			t.Errorf("Set should contain %d", expected)
		}
	}
}

func TestImmutableSetContains(t *testing.T) {
	set := SetOf(1, 2, 3, 4, 5)
	
	if !set.Contains(3) {
		t.Error("Set should contain 3")
	}
	
	if set.Contains(6) {
		t.Error("Set should not contain 6")
	}
}

func TestImmutableSetWithAdd(t *testing.T) {
	original := SetOf(1, 2, 3)
	newSet := original.WithAdd(4)
	
	// Original should be unchanged
	if original.Size() != 3 {
		t.Errorf("Original set size should remain 3, got %d", original.Size())
	}
	
	// New set should have the added element
	if newSet.Size() != 4 {
		t.Errorf("New set size should be 4, got %d", newSet.Size())
	}
	
	if !newSet.Contains(4) {
		t.Error("New set should contain 4")
	}
	
	// Test adding duplicate element
	sameSet := original.WithAdd(2)
	if sameSet != original {
		t.Error("Adding existing element should return same instance")
	}
}

func TestImmutableSetWithRemove(t *testing.T) {
	original := SetOf(1, 2, 3, 4)
	newSet := original.WithRemove(2)
	
	// Original should be unchanged
	if original.Size() != 4 {
		t.Errorf("Original set size should remain 4, got %d", original.Size())
	}
	
	if !original.Contains(2) {
		t.Error("Original set should still contain 2")
	}
	
	// New set should have the element removed
	if newSet.Size() != 3 {
		t.Errorf("New set size should be 3, got %d", newSet.Size())
	}
	
	if newSet.Contains(2) {
		t.Error("New set should not contain 2")
	}
	
	// Test removing non-existent element
	sameSet := original.WithRemove(5)
	if sameSet != original {
		t.Error("Removing non-existent element should return same instance")
	}
}

func TestImmutableSetWithClear(t *testing.T) {
	original := SetOf(1, 2, 3)
	newSet := original.WithClear()
	
	// Original should be unchanged
	if original.Size() != 3 {
		t.Errorf("Original set size should remain 3, got %d", original.Size())
	}
	
	// New set should be empty
	if !newSet.IsEmpty() {
		t.Error("New set should be empty")
	}
}

func TestImmutableSetUnion(t *testing.T) {
	set1 := SetOf(1, 2, 3)
	set2 := SetOf(3, 4, 5)
	
	union := set1.Union(set2)
	
	expectedSize := 5 // {1, 2, 3, 4, 5}
	if union.Size() != expectedSize {
		t.Errorf("Expected union size %d, got %d", expectedSize, union.Size())
	}
	
	expectedElements := []int{1, 2, 3, 4, 5}
	for _, element := range expectedElements {
		if !union.Contains(element) {
			t.Errorf("Union should contain %d", element)
		}
	}
}

func TestImmutableSetIntersection(t *testing.T) {
	set1 := SetOf(1, 2, 3, 4)
	set2 := SetOf(3, 4, 5, 6)
	
	intersection := set1.Intersection(set2)
	
	expectedSize := 2 // {3, 4}
	if intersection.Size() != expectedSize {
		t.Errorf("Expected intersection size %d, got %d", expectedSize, intersection.Size())
	}
	
	expectedElements := []int{3, 4}
	for _, element := range expectedElements {
		if !intersection.Contains(element) {
			t.Errorf("Intersection should contain %d", element)
		}
	}
	
	unexpectedElements := []int{1, 2, 5, 6}
	for _, element := range unexpectedElements {
		if intersection.Contains(element) {
			t.Errorf("Intersection should not contain %d", element)
		}
	}
}

func TestImmutableSetDifference(t *testing.T) {
	set1 := SetOf(1, 2, 3, 4)
	set2 := SetOf(3, 4, 5, 6)
	
	difference := set1.Difference(set2)
	
	expectedSize := 2 // {1, 2}
	if difference.Size() != expectedSize {
		t.Errorf("Expected difference size %d, got %d", expectedSize, difference.Size())
	}
	
	expectedElements := []int{1, 2}
	for _, element := range expectedElements {
		if !difference.Contains(element) {
			t.Errorf("Difference should contain %d", element)
		}
	}
	
	unexpectedElements := []int{3, 4, 5, 6}
	for _, element := range unexpectedElements {
		if difference.Contains(element) {
			t.Errorf("Difference should not contain %d", element)
		}
	}
}

func TestImmutableSetIsSubsetOf(t *testing.T) {
	set1 := SetOf(1, 2)
	set2 := SetOf(1, 2, 3, 4)
	set3 := SetOf(1, 5)
	
	if !set1.IsSubsetOf(set2) {
		t.Error("set1 should be a subset of set2")
	}
	
	if set1.IsSubsetOf(set3) {
		t.Error("set1 should not be a subset of set3")
	}
	
	if !set1.IsSubsetOf(set1) {
		t.Error("set should be a subset of itself")
	}
}

func TestImmutableSetIsSupersetOf(t *testing.T) {
	set1 := SetOf(1, 2, 3, 4)
	set2 := SetOf(1, 2)
	set3 := SetOf(1, 5)
	
	if !set1.IsSupersetOf(set2) {
		t.Error("set1 should be a superset of set2")
	}
	
	if set1.IsSupersetOf(set3) {
		t.Error("set1 should not be a superset of set3")
	}
	
	if !set1.IsSupersetOf(set1) {
		t.Error("set should be a superset of itself")
	}
}

func TestImmutableSetToSlice(t *testing.T) {
	set := SetOf(1, 2, 3)
	slice := set.ToSlice()
	
	if len(slice) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(slice))
	}
	
	// Check that all elements are present (order doesn't matter in sets)
	elementMap := make(map[int]bool)
	for _, element := range slice {
		elementMap[element] = true
	}
	
	expectedElements := []int{1, 2, 3}
	for _, expected := range expectedElements {
		if !elementMap[expected] {
			t.Errorf("Slice should contain %d", expected)
		}
	}
}

func TestImmutableSetIterator(t *testing.T) {
	set := SetOf(1, 2, 3)
	iterator := set.Iterator()
	
	elementMap := make(map[int]bool)
	count := 0
	
	for iterator.HasNext() {
		val, ok := iterator.Next()
		if !ok {
			t.Error("Iterator.Next() should return true")
		}
		
		elementMap[val] = true
		count++
	}
	
	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}
	
	expectedElements := []int{1, 2, 3}
	for _, expected := range expectedElements {
		if !elementMap[expected] {
			t.Errorf("Iterator should have returned %d", expected)
		}
	}
}

func TestImmutableSetForEach(t *testing.T) {
	set := SetOf(1, 2, 3)
	sum := 0
	
	set.ForEach(func(element int) {
		sum += element
	})
	
	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

func TestImmutableSetString(t *testing.T) {
	emptySet := NewImmutableSet[int]()
	if emptySet.String() != "{}" {
		t.Errorf("Expected '{}' for empty set, got '%s'", emptySet.String())
	}
	
	set := SetOf(1)
	str := set.String()
	expected := "{1}"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestImmutableSetImmutability(t *testing.T) {
	set := SetOf(1, 2, 3)
	
	// Test that modification methods don't change the original
	set.Add(4)
	set.Remove(2)
	set.Clear()
	
	// Original should remain unchanged
	if set.Size() != 3 {
		t.Errorf("Original set should remain unchanged, expected size 3, got %d", set.Size())
	}
	
	expectedElements := []int{1, 2, 3}
	for _, expected := range expectedElements {
		if !set.Contains(expected) {
			t.Errorf("Original set should remain unchanged, should contain %d", expected)
		}
	}
}