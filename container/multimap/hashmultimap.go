package multimap

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chenjianyu/collections/container/set"
)

// HashMultimap is a multimap implementation that uses HashSet to store multiple values for a key
type HashMultimap[K comparable, V comparable] struct {
	data   map[K]set.Set[V]
	size   int
	mutex  sync.RWMutex
}

// NewHashMultimap creates a new HashMultimap
func NewHashMultimap[K comparable, V comparable]() *HashMultimap[K, V] {
	return &HashMultimap[K, V]{
		data: make(map[K]set.Set[V]),
		size: 0,
	}
}

// Put adds a key-value mapping to this multimap
func (m *HashMultimap[K, V]) Put(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		values = set.New[V]()
		m.data[key] = values
	}

	result := values.Add(value)
	if result {
		m.size++
	}

	return result
}

// PutAll adds all key-value mappings from the specified multimap to this multimap
func (m *HashMultimap[K, V]) PutAll(multimap Multimap[K, V]) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	changed := false
	multimap.ForEach(func(key K, value V) {
		values, exists := m.data[key]
		if !exists {
			values = set.New[V]()
			m.data[key] = values
		}

		result := values.Add(value)
		if result {
			m.size++
			changed = true
		}
	})

	return changed
}

// ReplaceValues replaces all values for a key with the specified collection of values
func (m *HashMultimap[K, V]) ReplaceValues(key K, values []V) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldValues, exists := m.data[key]
	if exists {
		oldValuesSlice := oldValues.ToSlice()
		m.size -= oldValues.Size()
		delete(m.data, key)

		if len(values) > 0 {
			newValues := set.New[V]()
			for _, value := range values {
				newValues.Add(value)
			}
			m.data[key] = newValues
			m.size += newValues.Size()
		}

		return oldValuesSlice
	} else if len(values) > 0 {
		newValues := set.New[V]()
		for _, value := range values {
			newValues.Add(value)
		}
		m.data[key] = newValues
		m.size += newValues.Size()
	}

	return nil
}

// Remove removes a key-value mapping from this multimap
func (m *HashMultimap[K, V]) Remove(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		return false
	}

	result := values.Remove(value)
	if result {
		m.size--
		if values.IsEmpty() {
			delete(m.data, key)
		}
	}

	return result
}

// RemoveAll removes all values associated with a key
func (m *HashMultimap[K, V]) RemoveAll(key K) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	result := values.ToSlice()
	m.size -= values.Size()
	delete(m.data, key)

	return result
}

// ContainsKey returns true if this multimap contains at least one key-value mapping with the specified key
func (m *HashMultimap[K, V]) ContainsKey(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.data[key]
	return exists
}

// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
func (m *HashMultimap[K, V]) ContainsValue(value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, values := range m.data {
		if values.Contains(value) {
			return true
		}
	}

	return false
}

// ContainsEntry returns true if this multimap contains the specified key-value mapping
func (m *HashMultimap[K, V]) ContainsEntry(key K, value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return false
	}

	return values.Contains(value)
}

// Get returns all values associated with the specified key
func (m *HashMultimap[K, V]) Get(key K) []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	return values.ToSlice()
}

// Keys returns all distinct keys in this multimap
func (m *HashMultimap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}

	return keys
}

// Values returns all values in this multimap
func (m *HashMultimap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, m.size)
	for _, valueSet := range m.data {
		values = append(values, valueSet.ToSlice()...)
	}

	return values
}

// Entries returns all key-value pairs in this multimap
func (m *HashMultimap[K, V]) Entries() []Entry[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entries := make([]Entry[K, V], 0, m.size)
	for key, valueSet := range m.data {
		for _, value := range valueSet.ToSlice() {
			entries = append(entries, Entry[K, V]{Key: key, Value: value})
		}
	}

	return entries
}

// KeySet returns a set view of the distinct keys in this multimap
func (m *HashMultimap[K, V]) KeySet() []K {
	return m.Keys()
}

// AsMap returns a map view of this multimap, mapping each key to its collection of values
func (m *HashMultimap[K, V]) AsMap() map[K][]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[K][]V, len(m.data))
	for key, valueSet := range m.data {
		result[key] = valueSet.ToSlice()
	}

	return result
}

// ForEach executes the given function for each key-value pair in this multimap
func (m *HashMultimap[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for key, valueSet := range m.data {
		for _, value := range valueSet.ToSlice() {
			f(key, value)
		}
	}
}

// Size returns the number of key-value mappings in this multimap
func (m *HashMultimap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size
}

// IsEmpty returns true if this multimap contains no key-value mappings
func (m *HashMultimap[K, V]) IsEmpty() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size == 0
}

// Clear removes all key-value mappings from this multimap
func (m *HashMultimap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[K]set.Set[V])
	m.size = 0
}

// Contains returns true if this multimap contains the specified element
func (m *HashMultimap[K, V]) Contains(key K) bool {
	return m.ContainsKey(key)
}

// String returns a string representation of this multimap
func (m *HashMultimap[K, V]) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.IsEmpty() {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")

	first := true
	for key, valueSet := range m.data {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		builder.WriteString(fmt.Sprintf("%v=[%v]", key, formatValues(valueSet.ToSlice())))
	}

	builder.WriteString("}")
	return builder.String()
}