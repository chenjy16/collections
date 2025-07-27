package list

import (
	"testing"
)

func TestArrayList_New(t *testing.T) {
	list := New[int]()
	if list == nil {
		t.Error("New should return a non-nil ArrayList")
	}
	if !list.IsEmpty() {
		t.Error("New ArrayList should be empty")
	}
	if list.Size() != 0 {
		t.Errorf("Expected size 0, got %d", list.Size())
	}
}

func TestArrayList_WithCapacity(t *testing.T) {
	// 测试正常容量
	list := WithCapacity[int](10)
	if list == nil {
		t.Error("WithCapacity should return a non-nil ArrayList")
	}
	if !list.IsEmpty() {
		t.Error("New ArrayList should be empty regardless of capacity")
	}

	// 测试负容量（应该被处理为0）
	list = WithCapacity[int](-5)
	if list == nil {
		t.Error("WithCapacity should handle negative capacity")
	}
	if !list.IsEmpty() {
		t.Error("New ArrayList should be empty")
	}
}

func TestArrayList_FromSlice(t *testing.T) {
	// 测试空切片
	emptySlice := []int{}
	list := FromSlice(emptySlice)
	if !list.IsEmpty() {
		t.Error("ArrayList from empty slice should be empty")
	}

	// 测试非空切片
	slice := []int{1, 2, 3}
	list = FromSlice(slice)
	if list.Size() != len(slice) {
		t.Errorf("Expected size %d, got %d", len(slice), list.Size())
	}

	// 验证元素顺序
	for i, expected := range slice {
		if val, err := list.Get(i); err != nil || val != expected {
			t.Errorf("Expected %d at index %d, got %d with error: %v", expected, i, val, err)
		}
	}

	// 验证是否是副本（修改原切片不应影响列表）
	slice[0] = 100
	if val, _ := list.Get(0); val == 100 {
		t.Error("FromSlice should create a copy, not reference the original slice")
	}
}

func TestArrayList_Add(t *testing.T) {
	list := New[int]()

	// 测试添加元素
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.Size() != 3 {
		t.Errorf("Expected size 3, got %d", list.Size())
	}

	// 验证元素顺序
	expected := []int{1, 2, 3}
	for i, v := range expected {
		if val, err := list.Get(i); err != nil || val != v {
			t.Errorf("Expected %d at index %d, got %d with error: %v", v, i, val, err)
		}
	}
}

func TestArrayList_Insert(t *testing.T) {
	list := New[int]()

	// 测试在指定位置添加元素
	list.Add(1)
	list.Add(3)

	// 在中间插入
	err := list.Insert(1, 2)
	if err != nil {
		t.Errorf("Insert at valid index should not return error: %v", err)
	}

	// 在开头插入
	err = list.Insert(0, 0)
	if err != nil {
		t.Errorf("Insert at beginning should not return error: %v", err)
	}

	// 在末尾插入
	err = list.Insert(list.Size(), 4)
	if err != nil {
		t.Errorf("Insert at end should not return error: %v", err)
	}

	// 验证元素顺序
	expected := []int{0, 1, 2, 3, 4}
	for i, v := range expected {
		if val, err := list.Get(i); err != nil || val != v {
			t.Errorf("Expected %d at index %d, got %d with error: %v", v, i, val, err)
		}
	}

	// 测试无效索引
	if err := list.Insert(-1, 99); err == nil {
		t.Error("Insert should return error for negative index")
	}

	if err := list.Insert(list.Size()+1, 99); err == nil {
		t.Error("Insert should return error for index > size")
	}
}

