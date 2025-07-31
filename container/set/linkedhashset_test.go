package set

import (
	"testing"
)

func TestLinkedHashSet_New(t *testing.T) {
	set := NewLinkedHashSet[int]()
	if set == nil {
		t.Error("NewLinkedHashSet should return a non-nil LinkedHashSet")
	}
	if !set.IsEmpty() {
		t.Error("New LinkedHashSet should be empty")
	}
	if set.Size() != 0 {
		t.Error("New LinkedHashSet should have size 0")
	}
}

func TestLinkedHashSet_FromSlice(t *testing.T) {
	slice := []int{3, 1, 4, 1, 5, 9, 2, 6}
	set := LinkedHashSetFromSlice(slice)

	// Should maintain insertion order and remove duplicates
	expected := []int{3, 1, 4, 5, 9, 2, 6}
	if set.Size() != len(expected) {
		t.Errorf("Size should be %d, got %d", len(expected), set.Size())
	}

	// Check order is maintained
	result := set.ToSlice()
	for i, val := range expected {
		if result[i] != val {
			t.Errorf("Element at index %d should be %d, got %d", i, val, result[i])
		}
	}
}

func TestLinkedHashSet_Add(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Test adding elements
	if !set.Add(5) {
		t.Error("Add should return true for new element")
	}
	if !set.Add(3) {
		t.Error("Add should return true for new element")
	}
	if !set.Add(8) {
		t.Error("Add should return true for new element")
	}

	// Test adding duplicate
	if set.Add(5) {
		t.Error("Add should return false for duplicate element")
	}

	if set.Size() != 3 {
		t.Errorf("Size should be 3, got %d", set.Size())
	}

	// Check insertion order is maintained
	result := set.ToSlice()
	expected := []int{5, 3, 8}
	for i, val := range expected {
		if result[i] != val {
			t.Errorf("Element at index %d should be %d, got %d", i, val, result[i])
		}
	}
}

func TestLinkedHashSet_Remove(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Add elements
	elements := []int{1, 2, 3, 4, 5}
	for _, e := range elements {
		set.Add(e)
	}

	// Test removing existing element
	if !set.Remove(3) {
		t.Error("Remove should return true for existing element")
	}
	if set.Size() != 4 {
		t.Errorf("Size should be 4 after removal, got %d", set.Size())
	}
	if set.Contains(3) {
		t.Error("Set should not contain removed element")
	}

	// Check order is maintained after removal
	result := set.ToSlice()
	expected := []int{1, 2, 4, 5}
	for i, val := range expected {
		if result[i] != val {
			t.Errorf("Element at index %d should be %d, got %d", i, val, result[i])
		}
	}

	// Test removing non-existing element
	if set.Remove(10) {
		t.Error("Remove should return false for non-existing element")
	}
}

func TestLinkedHashSet_Contains(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Test empty set
	if set.Contains(1) {
		t.Error("Empty set should not contain any element")
	}

	// Add elements
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// Test existing elements
	if !set.Contains(1) {
		t.Error("Set should contain 1")
	}
	if !set.Contains(2) {
		t.Error("Set should contain 2")
	}
	if !set.Contains(3) {
		t.Error("Set should contain 3")
	}

	// Test non-existing element
	if set.Contains(4) {
		t.Error("Set should not contain 4")
	}
}

func TestLinkedHashSet_Clear(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Add elements
	for i := 1; i <= 5; i++ {
		set.Add(i)
	}

	// Clear the set
	set.Clear()

	if !set.IsEmpty() {
		t.Error("Set should be empty after clear")
	}
	if set.Size() != 0 {
		t.Error("Set size should be 0 after clear")
	}
}

