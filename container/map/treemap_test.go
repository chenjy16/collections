package maps

import (
	"strconv"
	"testing"
)

func TestTreeMapNew(t *testing.T) {
	m := NewTreeMap[int, string]()
	if m == nil {
		t.Error("NewTreeMap should not return nil")
	}
	if !m.IsEmpty() {
		t.Error("New TreeMap should be empty")
	}
	if m.Size() != 0 {
		t.Error("New TreeMap size should be 0")
	}
}

func TestTreeMapNewWithComparator(t *testing.T) {
	// 自定义比较器：降序
	comparator := func(a, b int) int {
		if a > b {
			return -1
		} else if a < b {
			return 1
		}
		return 0
	}

	m := NewTreeMapWithComparator[int, string](comparator)
	if m == nil {
		t.Error("NewTreeMapWithComparator should not return nil")
	}

	// 测试自定义比较器是否生效
	m.Put(1, "one")
	m.Put(2, "two")
	m.Put(3, "three")

	keys := m.Keys()
	expected := []int{3, 2, 1} // 降序
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key %d at index %d, got %d", expected[i], i, key)
		}
	}
}

func TestTreeMapPutAndGet(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试插入新键
	oldValue, existed := m.Put(1, "one")
	if existed {
		t.Error("Put should return false for new key")
	}
	if oldValue != "" {
		t.Error("Put should return zero value for new key")
	}

	// 测试获取值
	value, found := m.Get(1)
	if !found {
		t.Error("Get should find existing key")
	}
	if value != "one" {
		t.Errorf("Expected 'one', got '%s'", value)
	}

	// 测试更新现有键
	oldValue, existed = m.Put(1, "ONE")
	if !existed {
		t.Error("Put should return true for existing key")
	}
	if oldValue != "one" {
		t.Errorf("Expected old value 'one', got '%s'", oldValue)
	}

	// 验证值已更新
	value, found = m.Get(1)
	if !found {
		t.Error("Get should find existing key")
	}
	if value != "ONE" {
		t.Errorf("Expected 'ONE', got '%s'", value)
	}
}

func TestTreeMapContainsKey(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	if m.ContainsKey(1) {
		t.Error("Empty map should not contain any key")
	}

	// 添加元素
	m.Put(1, "one")
	m.Put(2, "two")

	// 测试存在的键
	if !m.ContainsKey(1) {
		t.Error("Map should contain key 1")
	}
	if !m.ContainsKey(2) {
		t.Error("Map should contain key 2")
	}

	// 测试不存在的键
	if m.ContainsKey(3) {
		t.Error("Map should not contain key 3")
	}
}

func TestTreeMapRemove(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试移除不存在的键
	value, removed := m.Remove(1)
	if removed {
		t.Error("Remove should return false for non-existing key")
	}
	if value != "" {
		t.Error("Remove should return zero value for non-existing key")
	}

	// 添加元素
	m.Put(1, "one")
	m.Put(2, "two")
	m.Put(3, "three")

	// 测试移除存在的键
	value, removed = m.Remove(2)
	if !removed {
		t.Error("Remove should return true for existing key")
	}
	if value != "two" {
		t.Errorf("Expected 'two', got '%s'", value)
	}

	// 验证键已被移除
	if m.ContainsKey(2) {
		t.Error("Key 2 should be removed")
	}
	if m.Size() != 2 {
		t.Errorf("Expected size 2, got %d", m.Size())
	}

	// 验证其他键仍然存在
	if !m.ContainsKey(1) {
		t.Error("Key 1 should still exist")
	}
	if !m.ContainsKey(3) {
		t.Error("Key 3 should still exist")
	}
}

func TestTreeMapSize(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}

	// 添加元素
	for i := 1; i <= 5; i++ {
		m.Put(i, strconv.Itoa(i))
		if m.Size() != i {
			t.Errorf("Expected size %d, got %d", i, m.Size())
		}
	}

	// 移除元素
	for i := 5; i >= 1; i-- {
		m.Remove(i)
		if m.Size() != i-1 {
			t.Errorf("Expected size %d, got %d", i-1, m.Size())
		}
	}
}

func TestTreeMapIsEmpty(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	if !m.IsEmpty() {
		t.Error("New map should be empty")
	}

	// 添加元素
	m.Put(1, "one")
	if m.IsEmpty() {
		t.Error("Map with elements should not be empty")
	}

	// 移除元素
	m.Remove(1)
	if !m.IsEmpty() {
		t.Error("Map should be empty after removing all elements")
	}
}

func TestTreeMapClear(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 添加元素
	for i := 1; i <= 10; i++ {
		m.Put(i, strconv.Itoa(i))
	}

	if m.Size() != 10 {
		t.Errorf("Expected size 10, got %d", m.Size())
	}

	// 清空映射
	m.Clear()

	if !m.IsEmpty() {
		t.Error("Map should be empty after Clear")
	}
	if m.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", m.Size())
	}

	// 验证所有键都被移除
	for i := 1; i <= 10; i++ {
		if m.ContainsKey(i) {
			t.Errorf("Key %d should not exist after Clear", i)
		}
	}
}