func TestArrayList_Get(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if _, err := list.Get(0); err == nil {
		t.Error("Get should return error for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试有效索引
	if val, err := list.Get(1); err != nil || val != 20 {
		t.Errorf("Expected 20 at index 1, got %d with error: %v", val, err)
	}

	// 测试无效索引
	if _, err := list.Get(-1); err == nil {
		t.Error("Get should return error for negative index")
	}

	if _, err := list.Get(list.Size()); err == nil {
		t.Error("Get should return error for index >= size")
	}
}

func TestArrayList_Set(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if _, err := list.Set(0, 100); err == nil {
		t.Error("Set should return error for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试有效索引
	if oldVal, err := list.Set(1, 25); err != nil || oldVal != 20 {
		t.Errorf("Expected old value 20, got %d with error: %v", oldVal, err)
	}

	// 验证更新后的值
	if val, err := list.Get(1); err != nil || val != 25 {
		t.Errorf("Expected 25 at index 1 after Set, got %d with error: %v", val, err)
	}

	// 测试无效索引
	if _, err := list.Set(-1, 100); err == nil {
		t.Error("Set should return error for negative index")
	}

	if _, err := list.Set(list.Size(), 100); err == nil {
		t.Error("Set should return error for index >= size")
	}
}

func TestArrayList_RemoveAt(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if _, err := list.RemoveAt(0); err == nil {
		t.Error("RemoveAt should return error for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试移除中间元素
	if val, err := list.RemoveAt(1); err != nil || val != 20 {
		t.Errorf("Expected removed value 20, got %d with error: %v", val, err)
	}

	// 验证移除后的列表
	expected := []int{10, 30}
	for i, v := range expected {
		if val, err := list.Get(i); err != nil || val != v {
			t.Errorf("Expected %d at index %d after RemoveAt, got %d with error: %v", v, i, val, err)
		}
	}

	// 测试无效索引
	if _, err := list.RemoveAt(-1); err == nil {
		t.Error("RemoveAt should return error for negative index")
	}

	if _, err := list.RemoveAt(list.Size()); err == nil {
		t.Error("RemoveAt should return error for index >= size")
	}
}

func TestArrayList_Remove(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if list.Remove(10) {
		t.Error("Remove should return false for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)
	list.Add(20) // 重复元素

	// 测试移除存在的元素
	if !list.Remove(20) {
		t.Error("Remove should return true for existing element")
	}

	// 验证只移除了第一个匹配的元素
	expected := []int{10, 30, 20}
	for i, v := range expected {
		if val, err := list.Get(i); err != nil || val != v {
			t.Errorf("Expected %d at index %d after Remove, got %d with error: %v", v, i, val, err)
		}
	}

	// 测试移除不存在的元素
	if list.Remove(50) {
		t.Error("Remove should return false for non-existing element")
	}
}

func TestArrayList_Contains(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if list.Contains(10) {
		t.Error("Contains should return false for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试包含的元素
	if !list.Contains(20) {
		t.Error("Contains should return true for existing element")
	}

	// 测试不包含的元素
	if list.Contains(50) {
		t.Error("Contains should return false for non-existing element")
	}
}

func TestArrayList_IndexOf(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if list.IndexOf(10) != -1 {
		t.Error("IndexOf should return -1 for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)
	list.Add(20) // 重复元素

	// 测试存在的元素
	if list.IndexOf(20) != 1 {
		t.Errorf("IndexOf should return 1 for first occurrence of 20, got %d", list.IndexOf(20))
	}

	// 测试不存在的元素
	if list.IndexOf(50) != -1 {
		t.Error("IndexOf should return -1 for non-existing element")
	}
}

func TestArrayList_LastIndexOf(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if list.LastIndexOf(10) != -1 {
		t.Error("LastIndexOf should return -1 for empty list")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)
	list.Add(20) // 重复元素

	// 测试存在的元素
	if list.LastIndexOf(20) != 3 {
		t.Errorf("LastIndexOf should return 3 for last occurrence of 20, got %d", list.LastIndexOf(20))
	}

	// 测试不存在的元素
	if list.LastIndexOf(50) != -1 {
		t.Error("LastIndexOf should return -1 for non-existing element")
	}
}

func TestArrayList_Clear(t *testing.T) {
	list := New[int]()

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 清空列表
	list.Clear()

	// 验证列表为空
	if !list.IsEmpty() {
		t.Error("List should be empty after Clear")
	}

	if list.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", list.Size())
	}
}

func TestArrayList_ToSlice(t *testing.T) {
	list := New[int]()

	// 测试空列表
	slice := list.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice for empty list, got %v", slice)
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试非空列表
	slice = list.ToSlice()
	expected := []int{10, 20, 30}

	if len(slice) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(slice))
	}

	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d in slice, got %d", v, i, slice[i])
		}
	}

	// 验证返回的切片是原列表的副本
	slice[0] = 100
	if val, _ := list.Get(0); val == 100 {
		t.Error("ToSlice should return a copy, not a reference to the internal slice")
	}
}

func TestArrayList_SubList(t *testing.T) {
	list := New[int]()

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)
	list.Add(40)
	list.Add(50)

	// 测试有效范围
	subList, err := list.SubList(1, 4)
	if err != nil {
		t.Errorf("SubList with valid range should not return error: %v", err)
	}

	if subList.Size() != 3 {
		t.Errorf("Expected sublist size 3, got %d", subList.Size())
	}

	expected := []int{20, 30, 40}
	for i, v := range expected {
		if val, err := subList.Get(i); err != nil || val != v {
			t.Errorf("Expected %d at index %d in sublist, got %d with error: %v", v, i, val, err)
		}
	}

	// 测试空范围
	emptySubList, err := list.SubList(2, 2)
	if err != nil {
		t.Errorf("SubList with empty range should not return error: %v", err)
	}
	if emptySubList.Size() != 0 {
		t.Errorf("Expected empty sublist, got size %d", emptySubList.Size())
	}

	// 测试无效范围
	_, err = list.SubList(-1, 3)
	if err == nil {
		t.Error("SubList should return error for negative fromIndex")
	}

	_, err = list.SubList(1, 6)
	if err == nil {
		t.Error("SubList should return error for toIndex > size")
	}

	_, err = list.SubList(3, 2)
	if err == nil {
		t.Error("SubList should return error for fromIndex > toIndex")
	}
}

func TestArrayList_ForEach(t *testing.T) {
	list := New[int]()

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试ForEach
	sum := 0
	list.ForEach(func(e int) {
		sum += e
	})

	expectedSum := 60
	if sum != expectedSum {
		t.Errorf("Expected sum %d after ForEach, got %d", expectedSum, sum)
	}
}

func TestArrayList_String(t *testing.T) {
	list := New[int]()

	// 测试空列表
	if list.String() != "[]" {
		t.Errorf("Expected \"[]\" for empty list, got %s", list.String())
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试非空列表
	expected := "[10, 20, 30]"
	if list.String() != expected {
		t.Errorf("Expected %s, got %s", expected, list.String())
	}
}

func TestArrayList_Iterator(t *testing.T) {
	list := New[int]()

	// 测试空列表的迭代器
	it := list.Iterator()
	if it.HasNext() {
		t.Error("Iterator of empty list should not have next element")
	}

	// 添加元素
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// 测试迭代
	it = list.Iterator()
	expected := []int{10, 20, 30}
	index := 0

	for it.HasNext() {
		val, ok := it.Next()
		if !ok {
			t.Error("Next should return true when HasNext is true")
		}
		if val != expected[index] {
			t.Errorf("Expected %d at iteration %d, got %d", expected[index], index, val)
		}
		index++
	}

	if index != len(expected) {
		t.Errorf("Iterator should iterate through all elements, expected %d iterations, got %d", len(expected), index)
	}

	// 测试迭代器的Remove方法
	it = list.Iterator()
	it.Next() // 移动到第一个元素
	if !it.Remove() {
		t.Error("Remove should return true after Next")
	}

	// 验证元素被移除
	if list.Size() != 2 {
		t.Errorf("Expected size 2 after iterator remove, got %d", list.Size())
	}
	if val, _ := list.Get(0); val != 20 {
		t.Errorf("Expected 20 at index 0 after iterator remove, got %d", val)
	}

	// 测试在调用Next之前调用Remove
	it = list.Iterator()
	if it.Remove() {
		t.Error("Remove should return false before Next")
	}

	// 测试在同一个Next之后多次调用Remove
	it.Next()
	if !it.Remove() {
		t.Error("First Remove after Next should return true")
	}
	if it.Remove() {
		t.Error("Second Remove after same Next should return false")
	}
}
