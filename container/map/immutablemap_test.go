package maps

import (
	"testing"
)

func TestNewImmutableMap(t *testing.T) {
	m := NewImmutableMap[string, int]()
	
	if !m.IsEmpty() {
		t.Error("New map should be empty")
	}
	
	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}
}

func TestNewImmutableMapFromMap(t *testing.T) {
	original := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	
	m := NewImmutableMapFromMap(original)
	
	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}
	
	for key, expectedValue := range original {
		if value, exists := m.Get(key); !exists || value != expectedValue {
			t.Errorf("Expected %d for key %s, got %d (exists: %t)", expectedValue, key, value, exists)
		}
	}
}

func TestMapOf(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
		NewPair("three", 3),
	)
	
	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}
	
	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	
	for key, expectedValue := range expected {
		if value, exists := m.Get(key); !exists || value != expectedValue {
			t.Errorf("Expected %d for key %s, got %d (exists: %t)", expectedValue, key, value, exists)
		}
	}
}

func TestImmutableMapGet(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	// Existing keys
	if value, exists := m.Get("one"); !exists || value != 1 {
		t.Errorf("Expected (1, true) for key 'one', got (%d, %t)", value, exists)
	}
	
	if value, exists := m.Get("two"); !exists || value != 2 {
		t.Errorf("Expected (2, true) for key 'two', got (%d, %t)", value, exists)
	}
	
	// Non-existing key
	if value, exists := m.Get("three"); exists {
		t.Errorf("Expected (0, false) for non-existing key, got (%d, %t)", value, exists)
	}
}

func TestImmutableMapContainsKey(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	if !m.ContainsKey("one") {
		t.Error("Map should contain key 'one'")
	}
	
	if m.ContainsKey("three") {
		t.Error("Map should not contain key 'three'")
	}
}

func TestImmutableMapContainsValue(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	if !m.ContainsValue(1) {
		t.Error("Map should contain value 1")
	}
	
	if m.ContainsValue(3) {
		t.Error("Map should not contain value 3")
	}
}

func TestImmutableMapWithPut(t *testing.T) {
	original := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	// Add new key-value pair
	newMap := original.WithPut("three", 3)
	
	// Original should be unchanged
	if original.Size() != 2 {
		t.Errorf("Original map size should remain 2, got %d", original.Size())
	}
	
	if original.ContainsKey("three") {
		t.Error("Original map should not contain key 'three'")
	}
	
	// New map should have the added pair
	if newMap.Size() != 3 {
		t.Errorf("New map size should be 3, got %d", newMap.Size())
	}
	
	if value, exists := newMap.Get("three"); !exists || value != 3 {
		t.Errorf("New map should contain 'three' -> 3, got %d (exists: %t)", value, exists)
	}
	
	// Update existing key
	updatedMap := original.WithPut("one", 10)
	
	if value, exists := updatedMap.Get("one"); !exists || value != 10 {
		t.Errorf("Updated map should contain 'one' -> 10, got %d (exists: %t)", value, exists)
	}
	
	// Original should still have old value
	if value, exists := original.Get("one"); !exists || value != 1 {
		t.Errorf("Original map should still contain 'one' -> 1, got %d (exists: %t)", value, exists)
	}
}

func TestImmutableMapWithRemove(t *testing.T) {
	original := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
		NewPair("three", 3),
	)
	
	newMap := original.WithRemove("two")
	
	// Original should be unchanged
	if original.Size() != 3 {
		t.Errorf("Original map size should remain 3, got %d", original.Size())
	}
	
	if !original.ContainsKey("two") {
		t.Error("Original map should still contain key 'two'")
	}
	
	// New map should have the key removed
	if newMap.Size() != 2 {
		t.Errorf("New map size should be 2, got %d", newMap.Size())
	}
	
	if newMap.ContainsKey("two") {
		t.Error("New map should not contain key 'two'")
	}
	
	// Test removing non-existent key
	sameMap := original.WithRemove("four")
	if sameMap != original {
		t.Error("Removing non-existent key should return same instance")
	}
}

