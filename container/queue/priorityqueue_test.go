package queue

import (
	"testing"
)

// 为测试创建一个简单的可比较类型
type IntComparable int

func (i IntComparable) CompareTo(other interface{}) int {
	otherInt, ok := other.(IntComparable)
	if !ok {
		// 返回错误值而不是 panic
		return 0 // 或者可以返回一个特殊值表示错误
	}
	if i < otherInt {
		return -1
	} else if i > otherInt {
		return 1
	}
	return 0
}

func TestPriorityQueue_NewPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()
	if pq.Size() != 0 {
		t.Errorf("Expected empty queue, got size %d", pq.Size())
	}
	if !pq.IsEmpty() {
		t.Error("Expected IsEmpty() to return true for new queue")
	}
}

func TestPriorityQueue_NewPriorityQueueWithComparator(t *testing.T) {
	// 创建一个最大堆（默认是最小堆）
	pq := NewPriorityQueueWithComparator(func(a, b int) int {
		if a > b {
			return -1 // 反转比较结果
		} else if a < b {
			return 1
		}
		return 0
	})

	if pq.Size() != 0 {
		t.Errorf("Expected empty queue, got size %d", pq.Size())
	}

	// 添加元素
	_ = pq.Add(1)
	_ = pq.Add(5)
	_ = pq.Add(3)

	// 验证最大值在顶部
	val, ok := pq.Peek()
	if !ok {
		t.Error("Peek failed on non-empty queue")
	}
	if val != 5 {
		t.Errorf("Expected 5 (max value) at top, got %d", val)
	}
}

func TestPriorityQueue_WithCapacity(t *testing.T) {
	pq := WithCapacity[IntComparable](3)
	if pq.Size() != 0 {
		t.Errorf("Expected empty queue, got size %d", pq.Size())
	}

	// 测试容量限制
	for i := 0; i < 3; i++ {
		err := pq.Add(IntComparable(i))
		if err != nil {
			t.Errorf("Unexpected error adding element %d: %v", i, err)
		}
	}

	// 添加第4个元素应该返回错误
	err := pq.Add(IntComparable(3))
	if err != ErrFullQueue {
		t.Errorf("Expected ErrFullQueue, got %v", err)
	}
}

func TestPriorityQueue_FromSlice(t *testing.T) {
	slice := []IntComparable{5, 3, 1, 4, 2}
	pq := NewFromSlice(slice)

	if pq.Size() != len(slice) {
		t.Errorf("Expected size %d, got %d", len(slice), pq.Size())
	}

	// 验证元素按优先级排序
	expectedOrder := []int{1, 2, 3, 4, 5} // 最小堆
	for _, expected := range expectedOrder {
		val, err := pq.Remove()
		if err != nil {
			t.Errorf("Unexpected error removing element: %v", err)
		}
		if int(val) != expected {
			t.Errorf("Expected %v, got %v", expected, val)
		}
	}

	if !pq.IsEmpty() {
		t.Error("Queue should be empty after removing all elements")
	}
}

func TestPriorityQueue_FromSliceWithComparator(t *testing.T) {
	slice := []int{5, 3, 1, 4, 2}
	// 创建一个最大堆
	pq := NewFromSliceWithComparator(slice, func(a, b int) int {
		if a > b {
			return -1 // 反转比较结果
		} else if a < b {
			return 1
		}
		return 0
	})

	if pq.Size() != len(slice) {
		t.Errorf("Expected size %d, got %d", len(slice), pq.Size())
	}

	// 验证元素按优先级排序（最大堆）
	expectedOrder := []int{5, 4, 3, 2, 1}
	for _, expected := range expectedOrder {
		val, err := pq.Remove()
		if err != nil {
			t.Errorf("Unexpected error removing element: %v", err)
		}
		if val != expected {
			t.Errorf("Expected %v, got %v", expected, val)
		}
	}
}

func TestPriorityQueue_Add(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		err := pq.Add(IntComparable(e))
		if err != nil {
			t.Errorf("Unexpected error adding element %d: %v", e, err)
		}
	}

	if pq.Size() != len(elements) {
		t.Errorf("Expected size %d, got %d", len(elements), pq.Size())
	}

	// 验证最小元素在顶部
	val, ok := pq.Peek()
	if !ok {
		t.Error("Peek failed on non-empty queue")
	}
	if val != IntComparable(1) {
		t.Errorf("Expected 1 (min value) at top, got %d", val)
	}
}

func TestPriorityQueue_Offer(t *testing.T) {
	pq := WithCapacity[IntComparable](3)

	// 添加3个元素
	for i := 0; i < 3; i++ {
		success := pq.Offer(IntComparable(i))
		if !success {
			t.Errorf("Expected Offer to return true for element %d", i)
		}
	}

	// 第4个元素应该返回false
	success := pq.Offer(IntComparable(3))
	if success {
		t.Error("Expected Offer to return false when queue is full")
	}

	if pq.Size() != 3 {
		t.Errorf("Expected size 3, got %d", pq.Size())
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 从空队列移除应该返回错误
	_, err := pq.Remove()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素（无序）
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 验证元素按优先级顺序移除（最小堆）
	expectedOrder := []int{1, 2, 3, 4, 5}
	for _, expected := range expectedOrder {
		val, err := pq.Remove()
		if err != nil {
			t.Errorf("Unexpected error removing element: %v", err)
		}
		if int(val) != expected {
			t.Errorf("Expected %v, got %v", expected, val)
		}
	}

	// 队列应该为空
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after removing all elements")
	}
}

