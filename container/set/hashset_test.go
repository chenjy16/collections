package set

import (
	"fmt"
	"strings"
	"testing"
)

func TestHashSet_New(t *testing.T) {
	set := New[int]()
	if set == nil {
		t.Error("New should return a non-nil HashSet")
	}
	if !set.IsEmpty() {
		t.Error("New HashSet should be empty")
	}
	if set.Size() != 0 {
		t.Errorf("Expected size 0, got %d", set.Size())
	}
}

func TestHashSet_FromSlice(t *testing.T) {
	// 测试空切片
	emptySlice := []int{}
	set := FromSlice(emptySlice)
	if !set.IsEmpty() {
		t.Error("HashSet from empty slice should be empty")
	}

	// 测试非空切片
	slice := []int{1, 2, 3, 2, 1} // 包含重复元素
	set = FromSlice(slice)
	if set.Size() != 3 { // 应该只有3个唯一元素
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	// 验证元素存在
	for _, e := range []int{1, 2, 3} {
		if !set.Contains(e) {
			t.Errorf("Set should contain %d", e)
		}
	}
}

func TestHashSet_Add(t *testing.T) {
	set := New[int]()

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

func TestHashSet_Remove(t *testing.T) {
	set := New[int]()

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

func TestHashSet_Contains(t *testing.T) {
	set := New[int]()

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

func TestHashSet_Clear(t *testing.T) {
	set := New[int]()

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

func TestHashSet_ToSlice(t *testing.T) {
	set := New[int]()

	// 测试空集合
	slice := set.ToSlice()
	if len(slice) != 0 {
		t.Errorf("Expected empty slice for empty set, got %v", slice)
	}

	// 添加元素
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

	// 验证所有元素都存在
	sliceMap := make(map[int]bool)
	for _, v := range slice {
		sliceMap[v] = true
	}

	for _, e := range elements {
		if !sliceMap[e] {
			t.Errorf("Element %d should be in the slice", e)
		}
	}
}

func TestHashSet_ForEach(t *testing.T) {
	set := New[int]()

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

	// 验证没有重复元素
	uniqueElements := make(map[int]bool)
	for _, e := range elements {
		uniqueElements[e] = true
	}

	if len(result) != len(uniqueElements) {
		t.Errorf("Expected %d unique elements, got %d", len(uniqueElements), len(result))
	}

	// 验证所有元素都被遍历到
	resultMap := make(map[int]bool)
	for _, v := range result {
		resultMap[v] = true
	}

	for e := range uniqueElements {
		if !resultMap[e] {
			t.Errorf("Element %d should be visited by ForEach", e)
		}
	}
}

func TestHashSet_SetOperations(t *testing.T) {
	set1 := New[int]()
	set2 := New[int]()

	// 初始化集合
	for _, e := range []int{1, 2, 3, 4, 5} {
		set1.Add(e)
	}

	for _, e := range []int{3, 4, 5, 6, 7} {
		set2.Add(e)
	}

	// 测试并集
	union := set1.Union(set2)
	expectedUnion := []int{1, 2, 3, 4, 5, 6, 7}
	if union.Size() != len(expectedUnion) {
		t.Errorf("Expected union size %d, got %d", len(expectedUnion), union.Size())
	}

	for _, e := range expectedUnion {
		if !union.Contains(e) {
			t.Errorf("Union should contain %d", e)
		}
	}

	// 测试交集
	intersection := set1.Intersection(set2)
	expectedIntersection := []int{3, 4, 5}
	if intersection.Size() != len(expectedIntersection) {
		t.Errorf("Expected intersection size %d, got %d", len(expectedIntersection), intersection.Size())
	}

	for _, e := range expectedIntersection {
		if !intersection.Contains(e) {
			t.Errorf("Intersection should contain %d", e)
		}
	}

	// 测试差集
	difference := set1.Difference(set2)
	expectedDifference := []int{1, 2}
	if difference.Size() != len(expectedDifference) {
		t.Errorf("Expected difference size %d, got %d", len(expectedDifference), difference.Size())
	}

	for _, e := range expectedDifference {
		if !difference.Contains(e) {
			t.Errorf("Difference should contain %d", e)
		}
	}

	// 测试子集关系
	subset := New[int]()
	for _, e := range []int{1, 2} {
		subset.Add(e)
	}

	if !subset.IsSubsetOf(set1) {
		t.Error("subset should be a subset of set1")
	}

	if subset.IsSubsetOf(set2) {
		t.Error("subset should not be a subset of set2")
	}

	// 测试超集关系
	if !set1.IsSupersetOf(subset) {
		t.Error("set1 should be a superset of subset")
	}

	if set2.IsSupersetOf(subset) {
		t.Error("set2 should not be a superset of subset")
	}

	// 测试空集
	emptySet := New[int]()
	if !emptySet.IsSubsetOf(set1) {
		t.Error("Empty set should be a subset of any set")
	}

	if emptySet.IsSupersetOf(set1) {
		t.Error("Empty set should not be a superset of non-empty set")
	}

	if !set1.IsSupersetOf(emptySet) {
		t.Error("Any set should be a superset of empty set")
	}
}

func TestHashSet_String(t *testing.T) {
	set := New[int]()

	// 测试空集合
	if set.String() != "[]" {
		t.Errorf("Expected \"[]\" for empty set, got %s", set.String())
	}

	// 添加元素
	set.Add(3)
	set.Add(1)
	set.Add(2)

	// 测试非空集合
	// 由于哈希集合的顺序不确定，我们需要检查字符串包含所有元素
	str := set.String()
	for _, e := range []int{1, 2, 3} {
		if !strings.Contains(str, fmt.Sprintf("%d", e)) {
			t.Errorf("String representation should contain %d", e)
		}
	}
}

func TestHashSet_Iterator(t *testing.T) {
	set := New[int]()

	// 测试空集合的迭代器
	it := set.Iterator()
	if it.HasNext() {
		t.Error("Iterator of empty set should not have next element")
	}

	// 添加元素
	set.Add(10)
	set.Add(20)
	set.Add(30)

	// 测试迭代
	it = set.Iterator()
	expected := []int{10, 20, 30}
	visited := make(map[int]bool)

	for it.HasNext() {
		val, ok := it.Next()
		if !ok {
			t.Error("Next should return true when HasNext is true")
		}
		visited[val] = true
	}

	for _, e := range expected {
		if !visited[e] {
			t.Errorf("Iterator should visit element %d", e)
		}
	}

	// 测试迭代器的Remove方法
	set.Clear()
	for _, e := range expected {
		set.Add(e)
	}

	it = set.Iterator()
	it.Next() // 移动到第一个元素
	if !it.Remove() {
		t.Error("Remove should return true after Next")
	}

	// 验证元素被移除
	if set.Size() != 2 {
		t.Errorf("Expected size 2 after iterator remove, got %d", set.Size())
	}

	// 测试在调用Next之前调用Remove
	it = set.Iterator()
	if it.Remove() {
		t.Error("Remove should return false before Next")
	}
}
