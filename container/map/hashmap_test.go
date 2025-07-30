package maps

import (
	"fmt"
	"testing"
)

func TestLinkedHashMapBasicOperations(t *testing.T) {
	// Create a new LinkedHashMap
	m := NewLinkedHashMap[string, int]()

	// Test Put and Get
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

	// Test non-existent key
	val, found = m.Get("four")
	if found {
		t.Errorf("Get(\"four\") = %v, %v; want 0, false", val, found)
	}

	// Test Size
	if size := m.Size(); size != 3 {
		t.Errorf("Size() = %v; want 3", size)
	}

	// Test ContainsKey
	if !m.ContainsKey("one") {
		t.Errorf("ContainsKey(\"one\") = false; want true")
	}

	if m.ContainsKey("four") {
		t.Errorf("ContainsKey(\"four\") = true; want false")
	}

	// Test Remove
	val, found = m.Remove("two")
	if !found || val != 2 {
		t.Errorf("Remove(\"two\") = %v, %v; want 2, true", val, found)
	}

	// Confirm deletion
	if m.ContainsKey("two") {
		t.Errorf("After Remove, ContainsKey(\"two\") = true; want false")
	}

	// Test Size update
	if size := m.Size(); size != 2 {
		t.Errorf("After Remove, Size() = %v; want 2", size)
	}

	// Test Clear
	m.Clear()
	if !m.IsEmpty() {
		t.Errorf("After Clear, IsEmpty() = false; want true")
	}

	if size := m.Size(); size != 0 {
		t.Errorf("After Clear, Size() = %v; want 0", size)
	}
}

func TestLinkedHashMapCollisionHandling(t *testing.T) {
	// Create a small capacity LinkedHashMap to easily trigger collisions
	m := NewLinkedHashMapWithCapacity[string, int](4)

	// Add enough elements to trigger collisions and treeification
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Put(key, i)
	}

	// Verify all elements can be retrieved correctly
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key%d", i)
		val, found := m.Get(key)
		if !found || val != i {
			t.Errorf("Get(%q) = %v, %v; want %d, true", key, val, found, i)
		}
	}

	// Test removing some elements
	for i := 0; i < 10; i += 2 {
		key := fmt.Sprintf("key%d", i)
		val, found := m.Remove(key)
		if !found || val != i {
			t.Errorf("Remove(%q) = %v, %v; want %d, true", key, val, found, i)
		}
	}

	// Verify removed elements don't exist, non-removed elements exist
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		_, found := m.Get(key)
		expected := i%2 != 0 // Elements with odd indices should exist
		if found != expected {
			t.Errorf("After removal, Get(%q) found = %v; want %v", key, found, expected)
		}
	}

	// Test Size
	expectedSize := 15 // 20 - 5(removed elements)
	if size := m.Size(); size != expectedSize {
		t.Errorf("After removal, Size() = %v; want %v", size, expectedSize)
	}
}

func TestLinkedHashMapResizing(t *testing.T) {
	// Create a small capacity LinkedHashMap
	m := NewLinkedHashMapWithCapacity[int, string](4)

	// Add enough elements to trigger resizing
	for i := 0; i < 100; i++ {
		m.Put(i, fmt.Sprintf("value%d", i))
	}

	// Verify all elements can be retrieved correctly
	for i := 0; i < 100; i++ {
		val, found := m.Get(i)
		expected := fmt.Sprintf("value%d", i)
		if !found || val != expected {
			t.Errorf("Get(%d) = %v, %v; want %q, true", i, val, found, expected)
		}
	}

	// Test Size
	if size := m.Size(); size != 100 {
		t.Errorf("Size() = %v; want 100", size)
	}
}

