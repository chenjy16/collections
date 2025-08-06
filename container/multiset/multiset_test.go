package multiset

import (
	"testing"
)

// Test HashMultiset
func TestHashMultiset(t *testing.T) {
	testMultisetImplementation(t, func() Multiset[string] {
		return NewHashMultiset[string]()
	})
}

func TestHashMultisetFromSlice(t *testing.T) {
	elements := []string{"a", "b", "a", "c", "b", "a"}
	ms := NewHashMultisetFromSlice(elements)
	
	if ms.Count("a") != 3 {
		t.Errorf("Expected count of 'a' to be 3, got %d", ms.Count("a"))
	}
	if ms.Count("b") != 2 {
		t.Errorf("Expected count of 'b' to be 2, got %d", ms.Count("b"))
	}
	if ms.Count("c") != 1 {
		t.Errorf("Expected count of 'c' to be 1, got %d", ms.Count("c"))
	}
	if ms.TotalSize() != 6 {
		t.Errorf("Expected total size to be 6, got %d", ms.TotalSize())
	}
}

// Test TreeMultiset
func TestTreeMultiset(t *testing.T) {
	testMultisetImplementation(t, func() Multiset[string] {
		return NewTreeMultiset[string]()
	})
}

func TestTreeMultisetOrdering(t *testing.T) {
	ms := NewTreeMultiset[string]()
	elements := []string{"c", "a", "b", "a"}
	
	for _, elem := range elements {
		ms.Add(elem)
	}
	
	elementSet := ms.ElementSet()
	expected := []string{"a", "b", "c"}
	
	if len(elementSet) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(elementSet))
	}
	
	for i, elem := range elementSet {
		if elem != expected[i] {
			t.Errorf("Expected element at index %d to be %s, got %s", i, expected[i], elem)
		}
	}
}

// Test LinkedHashMultiset
func TestLinkedHashMultiset(t *testing.T) {
	testMultisetImplementation(t, func() Multiset[string] {
		return NewLinkedHashMultiset[string]()
	})
}

func TestLinkedHashMultisetOrdering(t *testing.T) {
	ms := NewLinkedHashMultiset[string]()
	elements := []string{"c", "a", "b", "a"}
	
	for _, elem := range elements {
		ms.Add(elem)
	}
	
	elementSet := ms.ElementSet()
	expected := []string{"c", "a", "b"} // insertion order
	
	if len(elementSet) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(elementSet))
	}
	
	for i, elem := range elementSet {
		if elem != expected[i] {
			t.Errorf("Expected element at index %d to be %s, got %s", i, expected[i], elem)
		}
	}
}

// Test ConcurrentHashMultiset
func TestConcurrentHashMultiset(t *testing.T) {
	testMultisetImplementation(t, func() Multiset[string] {
		return NewConcurrentHashMultiset[string]()
	})
}

func TestConcurrentHashMultisetConcurrency(t *testing.T) {
	ms := NewConcurrentHashMultiset[int]()
	
	// Test concurrent access
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				ms.Add(id)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Check results
	if ms.TotalSize() != 1000 {
		t.Errorf("Expected total size to be 1000, got %d", ms.TotalSize())
	}
	
	for i := 0; i < 10; i++ {
		if ms.Count(i) != 100 {
			t.Errorf("Expected count of %d to be 100, got %d", i, ms.Count(i))
		}
	}
}

// Test ImmutableMultiset
func TestImmutableMultiset(t *testing.T) {
	// Test basic operations
	ms := NewImmutableMultiset[string]()
	
	if !ms.IsEmpty() {
		t.Error("New multiset should be empty")
	}
	
	ms2 := ms.WithAdd("a")
	if ms.Contains("a") {
		t.Error("Original multiset should not be modified")
	}
	if !ms2.Contains("a") {
		t.Error("New multiset should contain added element")
	}
	
	ms3, _ := ms2.WithAddCount("a", 2)
	if ms2.Count("a") != 1 {
		t.Error("Previous multiset should not be modified")
	}
	if ms3.Count("a") != 3 {
		t.Errorf("Expected count to be 3, got %d", ms3.Count("a"))
	}
}

func TestImmutableMultisetFromSlice(t *testing.T) {
	elements := []string{"a", "b", "a", "c"}
	ms := NewImmutableMultisetFromSlice(elements)
	
	if ms.Count("a") != 2 {
		t.Errorf("Expected count of 'a' to be 2, got %d", ms.Count("a"))
	}
	if ms.TotalSize() != 4 {
		t.Errorf("Expected total size to be 4, got %d", ms.TotalSize())
	}
}