func TestLinkedHashSet_ToSlice(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Test empty set
	slice := set.ToSlice()
	if len(slice) != 0 {
		t.Error("Empty set should return empty slice")
	}

	// Add elements in specific order
	elements := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, e := range elements {
		set.Add(e)
	}

	slice = set.ToSlice()
	expected := []int{3, 1, 4, 5, 9, 2, 6} // Duplicates removed, order maintained

	if len(slice) != len(expected) {
		t.Errorf("Slice length should be %d, got %d", len(expected), len(slice))
	}

	for i, val := range expected {
		if slice[i] != val {
			t.Errorf("Element at index %d should be %d, got %d", i, val, slice[i])
		}
	}
}

func TestLinkedHashSet_ForEach(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Add elements
	elements := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, e := range elements {
		set.Add(e)
	}

	// Test ForEach maintains insertion order
	var result []int
	set.ForEach(func(element int) {
		result = append(result, element)
	})

	expected := []int{3, 1, 4, 5, 9, 2, 6}
	if len(result) != len(expected) {
		t.Errorf("ForEach result length should be %d, got %d", len(expected), len(result))
	}

	for i, val := range expected {
		if result[i] != val {
			t.Errorf("Element at index %d should be %d, got %d", i, val, result[i])
		}
	}
}

func TestLinkedHashSet_SetOperations(t *testing.T) {
	set1 := NewLinkedHashSet[int]()
	set2 := NewLinkedHashSet[int]()

	// Initialize sets
	for i := 1; i <= 5; i++ {
		set1.Add(i)
	}
	for i := 3; i <= 7; i++ {
		set2.Add(i)
	}

	// Test Union - should maintain order from first set, then second set
	union := set1.Union(set2)
	expectedUnion := []int{1, 2, 3, 4, 5, 6, 7}
	if union.Size() != len(expectedUnion) {
		t.Errorf("Union size should be %d, got %d", len(expectedUnion), union.Size())
	}
	unionSlice := union.ToSlice()
	for i, val := range expectedUnion {
		if unionSlice[i] != val {
			t.Errorf("Union element at index %d should be %d, got %d", i, val, unionSlice[i])
		}
	}

	// Test Intersection - should maintain order from first set
	intersection := set1.Intersection(set2)
	expectedIntersection := []int{3, 4, 5}
	if intersection.Size() != len(expectedIntersection) {
		t.Errorf("Intersection size should be %d, got %d", len(expectedIntersection), intersection.Size())
	}
	intersectionSlice := intersection.ToSlice()
	for i, val := range expectedIntersection {
		if intersectionSlice[i] != val {
			t.Errorf("Intersection element at index %d should be %d, got %d", i, val, intersectionSlice[i])
		}
	}

	// Test Difference
	difference := set1.Difference(set2)
	expectedDifference := []int{1, 2}
	if difference.Size() != len(expectedDifference) {
		t.Errorf("Difference size should be %d, got %d", len(expectedDifference), difference.Size())
	}
	differenceSlice := difference.ToSlice()
	for i, val := range expectedDifference {
		if differenceSlice[i] != val {
			t.Errorf("Difference element at index %d should be %d, got %d", i, val, differenceSlice[i])
		}
	}

	// Test IsSubsetOf
	subset := NewLinkedHashSet[int]()
	subset.Add(2)
	subset.Add(3)
	if !subset.IsSubsetOf(set1) {
		t.Error("subset should be a subset of set1")
	}
	if subset.IsSubsetOf(set2) {
		t.Error("subset should not be a subset of set2")
	}

	// Test IsSupersetOf
	if !set1.IsSupersetOf(subset) {
		t.Error("set1 should be a superset of subset")
	}
}

func TestLinkedHashSet_String(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Test empty set
	if set.String() != "[]" {
		t.Errorf("Empty set string should be '[]', got '%s'", set.String())
	}

	// Add elements
	set.Add(1)
	set.Add(2)
	set.Add(3)

	str := set.String()
	expected := "[1, 2, 3]"
	if str != expected {
		t.Errorf("Set string should be '%s', got '%s'", expected, str)
	}
}

