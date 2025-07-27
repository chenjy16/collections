package list

import (
	"testing"
)

func TestLinkedListBasicOperations(t *testing.T) {
	list := NewLinkedList[int]()

	// 测试空列表
	if !list.IsEmpty() {
		t.Error("新创建的列表应该为空")
	}

	if list.Size() != 0 {
		t.Errorf("空列表的大小应该为0，实际为%d", list.Size())
	}

	// 测试添加元素
	if !list.Add(1) {
		t.Error("添加元素应该成功")
	}

	if !list.Add(2) {
		t.Error("添加元素应该成功")
	}

	if !list.Add(3) {
		t.Error("添加元素应该成功")
	}

	if list.Size() != 3 {
		t.Errorf("列表大小应该为3，实际为%d", list.Size())
	}

	if list.IsEmpty() {
		t.Error("非空列表不应该为空")
	}
}

func TestLinkedListGetAndSet(t *testing.T) {
	list := NewLinkedList[string]()
	list.Add("first")
	list.Add("second")
	list.Add("third")

	// 测试Get
	if val, err := list.Get(0); err != nil || val != "first" {
		t.Errorf("Get(0)应该返回'first'，实际返回'%s'，错误：%v", val, err)
	}

	if val, err := list.Get(1); err != nil || val != "second" {
		t.Errorf("Get(1)应该返回'second'，实际返回'%s'，错误：%v", val, err)
	}

	if val, err := list.Get(2); err != nil || val != "third" {
		t.Errorf("Get(2)应该返回'third'，实际返回'%s'，错误：%v", val, err)
	}

	// 测试越界访问
	if _, err := list.Get(-1); err == nil {
		t.Error("Get(-1)应该返回错误")
	}

	if _, err := list.Get(3); err == nil {
		t.Error("Get(3)应该返回错误")
	}

	// 测试Set
	if old, err := list.Set(1, "modified"); err != nil || old != "second" {
		t.Errorf("Set(1, 'modified')应该返回'second'，实际返回'%s'，错误：%v", old, err)
	}

	if val, _ := list.Get(1); val != "modified" {
		t.Errorf("Set后Get(1)应该返回'modified'，实际返回'%s'", val)
	}
}

func TestLinkedListInsertAndRemove(t *testing.T) {
	list := NewLinkedList[int]()

	// 测试Insert
	if err := list.Insert(0, 10); err != nil {
		t.Errorf("在空列表插入应该成功，错误：%v", err)
	}

	if err := list.Insert(0, 5); err != nil {
		t.Errorf("在头部插入应该成功，错误：%v", err)
	}

	if err := list.Insert(2, 15); err != nil {
		t.Errorf("在尾部插入应该成功，错误：%v", err)
	}

	if err := list.Insert(2, 12); err != nil {
		t.Errorf("在中间插入应该成功，错误：%v", err)
	}

	// 验证顺序：[5, 10, 12, 15]
	expected := []int{5, 10, 12, 15}
	for i, exp := range expected {
		if val, _ := list.Get(i); val != exp {
			t.Errorf("索引%d的值应该为%d，实际为%d", i, exp, val)
		}
	}

	// 测试RemoveAt
	if removed, err := list.RemoveAt(1); err != nil || removed != 10 {
		t.Errorf("RemoveAt(1)应该返回10，实际返回%d，错误：%v", removed, err)
	}

	if list.Size() != 3 {
		t.Errorf("移除后列表大小应该为3，实际为%d", list.Size())
	}

	// 测试Remove
	if !list.Remove(12) {
		t.Error("Remove(12)应该成功")
	}

	if list.Remove(100) {
		t.Error("Remove(100)应该失败")
	}
}

func TestLinkedListIndexOf(t *testing.T) {
	list := NewLinkedList[string]()
	list.Add("apple")
	list.Add("banana")
	list.Add("apple")
	list.Add("cherry")

	if index := list.IndexOf("apple"); index != 0 {
		t.Errorf("IndexOf('apple')应该返回0，实际返回%d", index)
	}

	if index := list.LastIndexOf("apple"); index != 2 {
		t.Errorf("LastIndexOf('apple')应该返回2，实际返回%d", index)
	}

	if index := list.IndexOf("grape"); index != -1 {
		t.Errorf("IndexOf('grape')应该返回-1，实际返回%d", index)
	}
}

