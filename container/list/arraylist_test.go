package list

import (
	"testing"
)

func TestArrayList_Basic(t *testing.T) {
	list := New[int]()

	// Test empty list
	if !list.IsEmpty() {
		t.Error("New list should be empty")
	}

	if list.Size() != 0 {
		t.Error("New list size should be 0")
	}

	// Test add elements
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.Size() != 3 {
		t.Errorf("List size should be 3, got %d", list.Size())
	}

	if list.IsEmpty() {
		t.Error("List should not be empty")
	}
}

func TestArrayList_WithCapacity(t *testing.T) {
	list := WithCapacity[int](2)

	// Test add within initial capacity
	list.Add(1)
	list.Add(2)

	if list.Size() != 2 {
		t.Errorf("List size should be 2, got %d", list.Size())
	}

	// Test add beyond initial capacity (should auto-expand)
	list.Add(3)

	if list.Size() != 3 {
		t.Errorf("List size should be 3, got %d", list.Size())
	}
}

func TestArrayList_FromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	list := FromSlice(slice)

	if list.Size() != len(slice) {
		t.Errorf("List size should be %d, got %d", len(slice), list.Size())
	}

	// Verify elements
	for i, expected := range slice {
		val, err := list.Get(i)
		if err != nil {
			t.Errorf("Get(%d) should not return error, got %v", i, err)
		}
		if val != expected {
			t.Errorf("Get(%d) should return %d, got %d", i, expected, val)
		}
	}
}

func TestArrayList_Get(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	// Test valid indices
	for i := 0; i < 3; i++ {
		val, err := list.Get(i)
		if err != nil {
			t.Errorf("Get(%d) should not return error, got %v", i, err)
		}
		if val != i+1 {
			t.Errorf("Get(%d) should return %d, got %d", i, i+1, val)
		}
	}

	// Test invalid indices
	_, err := list.Get(-1)
	if err == nil {
		t.Error("Get(-1) should return error")
	}

	_, err = list.Get(3)
	if err == nil {
		t.Error("Get(3) should return error")
	}
}

func TestArrayList_Set(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	// Test valid set
	_, success := list.Set(1, 10)
	if !success {
		t.Error("Set(1, 10) should return true")
	}

	val, _ := list.Get(1)
	if val != 10 {
		t.Errorf("Get(1) should return 10, got %d", val)
	}

	// Test invalid indices
	_, success = list.Set(-1, 100)
	if success {
		t.Error("Set(-1, 100) should return false")
	}

	_, success = list.Set(3, 100)
	if success {
		t.Error("Set(3, 100) should return false")
	}
}

func TestArrayList_Insert(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(3)

	// Test insert at middle
	err := list.Insert(1, 2)
	if err != nil {
		t.Errorf("Insert(1, 2) should not return error, got %v", err)
	}

	if list.Size() != 3 {
		t.Errorf("List size should be 3, got %d", list.Size())
	}

	// Verify order
	expected := []int{1, 2, 3}
	for i, exp := range expected {
		val, _ := list.Get(i)
		if val != exp {
			t.Errorf("Get(%d) should return %d, got %d", i, exp, val)
		}
	}

	// Test insert at beginning
	err = list.Insert(0, 0)
	if err != nil {
		t.Errorf("Insert(0, 0) should not return error, got %v", err)
	}

	// Test insert at end
	err = list.Insert(list.Size(), 4)
	if err != nil {
		t.Errorf("Insert at end should not return error, got %v", err)
	}

	// Test invalid index
	err = list.Insert(-1, 100)
	if err == nil {
		t.Error("Insert(-1, 100) should return error")
	}

	err = list.Insert(list.Size()+1, 100)
	if err == nil {
		t.Error("Insert beyond size should return error")
	}
}

func TestArrayList_RemoveAt(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(4)

	// Test remove at middle
	val, success := list.RemoveAt(1)
	if !success {
		t.Error("RemoveAt(1) should return true")
	}

	if val != 2 {
		t.Errorf("RemoveAt(1) should return 2, got %d", val)
	}

	if list.Size() != 3 {
		t.Errorf("List size should be 3, got %d", list.Size())
	}

	// Verify remaining elements
	expected := []int{1, 3, 4}
	for i, exp := range expected {
		val, _ := list.Get(i)
		if val != exp {
			t.Errorf("Get(%d) should return %d, got %d", i, exp, val)
		}
	}

	// Test invalid indices
	_, success = list.RemoveAt(-1)
	if success {
		t.Error("RemoveAt(-1) should return false")
	}

	_, success = list.RemoveAt(list.Size())
	if success {
		t.Error("RemoveAt beyond size should return false")
	}
}