func TestTreeMapKeys(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	keys := m.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(keys))
	}

	// 添加元素（无序插入）
	testData := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	for _, key := range testData {
		m.Put(key, strconv.Itoa(key))
	}

	keys = m.Keys()

	// 验证键的数量
	if len(keys) != len(testData) {
		t.Errorf("Expected %d keys, got %d", len(testData), len(keys))
	}

	// 验证键是有序的
	for i := 0; i < len(keys); i++ {
		if keys[i] != i+1 {
			t.Errorf("Expected key %d at index %d, got %d", i+1, i, keys[i])
		}
	}
}

func TestTreeMapValues(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	values := m.Values()
	if len(values) != 0 {
		t.Errorf("Expected 0 values, got %d", len(values))
	}

	// 添加元素
	testData := map[int]string{3: "three", 1: "one", 2: "two"}
	for key, value := range testData {
		m.Put(key, value)
	}

	values = m.Values()

	// 验证值的数量
	if len(values) != len(testData) {
		t.Errorf("Expected %d values, got %d", len(testData), len(values))
	}

	// 验证值的顺序（应该按键的顺序）
	expected := []string{"one", "two", "three"}
	for i, value := range values {
		if value != expected[i] {
			t.Errorf("Expected value '%s' at index %d, got '%s'", expected[i], i, value)
		}
	}
}

func TestTreeMapEntries(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	entries := m.Entries()
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}

	// 添加元素
	testData := map[int]string{3: "three", 1: "one", 2: "two"}
	for key, value := range testData {
		m.Put(key, value)
	}

	entries = m.Entries()

	// 验证条目的数量
	if len(entries) != len(testData) {
		t.Errorf("Expected %d entries, got %d", len(testData), len(entries))
	}

	// 验证条目的顺序（应该按键的顺序）
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}

	for i, entry := range entries {
		if entry.Key != expectedKeys[i] {
			t.Errorf("Expected key %d at index %d, got %d", expectedKeys[i], i, entry.Key)
		}
		if entry.Value != expectedValues[i] {
			t.Errorf("Expected value '%s' at index %d, got '%s'", expectedValues[i], i, entry.Value)
		}
	}
}

func TestTreeMapForEach(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 添加元素
	testData := map[int]string{3: "three", 1: "one", 2: "two"}
	for key, value := range testData {
		m.Put(key, value)
	}

	// 测试ForEach
	var keys []int
	var values []string

	m.ForEach(func(key int, value string) {
		keys = append(keys, key)
		values = append(values, value)
	})

	// 验证顺序
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}

	for i := range keys {
		if keys[i] != expectedKeys[i] {
			t.Errorf("Expected key %d at index %d, got %d", expectedKeys[i], i, keys[i])
		}
		if values[i] != expectedValues[i] {
			t.Errorf("Expected value '%s' at index %d, got '%s'", expectedValues[i], i, values[i])
		}
	}
}

func TestTreeMapString(t *testing.T) {
	m := NewTreeMap[int, string]()

	// 测试空映射
	str := m.String()
	if str != "{}" {
		t.Errorf("Expected '{}', got '%s'", str)
	}

	// 添加一个元素
	m.Put(1, "one")
	str = m.String()
	expected := "{1=one}"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}

	// 添加多个元素
	m.Put(2, "two")
	m.Put(3, "three")
	str = m.String()
	expected = "{1=one, 2=two, 3=three}"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestTreeMapLargeDataSet(t *testing.T) {
	m := NewTreeMap[int, int]()

	// 插入大量数据
	n := 1000
	for i := 0; i < n; i++ {
		m.Put(i, i*2)
	}

	// 验证大小
	if m.Size() != n {
		t.Errorf("Expected size %d, got %d", n, m.Size())
	}

	// 验证所有数据
	for i := 0; i < n; i++ {
		value, found := m.Get(i)
		if !found {
			t.Errorf("Key %d should exist", i)
		}
		if value != i*2 {
			t.Errorf("Expected value %d for key %d, got %d", i*2, i, value)
		}
	}

	// 验证键的顺序
	keys := m.Keys()
	for i, key := range keys {
		if key != i {
			t.Errorf("Expected key %d at index %d, got %d", i, i, key)
		}
	}
}

// 基准测试
func BenchmarkTreeMapPut(b *testing.B) {
	m := NewTreeMap[int, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}
}

func BenchmarkTreeMapGet(b *testing.B) {
	m := NewTreeMap[int, int]()

	// 预填充数据
	for i := 0; i < 1000; i++ {
		m.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

func BenchmarkTreeMapRemove(b *testing.B) {
	m := NewTreeMap[int, int]()

	// 预填充数据
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(i)
	}
}
