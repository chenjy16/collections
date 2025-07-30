package maps

import (
	"fmt"
	"github.com/chenjianyu/collections/container/common"
	"strings"
	"sync"
)

// CopyOnWriteMap is a thread-safe map implementation using copy-on-write strategy
// Suitable for scenarios with more reads than writes
type CopyOnWriteMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewCopyOnWriteMap creates a new CopyOnWriteMap
func NewCopyOnWriteMap[K comparable, V any]() *CopyOnWriteMap[K, V] {
	return NewCopyOnWriteMapWithCapacity[K, V](0)
}

// NewCopyOnWriteMapWithCapacity creates a new CopyOnWriteMap with initial capacity
func NewCopyOnWriteMapWithCapacity[K comparable, V any](capacity int) *CopyOnWriteMap[K, V] {
	return &CopyOnWriteMap[K, V]{
		data: make(map[K]V, capacity),
	}
}

// CopyOnWriteMapFromMap creates a new CopyOnWriteMap from a Go map
func CopyOnWriteMapFromMap[K comparable, V any](source map[K]V) *CopyOnWriteMap[K, V] {
	m := NewCopyOnWriteMapWithCapacity[K, V](len(source))
	for k, v := range source {
		m.Put(k, v)
	}
	return m
}

// Put associates the specified value with the specified key in this map
func (m *CopyOnWriteMap[K, V]) Put(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	oldValue, existed := m.data[key]

	// Copy-on-write: create a new map copy
	newData := make(map[K]V, len(m.data)+1)
	for k, v := range m.data {
		newData[k] = v
	}
	newData[key] = value
	m.data = newData

	return oldValue, existed
}

// Get returns the value mapped to the specified key
func (m *CopyOnWriteMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.data[key]
	return value, exists
}

// Remove if exists, removes mapping relationship for the key from this map
func (m *CopyOnWriteMap[K, V]) Remove(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	oldValue, existed := m.data[key]
	if !existed {
		return oldValue, false
	}

	// Copy-on-write: create a new map copy
	newData := make(map[K]V, len(m.data)-1)
	for k, v := range m.data {
		if k != key {
			newData[k] = v
		}
	}
	m.data = newData

	return oldValue, true
}

// ContainsKey if this map contains mapping relationship for the specified key, returns true
func (m *CopyOnWriteMap[K, V]) ContainsKey(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists
}

// ContainsValue if this map maps one or more keys to the specified value, returns true
func (m *CopyOnWriteMap[K, V]) ContainsValue(value V) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, v := range m.data {
		if Equal(v, value) {
			return true
		}
	}
	return false
}

// Size returns the number of key-value mapping relationships in this map
func (m *CopyOnWriteMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data)
}

// IsEmpty if this map does not contain key-value mapping relationships, returns true
func (m *CopyOnWriteMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear removes all mapping relationships from this map
func (m *CopyOnWriteMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[K]V)
}

// Keys returns the keys contained in this map
func (m *CopyOnWriteMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values contained in this map
func (m *CopyOnWriteMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	values := make([]V, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, v)
	}
	return values
}

// Entries returns the mapping relationships contained in this map
func (m *CopyOnWriteMap[K, V]) Entries() []Pair[K, V] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := make([]Pair[K, V], 0, len(m.data))
	for k, v := range m.data {
		entries = append(entries, NewPair(k, v))
	}
	return entries
}

// ForEach executes the given operation for each entry in this map
func (m *CopyOnWriteMap[K, V]) ForEach(f func(K, V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.data {
		f(k, v)
	}
}

// String returns the string representation of the map
func (m *CopyOnWriteMap[K, V]) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.data) == 0 {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")
	first := true
	for k, v := range m.data {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v=%v", k, v))
		first = false
	}
	builder.WriteString("}")
	return builder.String()
}

// PutAll copies all mapping relationships from the specified map to this map
func (m *CopyOnWriteMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(k K, v V) {
		m.Put(k, v)
	})
}

// PutAllFromMap copies all mappings from a Go map to this map
func (m *CopyOnWriteMap[K, V]) PutAllFromMap(source map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	newData := make(map[K]V, len(m.data)+len(source))
	for k, v := range m.data {
		newData[k] = v
	}
	for k, v := range source {
		newData[k] = v
	}
	m.data = newData
}

// ToMap returns a copy of the map as a Go map
func (m *CopyOnWriteMap[K, V]) ToMap() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copy := make(map[K]V, len(m.data))
	for k, v := range m.data {
		copy[k] = v
	}
	return copy
}

// Snapshot returns a snapshot of the current map
func (m *CopyOnWriteMap[K, V]) Snapshot() map[K]V {
	return m.ToMap()
}

// PutIfAbsent puts the value if the key is not present
func (m *CopyOnWriteMap[K, V]) PutIfAbsent(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	oldValue, exists := m.data[key]
	if exists {
		return oldValue, false
	}
	newData := make(map[K]V, len(m.data)+1)
	for k, v := range m.data {
		newData[k] = v
	}
	newData[key] = value
	m.data = newData
	return oldValue, true
}

// Replace replaces the value for the key if present
func (m *CopyOnWriteMap[K, V]) Replace(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	oldValue, exists := m.data[key]
	if !exists {
		return oldValue, false
	}
	newData := make(map[K]V, len(m.data))
	for k, v := range m.data {
		newData[k] = v
	}
	newData[key] = value
	m.data = newData
	return oldValue, true
}

// ReplaceIf replaces the value for the key if it matches the old value
func (m *CopyOnWriteMap[K, V]) ReplaceIf(key K, oldValue V, newValue V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	current, exists := m.data[key]
	if !exists || !common.Equal(current, oldValue) {
		return false
	}
	newData := make(map[K]V, len(m.data))
	for k, v := range m.data {
		newData[k] = v
	}
	newData[key] = newValue
	m.data = newData
	return true
}