func TestPriorityQueue_Poll(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 从空队列Poll应该返回false
	_, success := pq.Poll()
	if success {
		t.Error("Expected Poll to return false for empty queue")
	}

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 验证元素按优先级顺序移除
	expectedOrder := []int{1, 2, 3, 4, 5}
	for _, expected := range expectedOrder {
		val, success := pq.Poll()
		if !success {
			t.Errorf("Unexpected failure polling element")
		}
		if int(val) != expected {
			t.Errorf("Expected %v, got %v", expected, val)
		}
	}

	// 队列应该为空
	if !pq.IsEmpty() {
		t.Error("Queue should be empty after polling all elements")
	}
}

func TestPriorityQueue_Element(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 从空队列获取元素应该返回错误
	_, err := pq.Element()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 获取顶部元素
	val, err := pq.Element()
	if err != nil {
		t.Errorf("Unexpected error getting element: %v", err)
	}
	if val != IntComparable(1) {
		t.Errorf("Expected 1, got %d", val)
	}

	// 确认元素没有被移除
	if pq.Size() != len(elements) {
		t.Errorf("Expected size %d, got %d", len(elements), pq.Size())
	}
}

func TestPriorityQueue_Peek(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 从空队列Peek应该返回false
	_, success := pq.Peek()
	if success {
		t.Error("Expected Peek to return false for empty queue")
	}

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// Peek顶部元素
	val, success := pq.Peek()
	if !success {
		t.Error("Unexpected failure peeking element")
	}
	if val != IntComparable(1) {
		t.Errorf("Expected 1, got %d", val)
	}

	// 确认元素没有被移除
	if pq.Size() != len(elements) {
		t.Errorf("Expected size %d, got %d", len(elements), pq.Size())
	}
}

func TestPriorityQueue_ToSlice(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 空队列应该返回空切片
	slice := pq.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(slice))
	}

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 转换为切片
	slice = pq.ToSlice()
	if len(slice) != len(elements) {
		t.Errorf("Expected slice length %d, got %d", len(elements), len(slice))
	}

	// 确认原队列没有被修改
	if pq.Size() != len(elements) {
		t.Errorf("Expected queue size %d, got %d", len(elements), pq.Size())
	}

	// 修改切片不应影响原队列
	slice[0] = IntComparable(100)
	val, _ := pq.Peek()
	if val == IntComparable(100) {
		t.Error("Modifying slice should not affect the original queue")
	}
}

func TestPriorityQueue_ToSortedSlice(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 空队列应该返回空切片
	slice := pq.ToSortedSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(slice))
	}

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 转换为排序切片
	slice = pq.ToSortedSlice()
	if len(slice) != len(elements) {
		t.Errorf("Expected slice length %d, got %d", len(elements), len(slice))
	}

	// 验证切片是按优先级排序的
	expectedOrder := []int{1, 2, 3, 4, 5}
	for i, expected := range expectedOrder {
		if int(slice[i]) != expected {
			t.Errorf("Expected %d at index %d, got %d", expected, i, slice[i])
		}
	}

	// 确认原队列没有被修改
	if pq.Size() != len(elements) {
		t.Errorf("Expected queue size %d, got %d", len(elements), pq.Size())
	}
}

func TestPriorityQueue_Clear(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	if pq.Size() != len(elements) {
		t.Errorf("Expected size %d, got %d", len(elements), pq.Size())
	}

	// 清空队列
	pq.Clear()

	if !pq.IsEmpty() {
		t.Error("Expected queue to be empty after Clear()")
	}
	if pq.Size() != 0 {
		t.Errorf("Expected size 0, got %d", pq.Size())
	}

	// 从空队列移除应该返回错误
	_, err := pq.Remove()
	if err != ErrEmptyQueue {
		t.Errorf("Expected ErrEmptyQueue, got %v", err)
	}
}

func TestPriorityQueue_Contains(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 空队列不应包含任何元素
	if pq.Contains(IntComparable(1)) {
		t.Error("Empty queue should not contain any elements")
	}

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 检查包含的元素
	for _, e := range elements {
		if !pq.Contains(IntComparable(e)) {
			t.Errorf("Queue should contain element %d", e)
		}
	}

	// 检查不包含的元素
	if pq.Contains(IntComparable(10)) {
		t.Error("Queue should not contain element 10")
	}
}

func TestPriorityQueue_ForEach(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 添加元素
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = pq.Add(IntComparable(e))
	}

	// 使用ForEach遍历元素
	sum := 0
	pq.ForEach(func(e IntComparable) {
		sum += int(e)
	})

	// 验证结果
	expectedSum := 0
	for _, e := range elements {
		expectedSum += e
	}

	if sum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, sum)
	}

	// 确认原队列没有被修改
	if pq.Size() != len(elements) {
		t.Errorf("Expected queue size %d, got %d", len(elements), pq.Size())
	}
}

func TestPriorityQueue_String(t *testing.T) {
	pq := NewPriorityQueue[IntComparable]()

	// 空队列
	if pq.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", pq.String())
	}

	// 添加一个元素
	_ = pq.Add(IntComparable(1))
	if pq.String() != "[1]" {
		t.Errorf("Expected '[1]', got '%s'", pq.String())
	}

	// 添加更多元素
	_ = pq.Add(IntComparable(2))
	_ = pq.Add(IntComparable(3))

	// 注意：String()方法返回的是内部数组的表示，不一定是按优先级排序的
	// 所以这里只检查长度和格式
	str := pq.String()
	if len(str) < 7 { // 至少应该是"[1, 2, 3]"的长度
		t.Errorf("Expected string representation with 3 elements, got '%s'", str)
	}
}
