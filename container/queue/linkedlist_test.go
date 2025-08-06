package queue

import (
	"errors"
	"testing"

	"github.com/chenjianyu/collections/container/common"
)

func TestLinkedListQueue_New(t *testing.T) {
	q := New[int]()
	if q == nil {
		t.Error("New() should not return nil")
	}
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if q.Size() != 0 {
		t.Errorf("New queue size should be 0, got %d", q.Size())
	}
}

func TestLinkedListQueue_WithCapacity(t *testing.T) {
	q := WithCapacity[int](3)
	if q == nil {
		t.Error("WithCapacity() should not return nil")
	}
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if q.Size() != 0 {
		t.Errorf("New queue size should be 0, got %d", q.Size())
	}
}

func TestLinkedListQueue_FromSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	q := FromSlice(slice)
	if q.Size() != 3 {
		t.Errorf("Queue size should be 3, got %d", q.Size())
	}
	if q.IsEmpty() {
		t.Error("Queue should not be empty")
	}
}

func TestLinkedListQueue_Add(t *testing.T) {
	q := New[int]()
	
	// Test adding to empty queue
	err := q.Add(1)
	if err != nil {
		t.Errorf("Add should not return error, got %v", err)
	}
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}
	if q.IsEmpty() {
		t.Error("Queue should not be empty")
	}
	
	// Test adding multiple elements
	err = q.Add(2)
	if err != nil {
		t.Errorf("Add should not return error, got %v", err)
	}
	err = q.Add(3)
	if err != nil {
		t.Errorf("Add should not return error, got %v", err)
	}
	if q.Size() != 3 {
		t.Errorf("Queue size should be 3, got %d", q.Size())
	}
}

func TestLinkedListQueue_AddWithCapacity(t *testing.T) {
	q := WithCapacity[int](2)
	
	// Add within capacity
	err := q.Add(1)
	if err != nil {
		t.Errorf("Add should not return error, got %v", err)
	}
	err = q.Add(2)
	if err != nil {
		t.Errorf("Add should not return error, got %v", err)
	}
	
	// Add beyond capacity
	err = q.Add(3)
	if err == nil {
		t.Error("Add should return error when queue is full")
	}
	if !errors.Is(err, common.ErrFullContainer) {
		t.Errorf("Expected FullContainerError, got %T", err)
	}
}

func TestLinkedListQueue_Offer(t *testing.T) {
	q := New[int]()
	
	// Test offering to empty queue
	success := q.Offer(1)
	if !success {
		t.Error("Offer should return true")
	}
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}
	
	// Test offering with capacity
	q = WithCapacity[int](1)
	success = q.Offer(1)
	if !success {
		t.Error("Offer should return true")
	}
	success = q.Offer(2)
	if success {
		t.Error("Offer should return false when queue is full")
	}
}

func TestLinkedListQueue_Remove(t *testing.T) {
	q := New[int]()
	
	// Test removing from empty queue
	_, err := q.Remove()
	if err == nil {
		t.Error("Remove should return error when queue is empty")
	}
	if !errors.Is(err, common.ErrEmptyContainer) {
		t.Errorf("Expected EmptyContainerError, got %T", err)
	}
	
	// Test removing from non-empty queue
	q.Add(1)
	q.Add(2)
	val, err := q.Remove()
	if err != nil {
		t.Errorf("Remove should not return error, got %v", err)
	}
	if val != 1 {
		t.Errorf("Remove should return 1, got %d", val)
	}
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}
}

func TestLinkedListQueue_Poll(t *testing.T) {
	q := New[int]()
	
	// Test polling from empty queue
	_, success := q.Poll()
	if success {
		t.Error("Poll should return false when queue is empty")
	}
	
	// Test polling from non-empty queue
	q.Add(1)
	q.Add(2)
	val, success := q.Poll()
	if !success {
		t.Error("Poll should return true")
	}
	if val != 1 {
		t.Errorf("Poll should return 1, got %d", val)
	}
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}
}

func TestLinkedListQueue_Element(t *testing.T) {
	q := New[int]()
	
	// Test element from empty queue
	_, err := q.Element()
	if err == nil {
		t.Error("Element should return error when queue is empty")
	}
	
	// Test element from non-empty queue
	q.Add(1)
	q.Add(2)
	val, err := q.Element()
	if err != nil {
		t.Errorf("Element should not return error, got %v", err)
	}
	if val != 1 {
		t.Errorf("Element should return 1, got %d", val)
	}
	if q.Size() != 2 {
		t.Errorf("Queue size should remain 2, got %d", q.Size())
	}
}

func TestLinkedListQueue_Peek(t *testing.T) {
	q := New[int]()
	
	// Test peek from empty queue
	_, success := q.Peek()
	if success {
		t.Error("Peek should return false when queue is empty")
	}
	
	// Test peek from non-empty queue
	q.Add(1)
	q.Add(2)
	val, success := q.Peek()
	if !success {
		t.Error("Peek should return true")
	}
	if val != 1 {
		t.Errorf("Peek should return 1, got %d", val)
	}
	if q.Size() != 2 {
		t.Errorf("Queue size should remain 2, got %d", q.Size())
	}
}

