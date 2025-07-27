package queue

import (
	"testing"
)

func TestLinkedList_New(t *testing.T) {
	list := New[int]()
	if list.Size() != 0 {
		t.Errorf("Expected empty list, got size %d", list.Size())
	}
	if !list.IsEmpty() {
		t.Error("Expected IsEmpty() to return true for new list")
	}
}

func TestLinkedList_WithCapacity(t *testing.T) {
	list := WithCapacity[int](5)
	if list.Size() != 0 {
		t.Errorf("Expected empty list, got size %d", list.Size())
	}
	if !list.IsEmpty() {
		t.Error("Expected IsEmpty() to return true for new list")
	}

	// 测试容量限制
	for i := 0; i < 5; i++ {
		err := list.Add(i)
		if err != nil {
			t.Errorf("Unexpected error adding element %d: %v", i, err)
		}
	}

	// 添加第6个元素应该返回错误
	err := list.Add(5)
	if err != ErrFullQueue {
		t.Errorf("Expected ErrFullQueue, got %v", err)
	}
}

func TestLinkedList_FromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	list := FromSlice(slice)

	if list.Size() != len(slice) {
		t.Errorf("Expected size %d, got %d", len(slice), list.Size())
	}

	// 验证元素顺序
	for i, expected := range slice {
		val, err := list.Remove()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != expected {
			t.Errorf("Expected %v at position %d, got %v", expected, i, val)
		}
	}

	if !list.IsEmpty() {
		t.Error("List should be empty after removing all elements")
	}
}

func TestLinkedList_Add(t *testing.T) {
	list := New[int]()

	// 添加元素
	for i := 0; i < 5; i++ {
		err := list.Add(i)
		if err != nil {
			t.Errorf("Unexpected error adding element %d: %v", i, err)
		}
	}

	if list.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list.Size())
	}

	// 验证元素顺序 (FIFO)
	for i := 0; i < 5; i++ {
		val, err := list.Remove()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}
}

func TestLinkedList_Offer(t *testing.T) {
	list := WithCapacity[int](3)

	// 添加3个元素
	for i := 0; i < 3; i++ {
		success := list.Offer(i)
		if !success {
			t.Errorf("Expected Offer to return true for element %d", i)
		}
	}

	// 第4个元素应该返回false
	success := list.Offer(3)
	if success {
		t.Error("Expected Offer to return false when queue is full")
	}

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}
}

func TestLinkedList_Remove(t *testing.T) {
	list := New[int]()

	// 从空队列移除应该返回错误
	_, err := list.Remove()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// 移除元素
	for i := 0; i < 3; i++ {
		val, err := list.Remove()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	// 队列应该为空
	if !list.IsEmpty() {
		t.Error("Queue should be empty after removing all elements")
	}
}

func TestLinkedList_Poll(t *testing.T) {
	list := New[int]()

	// 从空队列Poll应该返回false
	_, success := list.Poll()
	if success {
		t.Error("Expected Poll to return false for empty queue")
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// Poll元素
	for i := 0; i < 3; i++ {
		val, success := list.Poll()
		if !success {
			t.Errorf("Unexpected failure polling element %d", i)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	// 队列应该为空
	if !list.IsEmpty() {
		t.Error("Queue should be empty after polling all elements")
	}
}

func TestLinkedList_Element(t *testing.T) {
	list := New[int]()

	// 从空队列获取元素应该返回错误
	_, err := list.Element()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// 获取头部元素
	val, err := list.Element()
	if err != nil {
		t.Errorf("Unexpected error getting element: %v", err)
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}

	// 确认元素没有被移除
	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}
}

func TestLinkedList_Peek(t *testing.T) {
	list := New[int]()

	// 从空队列Peek应该返回false
	_, success := list.Peek()
	if success {
		t.Error("Expected Peek to return false for empty queue")
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// Peek头部元素
	val, success := list.Peek()
	if !success {
		t.Error("Unexpected failure peeking element")
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}

	// 确认元素没有被移除
	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}
}

func TestLinkedList_AddFirst(t *testing.T) {
	list := New[int]()

	// 添加元素到头部
	for i := 0; i < 3; i++ {
		err := list.AddFirst(i)
		if err != nil {
			t.Errorf("Unexpected error adding element %d: %v", i, err)
		}
	}

	// 验证元素顺序 (LIFO)
	for i := 0; i < 3; i++ {
		val, err := list.RemoveFirst()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != 2-i {
			t.Errorf("Expected %d, got %d", 2-i, val)
		}
	}
}

func TestLinkedList_AddLast(t *testing.T) {
	list := New[int]()

	// 添加元素到尾部
	for i := 0; i < 3; i++ {
		err := list.AddLast(i)
		if err != nil {
			t.Errorf("Unexpected error adding element %d: %v", i, err)
		}
	}

	// 验证元素顺序 (FIFO)
	for i := 0; i < 3; i++ {
		val, err := list.RemoveFirst()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}
}

func TestLinkedList_RemoveFirst(t *testing.T) {
	list := New[int]()

	// 从空队列移除应该返回错误
	_, err := list.RemoveFirst()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// 移除头部元素
	for i := 0; i < 3; i++ {
		val, err := list.RemoveFirst()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}
}

func TestLinkedList_RemoveLast(t *testing.T) {
	list := New[int]()

	// 从空队列移除应该返回错误
	_, err := list.RemoveLast()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// 移除尾部元素
	for i := 0; i < 3; i++ {
		val, err := list.RemoveLast()
		if err != nil {
			t.Errorf("Unexpected error removing element %d: %v", i, err)
		}
		if val != 2-i {
			t.Errorf("Expected %d, got %d", 2-i, val)
		}
	}
}

func TestLinkedList_GetFirst(t *testing.T) {
	list := New[int]()

	// 从空队列获取元素应该返回错误
	_, err := list.GetFirst()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// 获取头部元素
	val, err := list.GetFirst()
	if err != nil {
		t.Errorf("Unexpected error getting element: %v", err)
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}

	// 确认元素没有被移除
	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}
}

func TestLinkedList_GetLast(t *testing.T) {
	list := New[int]()

	// 从空队列获取元素应该返回错误
	_, err := list.GetLast()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = list.Add(i)
	}

	// 获取尾部元素
	val, err := list.GetLast()
	if err != nil {
		t.Errorf("Unexpected error getting element: %v", err)
	}
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	// 确认元素没有被移除
	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}
}