func TestLinkedListSubList(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 0; i < 5; i++ {
		list.Add(i)
	}

	subList, err := list.SubList(1, 4)
	if err != nil {
		t.Errorf("SubList(1, 4)应该成功，错误：%v", err)
	}

	if subList.Size() != 3 {
		t.Errorf("子列表大小应该为3，实际为%d", subList.Size())
	}

	expected := []int{1, 2, 3}
	for i, exp := range expected {
		if val, _ := subList.Get(i); val != exp {
			t.Errorf("子列表索引%d的值应该为%d，实际为%d", i, exp, val)
		}
	}
}

func TestLinkedListIterator(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 3; i++ {
		list.Add(i)
	}

	iter := list.Iterator()
	values := []int{}

	for iter.HasNext() {
		val, ok := iter.Next()
		if !ok {
			t.Error("Next()应该成功")
		}
		values = append(values, val)
	}

	expected := []int{1, 2, 3}
	if len(values) != len(expected) {
		t.Errorf("迭代器返回的元素数量应该为%d，实际为%d", len(expected), len(values))
	}

	for i, exp := range expected {
		if values[i] != exp {
			t.Errorf("索引%d的值应该为%d，实际为%d", i, exp, values[i])
		}
	}
}

func TestLinkedListFirstLast(t *testing.T) {
	list := NewLinkedList[string]()

	// 测试空列表
	if _, err := list.GetFirst(); err == nil {
		t.Error("空列表GetFirst()应该返回错误")
	}

	if _, err := list.GetLast(); err == nil {
		t.Error("空列表GetLast()应该返回错误")
	}

	// 添加元素
	list.AddFirst("first")
	list.AddLast("last")
	list.AddFirst("new_first")

	// 验证顺序：[new_first, first, last]
	if val, _ := list.GetFirst(); val != "new_first" {
		t.Errorf("GetFirst()应该返回'new_first'，实际返回'%s'", val)
	}

	if val, _ := list.GetLast(); val != "last" {
		t.Errorf("GetLast()应该返回'last'，实际返回'%s'", val)
	}

	// 测试移除
	if removed, err := list.RemoveFirst(); err != nil || removed != "new_first" {
		t.Errorf("RemoveFirst()应该返回'new_first'，实际返回'%s'，错误：%v", removed, err)
	}

	if removed, err := list.RemoveLast(); err != nil || removed != "last" {
		t.Errorf("RemoveLast()应该返回'last'，实际返回'%s'，错误：%v", removed, err)
	}

	if list.Size() != 1 {
		t.Errorf("移除后列表大小应该为1，实际为%d", list.Size())
	}
}

func TestLinkedListFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	list := LinkedListFromSlice(slice)

	if list.Size() != len(slice) {
		t.Errorf("列表大小应该为%d，实际为%d", len(slice), list.Size())
	}

	for i, expected := range slice {
		if val, _ := list.Get(i); val != expected {
			t.Errorf("索引%d的值应该为%d，实际为%d", i, expected, val)
		}
	}
}

func TestLinkedListToSlice(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 5; i++ {
		list.Add(i)
	}

	slice := list.ToSlice()
	if len(slice) != list.Size() {
		t.Errorf("切片长度应该为%d，实际为%d", list.Size(), len(slice))
	}

	for i := 0; i < list.Size(); i++ {
		if val, _ := list.Get(i); slice[i] != val {
			t.Errorf("索引%d的值应该为%d，实际为%d", i, val, slice[i])
		}
	}
}

func TestLinkedListClear(t *testing.T) {
	list := NewLinkedList[int]()
	for i := 1; i <= 5; i++ {
		list.Add(i)
	}

	list.Clear()

	if !list.IsEmpty() {
		t.Error("清空后列表应该为空")
	}

	if list.Size() != 0 {
		t.Errorf("清空后列表大小应该为0，实际为%d", list.Size())
	}
}

func TestLinkedListContains(t *testing.T) {
	list := NewLinkedList[string]()
	list.Add("apple")
	list.Add("banana")
	list.Add("cherry")

	if !list.Contains("banana") {
		t.Error("列表应该包含'banana'")
	}

	if list.Contains("grape") {
		t.Error("列表不应该包含'grape'")
	}
}