func TestLinkedListQueue_DequeOperations(t *testing.T) {
	q := New[int]()
	
	// Test AddFirst
	err := q.AddFirst(1)
	if err != nil {
		t.Errorf("AddFirst should not return error, got %v", err)
	}
	err = q.AddFirst(2)
	if err != nil {
		t.Errorf("AddFirst should not return error, got %v", err)
	}
	
	// Test AddLast
	err = q.AddLast(3)
	if err != nil {
		t.Errorf("AddLast should not return error, got %v", err)
	}
	
	if q.Size() != 3 {
		t.Errorf("Queue size should be 3, got %d", q.Size())
	}
	
	// Test RemoveFirst
	val, err := q.RemoveFirst()
	if err != nil {
		t.Errorf("RemoveFirst should not return error, got %v", err)
	}
	if val != 2 {
		t.Errorf("RemoveFirst should return 2, got %d", val)
	}
	
	// Test RemoveLast
	val, err = q.RemoveLast()
	if err != nil {
		t.Errorf("RemoveLast should not return error, got %v", err)
	}
	if val != 3 {
		t.Errorf("RemoveLast should return 3, got %d", val)
	}
	
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}
}

func TestLinkedListQueue_OfferOperations(t *testing.T) {
	q := WithCapacity[int](2)
	
	// Test OfferFirst
	success := q.OfferFirst(1)
	if !success {
		t.Error("OfferFirst should return true")
	}
	success = q.OfferFirst(2)
	if !success {
		t.Error("OfferFirst should return true")
	}
	success = q.OfferFirst(3)
	if success {
		t.Error("OfferFirst should return false when queue is full")
	}
	
	q.Clear()
	
	// Test OfferLast
	success = q.OfferLast(1)
	if !success {
		t.Error("OfferLast should return true")
	}
	success = q.OfferLast(2)
	if !success {
		t.Error("OfferLast should return true")
	}
	success = q.OfferLast(3)
	if success {
		t.Error("OfferLast should return false when queue is full")
	}
}

func TestLinkedListQueue_PollOperations(t *testing.T) {
	q := New[int]()
	q.Add(1)
	q.Add(2)
	q.Add(3)
	
	// Test PollFirst
	val, success := q.PollFirst()
	if !success {
		t.Error("PollFirst should return true")
	}
	if val != 1 {
		t.Errorf("PollFirst should return 1, got %d", val)
	}
	
	// Test PollLast
	val, success = q.PollLast()
	if !success {
		t.Error("PollLast should return true")
	}
	if val != 3 {
		t.Errorf("PollLast should return 3, got %d", val)
	}
	
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}
}

func TestLinkedListQueue_Contains(t *testing.T) {
	q := New[int]()
	q.Add(1)
	q.Add(2)
	q.Add(3)
	
	if !q.Contains(2) {
		t.Error("Queue should contain 2")
	}
	if q.Contains(4) {
		t.Error("Queue should not contain 4")
	}
}

func TestLinkedListQueue_Clear(t *testing.T) {
	q := New[int]()
	q.Add(1)
	q.Add(2)
	q.Add(3)
	
	q.Clear()
	if !q.IsEmpty() {
		t.Error("Queue should be empty after clear")
	}
	if q.Size() != 0 {
		t.Errorf("Queue size should be 0 after clear, got %d", q.Size())
	}
}

func TestLinkedListQueue_ForEach(t *testing.T) {
	q := New[int]()
	q.Add(1)
	q.Add(2)
	q.Add(3)
	
	sum := 0
	q.ForEach(func(val int) {
		sum += val
	})
	
	if sum != 6 {
		t.Errorf("Sum should be 6, got %d", sum)
	}
}

func TestLinkedListQueue_String(t *testing.T) {
	q := New[int]()
	q.Add(1)
	q.Add(2)
	q.Add(3)
	
	str := q.String()
	if str == "" {
		t.Error("String should not be empty")
	}
}

func TestLinkedListQueue_EmptyOperations(t *testing.T) {
	q := New[int]()
	
	// Test RemoveFirst on empty queue
	_, err := q.RemoveFirst()
	if err == nil {
		t.Error("RemoveFirst should return error when queue is empty")
	}
	if !errors.Is(err, common.ErrEmptyContainer) {
		t.Errorf("Expected EmptyContainerError, got %T", err)
	}
	
	// Test RemoveLast on empty queue
	_, err = q.RemoveLast()
	if err == nil {
		t.Error("RemoveLast should return error when queue is empty")
	}
	if !errors.Is(err, common.ErrEmptyContainer) {
		t.Errorf("Expected EmptyContainerError, got %T", err)
	}
	
	// Test PollFirst on empty queue
	_, success := q.PollFirst()
	if success {
		t.Error("PollFirst should return false when queue is empty")
	}
	
	// Test PollLast on empty queue
	_, success = q.PollLast()
	if success {
		t.Error("PollLast should return false when queue is empty")
	}
}