func TestLinkedList_ToSlice(t *testing.T) {
	list := New[int]()

	// 空队列应该返回空切片
	slice := list.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(slice))
	}

	// 添加元素
	expected := []int{0, 1, 2}
	for _, v := range expected {
		_ = list.Add(v)
	}

	// 转换为切片
	slice = list.ToSlice()
	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}

	// 验证元素顺序
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, slice[i])
		}
	}

	// 确认原队列没有被修改
	if list.Size() != len(expected) {
		t.Errorf("Expected queue size %d, got %d", len(expected), list.Size())
	}

	// 修改切片不应影响原队列
	slice[0] = 100
	val, _ := list.GetFirst()
	if val == 100 {
		t.Error("Modifying slice should not affect the original queue")
	}
}

func TestLinkedList_Clear(t *testing.T) {
	list := New[int]()

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = list.Add(i)
	}

	if list.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list.Size())
	}

	// 清空队列
	list.Clear()

	if !list.IsEmpty() {
		t.Error("Expected queue to be empty after Clear()")
	}
	if list.Size() != 0 {
		t.Errorf("Expected size 0, got %d", list.Size())
	}

	// 从空队列移除应该返回错误
	_, err := list.Remove()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}
}

func TestLinkedList_Contains(t *testing.T) {
	list := New[int]()

	// 空队列不应包含任何元素
	if list.Contains(1) {
		t.Error("Empty queue should not contain any elements")
	}

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = list.Add(i)
	}

	// 检查包含的元素
	for i := 0; i < 5; i++ {
		if !list.Contains(i) {
			t.Errorf("Queue should contain element %d", i)
		}
	}

	// 检查不包含的元素
	if list.Contains(10) {
		t.Error("Queue should not contain element 10")
	}
}

func TestLinkedList_ForEach(t *testing.T) {
	list := New[int]()

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = list.Add(i)
	}

	// 使用ForEach遍历元素
	sum := 0
	list.ForEach(func(e int) {
		sum += e
	})

	// 验证结果
	expectedSum := 0 + 1 + 2 + 3 + 4
	if sum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, sum)
	}

	// 确认原队列没有被修改
	if list.Size() != 5 {
		t.Errorf("Expected queue size 5, got %d", list.Size())
	}
}

func TestLinkedList_String(t *testing.T) {
	list := New[int]()

	// 空队列
	if list.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", list.String())
	}

	// 添加一个元素
	_ = list.Add(1)
	if list.String() != "[1]" {
		t.Errorf("Expected '[1]', got '%s'", list.String())
	}

	// 添加更多元素
	_ = list.Add(2)
	_ = list.Add(3)
	if list.String() != "[1, 2, 3]" {
		t.Errorf("Expected '[1, 2, 3]', got '%s'", list.String())
	}
}