func TestLinkedHashSet_Iterator(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Test empty set iterator
	it := set.Iterator()
	if it.HasNext() {
		t.Error("Empty set iterator should not have next")
	}
	if _, ok := it.Next(); ok {
		t.Error("Empty set iterator Next should return false")
	}

	// Add elements
	elements := []int{3, 1, 4, 5, 9}
	for _, e := range elements {
		set.Add(e)
	}

	// Test iterator maintains insertion order
	it = set.Iterator()
	var result []int
	for it.HasNext() {
		if val, ok := it.Next(); ok {
			result = append(result, val)
		}
	}

	if len(result) != len(elements) {
		t.Errorf("Iterator result length should be %d, got %d", len(elements), len(result))
	}

	for i, val := range elements {
		if result[i] != val {
			t.Errorf("Iterator element at index %d should be %d, got %d", i, val, result[i])
		}
	}

	// Test iterator Remove
	it = set.Iterator()
	if it.HasNext() {
		val, _ := it.Next()
		if !it.Remove() {
			t.Error("Iterator Remove should return true after Next")
		}
		if set.Contains(val) {
			t.Error("Set should not contain removed element")
		}
		if set.Size() != len(elements)-1 {
			t.Errorf("Set size should be %d after iterator remove, got %d", len(elements)-1, set.Size())
		}
	}

	// Test Remove without Next
	it = set.Iterator()
	if it.Remove() {
		t.Error("Iterator Remove should return false without Next")
	}
}

func TestLinkedHashSet_InsertionOrder(t *testing.T) {
	set := NewLinkedHashSet[string]()

	// Add elements in specific order
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, word := range words {
		set.Add(word)
	}

	// Verify order is maintained
	result := set.ToSlice()
	for i, word := range words {
		if result[i] != word {
			t.Errorf("Element at index %d should be %s, got %s", i, word, result[i])
		}
	}

	// Remove middle element and verify order is still maintained
	set.Remove("cherry")
	expected := []string{"apple", "banana", "date", "elderberry"}
	result = set.ToSlice()
	for i, word := range expected {
		if result[i] != word {
			t.Errorf("After removal, element at index %d should be %s, got %s", i, word, result[i])
		}
	}

	// Add new element and verify it goes to the end
	set.Add("fig")
	expected = []string{"apple", "banana", "date", "elderberry", "fig"}
	result = set.ToSlice()
	for i, word := range expected {
		if result[i] != word {
			t.Errorf("After addition, element at index %d should be %s, got %s", i, word, result[i])
		}
	}
}

func TestLinkedHashSet_LargeDataset(t *testing.T) {
	set := NewLinkedHashSet[int]()

	// Add many elements to test resize functionality
	n := 1000
	for i := 0; i < n; i++ {
		if !set.Add(i) {
			t.Errorf("Add should return true for element %d", i)
		}
	}

	if set.Size() != n {
		t.Errorf("Set size should be %d, got %d", n, set.Size())
	}

	// Verify all elements are present and in order
	result := set.ToSlice()
	for i := 0; i < n; i++ {
		if result[i] != i {
			t.Errorf("Element at index %d should be %d, got %d", i, i, result[i])
		}
		if !set.Contains(i) {
			t.Errorf("Set should contain element %d", i)
		}
	}

	// Remove half the elements
	for i := 0; i < n; i += 2 {
		if !set.Remove(i) {
			t.Errorf("Remove should return true for element %d", i)
		}
	}

	if set.Size() != n/2 {
		t.Errorf("Set size should be %d after removals, got %d", n/2, set.Size())
	}

	// Verify remaining elements are still in order
	result = set.ToSlice()
	expectedIndex := 0
	for i := 1; i < n; i += 2 {
		if result[expectedIndex] != i {
			t.Errorf("Element at index %d should be %d, got %d", expectedIndex, i, result[expectedIndex])
		}
		expectedIndex++
	}
}