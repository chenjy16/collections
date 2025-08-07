package list

import (
	"testing"
)

func TestNewImmutableList(t *testing.T) {
	list := NewImmutableList[int]()
	
	if !list.IsEmpty() {
		t.Error("New list should be empty")
	}
	
	if list.Size() != 0 {
		t.Errorf("Expected size 0, got %d", list.Size())
	}
}

func TestNewImmutableListFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	list := NewImmutableListFromSlice(slice)
	
	if list.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list.Size())
	}
	
	for i, expected := range slice {
		if actual, _ := list.Get(i); actual != expected {
			t.Errorf("Expected %d at index %d, got %d", expected, i, actual)
		}
	}
}

func TestListOf(t *testing.T) {
	list := Of(1, 2, 3, 4, 5)
	
	if list.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list.Size())
	}
	
	expected := []int{1, 2, 3, 4, 5}
	for i, exp := range expected {
		if actual, err := list.Get(i); err != nil || actual != exp {
			t.Errorf("Expected %d at index %d, got %d (error: %v)", exp, i, actual, err)
		}
	}
}

func TestImmutableListGet(t *testing.T) {
	list := Of(10, 20, 30)
	
	// Valid indices
	if val, err := list.Get(0); err != nil || val != 10 {
		t.Errorf("Expected (10, nil), got (%d, %v)", val, err)
	}
	
	if val, err := list.Get(1); err != nil || val != 20 {
		t.Errorf("Expected (20, nil), got (%d, %v)", val, err)
	}
	
	if val, err := list.Get(2); err != nil || val != 30 {
		t.Errorf("Expected (30, nil), got (%d, %v)", val, err)
	}
	
	// Invalid indices
	if val, err := list.Get(-1); err == nil {
		t.Errorf("Expected error for negative index, got (%d, %v)", val, err)
	}
	
	if val, err := list.Get(3); err == nil {
		t.Errorf("Expected error for out of bounds index, got (%d, %v)", val, err)
	}
}

func TestImmutableListContains(t *testing.T) {
	list := Of(1, 2, 3, 4, 5)
	
	if !list.Contains(3) {
		t.Error("List should contain 3")
	}
	
	if list.Contains(6) {
		t.Error("List should not contain 6")
	}
}

func TestImmutableListIndexOf(t *testing.T) {
	list := Of(1, 2, 3, 2, 4)
	
	if index := list.IndexOf(2); index != 1 {
		t.Errorf("Expected index 1 for first occurrence of 2, got %d", index)
	}
	
	if index := list.IndexOf(5); index != -1 {
		t.Errorf("Expected index -1 for non-existent element, got %d", index)
	}
}

func TestImmutableListLastIndexOf(t *testing.T) {
	list := Of(1, 2, 3, 2, 4)
	
	if index := list.LastIndexOf(2); index != 3 {
		t.Errorf("Expected index 3 for last occurrence of 2, got %d", index)
	}
	
	if index := list.LastIndexOf(5); index != -1 {
		t.Errorf("Expected index -1 for non-existent element, got %d", index)
	}
}

func TestImmutableListWithAdd(t *testing.T) {
	original := Of(1, 2, 3)
	newList := original.WithAdd(4)
	
	// Original should be unchanged
	if original.Size() != 3 {
		t.Errorf("Original list size should remain 3, got %d", original.Size())
	}
	
	// New list should have the added element
	if newList.Size() != 4 {
		t.Errorf("New list size should be 4, got %d", newList.Size())
	}
	
	if val, err := newList.Get(3); err != nil || val != 4 {
		t.Errorf("Expected 4 at index 3, got %d (error: %v)", val, err)
	}
}

func TestImmutableListWithInsert(t *testing.T) {
	original := Of(1, 3, 4)
	newList, err := original.WithInsert(1, 2)
	if err != nil {
		t.Fatalf("WithInsert failed: %v", err)
	}
	
	// Original should be unchanged
	if original.Size() != 3 {
		t.Errorf("Original list size should remain 3, got %d", original.Size())
	}
	
	// New list should have the inserted element
	if newList.Size() != 4 {
		t.Errorf("New list size should be 4, got %d", newList.Size())
	}
	
	expected := []int{1, 2, 3, 4}
	for i, exp := range expected {
		if val, err := newList.Get(i); err != nil || val != exp {
			t.Errorf("Expected %d at index %d, got %d (error: %v)", exp, i, val, err)
		}
	}
}

func TestImmutableListWithSet(t *testing.T) {
	original := Of(1, 2, 3)
	newList, err := original.WithSet(1, 10)
	if err != nil {
		t.Fatalf("WithSet failed: %v", err)
	}
	
	// Original should be unchanged
	if val, err := original.Get(1); err != nil || val != 2 {
		t.Errorf("Original list should remain unchanged, expected 2 at index 1, got %d (error: %v)", val, err)
	}
	
	// New list should have the updated element
	if val, err := newList.Get(1); err != nil || val != 10 {
		t.Errorf("Expected 10 at index 1, got %d (error: %v)", val, err)
	}
}

