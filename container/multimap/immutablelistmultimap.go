package multimap

import (
	"fmt"
	"strings"
)

// ImmutableListMultimap is an immutable implementation of a multimap that preserves duplicate values and insertion order
type ImmutableListMultimap[K comparable, V comparable] struct {
	entries []Entry[K, V]
	data    map[K][]V
}

// NewImmutableListMultimap creates a new ImmutableListMultimap from the given entries
func NewImmutableListMultimap[K comparable, V comparable](entries []Entry[K, V]) *ImmutableListMultimap[K, V] {
	data := make(map[K][]V)
	
	// Group values by key while preserving order and duplicates
	for _, entry := range entries {
		data[entry.Key] = append(data[entry.Key], entry.Value)
	}
	
	return &ImmutableListMultimap[K, V]{
		entries: append([]Entry[K, V]{}, entries...), // Create a copy of entries
		data:    data,
	}
}

// ListOf creates a new ImmutableListMultimap from the given key-value pairs
func ListOf[K comparable, V comparable](pairs ...interface{}) *ImmutableListMultimap[K, V] {
	if len(pairs)%2 != 0 {
		panic("ImmutableListMultimap.ListOf requires an even number of arguments")
	}
	
	entries := make([]Entry[K, V], 0, len(pairs)/2)
	
	for i := 0; i < len(pairs); i += 2 {
		key, ok1 := pairs[i].(K)
		value, ok2 := pairs[i+1].(V)
		
		if !ok1 || !ok2 {
			panic("ImmutableListMultimap.ListOf: invalid type for key or value")
		}
		
		entries = append(entries, Entry[K, V]{Key: key, Value: value})
	}
	
	return NewImmutableListMultimap(entries)
}

// FromArrayListMultimap creates a new ImmutableListMultimap from the given ArrayListMultimap
func FromArrayListMultimap[K comparable, V comparable](multimap *ArrayListMultimap[K, V]) *ImmutableListMultimap[K, V] {
	return NewImmutableListMultimap(multimap.Entries())
}

// FromMultimap creates a new ImmutableListMultimap from any multimap
func FromMultimapToList[K comparable, V comparable](multimap Multimap[K, V]) *ImmutableListMultimap[K, V] {
	return NewImmutableListMultimap(multimap.Entries())
}

// Put is not supported for ImmutableListMultimap and will panic
func (m *ImmutableListMultimap[K, V]) Put(key K, value V) bool {
	panic("ImmutableListMultimap is immutable")
}

// PutAll is not supported for ImmutableListMultimap and will panic
func (m *ImmutableListMultimap[K, V]) PutAll(multimap Multimap[K, V]) bool {
	panic("ImmutableListMultimap is immutable")
}

// ReplaceValues is not supported for ImmutableListMultimap and will panic
func (m *ImmutableListMultimap[K, V]) ReplaceValues(key K, values []V) []V {
	panic("ImmutableListMultimap is immutable")
}

// Remove is not supported for ImmutableListMultimap and will panic
func (m *ImmutableListMultimap[K, V]) Remove(key K, value V) bool {
	panic("ImmutableListMultimap is immutable")
}

// RemoveAll is not supported for ImmutableListMultimap and will panic
func (m *ImmutableListMultimap[K, V]) RemoveAll(key K) []V {
	panic("ImmutableListMultimap is immutable")
}

// ContainsKey returns true if this multimap contains at least one key-value mapping with the specified key
func (m *ImmutableListMultimap[K, V]) ContainsKey(key K) bool {
	_, exists := m.data[key]
	return exists
}

// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
func (m *ImmutableListMultimap[K, V]) ContainsValue(value V) bool {
	for _, values := range m.data {
		for _, v := range values {
			if v == value {
				return true
			}
		}
	}
	return false
}

// ContainsEntry returns true if this multimap contains the specified key-value mapping
func (m *ImmutableListMultimap[K, V]) ContainsEntry(key K, value V) bool {
	values, exists := m.data[key]
	if !exists {
		return false
	}
	
	for _, v := range values {
		if v == value {
			return true
		}
	}
	
	return false
}

// Get returns all values associated with the specified key, preserving duplicates and order
func (m *ImmutableListMultimap[K, V]) Get(key K) []V {
	values, exists := m.data[key]
	if !exists {
		return nil
	}
	
	// Return a copy to maintain immutability
	result := make([]V, len(values))
	copy(result, values)
	return result
}

// Keys returns all keys in this multimap, including duplicates for each value
func (m *ImmutableListMultimap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.entries))
	for _, entry := range m.entries {
		keys = append(keys, entry.Key)
	}
	return keys
}

// Values returns all values in this multimap in insertion order
func (m *ImmutableListMultimap[K, V]) Values() []V {
	values := make([]V, 0, len(m.entries))
	for _, entry := range m.entries {
		values = append(values, entry.Value)
	}
	return values
}

// Entries returns all key-value pairs in this multimap in insertion order
func (m *ImmutableListMultimap[K, V]) Entries() []Entry[K, V] {
	// Return a copy to maintain immutability
	result := make([]Entry[K, V], len(m.entries))
	copy(result, m.entries)
	return result
}

// KeySet returns a set view of the distinct keys in this multimap
func (m *ImmutableListMultimap[K, V]) KeySet() []K {
	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

// AsMap returns a map view of this multimap, mapping each key to its collection of values
func (m *ImmutableListMultimap[K, V]) AsMap() map[K][]V {
	// Return a deep copy to maintain immutability
	result := make(map[K][]V, len(m.data))
	for key, values := range m.data {
		valuesCopy := make([]V, len(values))
		copy(valuesCopy, values)
		result[key] = valuesCopy
	}
	return result
}

// ForEach executes the given function for each key-value pair in this multimap in insertion order
func (m *ImmutableListMultimap[K, V]) ForEach(f func(K, V)) {
	for _, entry := range m.entries {
		f(entry.Key, entry.Value)
	}
}

// Size returns the number of key-value mappings in this multimap
func (m *ImmutableListMultimap[K, V]) Size() int {
	return len(m.entries)
}

// IsEmpty returns true if this multimap contains no key-value mappings
func (m *ImmutableListMultimap[K, V]) IsEmpty() bool {
	return len(m.entries) == 0
}

// Clear is not supported for ImmutableListMultimap and will panic
func (m *ImmutableListMultimap[K, V]) Clear() {
	panic("ImmutableListMultimap is immutable")
}

// Contains returns true if this multimap contains the specified element
func (m *ImmutableListMultimap[K, V]) Contains(key K) bool {
	return m.ContainsKey(key)
}

// String returns a string representation of this multimap
func (m *ImmutableListMultimap[K, V]) String() string {
	if m.IsEmpty() {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")

	first := true
	for key, values := range m.data {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		builder.WriteString(fmt.Sprintf("%v=[%v]", key, formatValues(values)))
	}

	builder.WriteString("}")
	return builder.String()
}