package multimap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Common test suite for all Multimap implementations
func testMultimapBasicOperations[K comparable, V comparable](t *testing.T, m Multimap[K, V], key1 K, key2 K, value1 V, value2 V, value3 V) {
	// Test initial state
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Size())

	// Test Put
	assert.True(t, m.Put(key1, value1))
	assert.Equal(t, 1, m.Size())
	assert.False(t, m.IsEmpty())
	assert.True(t, m.ContainsKey(key1))
	assert.True(t, m.ContainsValue(value1))
	assert.True(t, m.ContainsEntry(key1, value1))

	// Test Get
	values := m.Get(key1)
	assert.Equal(t, 1, len(values))
	assert.Contains(t, values, value1)

	// Test Put with same key, different value
	assert.True(t, m.Put(key1, value2))
	assert.Equal(t, 2, m.Size())
	values = m.Get(key1)
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, value1)
	assert.Contains(t, values, value2)

	// Test Put with different key
	assert.True(t, m.Put(key2, value3))
	assert.Equal(t, 3, m.Size())
	assert.True(t, m.ContainsKey(key2))
	values = m.Get(key2)
	assert.Equal(t, 1, len(values))
	assert.Contains(t, values, value3)

	// Test Keys
	keys := m.Keys()
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, key1)
	assert.Contains(t, keys, key2)

	// Test Values
	allValues := m.Values()
	assert.Equal(t, 3, len(allValues))
	assert.Contains(t, allValues, value1)
	assert.Contains(t, allValues, value2)
	assert.Contains(t, allValues, value3)

	// Test Entries
	entries := m.Entries()
	assert.Equal(t, 3, len(entries))
	assertContainsEntry(t, entries, key1, value1)
	assertContainsEntry(t, entries, key1, value2)
	assertContainsEntry(t, entries, key2, value3)

	// Test Remove
	assert.True(t, m.Remove(key1, value1))
	assert.Equal(t, 2, m.Size())
	assert.True(t, m.ContainsKey(key1)) // key1 still has value2
	values = m.Get(key1)
	assert.Equal(t, 1, len(values))
	assert.Contains(t, values, value2)

	// Test RemoveAll
	removed := m.RemoveAll(key1)
	assert.Equal(t, 1, len(removed))
	assert.Contains(t, removed, value2)
	assert.Equal(t, 1, m.Size())
	assert.False(t, m.ContainsKey(key1))

	// Test Clear
	m.Clear()
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Size())
}

// Helper function to check if entries contains a specific key-value pair
func assertContainsEntry[K comparable, V comparable](t *testing.T, entries []Entry[K, V], key K, value V) {
	for _, entry := range entries {
		if entry.Key == key && entry.Value == value {
			return
		}
	}
	t.Errorf("Entry with key %v and value %v not found", key, value)
}

func TestArrayListMultimap(t *testing.T) {
	m := NewArrayListMultimap[string, int]()
	testMultimapBasicOperations[string, int](t, m, "key1", "key2", 1, 2, 3)

	// Test duplicate values (ArrayList allows duplicates)
	m.Clear()
	m.Put("key", 1)
	m.Put("key", 1) // Duplicate value
	assert.Equal(t, 2, m.Size())
	values := m.Get("key")
	assert.Equal(t, 2, len(values))

	// Test ReplaceValues
	m.Clear()
	m.Put("key", 1)
	m.Put("key", 2)
	oldValues := m.ReplaceValues("key", []int{3, 4})
	assert.Equal(t, 2, len(oldValues))
	assert.Contains(t, oldValues, 1)
	assert.Contains(t, oldValues, 2)
	values = m.Get("key")
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, 3)
	assert.Contains(t, values, 4)

	// Test ForEach
	m.Clear()
	m.Put("key1", 1)
	m.Put("key1", 2)
	m.Put("key2", 3)
	sum := 0
	m.ForEach(func(key string, value int) {
		sum += value
	})
	assert.Equal(t, 6, sum)
}

