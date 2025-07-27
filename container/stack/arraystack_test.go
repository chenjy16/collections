package stack

import (
	"testing"
)

func TestArrayStack_New(t *testing.T) {
	stack := New[int]()
	if stack.Size() != 0 {
		t.Errorf("Expected empty stack, got size %d", stack.Size())
	}
	if !stack.IsEmpty() {
		t.Error("Expected IsEmpty() to return true for new stack")
	}
}

func TestArrayStack_WithCapacity(t *testing.T) {
	stack := WithCapacity[int](3)
	if stack.Size() != 0 {
		t.Errorf("Expected empty stack, got size %d", stack.Size())
	}

	// 测试容量限制
	for i := 0; i < 3; i++ {
		err := stack.Push(i)
		if err != nil {
			t.Errorf("Unexpected error pushing element %d: %v", i, err)
		}
	}

	// 添加第4个元素应该返回错误
	err := stack.Push(3)
	if err != ErrFullStack {
		t.Errorf("Expected ErrFullStack, got %v", err)
	}
}

func TestArrayStack_FromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stack := FromSlice(slice)

	if stack.Size() != len(slice) {
		t.Errorf("Expected size %d, got %d", len(slice), stack.Size())
	}

	// 验证元素顺序（栈顶应该是切片的最后一个元素）
	val, err := stack.Peek()
	if err != nil {
		t.Errorf("Unexpected error peeking: %v", err)
	}
	if val != slice[len(slice)-1] {
		t.Errorf("Expected top element to be %d, got %d", slice[len(slice)-1], val)
	}

	// 验证所有元素
	for i := len(slice) - 1; i >= 0; i-- {
		val, err := stack.Pop()
		if err != nil {
			t.Errorf("Unexpected error popping element %d: %v", i, err)
		}
		if val != slice[i] {
			t.Errorf("Expected %d, got %d", slice[i], val)
		}
	}

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping all elements")
	}
}

func TestArrayStack_Push(t *testing.T) {
	stack := New[int]()

	// 添加元素
	for i := 0; i < 5; i++ {
		err := stack.Push(i)
		if err != nil {
			t.Errorf("Unexpected error pushing element %d: %v", i, err)
		}
	}

	if stack.Size() != 5 {
		t.Errorf("Expected size 5, got %d", stack.Size())
	}

	// 验证栈顶元素
	val, err := stack.Peek()
	if err != nil {
		t.Errorf("Unexpected error peeking: %v", err)
	}
	if val != 4 {
		t.Errorf("Expected top element to be 4, got %d", val)
	}
}

