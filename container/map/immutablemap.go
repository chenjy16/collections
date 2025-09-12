package maps

import (
	"fmt"
	"strings"
)

// ImmutableMap is an immutable implementation of Map
// Once created, it cannot be modified. All modification operations return new instances
type ImmutableMap[K comparable, V any] struct {
	entries map[K]V
}

// NewImmutableMap creates a new empty ImmutableMap
func NewImmutableMap[K comparable, V any]() *ImmutableMap[K, V] {
	return &ImmutableMap[K, V]{
		entries: make(map[K]V),
	}
}

// NewImmutableMapFromMap creates a new ImmutableMap from an existing map
func NewImmutableMapFromMap[K comparable, V any](m map[K]V) *ImmutableMap[K, V] {
	entries := make(map[K]V, len(m))
	for k, v := range m {
		entries[k] = v
	}

	return &ImmutableMap[K, V]{
		entries: entries,
	}
}

// MapOf creates a new ImmutableMap from the given key-value pairs
func MapOf[K comparable, V any](pairs ...Pair[K, V]) *ImmutableMap[K, V] {
	entries := make(map[K]V, len(pairs))
	for _, pair := range pairs {
		entries[pair.Key] = pair.Value
	}

	return &ImmutableMap[K, V]{
		entries: entries,
	}
}

// copyEntries creates a copy of the entries map
func (im *ImmutableMap[K, V]) copyEntries() map[K]V {
	entriesCopy := make(map[K]V, len(im.entries))
	for k, v := range im.entries {
		entriesCopy[k] = v
	}
	return entriesCopy
}

// Size returns the number of key-value pairs in the map
func (im *ImmutableMap[K, V]) Size() int {
	return len(im.entries)
}

// IsEmpty returns true if the map is empty
func (im *ImmutableMap[K, V]) IsEmpty() bool {
	return len(im.entries) == 0
}

// Clear is a no-op as ImmutableMap is immutable
func (im *ImmutableMap[K, V]) Clear() {
	// No-op for immutable collections
}

// Put returns the zero value and false as ImmutableMap is immutable
func (im *ImmutableMap[K, V]) Put(key K, value V) (V, bool) {
	var zero V
	// Return zero value and false to indicate the operation failed
	return zero, false
}

// Get returns the value associated with the key
func (im *ImmutableMap[K, V]) Get(key K) (V, bool) {
	value, exists := im.entries[key]
	return value, exists
}

// Remove returns the zero value and false as ImmutableMap is immutable
func (im *ImmutableMap[K, V]) Remove(key K) (V, bool) {
	var zero V
	// Return zero value and false to indicate the operation failed
	return zero, false
}

// ContainsKey returns true if the map contains the specified key
func (im *ImmutableMap[K, V]) ContainsKey(key K) bool {
	_, exists := im.entries[key]
	return exists
}

// ContainsValue returns true if the map contains the specified value
func (im *ImmutableMap[K, V]) ContainsValue(value V) bool {
	for _, v := range im.entries {
		if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", value) {
			return true
		}
	}
	return false
}

// Keys returns a slice of all keys in the map
func (im *ImmutableMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(im.entries))
	for k := range im.entries {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values in the map
func (im *ImmutableMap[K, V]) Values() []V {
	values := make([]V, 0, len(im.entries))
	for _, v := range im.entries {
		values = append(values, v)
	}
	return values
}

// Entries returns a slice of all key-value pairs in the map
func (im *ImmutableMap[K, V]) Entries() []Pair[K, V] {
	entries := make([]Pair[K, V], 0, len(im.entries))
	for k, v := range im.entries {
		entries = append(entries, Pair[K, V]{Key: k, Value: v})
	}
	return entries
}

// WithPut returns a new ImmutableMap with the key-value pair added/updated
func (im *ImmutableMap[K, V]) WithPut(key K, value V) *ImmutableMap[K, V] {
	newEntries := im.copyEntries()
	newEntries[key] = value

	return &ImmutableMap[K, V]{entries: newEntries}
}

// WithRemove returns a new ImmutableMap with the key removed
func (im *ImmutableMap[K, V]) WithRemove(key K) *ImmutableMap[K, V] {
	if !im.ContainsKey(key) {
		return im // Key doesn't exist, return same instance
	}

	newEntries := im.copyEntries()
	delete(newEntries, key)

	return &ImmutableMap[K, V]{entries: newEntries}
}

// WithClear returns a new empty ImmutableMap
func (im *ImmutableMap[K, V]) WithClear() *ImmutableMap[K, V] {
	return NewImmutableMap[K, V]()
}

// PutAll is not supported for ImmutableMap - this is a no-op
func (im *ImmutableMap[K, V]) PutAll(other Map[K, V]) {
	// No-op for immutable collections
}

// WithPutAll returns a new ImmutableMap with all key-value pairs from another map added
func (im *ImmutableMap[K, V]) WithPutAll(other Map[K, V]) *ImmutableMap[K, V] {
	newEntries := im.copyEntries()

	other.ForEach(func(key K, value V) {
		newEntries[key] = value
	})

	return &ImmutableMap[K, V]{entries: newEntries}
}

// ForEach executes the given function for each key-value pair in the map
func (im *ImmutableMap[K, V]) ForEach(fn func(K, V)) {
	for k, v := range im.entries {
		fn(k, v)
	}
}

// String returns a string representation of the map
func (im *ImmutableMap[K, V]) String() string {
	if im.IsEmpty() {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")

	first := true
	for k, v := range im.entries {
		if !first {
			builder.WriteString(", ")
		}
		first = false
		builder.WriteString(fmt.Sprintf("%v=%v", k, v))
	}

	builder.WriteString("}")
	return builder.String()
}
