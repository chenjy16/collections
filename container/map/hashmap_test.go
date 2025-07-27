package maps

import (
	"fmt"
	"testing"
)

func TestLinkedHashMapBasicOperations(t *testing.T) {
	// 创建一个新的LinkedHashMap
	m := NewLinkedHashMap[string, int]()

	// 测试Put和Get
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	val, found := m.Get("one")
	if !found || val != 1 {
		t.Errorf("Get(\"one\") = %v, %v; want 1, true", val, found)
	}

	val, found = m.Get("two")
	if !found || val != 2 {
		t.Errorf("Get(\"two\") = %v, %v; want 2, true", val, found)
	}

	val, found = m.Get("three")
	if !found || val != 3 {
		t.Errorf("Get(\"three\") = %v, %v; want 3, true", val, found)
	}

	// 测试不存在的键
	val, found = m.Get("four")
	if found {
		t.Errorf("Get(\"four\") = %v, %v; want 0, false", val, found)
	}

	// 测试Size
	if size := m.Size(); size != 3 {
		t.Errorf("Size() = %v; want 3", size)
	}

	// 测试ContainsKey
	if !m.ContainsKey("one") {
		t.Errorf("ContainsKey(\"one\") = false; want true")
	}

	if m.ContainsKey("four") {
		t.Errorf("ContainsKey(\"four\") = true; want false")
	}

	// 测试Remove
	val, found = m.Remove("two")
	if !found || val != 2 {
		t.Errorf("Remove(\"two\") = %v, %v; want 2, true", val, found)
	}

	// 确认已删除
	if m.ContainsKey("two") {
		t.Errorf("After Remove, ContainsKey(\"two\") = true; want false")
	}

	// 测试Size更新
	if size := m.Size(); size != 2 {
		t.Errorf("After Remove, Size() = %v; want 2", size)
	}

	// 测试Clear
	m.Clear()
	if !m.IsEmpty() {
		t.Errorf("After Clear, IsEmpty() = false; want true")
	}

	if size := m.Size(); size != 0 {
		t.Errorf("After Clear, Size() = %v; want 0", size)
	}
}

func TestLinkedHashMapCollisionHandling(t *testing.T) {
	// 创建一个小容量的LinkedHashMap，以便更容易触发冲突
	m := NewLinkedHashMapWithCapacity[string, int](4)

	// 添加足够多的元素以触发冲突和树化
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Put(key, i)
	}

	// 验证所有元素都能正确获取
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key%d", i)
		val, found := m.Get(key)
		if !found || val != i {
			t.Errorf("Get(%q) = %v, %v; want %d, true", key, val, found, i)
		}
	}

	// 测试删除一些元素
	for i := 0; i < 10; i += 2 {
		key := fmt.Sprintf("key%d", i)
		val, found := m.Remove(key)
		if !found || val != i {
			t.Errorf("Remove(%q) = %v, %v; want %d, true", key, val, found, i)
		}
	}

	// 验证删除的元素不存在，未删除的元素存在
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		_, found := m.Get(key)
		expected := i%2 != 0 // 奇数索引的元素应该存在
		if found != expected {
			t.Errorf("After removal, Get(%q) found = %v; want %v", key, found, expected)
		}
	}

	// 测试Size
	expectedSize := 15 // 20 - 5(删除的元素)
	if size := m.Size(); size != expectedSize {
		t.Errorf("After removal, Size() = %v; want %v", size, expectedSize)
	}
}

func TestLinkedHashMapResizing(t *testing.T) {
	// 创建一个小容量的LinkedHashMap
	m := NewLinkedHashMapWithCapacity[int, string](4)

	// 添加足够多的元素以触发扩容
	for i := 0; i < 100; i++ {
		m.Put(i, fmt.Sprintf("value%d", i))
	}

	// 验证所有元素都能正确获取
	for i := 0; i < 100; i++ {
		val, found := m.Get(i)
		expected := fmt.Sprintf("value%d", i)
		if !found || val != expected {
			t.Errorf("Get(%d) = %v, %v; want %q, true", i, val, found, expected)
		}
	}

	// 测试Size
	if size := m.Size(); size != 100 {
		t.Errorf("Size() = %v; want 100", size)
	}
}