func TestImmutableListWithRemoveAt(t *testing.T) {
	original := Of(1, 2, 3, 4)
	newList, err := original.WithRemoveAt(1)
	if err != nil {
		t.Fatalf("WithRemoveAt failed: %v", err)
	}
	
	// Original should be unchanged
	if original.Size() != 4 {
		t.Errorf("Original list size should remain 4, got %d", original.Size())
	}
	
	// New list should have the element removed
	if newList.Size() != 3 {
		t.Errorf("New list size should be 3, got %d", newList.Size())
	}
	
	expected := []int{1, 3, 4}
	for i, exp := range expected {
		if val, err := newList.Get(i); err != nil || val != exp {
			t.Errorf("Expected %d at index %d, got %d (error: %v)", exp, i, val, err)
		}
	}
}

func TestImmutableListWithRemove(t *testing.T) {
	original := Of(1, 2, 3, 2, 4)
	newList := original.WithRemove(2)
	
	// Original should be unchanged
	if original.Size() != 5 {
		t.Errorf("Original list size should remain 5, got %d", original.Size())
	}
	
	// New list should have the first occurrence removed
	if newList.Size() != 4 {
		t.Errorf("New list size should be 4, got %d", newList.Size())
	}
	
	expected := []int{1, 3, 2, 4}
	for i, exp := range expected {
		if val, err := newList.Get(i); err != nil || val != exp {
			t.Errorf("Expected %d at index %d, got %d (error: %v)", exp, i, val, err)
		}
	}
}

func TestImmutableListWithClear(t *testing.T) {
	original := Of(1, 2, 3)
	newList := original.WithClear()
	
	// Original should be unchanged
	if original.Size() != 3 {
		t.Errorf("Original list size should remain 3, got %d", original.Size())
	}
	
	// New list should be empty
	if !newList.IsEmpty() {
		t.Error("New list should be empty")
	}
}

func TestImmutableListSubList(t *testing.T) {
	list := Of(1, 2, 3, 4, 5)
	subList, err := list.SubList(1, 4)
	if err != nil {
		t.Fatalf("SubList failed: %v", err)
	}
	
	if subList.Size() != 3 {
		t.Errorf("Expected sublist size 3, got %d", subList.Size())
	}
	
	expected := []int{2, 3, 4}
	for i, exp := range expected {
		if val, err := subList.Get(i); err != nil || val != exp {
			t.Errorf("Expected %d at index %d, got %d (error: %v)", exp, i, val, err)
		}
	}
}

func TestImmutableListToSlice(t *testing.T) {
	list := Of(1, 2, 3, 4, 5)
	slice := list.ToSlice()
	
	expected := []int{1, 2, 3, 4, 5}
	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}
	
	for i, exp := range expected {
		if slice[i] != exp {
			t.Errorf("Expected %d at index %d, got %d", exp, i, slice[i])
		}
	}
}

func TestImmutableListIterator(t *testing.T) {
	list := Of(1, 2, 3)
	iterator := list.Iterator()
	
	expected := []int{1, 2, 3}
	index := 0
	
	for iterator.HasNext() {
		val, ok := iterator.Next()
		if !ok {
			t.Error("Iterator.Next() should return true")
		}
		
		if val != expected[index] {
			t.Errorf("Expected %d, got %d", expected[index], val)
		}
		index++
	}
	
	if index != len(expected) {
		t.Errorf("Expected to iterate %d times, got %d", len(expected), index)
	}
}

func TestImmutableListForEach(t *testing.T) {
	list := Of(1, 2, 3)
	sum := 0
	
	list.ForEach(func(element int) {
		sum += element
	})
	
	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

func TestImmutableListString(t *testing.T) {
	emptyList := NewImmutableList[int]()
	if emptyList.String() != "[]" {
		t.Errorf("Expected '[]' for empty list, got '%s'", emptyList.String())
	}
	
	list := Of(1, 2, 3)
	str := list.String()
	expected := "[1, 2, 3]"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestImmutableListImmutability(t *testing.T) {
	list := Of(1, 2, 3)
	
	// Test that modification methods don't change the original
	list.Add(4)
	list.Insert(0, 0)
	list.Set(0, 10)
	list.RemoveAt(0)
	list.Remove(2)
	list.Clear()
	
	// Original should remain unchanged
	if list.Size() != 3 {
		t.Errorf("Original list should remain unchanged, expected size 3, got %d", list.Size())
	}
	
	expected := []int{1, 2, 3}
	for i, exp := range expected {
		if val, err := list.Get(i); err != nil || val != exp {
			t.Errorf("Original list should remain unchanged, expected %d at index %d, got %d (error: %v)", exp, i, val, err)
		}
	}
}