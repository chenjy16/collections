package set

import (
	"strings"
	"testing"
)

func TestHashSet_New(t *testing.T) {
	set := New[int]()
	if set == nil {
		t.Error("New should return a non-nil HashSet")
	}
	if !set.IsEmpty() {
		t.Error("New HashSet should be empty")
	}
	if set.Size() != 0 {
		t.Error("New HashSet size should be 0")
	}
}

func TestHashSet_FromSlice(t *testing.T) {
	// Test empty slice
	set := FromSlice([]int{})
	if !set.IsEmpty() {
		t.Error("HashSet from empty slice should be empty")
	}

	// Test non-empty slice
	slice := []int{1, 2, 3, 2, 1} // Contains duplicate elements
	set = FromSlice(slice)
	if set.Size() != 3 { // Should only have 3 unique elements
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	// Verify elements exist
	if !set.Contains(1) || !set.Contains(2) || !set.Contains(3) {
		t.Error("Set should contain all unique elements from slice")
	}
}

func TestHashSet_Add(t *testing.T) {
	set := New[int]()

	// Test adding elements
	if !set.Add(5) {
		t.Error("Add should return true for new element")
	}
	if !set.Add(3) {
		t.Error("Add should return true for new element")
	}
	if !set.Add(7) {
		t.Error("Add should return true for new element")
	}

	// Test duplicate addition
	if set.Add(5) {
		t.Error("Add should return false for duplicate element")
	}

	// Verify size
	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	// Verify elements exist
	if !set.Contains(3) || !set.Contains(5) || !set.Contains(7) {
		t.Error("Set should contain all added elements")
	}
}

func TestHashSet_Remove(t *testing.T) {
	set := New[int]()

	// Add elements
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	set.Add(5)

	// Test removing existing element
	if !set.Remove(3) {
		t.Error("Remove should return true for existing element")
	}

	// Test removing non-existing element
	if set.Remove(10) {
		t.Error("Remove should return false for non-existing element")
	}

	// Verify size
	if set.Size() != 4 {
		t.Errorf("Expected size 4, got %d", set.Size())
	}

	// Verify element doesn't exist
	if set.Contains(3) {
		t.Error("Set should not contain removed element")
	}

	// Verify other elements still exist
	if !set.Contains(1) || !set.Contains(2) || !set.Contains(4) || !set.Contains(5) {
		t.Error("Set should contain non-removed elements")
	}
}

func TestHashSet_Contains(t *testing.T) {
	set := New[int]()

	// Test empty set
	if set.Contains(1) {
		t.Error("Empty set should not contain any element")
	}

	// Add elements
	set.Add(10)
	set.Add(20)
	set.Add(30)

	// Test contained elements
	if !set.Contains(20) {
		t.Error("Set should contain added element")
	}

	// Test non-contained elements
	if set.Contains(25) {
		t.Error("Set should not contain non-added element")
	}
}

func TestHashSet_Clear(t *testing.T) {
	set := New[int]()

	// Add elements
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// Clear set
	set.Clear()

	// Verify set is empty
	if !set.IsEmpty() {
		t.Error("Set should be empty after Clear")
	}

	if set.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", set.Size())
	}

	// Verify elements don't exist
	if set.Contains(1) || set.Contains(2) || set.Contains(3) {
		t.Error("Set should not contain any elements after Clear")
	}
}

func TestHashSet_ToSlice(t *testing.T) {
	set := New[int]()

	// Test empty set
	slice := set.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice for empty set, got %v", slice)
	}

	// Add elements
	elements := []int{5, 2, 8, 1, 9, 3}
	for _, e := range elements {
		set.Add(e)
	}

	// Get slice
	slice = set.ToSlice()

	// Verify length
	if len(slice) != len(elements) {
		t.Errorf("Expected slice length %d, got %d", len(elements), len(slice))
	}

	// Verify all elements exist
	elementMap := make(map[int]bool)
	for _, e := range slice {
		elementMap[e] = true
	}

	for _, e := range elements {
		if !elementMap[e] {
			t.Errorf("Slice should contain element %d", e)
		}
	}
}

func TestHashSet_ForEach(t *testing.T) {
	set := New[int]()

	// Add elements
	elements := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, e := range elements {
		set.Add(e)
	}

	// Test ForEach
	var result []int
	set.ForEach(func(e int) {
		result = append(result, e)
	})

	// Verify no duplicate elements
	uniqueElements := make(map[int]bool)
	for _, e := range elements {
		uniqueElements[e] = true
	}

	if len(result) != len(uniqueElements) {
		t.Errorf("Expected %d unique elements, got %d", len(uniqueElements), len(result))
	}

	// Verify all elements were traversed
	resultMap := make(map[int]bool)
	for _, e := range result {
		resultMap[e] = true
	}

	for e := range uniqueElements {
		if !resultMap[e] {
			t.Errorf("ForEach should traverse element %d", e)
		}
	}
}