func TestLinkedHashMapKeysValuesEntries(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Add some elements
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// Test Keys
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("len(Keys()) = %v; want 3", len(keys))
	}

	// Check if all keys exist
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	for _, k := range []string{"one", "two", "three"} {
		if !keyMap[k] {
			t.Errorf("Keys() does not contain %q", k)
		}
	}

	// Test Values
	values := m.Values()
	if len(values) != 3 {
		t.Errorf("len(Values()) = %v; want 3", len(values))
	}

	// Check if all values exist
	valueMap := make(map[int]bool)
	for _, v := range values {
		valueMap[v] = true
	}

	for _, v := range []int{1, 2, 3} {
		if !valueMap[v] {
			t.Errorf("Values() does not contain %d", v)
		}
	}

	// Test Entries
	entries := m.Entries()
	if len(entries) != 3 {
		t.Errorf("len(Entries()) = %v; want 3", len(entries))
	}

	// Check if all key-value pairs exist
	entryMap := make(map[string]int)
	for _, e := range entries {
		entryMap[e.Key] = e.Value
	}

	expectedEntries := map[string]int{"one": 1, "two": 2, "three": 3}
	for k, v := range expectedEntries {
		if entryMap[k] != v {
			t.Errorf("Entries() does not contain correct pair %q: %d", k, v)
		}
	}
}

func TestLinkedHashMapForEach(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Add some elements
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// Test ForEach
	visited := make(map[string]int)
	m.ForEach(func(k string, v int) {
		visited[k] = v
	})

	expectedEntries := map[string]int{"one": 1, "two": 2, "three": 3}
	if len(visited) != len(expectedEntries) {
		t.Errorf("ForEach visited %d entries; want %d", len(visited), len(expectedEntries))
	}

	for k, v := range expectedEntries {
		if visited[k] != v {
			t.Errorf("ForEach did not visit correct pair %q: %d", k, v)
		}
	}
}

func TestLinkedHashMapContainsValue(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Add some elements
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// Test ContainsValue
	if !m.ContainsValue(1) {
		t.Errorf("ContainsValue(1) = false; want true")
	}

	if !m.ContainsValue(2) {
		t.Errorf("ContainsValue(2) = false; want true")
	}

	if !m.ContainsValue(3) {
		t.Errorf("ContainsValue(3) = false; want true")
	}

	if m.ContainsValue(4) {
		t.Errorf("ContainsValue(4) = true; want false")
	}
}

func TestLinkedHashMapString(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Test empty map
	str := m.String()
	if str != "{}" {
		t.Errorf("String() for empty map = %q; want \"{}\"", str)
	}

	// Add some elements
	m.Put("one", 1)
	m.Put("two", 2)

	// Test non-empty map
	str = m.String()
	// The exact format may vary, but it should contain the key-value pairs
	if len(str) < 10 { // A reasonable minimum length
		t.Errorf("String() for non-empty map = %q; want longer string", str)
	}
}

func TestLinkedHashMapPutAll(t *testing.T) {
	m1 := NewLinkedHashMap[string, int]()
	m1.Put("one", 1)
	m1.Put("two", 2)

	m2 := NewLinkedHashMap[string, int]()
	m2.Put("three", 3)
	m2.Put("four", 4)

	// Test PutAll
	m1.PutAll(m2)

	// Verify all elements exist in m1
	expectedEntries := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
	if m1.Size() != len(expectedEntries) {
		t.Errorf("After PutAll, Size() = %v; want %v", m1.Size(), len(expectedEntries))
	}

	for k, v := range expectedEntries {
		val, found := m1.Get(k)
		if !found || val != v {
			t.Errorf("After PutAll, Get(%q) = %v, %v; want %d, true", k, val, found, v)
		}
	}
}

func TestLinkedHashMapUpdateExistingKey(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Put initial value
	oldVal, existed := m.Put("key", 1)
	if existed {
		t.Errorf("Put(\"key\", 1) returned existed = true; want false")
	}

	// Update existing key
	oldVal, existed = m.Put("key", 2)
	if !existed || oldVal != 1 {
		t.Errorf("Put(\"key\", 2) = %v, %v; want 1, true", oldVal, existed)
	}

	// Verify updated value
	val, found := m.Get("key")
	if !found || val != 2 {
		t.Errorf("Get(\"key\") = %v, %v; want 2, true", val, found)
	}

	// Verify size didn't change
	if size := m.Size(); size != 1 {
		t.Errorf("Size() = %v; want 1", size)
	}
}
