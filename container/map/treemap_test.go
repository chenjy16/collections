package maps

import (
	"fmt"
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

func TestTreeMapWithComparator(t *testing.T) {
	// Custom comparator: descending order
	comparator := func(a, b int) int {
		if a > b {
			return -1
		} else if a < b {
			return 1
		}
		return 0
	}

	tm := NewTreeMapWithComparator[int, string](comparator)

	tm.Put(1, "one")
	tm.Put(3, "three")
	tm.Put(2, "two")

	// Test if custom comparator is effective
	keys := tm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	expected := []int{3, 2, 1} // Descending order
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key %d at position %d, got %d", expected[i], i, key)
		}
	}
}

func TestTreeMapPut(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test inserting new key
	oldVal, existed := tm.Put(1, "one")
	if existed {
		t.Error("Key should not exist initially")
	}
	if oldVal != "" {
		t.Error("Old value should be empty for new key")
	}

	// Test getting value
	val, ok := tm.Get(1)
	if !ok {
		t.Error("Key should exist after Put")
	}
	if val != "one" {
		t.Errorf("Expected 'one', got '%s'", val)
	}

	// Test updating existing key
	oldVal, existed = tm.Put(1, "ONE")
	if !existed {
		t.Error("Key should exist")
	}
	if oldVal != "one" {
		t.Errorf("Expected old value 'one', got '%s'", oldVal)
	}

	// Verify value is updated
	val, ok = tm.Get(1)
	if !ok || val != "ONE" {
		t.Error("Failed to update existing key")
	}
}

func TestTreeMapGet(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	_, ok := tm.Get(1)
	if ok {
		t.Error("Empty map should not contain any key")
	}

	// Add elements
	tm.Put(1, "one")
	tm.Put(2, "two")

	// Test existing key
	val, ok := tm.Get(1)
	if !ok {
		t.Error("Key should exist")
	}
	if val != "one" {
		t.Errorf("Expected 'one', got '%s'", val)
	}

	// Test non-existing key
	_, ok = tm.Get(3)
	if ok {
		t.Error("Non-existing key should not be found")
	}
}

func TestTreeMapRemove(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test removing non-existing key
	_, existed := tm.Remove(1)
	if existed {
		t.Error("Non-existing key should not be found")
	}

	// Add elements
	tm.Put(1, "one")
	tm.Put(2, "two")
	tm.Put(3, "three")

	// Test removing existing key
	oldVal, existed := tm.Remove(2)
	if !existed {
		t.Error("Key should exist")
	}
	if oldVal != "two" {
		t.Errorf("Expected old value 'two', got '%s'", oldVal)
	}

	// Verify key is removed
	_, ok := tm.Get(2)
	if ok {
		t.Error("Key should not exist after removal")
	}

	// Verify other keys still exist
	val, ok := tm.Get(1)
	if !ok || val != "one" {
		t.Error("Other keys should remain")
	}
	val, ok = tm.Get(3)
	if !ok || val != "three" {
		t.Error("Other keys should remain")
	}
}

func TestTreeMapContainsKey(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	if tm.ContainsKey(1) {
		t.Error("Empty map should not contain any key")
	}

	// Add elements
	tm.Put(1, "one")
	tm.Put(2, "two")

	// Test existing key
	if !tm.ContainsKey(1) {
		t.Error("Map should contain key 1")
	}
	if !tm.ContainsKey(2) {
		t.Error("Map should contain key 2")
	}

	// Test non-existing key
	if tm.ContainsKey(3) {
		t.Error("Map should not contain key 3")
	}
}

func TestTreeMapSize(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	if tm.Size() != 0 {
		t.Error("Empty map size should be 0")
	}

	// Add elements
	tm.Put(1, "one")
	tm.Put(2, "two")
	tm.Put(3, "three")

	if tm.Size() != 3 {
		t.Errorf("Expected size 3, got %d", tm.Size())
	}

	// Remove elements
	tm.Remove(2)
	if tm.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", tm.Size())
	}
}

func TestTreeMapIsEmpty(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	if !tm.IsEmpty() {
		t.Error("New map should be empty")
	}

	// Add elements
	tm.Put(1, "one")
	if tm.IsEmpty() {
		t.Error("Map should not be empty after adding element")
	}

	// Remove elements
	tm.Remove(1)
	if !tm.IsEmpty() {
		t.Error("Map should be empty after removing all elements")
	}
}