func TestArrayStack_Pop(t *testing.T) {
	stack := New[int]()

	// 从空栈弹出应该返回错误
	_, err := stack.Pop()
	if err != ErrEmptyStack {
		t.Errorf("Expected ErrEmptyStack, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = stack.Push(i)
	}

	// 弹出元素
	for i := 2; i >= 0; i-- {
		val, err := stack.Pop()
		if err != nil {
			t.Errorf("Unexpected error popping element %d: %v", i, err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	// 栈应该为空
	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping all elements")
	}

	// 再次弹出应该返回错误
	_, err = stack.Pop()
	if err != ErrEmptyStack {
		t.Errorf("Expected ErrEmptyStack, got %v", err)
	}
}

func TestArrayStack_Peek(t *testing.T) {
	stack := New[int]()

	// 从空栈获取元素应该返回错误
	_, err := stack.Peek()
	if err != ErrEmptyStack {
		t.Errorf("Expected ErrEmptyStack, got %v", err)
	}

	// 添加元素
	for i := 0; i < 3; i++ {
		_ = stack.Push(i)
	}

	// 获取栈顶元素
	val, err := stack.Peek()
	if err != nil {
		t.Errorf("Unexpected error peeking: %v", err)
	}
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	// 确认元素没有被移除
	if stack.Size() != 3 {
		t.Errorf("Expected size 3, got %d", stack.Size())
	}

	// 多次Peek应该返回相同的值
	val2, _ := stack.Peek()
	if val != val2 {
		t.Errorf("Expected peek to return same value on multiple calls, got %d and %d", val, val2)
	}
}

func TestArrayStack_Search(t *testing.T) {
	stack := New[int]()

	// 空栈搜索应该返回-1
	pos := stack.Search(1)
	if pos != -1 {
		t.Errorf("Expected -1 for empty stack search, got %d", pos)
	}

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = stack.Push(i)
	}

	// 搜索存在的元素
	for i := 0; i < 5; i++ {
		pos := stack.Search(i)
		expectedPos := 5 - i // 栈顶是1，往下依次增加
		if pos != expectedPos {
			t.Errorf("Expected position %d for element %d, got %d", expectedPos, i, pos)
		}
	}

	// 搜索不存在的元素
	pos = stack.Search(10)
	if pos != -1 {
		t.Errorf("Expected -1 for non-existent element, got %d", pos)
	}
}

func TestArrayStack_ToSlice(t *testing.T) {
	stack := New[int]()

	// 空栈应该返回空切片
	slice := stack.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(slice))
	}

	// 添加元素
	expected := []int{0, 1, 2, 3, 4}
	for _, v := range expected {
		_ = stack.Push(v)
	}

	// 转换为切片
	slice = stack.ToSlice()
	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}

	// 验证元素顺序（栈底到栈顶）
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, slice[i])
		}
	}

	// 确认原栈没有被修改
	if stack.Size() != len(expected) {
		t.Errorf("Expected stack size %d, got %d", len(expected), stack.Size())
	}

	// 修改切片不应影响原栈
	slice[0] = 100
	val, _ := stack.Pop()
	if val == 100 {
		t.Error("Modifying slice should not affect the original stack")
	}
}

func TestArrayStack_Clear(t *testing.T) {
	stack := New[int]()

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = stack.Push(i)
	}

	if stack.Size() != 5 {
		t.Errorf("Expected size 5, got %d", stack.Size())
	}

	// 清空栈
	stack.Clear()

	if !stack.IsEmpty() {
		t.Error("Expected stack to be empty after Clear()")
	}
	if stack.Size() != 0 {
		t.Errorf("Expected size 0, got %d", stack.Size())
	}

	// 从空栈弹出应该返回错误
	_, err := stack.Pop()
	if err != ErrEmptyStack {
		t.Errorf("Expected ErrEmptyStack, got %v", err)
	}
}

func TestArrayStack_Contains(t *testing.T) {
	stack := New[int]()

	// 空栈不应包含任何元素
	if stack.Contains(1) {
		t.Error("Empty stack should not contain any elements")
	}

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = stack.Push(i)
	}

	// 检查包含的元素
	for i := 0; i < 5; i++ {
		if !stack.Contains(i) {
			t.Errorf("Stack should contain element %d", i)
		}
	}

	// 检查不包含的元素
	if stack.Contains(10) {
		t.Error("Stack should not contain element 10")
	}
}

func TestArrayStack_ForEach(t *testing.T) {
	stack := New[int]()

	// 添加元素
	for i := 0; i < 5; i++ {
		_ = stack.Push(i)
	}

	// 使用ForEach遍历元素
	sum := 0
	stack.ForEach(func(e int) {
		sum += e
	})

	// 验证结果
	expectedSum := 0 + 1 + 2 + 3 + 4
	if sum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, sum)
	}

	// 确认原栈没有被修改
	if stack.Size() != 5 {
		t.Errorf("Expected stack size 5, got %d", stack.Size())
	}
}

func TestArrayStack_String(t *testing.T) {
	stack := New[int]()

	// 空栈
	if stack.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", stack.String())
	}

	// 添加一个元素
	_ = stack.Push(1)
	if stack.String() != "[1]" {
		t.Errorf("Expected '[1]', got '%s'", stack.String())
	}

	// 添加更多元素
	_ = stack.Push(2)
	_ = stack.Push(3)
	if stack.String() != "[1, 2, 3]" {
		t.Errorf("Expected '[1, 2, 3]', got '%s'", stack.String())
	}
}
