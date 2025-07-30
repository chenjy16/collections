package queue

import (
	"testing"
)

// Simple comparable type for testing
type TestInt int

func (t TestInt) CompareTo(other interface{}) int {
	if o, ok := other.(TestInt); ok {
		if t < o {
			return -1
		} else if t > o {
			return 1
		}
		return 0
	}
	// Return error value instead of panic
	return 0 // or could return a special value to indicate error
}

func TestPriorityQueue_Basic(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	if !pq.IsEmpty() {
		t.Error("New queue should be empty")
	}

	if pq.Size() != 0 {
		t.Errorf("Empty queue size should be 0, got %d", pq.Size())
	}

	// Test adding elements
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	if pq.Size() != 4 {
		t.Errorf("Queue size should be 4, got %d", pq.Size())
	}

	// Create a max heap (default is min heap)
	maxPQ := NewPriorityQueueWithComparator[TestInt](func(a, b TestInt) int {
		// Reverse comparison result
		return -1 * a.CompareTo(b)
	})

	maxPQ.Add(TestInt(3))
	maxPQ.Add(TestInt(1))
	maxPQ.Add(TestInt(4))
	maxPQ.Add(TestInt(2))

	// Verify max value is at top
	if val, ok := maxPQ.Peek(); !ok || val != TestInt(4) {
		t.Errorf("Max heap peek should return 4, got %v", val)
	}
}

func TestPriorityQueue_WithCapacity(t *testing.T) {
	// Test capacity limit
	pq := WithCapacity[TestInt](3)

	pq.Add(TestInt(1))
	pq.Add(TestInt(2))
	pq.Add(TestInt(3))

	// Adding 4th element should return error
	if err := pq.Add(TestInt(4)); err == nil {
		t.Error("Adding to full queue should return error")
	}

	if pq.Size() != 3 {
		t.Errorf("Queue size should be 3, got %d", pq.Size())
	}
}

func TestPriorityQueue_Order(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Add elements in random order
	elements := []TestInt{5, 2, 8, 1, 9, 3}
	for _, elem := range elements {
		pq.Add(elem)
	}

	// Verify elements are sorted by priority
	expectedOrder := []int{1, 2, 3, 5, 8, 9} // min heap
	var result []int

	for !pq.IsEmpty() {
		if val, err := pq.Remove(); err == nil {
			result = append(result, int(val))
		}
	}

	for i, expected := range expectedOrder {
		if i >= len(result) || result[i] != expected {
			t.Errorf("Expected %d at position %d, got %d", expected, i, result[i])
		}
	}
}

func TestPriorityQueue_MaxHeap(t *testing.T) {
	// Create a max heap
	pq := NewPriorityQueueWithComparator[TestInt](func(a, b TestInt) int {
		// Reverse comparison result
		return -1 * a.CompareTo(b)
	})

	// Add elements
	elements := []TestInt{5, 2, 8, 1, 9, 3}
	for _, elem := range elements {
		pq.Add(elem)
	}

	// Verify elements are sorted by priority (max heap)
	expectedOrder := []int{9, 8, 5, 3, 2, 1}
	var result []int

	for !pq.IsEmpty() {
		if val, err := pq.Remove(); err == nil {
			result = append(result, int(val))
		}
	}

	for i, expected := range expectedOrder {
		if i >= len(result) || result[i] != expected {
			t.Errorf("Expected %d at position %d, got %d", expected, i, result[i])
		}
	}
}

func TestPriorityQueue_Offer(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Add elements
	if !pq.Offer(TestInt(3)) {
		t.Error("Offer should succeed")
	}
	if !pq.Offer(TestInt(1)) {
		t.Error("Offer should succeed")
	}
	if !pq.Offer(TestInt(2)) {
		t.Error("Offer should succeed")
	}

	// Verify min element is at top
	if val, ok := pq.Peek(); !ok || val != TestInt(1) {
		t.Errorf("Peek should return 1, got %v", val)
	}

	// Test with capacity
	cappedPQ := WithCapacity[TestInt](2)
	cappedPQ.Offer(TestInt(1))
	cappedPQ.Offer(TestInt(2))

	// 4th element should return false
	if cappedPQ.Offer(TestInt(3)) {
		t.Error("Offer to full queue should return false")
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Remove from empty queue should return error
	if _, err := pq.Remove(); err == nil {
		t.Error("Remove from empty queue should return error")
	}

	// Add elements (unordered)
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	// Verify elements are removed in priority order (min heap)
	expectedOrder := []TestInt{1, 2, 3, 4}
	for _, expected := range expectedOrder {
		if val, err := pq.Remove(); err != nil || val != expected {
			t.Errorf("Expected %v, got %v with error %v", expected, val, err)
		}
	}

	// Queue should be empty
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after removing all elements")
	}
}

func TestPriorityQueue_Poll(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Poll from empty queue should return false
	if _, ok := pq.Poll(); ok {
		t.Error("Poll from empty queue should return false")
	}

	// Add elements
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	// Verify elements are removed in priority order
	expectedOrder := []TestInt{1, 2, 3, 4}
	for _, expected := range expectedOrder {
		if val, ok := pq.Poll(); !ok || val != expected {
			t.Errorf("Expected %v, got %v with ok=%v", expected, val, ok)
		}
	}

	// Queue should be empty
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after polling all elements")
	}
}