func TestTreeMapClear(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Add elements
	tm.Put(1, "one")
	tm.Put(2, "two")
	tm.Put(3, "three")

	// Clear map
	tm.Clear()

	if !tm.IsEmpty() {
		t.Error("Map should be empty after Clear")
	}
	if tm.Size() != 0 {
		t.Error("Map size should be 0 after Clear")
	}

	// Verify all keys are removed
	if tm.ContainsKey(1) || tm.ContainsKey(2) || tm.ContainsKey(3) {
		t.Error("Map should not contain any key after Clear")
	}
}

func TestTreeMapKeys(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	keys := tm.Keys()
	if len(keys) != 0 {
		t.Error("Empty map should have no keys")
	}

	// Add elements (unordered insertion)
	tm.Put(3, "three")
	tm.Put(1, "one")
	tm.Put(4, "four")
	tm.Put(2, "two")

	// Verify key count
	keys = tm.Keys()
	if len(keys) != 4 {
		t.Errorf("Expected 4 keys, got %d", len(keys))
	}

	// Verify keys are ordered
	expected := []int{1, 2, 3, 4}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Expected key %d at position %d, got %d", expected[i], i, key)
		}
	}
}

func TestTreeMapValues(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	values := tm.Values()
	if len(values) != 0 {
		t.Error("Empty map should have no values")
	}

	// Add elements
	tm.Put(3, "three")
	tm.Put(1, "one")
	tm.Put(2, "two")

	// Verify value count
	values = tm.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Verify value order (should be by key order)
	expected := []string{"one", "two", "three"}
	for i, value := range values {
		if value != expected[i] {
			t.Errorf("Expected value '%s' at position %d, got '%s'", expected[i], i, value)
		}
	}
}

func TestTreeMapEntries(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	entries := tm.Entries()
	if len(entries) != 0 {
		t.Error("Empty map should have no entries")
	}

	// Add elements
	tm.Put(3, "three")
	tm.Put(1, "one")
	tm.Put(2, "two")

	// Verify entry count
	entries = tm.Entries()
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Verify entry order (should be by key order)
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}
	for i, entry := range entries {
		if entry.Key != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, entry.Key)
		}
		if entry.Value != expectedValues[i] {
			t.Errorf("Expected value '%s' at position %d, got '%s'", expectedValues[i], i, entry.Value)
		}
	}
}

func TestTreeMapForEach(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Add elements
	tm.Put(3, "three")
	tm.Put(1, "one")
	tm.Put(2, "two")

	// Test ForEach
	var keys []int
	var values []string
	tm.ForEach(func(k int, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})

	// Verify order
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}

	for i := range keys {
		if keys[i] != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, keys[i])
		}
		if values[i] != expectedValues[i] {
			t.Errorf("Expected value '%s' at position %d, got '%s'", expectedValues[i], i, values[i])
		}
	}
}

func TestTreeMapString(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Test empty map
	str := tm.String()
	if str != "{}" {
		t.Errorf("Expected '{}', got '%s'", str)
	}

	// Test single element
	tm.Put(1, "one")
	str = tm.String()
	if str != "{1=one}" {
		t.Errorf("Expected '{1=one}', got '%s'", str)
	}

	// Test multiple elements
	tm.Put(2, "two")
	tm.Put(3, "three")
	str = tm.String()
	expected := "{1=one, 2=two, 3=three}"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestTreeMapLargeDataset(t *testing.T) {
	tm := NewTreeMap[int, string]()

	// Insert large amount of data
	for i := 1000; i >= 1; i-- {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	// Verify size
	if tm.Size() != 1000 {
		t.Errorf("Expected size 1000, got %d", tm.Size())
	}

	// Verify all data
	for i := 1; i <= 1000; i++ {
		val, ok := tm.Get(i)
		if !ok {
			t.Errorf("Key %d should exist", i)
		}
		expected := fmt.Sprintf("value%d", i)
		if val != expected {
			t.Errorf("Expected value '%s' for key %d, got '%s'", expected, i, val)
		}
	}

	// Verify key order
	keys := tm.Keys()
	for i, key := range keys {
		if key != i+1 {
			t.Errorf("Expected key %d at position %d, got %d", i+1, i, key)
		}
	}
}

// Benchmark tests
func BenchmarkTreeMapPut(b *testing.B) {
	tm := NewTreeMap[int, string]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}
}

func BenchmarkTreeMapGet(b *testing.B) {
	tm := NewTreeMap[int, string]()

	// Pre-populate data
	for i := 0; i < 1000; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.Get(i % 1000)
	}
}

func BenchmarkTreeMapRemove(b *testing.B) {
	tm := NewTreeMap[int, string]()

	// Pre-populate data
	for i := 0; i < b.N; i++ {
		tm.Put(i, fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.Remove(i)
	}
}
