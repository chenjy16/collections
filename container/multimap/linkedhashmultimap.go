package multimap

import (
    "fmt"
    "strings"
    "sync"

    "github.com/chenjianyu/collections/container/common"
    "github.com/chenjianyu/collections/container/set"
)

// LinkedHashMultimap is a multimap implementation that maintains insertion order of keys and values
type LinkedHashMultimap[K comparable, V comparable] struct {
	data   map[K]set.Set[V]
	keys   []K                 // Maintains insertion order of keys
	values map[K]map[V]struct{} // Tracks insertion order of values for each key
	size   int
	mutex  sync.RWMutex
}

// NewLinkedHashMultimap creates a new LinkedHashMultimap
func NewLinkedHashMultimap[K comparable, V comparable]() *LinkedHashMultimap[K, V] {
	return &LinkedHashMultimap[K, V]{
		data:   make(map[K]set.Set[V]),
		keys:   make([]K, 0),
		values: make(map[K]map[V]struct{}),
		size:   0,
	}
}

// Put adds a key-value mapping to this multimap
func (m *LinkedHashMultimap[K, V]) Put(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		values = set.NewLinkedHashSet[V]()
		m.data[key] = values
		m.keys = append(m.keys, key)
		m.values[key] = make(map[V]struct{})
	}

	result := values.Add(value)
	if result {
		m.size++
		// Track insertion order of values
		m.values[key][value] = struct{}{}
	}

	return result
}

// PutAll adds all key-value mappings from the specified multimap to this multimap
func (m *LinkedHashMultimap[K, V]) PutAll(multimap Multimap[K, V]) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	changed := false
	
	// Preserve insertion order by iterating through entries
	for _, entry := range multimap.Entries() {
		key := entry.Key
		value := entry.Value
		
		values, exists := m.data[key]
		if !exists {
			values = set.NewLinkedHashSet[V]()
			m.data[key] = values
			m.keys = append(m.keys, key)
			m.values[key] = make(map[V]struct{})
		}

		result := values.Add(value)
		if result {
			m.size++
			m.values[key][value] = struct{}{}
			changed = true
		}
	}

	return changed
}

// ReplaceValues replaces all values for a key with the specified collection of values
func (m *LinkedHashMultimap[K, V]) ReplaceValues(key K, values []V) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldValues, exists := m.data[key]
	if exists {
		oldValuesSlice := oldValues.ToSlice()
		m.size -= oldValues.Size()
		delete(m.data, key)
		delete(m.values, key)

		if len(values) > 0 {
			newValues := set.NewLinkedHashSet[V]()
			m.values[key] = make(map[V]struct{})
			
			for _, value := range values {
				newValues.Add(value)
				m.values[key][value] = struct{}{}
			}
			
			m.data[key] = newValues
			m.size += newValues.Size()
		} else {
			// Remove key from keys slice if no values remain
			for i, k := range m.keys {
				if k == key {
					m.keys = append(m.keys[:i], m.keys[i+1:]...)
					break
				}
			}
		}

		return oldValuesSlice
	} else if len(values) > 0 {
		newValues := set.NewLinkedHashSet[V]()
		m.keys = append(m.keys, key)
		m.values[key] = make(map[V]struct{})
		
		for _, value := range values {
			newValues.Add(value)
			m.values[key][value] = struct{}{}
		}
		
		m.data[key] = newValues
		m.size += newValues.Size()
	}

	return nil
}

// Remove removes a key-value mapping from this multimap
func (m *LinkedHashMultimap[K, V]) Remove(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		return false
	}

	result := values.Remove(value)
	if result {
		m.size--
		delete(m.values[key], value)
		
		if values.IsEmpty() {
			delete(m.data, key)
			delete(m.values, key)
			
			// Remove key from keys slice
			for i, k := range m.keys {
				if k == key {
					m.keys = append(m.keys[:i], m.keys[i+1:]...)
					break
				}
			}
		}
	}

	return result
}

// RemoveAll removes all values associated with a key
func (m *LinkedHashMultimap[K, V]) RemoveAll(key K) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	result := values.ToSlice()
	m.size -= values.Size()
	delete(m.data, key)
	delete(m.values, key)
	
	// Remove key from keys slice
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}

	return result
}

// ContainsKey returns true if this multimap contains at least one key-value mapping with the specified key
func (m *LinkedHashMultimap[K, V]) ContainsKey(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.data[key]
	return exists
}

// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
func (m *LinkedHashMultimap[K, V]) ContainsValue(value V) bool {
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
func (m *LinkedHashMultimap[K, V]) ContainsEntry(key K, value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return false
	}

	return values.Contains(value)
}

// Get returns all values associated with the specified key
func (m *LinkedHashMultimap[K, V]) Get(key K) []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	return values.ToSlice()
}

// Keys returns all distinct keys in this multimap in insertion order
func (m *LinkedHashMultimap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return result
}

// Values returns all values in this multimap in insertion order
func (m *LinkedHashMultimap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, m.size)
	
	// Iterate through keys in insertion order
	for _, key := range m.keys {
		valueSet := m.data[key]
		values = append(values, valueSet.ToSlice()...)
	}

	return values
}

// Entries returns all key-value pairs in this multimap in insertion order
func (m *LinkedHashMultimap[K, V]) Entries() []common.Entry[K, V] {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    entries := make([]common.Entry[K, V], 0, m.size)
	
	// Iterate through keys in insertion order
	for _, key := range m.keys {
		valueSet := m.data[key]
		for _, value := range valueSet.ToSlice() {
            entries = append(entries, common.NewEntry[K, V](key, value))
        }
    }

    return entries
}

// KeySet returns a set view of the distinct keys in this multimap in insertion order
func (m *LinkedHashMultimap[K, V]) KeySet() []K {
	return m.Keys()
}

// AsMap returns a map view of this multimap, mapping each key to its collection of values
func (m *LinkedHashMultimap[K, V]) AsMap() map[K][]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[K][]V, len(m.data))
	for key, valueSet := range m.data {
		result[key] = valueSet.ToSlice()
	}

	return result
}

// ForEach executes the given function for each key-value pair in this multimap in insertion order
func (m *LinkedHashMultimap[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Iterate through keys in insertion order
	for _, key := range m.keys {
		valueSet := m.data[key]
		for _, value := range valueSet.ToSlice() {
			f(key, value)
		}
	}
}

// Size returns the number of key-value mappings in this multimap
func (m *LinkedHashMultimap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size
}

// IsEmpty returns true if this multimap contains no key-value mappings
func (m *LinkedHashMultimap[K, V]) IsEmpty() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size == 0
}

// Clear removes all key-value mappings from this multimap
func (m *LinkedHashMultimap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[K]set.Set[V])
	m.keys = make([]K, 0)
	m.values = make(map[K]map[V]struct{})
	m.size = 0
}

// Contains returns true if this multimap contains the specified element
func (m *LinkedHashMultimap[K, V]) Contains(key K) bool {
	return m.ContainsKey(key)
}

// String returns a string representation of this multimap
func (m *LinkedHashMultimap[K, V]) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.IsEmpty() {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")

	first := true
	for _, key := range m.keys {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		valueSet := m.data[key]
		builder.WriteString(fmt.Sprintf("%v=[%v]", key, formatValues(valueSet.ToSlice())))
	}

	builder.WriteString("}")
	return builder.String()
}