// Common test function for all multiset implementations
func testMultisetImplementation(t *testing.T, factory func() Multiset[string]) {
	ms := factory()
	
	// Test empty multiset
	if !ms.IsEmpty() {
		t.Error("New multiset should be empty")
	}
	if ms.Size() != 0 {
		t.Errorf("Expected size 0, got %d", ms.Size())
	}
	if ms.TotalSize() != 0 {
		t.Errorf("Expected total size 0, got %d", ms.TotalSize())
	}
	
	// Test adding elements
	prevCount := ms.Add("a")
	if prevCount != 0 {
		t.Errorf("Expected previous count 0, got %d", prevCount)
	}
	if ms.Count("a") != 1 {
		t.Errorf("Expected count 1, got %d", ms.Count("a"))
	}
	if ms.IsEmpty() {
		t.Error("Multiset should not be empty after adding element")
	}
	
	// Test adding same element again
	prevCount = ms.Add("a")
	if prevCount != 1 {
		t.Errorf("Expected previous count 1, got %d", prevCount)
	}
	if ms.Count("a") != 2 {
		t.Errorf("Expected count 2, got %d", ms.Count("a"))
	}
	
	// Test AddCount
	prevCount, _ = ms.AddCount("b", 3)
	if prevCount != 0 {
		t.Errorf("Expected previous count 0, got %d", prevCount)
	}
	if ms.Count("b") != 3 {
		t.Errorf("Expected count 3, got %d", ms.Count("b"))
	}
	
	// Test Contains
	if !ms.Contains("a") {
		t.Error("Multiset should contain 'a'")
	}
	if !ms.Contains("b") {
		t.Error("Multiset should contain 'b'")
	}
	if ms.Contains("c") {
		t.Error("Multiset should not contain 'c'")
	}
	
	// Test Size and TotalSize
	if ms.Size() != 2 {
		t.Errorf("Expected size 2, got %d", ms.Size())
	}
	if ms.TotalSize() != 5 {
		t.Errorf("Expected total size 5, got %d", ms.TotalSize())
	}
	
	// Test Remove
	prevCount = ms.Remove("a")
	if prevCount != 2 {
		t.Errorf("Expected previous count 2, got %d", prevCount)
	}
	if ms.Count("a") != 1 {
		t.Errorf("Expected count 1, got %d", ms.Count("a"))
	}
	
	// Test RemoveCount
	prevCount, _ = ms.RemoveCount("b", 2)
	if prevCount != 3 {
		t.Errorf("Expected previous count 3, got %d", prevCount)
	}
	if ms.Count("b") != 1 {
		t.Errorf("Expected count 1, got %d", ms.Count("b"))
	}
	
	// Test RemoveAll
	prevCount = ms.RemoveAll("a")
	if prevCount != 1 {
		t.Errorf("Expected previous count 1, got %d", prevCount)
	}
	if ms.Count("a") != 0 {
		t.Errorf("Expected count 0, got %d", ms.Count("a"))
	}
	if ms.Contains("a") {
		t.Error("Multiset should not contain 'a' after RemoveAll")
	}
	
	// Test SetCount
	prevCount, _ = ms.SetCount("c", 5)
	if prevCount != 0 {
		t.Errorf("Expected previous count 0, got %d", prevCount)
	}
	if ms.Count("c") != 5 {
		t.Errorf("Expected count 5, got %d", ms.Count("c"))
	}
	
	// Test SetCount to 0 (should remove element)
	_, _ = ms.SetCount("c", 0)
	if ms.Contains("c") {
		t.Error("Multiset should not contain 'c' after setting count to 0")
	}
	
	// Test ElementSet and EntrySet
	ms.AddCount("x", 2)
	ms.AddCount("y", 1)
	
	elementSet := ms.ElementSet()
	if len(elementSet) != 3 { // b, x, y
		t.Errorf("Expected 3 distinct elements, got %d", len(elementSet))
	}
	
	entrySet := ms.EntrySet()
	if len(entrySet) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entrySet))
	}
	
	// Test ToSlice
	slice := ms.ToSlice()
	if len(slice) != 4 { // b(1) + x(2) + y(1)
		t.Errorf("Expected slice length 4, got %d", len(slice))
	}
	
	// Test Iterator
	iterator := ms.Iterator()
	count := 0
	for iterator.HasNext() {
		_, ok := iterator.Next()
		if !ok {
			t.Error("Next() should return true when HasNext() is true")
		}
		count++
	}
	if count != ms.TotalSize() {
		t.Errorf("Iterator should iterate over all elements, expected %d, got %d", ms.TotalSize(), count)
	}
	
	// Test ForEach
	forEachCount := 0
	ms.ForEach(func(element string) {
		forEachCount++
	})
	if forEachCount != ms.TotalSize() {
		t.Errorf("ForEach should visit all elements, expected %d, got %d", ms.TotalSize(), forEachCount)
	}
	
	// Test Clear
	ms.Clear()
	if !ms.IsEmpty() {
		t.Error("Multiset should be empty after Clear")
	}
}

