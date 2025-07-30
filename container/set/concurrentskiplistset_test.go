package set

import (
	"sync"
	"testing"
)

func TestConcurrentSkipListSet_Basic(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// Test empty set
	if !set.IsEmpty() {
		t.Error("New set should be empty")
	}

	if set.Size() != 0 {
		t.Error("New set size should be 0")
	}

	// Test add elements
	if !set.Add(1) {
		t.Error("Add(1) should return true")
	}

	if !set.Add(2) {
		t.Error("Add(2) should return true")
	}

	if !set.Add(3) {
		t.Error("Add(3) should return true")
	}

	if set.Size() != 3 {
		t.Errorf("Set size should be 3, got %d", set.Size())
	}

	if set.IsEmpty() {
		t.Error("Set should not be empty")
	}

	// Test add duplicate
	if set.Add(2) {
		t.Error("Add(2) should return false for duplicate")
	}

	if set.Size() != 3 {
		t.Errorf("Set size should still be 3, got %d", set.Size())
	}
}

func TestConcurrentSkipListSet_Contains(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	if !set.Contains(2) {
		t.Error("Set should contain 2")
	}

	if set.Contains(4) {
		t.Error("Set should not contain 4")
	}
}

func TestConcurrentSkipListSet_Remove(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// Test remove existing element
	if !set.Remove(2) {
		t.Error("Remove(2) should return true")
	}

	if set.Size() != 2 {
		t.Errorf("Set size should be 2, got %d", set.Size())
	}

	if set.Contains(2) {
		t.Error("Set should not contain 2 after removal")
	}

	// Test remove non-existing element
	if set.Remove(4) {
		t.Error("Remove(4) should return false")
	}

	if set.Size() != 2 {
		t.Errorf("Set size should still be 2, got %d", set.Size())
	}
}

func TestConcurrentSkipListSet_Clear(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	set.Clear()

	if !set.IsEmpty() {
		t.Error("Set should be empty after clear")
	}

	if set.Size() != 0 {
		t.Errorf("Set size should be 0 after clear, got %d", set.Size())
	}
}

func TestConcurrentSkipListSet_ToSlice(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	set.Add(3)
	set.Add(1)
	set.Add(2)

	slice := set.ToSlice()

	if len(slice) != 3 {
		t.Errorf("Slice length should be 3, got %d", len(slice))
	}

	// Elements should be sorted
	expected := []int{1, 2, 3}
	for i, val := range slice {
		if val != expected[i] {
			t.Errorf("Slice[%d] should be %d, got %d", i, expected[i], val)
		}
	}
}

func TestConcurrentSkipListSet_ForEach(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	sum := 0
	set.ForEach(func(val int) {
		sum += val
	})

	if sum != 6 {
		t.Errorf("Sum should be 6, got %d", sum)
	}
}

func TestConcurrentSkipListSet_SetOperations(t *testing.T) {
	set1 := NewConcurrentSkipListSet[int]()
	set2 := NewConcurrentSkipListSet[int]()

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
	unionSlice := union.ToSlice()
	if len(unionSlice) != len(expectedUnion) {
		t.Errorf("Union size should be %d, got %d", len(expectedUnion), len(unionSlice))
	}
	for i, val := range expectedUnion {
		if unionSlice[i] != val {
			t.Errorf("Union[%d] should be %d, got %d", i, val, unionSlice[i])
		}
	}

	// Test Intersection
	intersection := set1.Intersection(set2)
	expectedIntersection := []int{3, 4, 5}
	intersectionSlice := intersection.ToSlice()
	if len(intersectionSlice) != len(expectedIntersection) {
		t.Errorf("Intersection size should be %d, got %d", len(expectedIntersection), len(intersectionSlice))
	}
	for i, val := range expectedIntersection {
		if intersectionSlice[i] != val {
			t.Errorf("Intersection[%d] should be %d, got %d", i, val, intersectionSlice[i])
		}
	}

	// Test Difference
	difference := set1.Difference(set2)
	expectedDifference := []int{1, 2}
	differenceSlice := difference.ToSlice()
	if len(differenceSlice) != len(expectedDifference) {
		t.Errorf("Difference size should be %d, got %d", len(expectedDifference), len(differenceSlice))
	}
	for i, val := range expectedDifference {
		if differenceSlice[i] != val {
			t.Errorf("Difference[%d] should be %d, got %d", i, val, differenceSlice[i])
		}
	}

	// Test IsSubsetOf
	subset := NewConcurrentSkipListSet[int]()
	subset.Add(1)
	subset.Add(2)
	if !subset.IsSubsetOf(set1) {
		t.Error("subset should be subset of set1")
	}

	// Test IsSupersetOf
	if !set1.IsSupersetOf(subset) {
		t.Error("set1 should be superset of subset")
	}
}

func TestConcurrentSkipListSet_String(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// Test empty set
	if set.String() != "{}" {
		t.Errorf("Empty set string should be '{}', got '%s'", set.String())
	}

	// Test non-empty set
	set.Add(1)
	set.Add(2)
	set.Add(3)

	str := set.String()
	if str != "{1, 2, 3}" {
		t.Errorf("Set string should be '{1, 2, 3}', got '%s'", str)
	}
}

func TestConcurrentSkipListSet_WithComparator(t *testing.T) {
	// Custom comparator for reverse order
	reverseComparator := func(a, b int) int {
		if a < b {
			return 1
		} else if a > b {
			return -1
		}
		return 0
	}

	set := NewConcurrentSkipListSetWithComparator(reverseComparator)
	set.Add(3)
	set.Add(1)
	set.Add(2)

	slice := set.ToSlice()
	expected := []int{3, 2, 1} // Reverse order

	if len(slice) != 3 {
		t.Errorf("Slice length should be 3, got %d", len(slice))
	}

	for i, val := range slice {
		if val != expected[i] {
			t.Errorf("Slice[%d] should be %d, got %d", i, expected[i], val)
		}
	}
}

func TestConcurrentSkipListSet_Concurrency(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup

	// Concurrent add operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				set.Add(start*numOperations + j)
			}
		}(i)
	}

	// Concurrent read operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				set.Contains(j)
				set.Size()
			}
		}()
	}

	wg.Wait()

	// Verify final state
	expectedSize := numGoroutines * numOperations
	if set.Size() != expectedSize {
		t.Errorf("Expected size %d after concurrent operations, got %d", expectedSize, set.Size())
	}
}