func TestHashMultimap(t *testing.T) {
	m := NewHashMultimap[string, int]()
	testMultimapBasicOperations[string, int](t, m, "key1", "key2", 1, 2, 3)

	// Test duplicate values (HashMultimap doesn't allow duplicates)
	m.Clear()
	m.Put("key", 1)
	assert.False(t, m.Put("key", 1)) // Duplicate value should not be added
	assert.Equal(t, 1, m.Size())
	values := m.Get("key")
	assert.Equal(t, 1, len(values))

	// Test ReplaceValues
	m.Clear()
	m.Put("key", 1)
	m.Put("key", 2)
	oldValues := m.ReplaceValues("key", []int{3, 4, 3}) // Duplicate in replacement
	assert.Equal(t, 2, len(oldValues))
	assert.Contains(t, oldValues, 1)
	assert.Contains(t, oldValues, 2)
	values = m.Get("key")
	assert.Equal(t, 2, len(values)) // Only unique values
	assert.Contains(t, values, 3)
	assert.Contains(t, values, 4)
}

func TestLinkedHashMultimap(t *testing.T) {
	m := NewLinkedHashMultimap[string, int]()
	testMultimapBasicOperations[string, int](t, m, "key1", "key2", 1, 2, 3)

	// Test insertion order of keys
	m.Clear()
	m.Put("c", 1)
	m.Put("a", 2)
	m.Put("b", 3)
	keys := m.Keys()
	assert.Equal(t, []string{"c", "a", "b"}, keys)

	// Test insertion order of values
	m.Clear()
	m.Put("key", 3)
	m.Put("key", 1)
	m.Put("key", 2)
	values := m.Get("key")
	assert.Equal(t, []int{3, 1, 2}, values)

	// Test order preservation after removal
	m.Remove("key", 1)
	values = m.Get("key")
	assert.Equal(t, []int{3, 2}, values)
}

// ComparableInt is a custom comparable type for testing
type ComparableInt int