func TestHashSet_SetOperations(t *testing.T) {
	set1 := New[int]()
	set2 := New[int]()

	// Initialize sets
	for i := 1; i <= 5; i++ {
		set1.Add(i)
	}
	for i := 3; i <= 7; i++ {
		set2.Add(i)
	}

	// Test Union
	union := set1.Union(set2)
	expectedUnion := []int{1, 2, 3, 4, 5, 6, 7}
	if union.Size() != len(expectedUnion) {
		t.Errorf("Union size should be %d, got %d", len(expectedUnion), union.Size())
	}
	for _, val := range expectedUnion {
		if !union.Contains(val) {
			t.Errorf("Union should contain %d", val)
		}
	}

	// Test Intersection
	intersection := set1.Intersection(set2)
	expectedIntersection := []int{3, 4, 5}
	if intersection.Size() != len(expectedIntersection) {
		t.Errorf("Intersection size should be %d, got %d", len(expectedIntersection), intersection.Size())
	}
	for _, val := range expectedIntersection {
		if !intersection.Contains(val) {
			t.Errorf("Intersection should contain %d", val)
		}
	}

	// Test Difference
	difference := set1.Difference(set2)
	expectedDifference := []int{1, 2}
	if difference.Size() != len(expectedDifference) {
		t.Errorf("Difference size should be %d, got %d", len(expectedDifference), difference.Size())
	}
	for _, val := range expectedDifference {
		if !difference.Contains(val) {
			t.Errorf("Difference should contain %d", val)
		}
	}

	// Test subset relationship
	subset := New[int]()
	subset.Add(1)
	subset.Add(2)
	if !subset.IsSubsetOf(set1) {
		t.Error("subset should be a subset of set1")
	}

	if subset.IsSubsetOf(set2) {
		t.Error("subset should not be a subset of set2")
	}

	// Test superset relationship
	if !set1.IsSupersetOf(subset) {
		t.Error("set1 should be a superset of subset")
	}

	if set2.IsSupersetOf(subset) {
		t.Error("set2 should not be a superset of subset")
	}
}

func TestHashSet_String(t *testing.T) {
	set := New[int]()

	// Test empty set
	if set.String() != "[]" {
		t.Errorf("Expected \"[]\" for empty set, got %s", set.String())
	}

	// Add elements
	set.Add(3)
	set.Add(1)
	set.Add(2)

	// Test non-empty set
	// Since hash set order is non-deterministic, we need to check if string contains all elements
	str := set.String()
	if !strings.Contains(str, "1") || !strings.Contains(str, "2") || !strings.Contains(str, "3") {
		t.Errorf("String representation should contain all elements, got %s", str)
	}
	if !strings.HasPrefix(str, "[") || !strings.HasSuffix(str, "]") {
		t.Errorf("String representation should be wrapped in brackets, got %s", str)
	}
}

func TestHashSet_Iterator(t *testing.T) {
	set := New[int]()

	// Test empty set iterator
	it := set.Iterator()
	if it.HasNext() {
		t.Error("Empty set iterator should not have next element")
	}

	// Add elements
	elements := []int{1, 2, 3, 4, 5}
	for _, e := range elements {
		set.Add(e)
	}

	// Test iteration
	it = set.Iterator()
	visited := make(map[int]bool)
	count := 0

	for it.HasNext() {
		val, ok := it.Next()
		if !ok {
			t.Error("Iterator should have next element")
		}
		if visited[val] {
			t.Errorf("Element %d visited twice", val)
		}
		visited[val] = true
		count++
	}

	if count != len(elements) {
		t.Errorf("Expected to visit %d elements, visited %d", len(elements), count)
	}

	// Test iterator's Remove method
	set = New[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	it = set.Iterator()
	it.Next() // Move to first element
	if !it.Remove() {
		t.Error("Remove should succeed")
	}

	// Verify element was removed
	if set.Size() != 2 {
		t.Errorf("Expected size 2 after remove, got %d", set.Size())
	}

	// Test calling Remove before Next
	it = set.Iterator()
	if it.Remove() {
		t.Error("Remove before Next should fail")
	}
}