func TestLinkedHashMapKeysValuesEntries(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// 添加一些元素
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// 测试Keys
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("len(Keys()) = %v; want 3", len(keys))
	}

	// 检查所有键是否存在
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	for _, k := range []string{"one", "two", "three"} {
		if !keyMap[k] {
			t.Errorf("Keys() does not contain %q", k)
		}
	}

	// 测试Values
	values := m.Values()
	if len(values) != 3 {
		t.Errorf("len(Values()) = %v; want 3", len(values))
	}

	// 检查所有值是否存在
	valueMap := make(map[int]bool)
	for _, v := range values {
		valueMap[v] = true
	}

	for _, v := range []int{1, 2, 3} {
		if !valueMap[v] {
			t.Errorf("Values() does not contain %d", v)
		}
	}

	// 测试Entries
	entries := m.Entries()
	if len(entries) != 3 {
		t.Errorf("len(Entries()) = %v; want 3", len(entries))
	}

	// 检查所有键值对是否存在
	entryMap := make(map[string]int)
	for _, e := range entries {
		entryMap[e.Key] = e.Value
	}

	expectedEntries := map[string]int{"one": 1, "two": 2, "three": 3}
	for k, v := range expectedEntries {
		if entryMap[k] != v {
			t.Errorf("Entries() does not contain %q=%d", k, v)
		}
	}
}

func TestLinkedHashMapForEach(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// 添加一些元素
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// 使用ForEach收集键值对
	collected := make(map[string]int)
	m.ForEach(func(k string, v int) {
		collected[k] = v
	})

	// 验证收集的键值对
	expected := map[string]int{"one": 1, "two": 2, "three": 3}
	for k, v := range expected {
		if collected[k] != v {
			t.Errorf("ForEach collected %q=%d; want %d", k, collected[k], v)
		}
	}

	if len(collected) != len(expected) {
		t.Errorf("ForEach collected %d entries; want %d", len(collected), len(expected))
	}
}

func TestLinkedHashMapPutAll(t *testing.T) {
	m1 := NewLinkedHashMap[string, int]()
	m1.Put("one", 1)
	m1.Put("two", 2)

	m2 := NewLinkedHashMap[string, int]()
	m2.Put("three", 3)
	m2.Put("four", 4)

	// 测试PutAll
	m1.PutAll(m2)

	// 验证m1包含所有键值对
	expected := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
	for k, v := range expected {
		val, found := m1.Get(k)
		if !found || val != v {
			t.Errorf("After PutAll, Get(%q) = %v, %v; want %d, true", k, val, found, v)
		}
	}

	// 测试Size
	if size := m1.Size(); size != 4 {
		t.Errorf("After PutAll, Size() = %v; want 4", size)
	}
}

func TestLinkedHashMapString(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// 测试空映射的字符串表示
	if s := m.String(); s != "{}" {
		t.Errorf("Empty map String() = %q; want \"{}\"", s)
	}

	// 添加一些元素
	m.Put("one", 1)
	m.Put("two", 2)

	// 测试非空映射的字符串表示
	s := m.String()
	// 由于映射的迭代顺序不确定，我们只检查字符串包含所有键值对
	for _, pair := range []string{"one=1", "two=2"} {
		if !contains(s, pair) {
			t.Errorf("String() = %q; should contain %q", s, pair)
		}
	}
}

// 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkLinkedHashMapPut(b *testing.B) {
	m := NewLinkedHashMap[int, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkLinkedHashMapGet(b *testing.B) {
	m := NewLinkedHashMap[int, int]()

	// 预先填充映射
	for i := 0; i < 1000; i++ {
		m.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 循环访问已有的键
		m.Get(i % 1000)
	}
}

func BenchmarkLinkedHashMapRemove(b *testing.B) {
	m := NewLinkedHashMap[int, int]()

	// 预先填充映射
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(i)
	}
}
