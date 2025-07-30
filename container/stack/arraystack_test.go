package stack

import (
	"testing"
)

func TestArrayStack_Basic(t *testing.T) {
	stack := New[int]()

	// Test empty stack
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}

	if stack.Size() != 0 {
		t.Error("New stack size should be 0")
	}

	// Test push elements
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Size() != 3 {
		t.Errorf("Stack size should be 3, got %d", stack.Size())
	}

	if stack.IsEmpty() {
		t.Error("Stack should not be empty")
	}
}

func TestArrayStack_WithCapacity(t *testing.T) {
	stack := WithCapacity[int](2)

	// Test push within capacity
	err := stack.Push(1)
	if err != nil {
		t.Error("Push should not return error")
	}

	err = stack.Push(2)
	if err != nil {
		t.Error("Push should not return error")
	}

	// Test push beyond capacity
	err = stack.Push(3)
	if err == nil {
		t.Error("Push should return error when capacity exceeded")
	}

	if stack.Size() != 2 {
		t.Errorf("Stack size should be 2, got %d", stack.Size())
	}
}

func TestArrayStack_Pop(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	// Test pop from non-empty stack
	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop should not return error, got %v", err)
	}

	if val != 3 {
		t.Errorf("Pop should return 3, got %d", val)
	}

	if stack.Size() != 2 {
		t.Errorf("Stack size should be 2, got %d", stack.Size())
	}

	// Test pop all elements
	stack.Pop()
	stack.Pop()

	// Test pop from empty stack
	_, err = stack.Pop()
	if err == nil {
		t.Error("Pop from empty stack should return error")
	}
}

func TestArrayStack_Peek(t *testing.T) {
	stack := New[int]()

	// Test peek from empty stack
	_, err := stack.Peek()
	if err == nil {
		t.Error("Peek from empty stack should return error")
	}

	// Test peek from non-empty stack
	stack.Push(1)
	stack.Push(2)

	val, err := stack.Peek()
	if err != nil {
		t.Errorf("Peek should not return error, got %v", err)
	}

	if val != 2 {
		t.Errorf("Peek should return 2, got %d", val)
	}

	// Verify peek doesn't remove the item
	if stack.Size() != 2 {
		t.Errorf("Stack size should still be 2, got %d", stack.Size())
	}
}

func TestArrayStack_Clear(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	stack.Clear()

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after clear")
	}

	if stack.Size() != 0 {
		t.Errorf("Stack size should be 0 after clear, got %d", stack.Size())
	}
}

func TestArrayStack_Contains(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if !stack.Contains(2) {
		t.Error("Stack should contain 2")
	}

	if stack.Contains(4) {
		t.Error("Stack should not contain 4")
	}
}

func TestArrayStack_ForEach(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	sum := 0
	stack.ForEach(func(val int) {
		sum += val
	})

	if sum != 6 {
		t.Errorf("Sum should be 6, got %d", sum)
	}
}

func TestArrayStack_ToSlice(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	slice := stack.ToSlice()

	if len(slice) != 3 {
		t.Errorf("Slice length should be 3, got %d", len(slice))
	}

	// Stack is LIFO, so slice should be in reverse order
	expected := []int{1, 2, 3} // ToSlice returns bottom to top order
	for i, val := range slice {
		if val != expected[i] {
			t.Errorf("Slice[%d] should be %d, got %d", i, expected[i], val)
		}
	}
}

func TestArrayStack_String(t *testing.T) {
	stack := New[int]()

	// Test empty stack
	if stack.String() != "[]" {
		t.Errorf("Empty stack string should be '[]', got '%s'", stack.String())
	}

	// Test non-empty stack
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	expected := "[1, 2, 3]" // String shows bottom to top order
	if stack.String() != expected {
		t.Errorf("Stack string should be '%s', got '%s'", expected, stack.String())
	}
}
