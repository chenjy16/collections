package multimap

import (
    "fmt"
    "strings"
    "sync"

    "github.com/chenjianyu/collections/container/common"
    "github.com/chenjianyu/collections/container/set"
)

// TreeMultimap is a multimap implementation that maintains keys in sorted order
type TreeMultimap[K comparable, V comparable] struct {
    data          map[K]set.Set[V]
    keys          []K // Maintains sorted order of keys
    size          int
    mutex         sync.RWMutex
    keyComparator func(a, b K) int
}

// NewTreeMultimap creates a new TreeMultimap
func NewTreeMultimap[K comparable, V comparable]() *TreeMultimap[K, V] {
    return &TreeMultimap[K, V]{
        data:          make(map[K]set.Set[V]),
        keys:          make([]K, 0),
        size:          0,
        keyComparator: defaultKeyComparator[K],
    }
}

// NewTreeMultimapWithComparator creates a new TreeMultimap with a custom key comparator
func NewTreeMultimapWithComparator[K comparable, V comparable](keyCmp func(a, b K) int) *TreeMultimap[K, V] {
    if keyCmp == nil {
        keyCmp = defaultKeyComparator[K]
    }
    return &TreeMultimap[K, V]{
        data:          make(map[K]set.Set[V]),
        keys:          make([]K, 0),
        size:          0,
        keyComparator: keyCmp,
    }
}

// sortKeys sorts the keys slice
func (m *TreeMultimap[K, V]) sortKeys() {
	// Sort keys using insertion sort for simplicity
	for i := 1; i < len(m.keys); i++ {
		key := m.keys[i]
		j := i - 1
        // Compare keys using injected key comparator
        for j >= 0 && m.keyComparator(m.keys[j], key) > 0 {
            m.keys[j+1] = m.keys[j]
            j--
        }
        m.keys[j+1] = key
    }
}

// defaultKeyComparator prefers Comparable.CompareTo; otherwise falls back to natural generic ordering
func defaultKeyComparator[K comparable](a, b K) int { return common.CompareNatural[K](a, b) }

// defaultValueComparator prefers Comparable.CompareTo; otherwise falls back to natural generic ordering
func defaultValueComparator[V comparable](a, b V) int { return common.CompareNatural[V](a, b) }

// findKeyIndex finds the index of a key in the sorted keys slice
func (m *TreeMultimap[K, V]) findKeyIndex(key K) int {
    for i, k := range m.keys {
        if m.keyComparator(k, key) == 0 {
            return i
        }
    }
    return -1
}

// Put adds a key-value mapping to this multimap
func (m *TreeMultimap[K, V]) Put(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

    values, exists := m.data[key]
    if !exists {
        values = set.NewTreeSetWithComparator[V](defaultValueComparator[V])
        m.data[key] = values
        m.keys = append(m.keys, key)
        m.sortKeys()
    }

	result := values.Add(value)
	if result {
		m.size++
	}

	return result
}

// PutAll adds all key-value mappings from the specified multimap to this multimap
func (m *TreeMultimap[K, V]) PutAll(multimap Multimap[K, V]) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	changed := false
    multimap.ForEach(func(key K, value V) {
        values, exists := m.data[key]
        if !exists {
            values = set.NewTreeSetWithComparator[V](defaultValueComparator[V])
            m.data[key] = values
            m.keys = append(m.keys, key)
            // We'll sort keys once at the end for efficiency
        }

		result := values.Add(value)
		if result {
			m.size++
			changed = true
		}
	})

	// Sort keys if any were added
	if changed {
		m.sortKeys()
	}

	return changed
}

// ReplaceValues replaces all values for a key with the specified collection of values
func (m *TreeMultimap[K, V]) ReplaceValues(key K, values []V) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldValues, exists := m.data[key]
    if exists {
        oldValuesSlice := oldValues.ToSlice()
        m.size -= oldValues.Size()
        delete(m.data, key)

        if len(values) > 0 {
            newValues := set.NewTreeSetWithComparator[V](defaultValueComparator[V])
            for _, value := range values {
                newValues.Add(value)
            }
            m.data[key] = newValues
            m.size += newValues.Size()
        } else {
			// Remove key from keys slice if no values remain
			index := m.findKeyIndex(key)
			if index >= 0 {
				m.keys = append(m.keys[:index], m.keys[index+1:]...)
			}
		}

		return oldValuesSlice
    } else if len(values) > 0 {
        newValues := set.NewTreeSetWithComparator[V](defaultValueComparator[V])
        for _, value := range values {
            newValues.Add(value)
        }
        m.data[key] = newValues
        m.keys = append(m.keys, key)
		m.sortKeys()
		m.size += newValues.Size()
	}

	return nil
}