// Test multiset operations
func TestMultisetOperations(t *testing.T) {
	ms1 := NewHashMultiset[string]()
	ms1.AddCount("a", 2)
	ms1.AddCount("b", 1)
	
	ms2 := NewHashMultiset[string]()
	ms2.AddCount("a", 1)
	ms2.AddCount("c", 3)
	
	// Test Union
	union := ms1.Union(ms2)
	if union.Count("a") != 2 { // max(2, 1)
		t.Errorf("Expected union count of 'a' to be 2, got %d", union.Count("a"))
	}
	if union.Count("b") != 1 {
		t.Errorf("Expected union count of 'b' to be 1, got %d", union.Count("b"))
	}
	if union.Count("c") != 3 {
		t.Errorf("Expected union count of 'c' to be 3, got %d", union.Count("c"))
	}
	
	// Test Intersection
	intersection := ms1.Intersection(ms2)
	if intersection.Count("a") != 1 { // min(2, 1)
		t.Errorf("Expected intersection count of 'a' to be 1, got %d", intersection.Count("a"))
	}
	if intersection.Count("b") != 0 { // min(1, 0)
		t.Errorf("Expected intersection count of 'b' to be 0, got %d", intersection.Count("b"))
	}
	if intersection.Count("c") != 0 { // min(0, 3)
		t.Errorf("Expected intersection count of 'c' to be 0, got %d", intersection.Count("c"))
	}
	
	// Test Difference
	difference := ms1.Difference(ms2)
	if difference.Count("a") != 1 { // 2 - 1
		t.Errorf("Expected difference count of 'a' to be 1, got %d", difference.Count("a"))
	}
	if difference.Count("b") != 1 { // 1 - 0
		t.Errorf("Expected difference count of 'b' to be 1, got %d", difference.Count("b"))
	}
	if difference.Count("c") != 0 { // 0 - 3 = 0 (can't be negative)
		t.Errorf("Expected difference count of 'c' to be 0, got %d", difference.Count("c"))
	}
	
	// Test IsSubsetOf
	ms3 := NewHashMultiset[string]()
	ms3.AddCount("a", 1)
	
	if !ms3.IsSubsetOf(ms1) {
		t.Error("ms3 should be a subset of ms1")
	}
	if ms1.IsSubsetOf(ms3) {
		t.Error("ms1 should not be a subset of ms3")
	}
	
	// Test IsSupersetOf
	if !ms1.IsSupersetOf(ms3) {
		t.Error("ms1 should be a superset of ms3")
	}
	if ms3.IsSupersetOf(ms1) {
		t.Error("ms3 should not be a superset of ms1")
	}
}

// Test error conditions
func TestMultisetErrorConditions(t *testing.T) {
	ms := NewHashMultiset[string]()
	
	// Test negative count returns error instead of panicking
	count, err := ms.AddCount("a", -1)
	if err == nil {
		t.Error("AddCount with negative count should return error")
	}
	if count != 0 {
		t.Errorf("AddCount with negative count should return 0, got %d", count)
	}
	
	count, err = ms.RemoveCount("a", -1)
	if err == nil {
		t.Error("RemoveCount with negative count should return error")
	}
	if count != 0 {
		t.Errorf("RemoveCount with negative count should return 0, got %d", count)
	}
	
	count, err = ms.SetCount("a", -1)
	if err == nil {
		t.Error("SetCount with negative count should return error")
	}
	if count != 0 {
		t.Errorf("SetCount with negative count should return 0, got %d", count)
	}
}