// CompareTo compares this ComparableInt with another object
// This method is required by the common.Comparable interface
func (a ComparableInt) CompareTo(other interface{}) int {
	if b, ok := other.(ComparableInt); ok {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	return 0 // Default for incompatible types
}

// ComparableString is a custom comparable type for testing
type ComparableString string

// CompareTo compares this ComparableString with another object
// This method is required by the common.Comparable interface
func (a ComparableString) CompareTo(other interface{}) int {
	if b, ok := other.(ComparableString); ok {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	return 0 // Default for incompatible types
}

func TestTreeMultimap(t *testing.T) {

	m := NewTreeMultimap[ComparableString, ComparableInt]()
	testMultimapBasicOperations[ComparableString, ComparableInt](t, m, ComparableString("key1"), ComparableString("key2"), ComparableInt(1), ComparableInt(2), ComparableInt(3))

	// Test sorted order of keys
	m.Clear()
	m.Put(ComparableString("c"), ComparableInt(1))
	m.Put(ComparableString("a"), ComparableInt(2))
	m.Put(ComparableString("b"), ComparableInt(3))
	keys := m.Keys()
	expected := []ComparableString{"a", "b", "c"}
	assert.Equal(t, expected, keys)

	// Test sorted order of values
	m.Clear()
	m.Put(ComparableString("key"), ComparableInt(3))
	m.Put(ComparableString("key"), ComparableInt(1))
	m.Put(ComparableString("key"), ComparableInt(2))
	values := m.Get(ComparableString("key"))
	expectedValues := []ComparableInt{1, 2, 3}
	assert.Equal(t, expectedValues, values)
}

func TestImmutableMultimap(t *testing.T) {
	// Create a mutable multimap and add some entries
	mutable := NewHashMultimap[string, int]()
	mutable.Put("key1", 1)
	mutable.Put("key1", 2)
	mutable.Put("key2", 3)

	// Create an immutable copy
	immutable := FromMultimap[string, int](mutable)

	// Test basic operations
	assert.Equal(t, 3, immutable.Size())
	assert.False(t, immutable.IsEmpty())
	assert.True(t, immutable.ContainsKey("key1"))
	assert.True(t, immutable.ContainsValue(2))
	assert.True(t, immutable.ContainsEntry("key2", 3))

	// Test Get
	values := immutable.Get("key1")
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, 1)
	assert.Contains(t, values, 2)

	// Test immutability - operations should return false/nil instead of panicking
	assert.False(t, immutable.Put("key3", 4))
	assert.False(t, immutable.Remove("key1", 1))
	immutable.Clear() // Should log warning but not panic

	// Test Of constructor
	immutable2 := Of[string, int]("key1", 1, "key1", 2, "key2", 3)
	assert.Equal(t, 3, immutable2.Size())
	assert.True(t, immutable2.ContainsEntry("key1", 1))
	assert.True(t, immutable2.ContainsEntry("key1", 2))
	assert.True(t, immutable2.ContainsEntry("key2", 3))
}

func TestImmutableListMultimap(t *testing.T) {
	// Create a mutable multimap and add some entries
	mutable := NewArrayListMultimap[string, int]()
	mutable.Put("key1", 1)
	mutable.Put("key1", 1) // Duplicate
	mutable.Put("key1", 2)
	mutable.Put("key2", 3)

	// Create an immutable copy
	immutable := FromArrayListMultimap(mutable)

	// Test basic operations
	assert.Equal(t, 4, immutable.Size())
	assert.False(t, immutable.IsEmpty())
	assert.True(t, immutable.ContainsKey("key1"))
	assert.True(t, immutable.ContainsValue(2))
	assert.True(t, immutable.ContainsEntry("key2", 3))

	// Test Get - should preserve duplicates
	values := immutable.Get("key1")
	assert.Equal(t, 3, len(values))
	assert.Equal(t, 1, values[0])
	assert.Equal(t, 1, values[1])
	assert.Equal(t, 2, values[2])

	// Test immutability - operations should return false/nil instead of panicking
	assert.False(t, immutable.Put("key3", 4))
	assert.False(t, immutable.Remove("key1", 1))
	immutable.Clear() // Should log warning but not panic

	// Test ListOf constructor
	immutable2 := ListOf[string, int]("key1", 1, "key1", 1, "key1", 2, "key2", 3)
	assert.Equal(t, 4, immutable2.Size())
	values = immutable2.Get("key1")
	assert.Equal(t, 3, len(values))
	assert.Equal(t, 1, values[0])
	assert.Equal(t, 1, values[1])
	assert.Equal(t, 2, values[2])
}

func TestImmutableSetMultimap(t *testing.T) {
	// Create a mutable multimap and add some entries
	mutable := NewHashMultimap[string, int]()
	mutable.Put("key1", 1)
	mutable.Put("key1", 1) // Duplicate that will be ignored
	mutable.Put("key1", 2)
	mutable.Put("key2", 3)

	// Create an immutable copy
	immutable := FromHashMultimap(mutable)

	// Test basic operations
	assert.Equal(t, 3, immutable.Size())
	assert.False(t, immutable.IsEmpty())
	assert.True(t, immutable.ContainsKey("key1"))
	assert.True(t, immutable.ContainsValue(2))
	assert.True(t, immutable.ContainsEntry("key2", 3))

	// Test Get - should eliminate duplicates
	values := immutable.Get("key1")
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, 1)
	assert.Contains(t, values, 2)

	// Test immutability - operations should return false/nil instead of panicking
	assert.False(t, immutable.Put("key3", 4))
	assert.False(t, immutable.Remove("key1", 1))
	immutable.Clear() // Should log warning but not panic

	// Test SetOf constructor
	immutable2 := SetOf[string, int]("key1", 1, "key1", 1, "key1", 2, "key2", 3)
	assert.Equal(t, 3, immutable2.Size())
	values = immutable2.Get("key1")
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, 1)
	assert.Contains(t, values, 2)
}