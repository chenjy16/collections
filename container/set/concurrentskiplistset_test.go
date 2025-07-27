package set

import (
	"sort"
	"sync"
	"testing"
)

func TestConcurrentSkipListSet_Add(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 测试添加元素
	if !set.Add(5) {
		t.Error("Add should return true for new element")
	}
	if !set.Add(3) {
		t.Error("Add should return true for new element")
	}
	if !set.Add(7) {
		t.Error("Add should return true for new element")
	}

	// 测试重复添加
	if set.Add(5) {
		t.Error("Add should return false for duplicate element")
	}

	// 验证大小
	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	// 验证元素存在
	if !set.Contains(3) || !set.Contains(5) || !set.Contains(7) {
		t.Error("Set should contain all added elements")
	}
}

func TestConcurrentSkipListSet_Remove(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 添加元素
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	set.Add(5)

	// 测试移除存在的元素
	if !set.Remove(3) {
		t.Error("Remove should return true for existing element")
	}

	// 测试移除不存在的元素
	if set.Remove(10) {
		t.Error("Remove should return false for non-existing element")
	}

	// 验证大小
	if set.Size() != 4 {
		t.Errorf("Expected size 4, got %d", set.Size())
	}

	// 验证元素不存在
	if set.Contains(3) {
		t.Error("Set should not contain removed element")
	}

	// 验证其他元素仍存在
	if !set.Contains(1) || !set.Contains(2) || !set.Contains(4) || !set.Contains(5) {
		t.Error("Set should contain non-removed elements")
	}
}

func TestConcurrentSkipListSet_Contains(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 测试空集合
	if set.Contains(1) {
		t.Error("Empty set should not contain any element")
	}

	// 添加元素
	set.Add(10)
	set.Add(20)
	set.Add(30)

	// 测试包含的元素
	if !set.Contains(20) {
		t.Error("Set should contain added element")
	}

	// 测试不包含的元素
	if set.Contains(25) {
		t.Error("Set should not contain non-added element")
	}
}

func TestConcurrentSkipListSet_Clear(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 添加元素
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// 清空集合
	set.Clear()

	// 验证集合为空
	if !set.IsEmpty() {
		t.Error("Set should be empty after Clear")
	}

	if set.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", set.Size())
	}

	// 验证元素不存在
	if set.Contains(1) || set.Contains(2) || set.Contains(3) {
		t.Error("Set should not contain any elements after Clear")
	}
}

func TestConcurrentSkipListSet_ToSlice(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 测试空集合
	slice := set.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice for empty set, got %v", slice)
	}

	// 添加元素（无序）
	elements := []int{5, 2, 8, 1, 9, 3}
	for _, e := range elements {
		set.Add(e)
	}

	// 获取切片
	slice = set.ToSlice()

	// 验证长度
	if len(slice) != len(elements) {
		t.Errorf("Expected slice length %d, got %d", len(elements), len(slice))
	}

	// 验证元素是否有序
	if !sort.IntsAreSorted(slice) {
		t.Errorf("Slice should be sorted, got %v", slice)
	}

	// 验证所有元素都存在
	expected := make([]int, len(elements))
	copy(expected, elements)
	sort.Ints(expected)

	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, slice[i])
		}
	}
}

func TestConcurrentSkipListSet_ForEach(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 添加元素
	elements := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, e := range elements {
		set.Add(e)
	}

	// 测试ForEach
	var result []int
	set.ForEach(func(e int) {
		result = append(result, e)
	})

	// 验证结果是有序的
	if !sort.IntsAreSorted(result) {
		t.Errorf("ForEach should iterate in sorted order, got %v", result)
	}

	// 验证没有重复元素
	uniqueElements := make(map[int]bool)
	for _, e := range elements {
		uniqueElements[e] = true
	}

	if len(result) != len(uniqueElements) {
		t.Errorf("Expected %d unique elements, got %d", len(uniqueElements), len(result))
	}
}

func TestConcurrentSkipListSet_SetOperations(t *testing.T) {
	set1 := NewConcurrentSkipListSet[int]()
	set2 := NewConcurrentSkipListSet[int]()

	// 初始化集合
	for _, e := range []int{1, 2, 3, 4, 5} {
		set1.Add(e)
	}
	for _, e := range []int{4, 5, 6, 7, 8} {
		set2.Add(e)
	}

	// 测试AddAll
	set3 := NewConcurrentSkipListSet[int]()
	set3.AddAll(set1)
	if set3.Size() != 5 {
		t.Errorf("Expected size 5 after AddAll, got %d", set3.Size())
	}

	// 测试ContainsAll
	if !set1.ContainsAll(NewConcurrentSkipListSetFromSlice([]int{1, 3, 5})) {
		t.Error("set1 should contain all elements {1, 3, 5}")
	}

	if set1.ContainsAll(set2) {
		t.Error("set1 should not contain all elements of set2")
	}

	// 测试RemoveAll
	set4 := NewConcurrentSkipListSet[int]()
	set4.AddAll(set1)
	set4.RemoveAll(set2)
	expected := []int{1, 2, 3}
	slice := set4.ToSlice()
	if len(slice) != len(expected) {
		t.Errorf("Expected %d elements after RemoveAll, got %d", len(expected), len(slice))
	}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d after RemoveAll, got %d", v, i, slice[i])
		}
	}

	// 测试RetainAll
	set5 := NewConcurrentSkipListSet[int]()
	set5.AddAll(set1)
	set5.RetainAll(set2)
	expected = []int{4, 5}
	slice = set5.ToSlice()
	if len(slice) != len(expected) {
		t.Errorf("Expected %d elements after RetainAll, got %d", len(expected), len(slice))
	}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected %d at index %d after RetainAll, got %d", v, i, slice[i])
		}
	}
}

