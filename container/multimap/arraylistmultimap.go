package multimap

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chenjianyu/collections/container/list"
)

// ArrayListMultimap is a multimap implementation that uses ArrayList to store multiple values for a key
type ArrayListMultimap[K comparable, V comparable] struct {
	data   map[K]list.List[V]
	size   int
	mutex  sync.RWMutex
}

// NewArrayListMultimap creates a new ArrayListMultimap
func NewArrayListMultimap[K comparable, V comparable]() *ArrayListMultimap[K, V] {
	return &ArrayListMultimap[K, V]{
		data: make(map[K]list.List[V]),
		size: 0,
	}
}

// Put adds a key-value mapping to this multimap
func (m *ArrayListMultimap[K, V]) Put(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		values = list.New[V]()
		m.data[key] = values
	}

	result := values.Add(value)
	if result {
		m.size++
	}

	return result
}

// PutAll adds all key-value mappings from the specified multimap to this multimap
func (m *ArrayListMultimap[K, V]) PutAll(multimap Multimap[K, V]) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	changed := false
	multimap.ForEach(func(key K, value V) {
		values, exists := m.data[key]
		if !exists {
			values = list.New[V]()
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
func (m *ArrayListMultimap[K, V]) ReplaceValues(key K, values []V) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldValues, exists := m.data[key]
	if exists {
		oldValuesSlice := oldValues.ToSlice()
		m.size -= oldValues.Size()
		delete(m.data, key)

		if len(values) > 0 {
			newValues := list.New[V]()
			for _, value := range values {
				newValues.Add(value)
			}
			m.data[key] = newValues
			m.size += newValues.Size()
		}

		return oldValuesSlice
	} else if len(values) > 0 {
		newValues := list.New[V]()
		for _, value := range values {
			newValues.Add(value)
		}
		m.data[key] = newValues
		m.size += newValues.Size()
	}

	return nil
}

// Remove removes a key-value mapping from this multimap
func (m *ArrayListMultimap[K, V]) Remove(key K, value V) bool {
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
func (m *ArrayListMultimap[K, V]) RemoveAll(key K) []V {
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
func (m *ArrayListMultimap[K, V]) ContainsKey(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.data[key]
	return exists
}

// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
func (m *ArrayListMultimap[K, V]) ContainsValue(value V) bool {
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
func (m *ArrayListMultimap[K, V]) ContainsEntry(key K, value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return false
	}

	return values.Contains(value)
}

// Get returns all values associated with the specified key
func (m *ArrayListMultimap[K, V]) Get(key K) []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	return values.ToSlice()
}

// Keys returns all distinct keys in this multimap
func (m *ArrayListMultimap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}

	return keys
}

// Values returns all values in this multimap
func (m *ArrayListMultimap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, m.size)
	for _, valueList := range m.data {
		values = append(values, valueList.ToSlice()...)
	}

	return values
}

// Entries returns all key-value pairs in this multimap
func (m *ArrayListMultimap[K, V]) Entries() []Entry[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entries := make([]Entry[K, V], 0, m.size)
	for key, valueList := range m.data {
		for _, value := range valueList.ToSlice() {
			entries = append(entries, Entry[K, V]{Key: key, Value: value})
		}
	}

	return entries
}

// KeySet returns a set view of the distinct keys in this multimap
func (m *ArrayListMultimap[K, V]) KeySet() []K {
	return m.Keys()
}

// AsMap returns a map view of this multimap, mapping each key to its collection of values
func (m *ArrayListMultimap[K, V]) AsMap() map[K][]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[K][]V, len(m.data))
	for key, valueList := range m.data {
		result[key] = valueList.ToSlice()
	}

	return result
}

// ForEach executes the given function for each key-value pair in this multimap
func (m *ArrayListMultimap[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for key, valueList := range m.data {
		for _, value := range valueList.ToSlice() {
			f(key, value)
		}
	}
}

// Size returns the number of key-value mappings in this multimap
func (m *ArrayListMultimap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size
}

// IsEmpty returns true if this multimap contains no key-value mappings
func (m *ArrayListMultimap[K, V]) IsEmpty() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size == 0
}

// Clear removes all key-value mappings from this multimap
func (m *ArrayListMultimap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[K]list.List[V])
	m.size = 0
}

// Contains returns true if this multimap contains the specified element
func (m *ArrayListMultimap[K, V]) Contains(key K) bool {
	return m.ContainsKey(key)
}

// String returns a string representation of this multimap
func (m *ArrayListMultimap[K, V]) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.IsEmpty() {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")

	first := true
	for key, valueList := range m.data {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		builder.WriteString(fmt.Sprintf("%v=[%v]", key, formatValues(valueList.ToSlice())))
	}

	builder.WriteString("}")
	return builder.String()
}

// formatValues formats a slice of values as a comma-separated string
func formatValues[V comparable](values []V) string {
	if len(values) == 0 {
		return ""
	}

	var builder strings.Builder
	for i, value := range values {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", value))
	}

	return builder.String()
}