// Remove removes a key-value mapping from this multimap
func (m *TreeMultimap[K, V]) Remove(key K, value V) bool {
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
			// Remove key from keys slice
			index := m.findKeyIndex(key)
			if index >= 0 {
				m.keys = append(m.keys[:index], m.keys[index+1:]...)
			}
		}
	}

	return result
}

// RemoveAll removes all values associated with a key
func (m *TreeMultimap[K, V]) RemoveAll(key K) []V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	result := values.ToSlice()
	m.size -= values.Size()
	delete(m.data, key)
	
	// Remove key from keys slice
	index := m.findKeyIndex(key)
	if index >= 0 {
		m.keys = append(m.keys[:index], m.keys[index+1:]...)
	}

	return result
}

// ContainsKey returns true if this multimap contains at least one key-value mapping with the specified key
func (m *TreeMultimap[K, V]) ContainsKey(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.data[key]
	return exists
}

// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
func (m *TreeMultimap[K, V]) ContainsValue(value V) bool {
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
func (m *TreeMultimap[K, V]) ContainsEntry(key K, value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return false
	}

	return values.Contains(value)
}

// Get returns all values associated with the specified key
func (m *TreeMultimap[K, V]) Get(key K) []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values, exists := m.data[key]
	if !exists {
		return nil
	}

	return values.ToSlice()
}

// Keys returns all distinct keys in this multimap in sorted order
func (m *TreeMultimap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return result
}

// Values returns all values in this multimap
func (m *TreeMultimap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, m.size)
	
	// Iterate through keys in sorted order
	for _, key := range m.keys {
		valueSet := m.data[key]
		values = append(values, valueSet.ToSlice()...)
	}

	return values
}

// Entries returns all key-value pairs in this multimap
func (m *TreeMultimap[K, V]) Entries() []common.Entry[K, V] {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    entries := make([]common.Entry[K, V], 0, m.size)
	
	// Iterate through keys in sorted order
	for _, key := range m.keys {
		valueSet := m.data[key]
        for _, value := range valueSet.ToSlice() {
            entries = append(entries, common.NewEntry[K, V](key, value))
        }
	}

    return entries
}

// KeySet returns a set view of the distinct keys in this multimap
func (m *TreeMultimap[K, V]) KeySet() []K {
	return m.Keys()
}

// AsMap returns a map view of this multimap, mapping each key to its collection of values
func (m *TreeMultimap[K, V]) AsMap() map[K][]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[K][]V, len(m.data))
	for key, valueSet := range m.data {
		result[key] = valueSet.ToSlice()
	}

	return result
}

// ForEach executes the given function for each key-value pair in this multimap
func (m *TreeMultimap[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Iterate through keys in sorted order
	for _, key := range m.keys {
		valueSet := m.data[key]
		for _, value := range valueSet.ToSlice() {
			f(key, value)
		}
	}
}

// Size returns the number of key-value mappings in this multimap
func (m *TreeMultimap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size
}

// IsEmpty returns true if this multimap contains no key-value mappings
func (m *TreeMultimap[K, V]) IsEmpty() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.size == 0
}

// Clear removes all key-value mappings from this multimap
func (m *TreeMultimap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[K]set.Set[V])
	m.keys = make([]K, 0)
	m.size = 0
}

// Contains returns true if this multimap contains the specified element
func (m *TreeMultimap[K, V]) Contains(key K) bool {
	return m.ContainsKey(key)
}

// String returns a string representation of this multimap
func (m *TreeMultimap[K, V]) String() string {
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
		builder.WriteString(fmt.Sprintf("%v=[%v]", key, formatValues[V](valueSet.ToSlice())))
	}

	builder.WriteString("}")
	return builder.String()
}