func TestImmutableMapWithClear(t *testing.T) {
	original := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	newMap := original.WithClear()
	
	// Original should be unchanged
	if original.Size() != 2 {
		t.Errorf("Original map size should remain 2, got %d", original.Size())
	}
	
	// New map should be empty
	if !newMap.IsEmpty() {
		t.Error("New map should be empty")
	}
}

func TestImmutableMapKeys(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
		NewPair("three", 3),
	)
	
	keys := m.Keys()
	
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
	
	expectedKeys := map[string]bool{
		"one":   true,
		"two":   true,
		"three": true,
	}
	
	for _, key := range keys {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key: %s", key)
		}
	}
}

func TestImmutableMapValues(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
		NewPair("three", 3),
	)
	
	values := m.Values()
	
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
	
	expectedValues := map[int]bool{
		1: true,
		2: true,
		3: true,
	}
	
	for _, value := range values {
		if !expectedValues[value] {
			t.Errorf("Unexpected value: %d", value)
		}
	}
}

func TestImmutableMapEntries(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	entries := m.Entries()
	
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
	
	expectedEntries := map[string]int{
		"one": 1,
		"two": 2,
	}
	
	for _, entry := range entries {
		if expectedValue, exists := expectedEntries[entry.Key]; !exists || expectedValue != entry.Value {
			t.Errorf("Unexpected entry: %s -> %d", entry.Key, entry.Value)
		}
	}
}

func TestImmutableMapWithPutAll(t *testing.T) {
	original := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	other := MapOf(
		NewPair("three", 3),
		NewPair("four", 4),
		NewPair("one", 10), // Override existing key
	)
	
	newMap := original.WithPutAll(other)
	
	// Original should be unchanged
	if original.Size() != 2 {
		t.Errorf("Original map size should remain 2, got %d", original.Size())
	}
	
	// New map should have all entries
	if newMap.Size() != 4 {
		t.Errorf("New map size should be 4, got %d", newMap.Size())
	}
	
	// Check specific values
	if value, exists := newMap.Get("one"); !exists || value != 10 {
		t.Errorf("Expected 'one' -> 10, got %d (exists: %t)", value, exists)
	}
	
	if value, exists := newMap.Get("three"); !exists || value != 3 {
		t.Errorf("Expected 'three' -> 3, got %d (exists: %t)", value, exists)
	}
}

func TestImmutableMapForEach(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
		NewPair("three", 3),
	)
	
	sum := 0
	keyCount := 0
	
	m.ForEach(func(key string, value int) {
		sum += value
		keyCount++
	})
	
	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
	
	if keyCount != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", keyCount)
	}
}

func TestImmutableMapString(t *testing.T) {
	emptyMap := NewImmutableMap[string, int]()
	if emptyMap.String() != "{}" {
		t.Errorf("Expected '{}' for empty map, got '%s'", emptyMap.String())
	}
	
	singleMap := MapOf(NewPair("one", 1))
	str := singleMap.String()
	expected := "{one=1}"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}

func TestImmutableMapImmutability(t *testing.T) {
	m := MapOf(
		NewPair("one", 1),
		NewPair("two", 2),
	)
	
	// Test that modification methods don't change the original
	m.Put("three", 3)
	m.Remove("one")
	m.Clear()
	m.PutAll(MapOf(NewPair("four", 4)))
	
	// Original should remain unchanged
	if m.Size() != 2 {
		t.Errorf("Original map should remain unchanged, expected size 2, got %d", m.Size())
	}
	
	expectedEntries := map[string]int{
		"one": 1,
		"two": 2,
	}
	
	for key, expectedValue := range expectedEntries {
		if value, exists := m.Get(key); !exists || value != expectedValue {
			t.Errorf("Original map should remain unchanged, expected %d for key %s, got %d (exists: %t)", expectedValue, key, value, exists)
		}
	}
}