func TestConcurrentSkipListSet_NavigationMethods(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 测试空集合
	if _, ok := set.First(); ok {
		t.Error("First should return false for empty set")
	}
	if _, ok := set.Last(); ok {
		t.Error("Last should return false for empty set")
	}

	// 添加元素
	elements := []int{5, 2, 8, 1, 9, 3, 7}
	for _, e := range elements {
		set.Add(e)
	}

	// 测试First和Last
	if first, ok := set.First(); !ok || first != 1 {
		t.Errorf("Expected first element 1, got %d", first)
	}
	if last, ok := set.Last(); !ok || last != 9 {
		t.Errorf("Expected last element 9, got %d", last)
	}

	// 测试Lower
	if lower, ok := set.Lower(5); !ok || lower != 3 {
		t.Errorf("Expected lower(5) = 3, got %d", lower)
	}
	if _, ok := set.Lower(1); ok {
		t.Error("Lower(1) should return false (no element lower than 1)")
	}

	// 测试Higher
	if higher, ok := set.Higher(5); !ok || higher != 7 {
		t.Errorf("Expected higher(5) = 7, got %d", higher)
	}
	if _, ok := set.Higher(9); ok {
		t.Error("Higher(9) should return false (no element higher than 9)")
	}

	// 测试Floor
	if floor, ok := set.Floor(5); !ok || floor != 5 {
		t.Errorf("Expected floor(5) = 5, got %d", floor)
	}
	if floor, ok := set.Floor(4); !ok || floor != 3 {
		t.Errorf("Expected floor(4) = 3, got %d", floor)
	}
	if _, ok := set.Floor(0); ok {
		t.Error("Floor(0) should return false")
	}

	// 测试Ceiling
	if ceiling, ok := set.Ceiling(5); !ok || ceiling != 5 {
		t.Errorf("Expected ceiling(5) = 5, got %d", ceiling)
	}
	if ceiling, ok := set.Ceiling(4); !ok || ceiling != 5 {
		t.Errorf("Expected ceiling(4) = 5, got %d", ceiling)
	}
	if _, ok := set.Ceiling(10); ok {
		t.Error("Ceiling(10) should return false")
	}
}

func TestConcurrentSkipListSet_PollMethods(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 测试空集合
	if _, ok := set.PollFirst(); ok {
		t.Error("PollFirst should return false for empty set")
	}
	if _, ok := set.PollLast(); ok {
		t.Error("PollLast should return false for empty set")
	}

	// 添加元素
	elements := []int{5, 2, 8, 1, 9, 3}
	for _, e := range elements {
		set.Add(e)
	}

	originalSize := set.Size()

	// 测试PollFirst
	if first, ok := set.PollFirst(); !ok || first != 1 {
		t.Errorf("Expected PollFirst() = 1, got %d", first)
	}
	if set.Size() != originalSize-1 {
		t.Errorf("Expected size %d after PollFirst, got %d", originalSize-1, set.Size())
	}
	if set.Contains(1) {
		t.Error("Set should not contain polled first element")
	}

	// 测试PollLast
	if last, ok := set.PollLast(); !ok || last != 9 {
		t.Errorf("Expected PollLast() = 9, got %d", last)
	}
	if set.Size() != originalSize-2 {
		t.Errorf("Expected size %d after PollLast, got %d", originalSize-2, set.Size())
	}
	if set.Contains(9) {
		t.Error("Set should not contain polled last element")
	}
}

func TestConcurrentSkipListSet_Concurrency(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()
	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup

	// 并发添加
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				set.Add(start*numOperations + j)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				set.Contains(j)
				set.Size()
			}
		}()
	}

	wg.Wait()

	// 验证最终状态
	expectedSize := numGoroutines * numOperations
	if set.Size() != expectedSize {
		t.Errorf("Expected size %d after concurrent operations, got %d", expectedSize, set.Size())
	}

	// 验证元素有序性
	slice := set.ToSlice()
	if !sort.IntsAreSorted(slice) {
		t.Error("Elements should be sorted after concurrent operations")
	}
}

func TestConcurrentSkipListSet_WithComparator(t *testing.T) {
	// 创建逆序比较器
	reverseComparator := func(a, b int) int {
		if a < b {
			return 1
		} else if a > b {
			return -1
		}
		return 0
	}

	set := NewConcurrentSkipListSetWithComparator(reverseComparator)

	// 添加元素
	elements := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, e := range elements {
		set.Add(e)
	}

	// 验证逆序
	slice := set.ToSlice()
	for i := 1; i < len(slice); i++ {
		if slice[i-1] < slice[i] {
			t.Errorf("Elements should be in reverse order, got %v", slice)
			break
		}
	}

	// 验证First和Last
	if first, ok := set.First(); !ok || first != 9 {
		t.Errorf("Expected first element 9 (max), got %d", first)
	}
	if last, ok := set.Last(); !ok || last != 1 {
		t.Errorf("Expected last element 1 (min), got %d", last)
	}
}

func TestConcurrentSkipListSet_String(t *testing.T) {
	set := NewConcurrentSkipListSet[int]()

	// 测试空集合
	if set.String() != "[]" {
		t.Errorf("Expected \"[]\" for empty set, got %s", set.String())
	}

	// 添加元素
	set.Add(3)
	set.Add(1)
	set.Add(2)

	// 测试非空集合
	expected := "[1, 2, 3]"
	if set.String() != expected {
		t.Errorf("Expected %s, got %s", expected, set.String())
	}
}