func TestPriorityQueue_Element(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Get element from empty queue should return error
	if _, err := pq.Element(); err == nil {
		t.Error("Element from empty queue should return error")
	}

	// Add elements
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	// Get top element
	if val, err := pq.Element(); err != nil || val != TestInt(1) {
		t.Errorf("Element should return 1, got %v with error %v", val, err)
	}

	// Confirm element was not removed
	if pq.Size() != 4 {
		t.Errorf("Size should still be 4, got %d", pq.Size())
	}
}

func TestPriorityQueue_Peek(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Peek from empty queue should return false
	if _, ok := pq.Peek(); ok {
		t.Error("Peek from empty queue should return false")
	}

	// Add elements
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	// Peek top element
	if val, ok := pq.Peek(); !ok || val != TestInt(1) {
		t.Errorf("Peek should return 1, got %v with ok=%v", val, ok)
	}

	// Confirm element was not removed
	if pq.Size() != 4 {
		t.Errorf("Size should still be 4, got %d", pq.Size())
	}
}

func TestPriorityQueue_ToSlice(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Empty queue should return empty slice
	slice := pq.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Empty queue slice should have length 0, got %d", len(slice))
	}

	// Add elements
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	// Convert to slice
	slice = pq.ToSlice()
	if len(slice) != 4 {
		t.Errorf("Slice should have length 4, got %d", len(slice))
	}

	// Confirm original queue was not modified
	if pq.Size() != 4 {
		t.Errorf("Original queue size should still be 4, got %d", pq.Size())
	}

	// Modifying slice should not affect original queue
	slice[0] = TestInt(99)
	if val, _ := pq.Peek(); val == TestInt(99) {
		t.Error("Modifying slice should not affect original queue")
	}
}

func TestPriorityQueue_ToSortedSlice(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Empty queue should return empty slice
	slice := pq.ToSortedSlice()
	if len(slice) != 0 {
		t.Errorf("Empty queue slice should have length 0, got %d", len(slice))
	}

	// Add elements
	pq.Add(TestInt(3))
	pq.Add(TestInt(1))
	pq.Add(TestInt(4))
	pq.Add(TestInt(2))

	// Convert to sorted slice
	slice = pq.ToSortedSlice()
	if len(slice) != 4 {
		t.Errorf("Slice should have length 4, got %d", len(slice))
	}

	// Verify slice is sorted by priority
	expectedOrder := []TestInt{1, 2, 3, 4}
	for i, expected := range expectedOrder {
		if slice[i] != expected {
			t.Errorf("Expected %v at position %d, got %v", expected, i, slice[i])
		}
	}

	// Confirm original queue was not modified
	if pq.Size() != 4 {
		t.Errorf("Original queue size should still be 4, got %d", pq.Size())
	}
}

func TestPriorityQueue_Clear(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Add elements
	pq.Add(TestInt(1))
	pq.Add(TestInt(2))
	pq.Add(TestInt(3))

	if pq.Size() != 3 {
		t.Errorf("Size should be 3, got %d", pq.Size())
	}

	// Clear queue
	pq.Clear()

	if !pq.IsEmpty() {
		t.Error("Queue should be empty after clear")
	}

	if pq.Size() != 0 {
		t.Errorf("Size should be 0 after clear, got %d", pq.Size())
	}

	// Remove from empty queue should return error
	if _, err := pq.Remove(); err == nil {
		t.Error("Remove from cleared queue should return error")
	}
}

func TestPriorityQueue_Contains(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Empty queue should not contain any elements
	if pq.Contains(TestInt(1)) {
		t.Error("Empty queue should not contain any elements")
	}

	// Add elements
	pq.Add(TestInt(1))
	pq.Add(TestInt(2))
	pq.Add(TestInt(3))

	// Check contained elements
	if !pq.Contains(TestInt(1)) {
		t.Error("Queue should contain 1")
	}
	if !pq.Contains(TestInt(2)) {
		t.Error("Queue should contain 2")
	}

	// Check non-contained elements
	if pq.Contains(TestInt(4)) {
		t.Error("Queue should not contain 4")
	}
}

func TestPriorityQueue_ForEach(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Add elements
	pq.Add(TestInt(1))
	pq.Add(TestInt(2))
	pq.Add(TestInt(3))

	var result []TestInt
	pq.ForEach(func(val TestInt) {
		result = append(result, val)
	})

	// Verify result
	if len(result) != 3 {
		t.Errorf("ForEach should visit 3 elements, got %d", len(result))
	}

	// Confirm original queue was not modified
	if pq.Size() != 3 {
		t.Errorf("Original queue size should still be 3, got %d", pq.Size())
	}
}

func TestPriorityQueue_String(t *testing.T) {
	pq := NewPriorityQueue[TestInt]()

	// Empty queue
	str := pq.String()
	if str != "[]" {
		t.Errorf("Empty queue string should be '[]', got '%s'", str)
	}

	// Add one element
	pq.Add(TestInt(1))
	str = pq.String()
	if str != "[1]" {
		t.Errorf("Single element queue string should be '[1]', got '%s'", str)
	}

	// Add more elements
	pq.Add(TestInt(2))
	pq.Add(TestInt(3))

	// Note: String() method returns the internal array representation, not necessarily sorted by priority
	// So here we only check length and format
	str = pq.String()
	if len(str) < 7 { // Should be at least "[1, 2, 3]" length
		t.Errorf("String representation seems too short: '%s'", str)
	}
}