func TestArrayList_Remove(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(2)

	// Test remove existing element
	removed := list.Remove(2)
	if !removed {
		t.Error("Remove(2) should return true")
	}

	if list.Size() != 3 {
		t.Errorf("List size should be 3, got %d", list.Size())
	}

	// Verify first occurrence was removed
	expected := []int{1, 3, 2}
	for i, exp := range expected {
		val, _ := list.Get(i)
		if val != exp {
			t.Errorf("Get(%d) should return %d, got %d", i, exp, val)
		}
	}

	// Test remove non-existing element
	removed = list.Remove(10)
	if removed {
		t.Error("Remove(10) should return false")
	}
}

func TestArrayList_IndexOf(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(2)

	// Test existing element
	index := list.IndexOf(2)
	if index != 1 {
		t.Errorf("IndexOf(2) should return 1, got %d", index)
	}

	// Test non-existing element
	index = list.IndexOf(10)
	if index != -1 {
		t.Errorf("IndexOf(10) should return -1, got %d", index)
	}
}

func TestArrayList_LastIndexOf(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(2)

	// Test existing element
	index := list.LastIndexOf(2)
	if index != 3 {
		t.Errorf("LastIndexOf(2) should return 3, got %d", index)
	}

	// Test non-existing element
	index = list.LastIndexOf(10)
	if index != -1 {
		t.Errorf("LastIndexOf(10) should return -1, got %d", index)
	}
}

func TestArrayList_SubList(t *testing.T) {
	list := New[int]()
	for i := 0; i < 5; i++ {
		list.Add(i)
	}

	// Test valid sublist
	sublist, err := list.SubList(1, 4)
	if err != nil {
		t.Errorf("SubList should not return error, got %v", err)
	}

	if sublist.Size() != 3 {
		t.Errorf("Sublist size should be 3, got %d", sublist.Size())
	}

	// Verify sublist elements
	expected := []int{1, 2, 3}
	for i, exp := range expected {
		val, _ := sublist.Get(i)
		if val != exp {
			t.Errorf("Sublist Get(%d) should return %d, got %d", i, exp, val)
		}
	}

	// Test invalid indices
	_, err = list.SubList(-1, 3)
	if err == nil {
		t.Error("SubList with invalid indices should return error")
	}
}

func TestArrayList_Clear(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	list.Clear()

	if !list.IsEmpty() {
		t.Error("List should be empty after clear")
	}

	if list.Size() != 0 {
		t.Errorf("List size should be 0 after clear, got %d", list.Size())
	}
}

func TestArrayList_Contains(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if !list.Contains(2) {
		t.Error("List should contain 2")
	}

	if list.Contains(4) {
		t.Error("List should not contain 4")
	}
}

func TestArrayList_ForEach(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	sum := 0
	list.ForEach(func(val int) {
		sum += val
	})

	if sum != 6 {
		t.Errorf("Sum should be 6, got %d", sum)
	}
}

func TestArrayList_ToSlice(t *testing.T) {
	list := New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	slice := list.ToSlice()

	if len(slice) != 3 {
		t.Errorf("Slice length should be 3, got %d", len(slice))
	}

	expected := []int{1, 2, 3}
	for i, val := range slice {
		if val != expected[i] {
			t.Errorf("Slice[%d] should be %d, got %d", i, expected[i], val)
		}
	}
}

func TestArrayList_String(t *testing.T) {
	list := New[int]()

	// Test empty list
	if list.String() != "[]" {
		t.Errorf("Empty list string should be '[]', got '%s'", list.String())
	}

	// Test non-empty list
	list.Add(1)
	list.Add(2)
	list.Add(3)

	expected := "[1, 2, 3]"
	if list.String() != expected {
		t.Errorf("List string should be '%s', got '%s'", expected, list.String())
	}
}
