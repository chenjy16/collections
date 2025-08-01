package multimap

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/set"
)

// ImmutableSetMultimap is an immutable implementation of a multimap that eliminates duplicate values
type ImmutableSetMultimap[K comparable, V comparable] struct {
	entries []Entry[K, V]
	data    map[K][]V
	sets    map[K]set.Set[V] // Used for quick lookups and deduplication
}

// NewImmutableSetMultimap creates a new ImmutableSetMultimap from the given entries
func NewImmutableSetMultimap[K comparable, V comparable](entries []Entry[K, V]) *ImmutableSetMultimap[K, V] {
	data := make(map[K][]V)
	sets := make(map[K]set.Set[V])
	dedupEntries := make([]Entry[K, V], 0)
	
	// Group values by key and eliminate duplicates
	for _, entry := range entries {
		key := entry.Key
		value := entry.Value
		
		// Initialize set for this key if it doesn't exist
		if _, exists := sets[key]; !exists {
			sets[key] = set.New[V]()
		}
		
		// Only add if not already present
		if !sets[key].Contains(value) {
			sets[key].Add(value)
			data[key] = append(data[key], value)
			dedupEntries = append(dedupEntries, entry)
		}
	}
	
	return &ImmutableSetMultimap[K, V]{
		entries: dedupEntries,
		data:    data,
		sets:    sets,
	}
}

// SetOf creates a new ImmutableSetMultimap from the given key-value pairs
func SetOf[K comparable, V comparable](pairs ...interface{}) *ImmutableSetMultimap[K, V] {
	if len(pairs)%2 != 0 {
		panic("ImmutableSetMultimap.SetOf requires an even number of arguments")
	}
	
	entries := make([]Entry[K, V], 0, len(pairs)/2)
	
	for i := 0; i < len(pairs); i += 2 {
		key, ok1 := pairs[i].(K)
		value, ok2 := pairs[i+1].(V)
		
		if !ok1 || !ok2 {
			panic("ImmutableSetMultimap.SetOf: invalid type for key or value")
		}
		
		entries = append(entries, Entry[K, V]{Key: key, Value: value})
	}
	
	return NewImmutableSetMultimap(entries)
}

// FromHashMultimap creates a new ImmutableSetMultimap from the given HashMultimap
func FromHashMultimap[K comparable, V comparable](multimap *HashMultimap[K, V]) *ImmutableSetMultimap[K, V] {
	return NewImmutableSetMultimap(multimap.Entries())
}

// FromMultimapToSet creates a new ImmutableSetMultimap from any multimap
func FromMultimapToSet[K comparable, V comparable](multimap Multimap[K, V]) *ImmutableSetMultimap[K, V] {
	return NewImmutableSetMultimap(multimap.Entries())
}

// Put is not supported for ImmutableSetMultimap and will panic
func (m *ImmutableSetMultimap[K, V]) Put(key K, value V) bool {
	panic("ImmutableSetMultimap is immutable")
}

// PutAll is not supported for ImmutableSetMultimap and will panic
func (m *ImmutableSetMultimap[K, V]) PutAll(multimap Multimap[K, V]) bool {
	panic("ImmutableSetMultimap is immutable")
}

// ReplaceValues is not supported for ImmutableSetMultimap and will panic
func (m *ImmutableSetMultimap[K, V]) ReplaceValues(key K, values []V) []V {
	panic("ImmutableSetMultimap is immutable")
}

// Remove is not supported for ImmutableSetMultimap and will panic
func (m *ImmutableSetMultimap[K, V]) Remove(key K, value V) bool {
	panic("ImmutableSetMultimap is immutable")
}

// RemoveAll is not supported for ImmutableSetMultimap and will panic
func (m *ImmutableSetMultimap[K, V]) RemoveAll(key K) []V {
	panic("ImmutableSetMultimap is immutable")
}

// ContainsKey returns true if this multimap contains at least one key-value mapping with the specified key
func (m *ImmutableSetMultimap[K, V]) ContainsKey(key K) bool {
	_, exists := m.data[key]
	return exists
}

// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
func (m *ImmutableSetMultimap[K, V]) ContainsValue(value V) bool {
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
func (m *ImmutableSetMultimap[K, V]) ContainsEntry(key K, value V) bool {
	set, exists := m.sets[key]
	if !exists {
		return false
	}
	
	return set.Contains(value)
}

// Get returns all values associated with the specified key, with duplicates removed
func (m *ImmutableSetMultimap[K, V]) Get(key K) []V {
	values, exists := m.data[key]
	if !exists {
		return nil
	}
	
	// Return a copy to maintain immutability
	result := make([]V, len(values))
	copy(result, values)
	return result
}

// Keys returns all distinct keys in this multimap
func (m *ImmutableSetMultimap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values in this multimap with duplicates removed
func (m *ImmutableSetMultimap[K, V]) Values() []V {
	values := make([]V, 0, len(m.entries))
	for _, entry := range m.entries {
		values = append(values, entry.Value)
	}
	return values
}

// Entries returns all key-value pairs in this multimap with duplicates removed
func (m *ImmutableSetMultimap[K, V]) Entries() []Entry[K, V] {
	// Return a copy to maintain immutability
	result := make([]Entry[K, V], len(m.entries))
	copy(result, m.entries)
	return result
}

// KeySet returns a set view of the distinct keys in this multimap
func (m *ImmutableSetMultimap[K, V]) KeySet() []K {
	return m.Keys()
}

// AsMap returns a map view of this multimap, mapping each key to its collection of values
func (m *ImmutableSetMultimap[K, V]) AsMap() map[K][]V {
	// Return a deep copy to maintain immutability
	result := make(map[K][]V, len(m.data))
	for key, values := range m.data {
		valuesCopy := make([]V, len(values))
		copy(valuesCopy, values)
		result[key] = valuesCopy
	}
	return result
}

// ForEach executes the given function for each key-value pair in this multimap
func (m *ImmutableSetMultimap[K, V]) ForEach(f func(K, V)) {
	for _, entry := range m.entries {
		f(entry.Key, entry.Value)
	}
}

// Size returns the number of key-value mappings in this multimap
func (m *ImmutableSetMultimap[K, V]) Size() int {
	return len(m.entries)
}

// IsEmpty returns true if this multimap contains no key-value mappings
func (m *ImmutableSetMultimap[K, V]) IsEmpty() bool {
	return len(m.entries) == 0
}

// Clear is not supported for ImmutableSetMultimap and will panic
func (m *ImmutableSetMultimap[K, V]) Clear() {
	panic("ImmutableSetMultimap is immutable")
}

// Contains returns true if this multimap contains the specified element
func (m *ImmutableSetMultimap[K, V]) Contains(key K) bool {
	return m.ContainsKey(key)
}

// String returns a string representation of this multimap
func (m *ImmutableSetMultimap[K, V]) String() string {
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