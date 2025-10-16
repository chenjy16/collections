// Package multimap provides multimap data structure implementations
package multimap

import (
    "github.com/chenjianyu/collections/container/common"
)

// Multimap represents a collection that maps keys to multiple values
type Multimap[K comparable, V comparable] interface {
    // Size returns the number of key-value mappings in this multimap
    Size() int

    // IsEmpty returns true if this multimap contains no key-value mappings
    IsEmpty() bool

    // Clear removes all key-value mappings from this multimap
    Clear()

    // Contains returns true if this multimap contains the specified key
    Contains(key K) bool

    // String returns the string representation of this multimap
    String() string

	// Put adds a key-value mapping to this multimap
	// Returns true if the multimap changed
	Put(key K, value V) bool

	// PutAll adds all key-value mappings from the specified multimap to this multimap
	// Returns true if the multimap changed
	PutAll(multimap Multimap[K, V]) bool

	// ReplaceValues replaces all values for a key with the specified collection of values
	// Returns the previous values associated with the key, or nil if none
	ReplaceValues(key K, values []V) []V

	// Remove removes a key-value mapping from this multimap
	// Returns true if the multimap changed
	Remove(key K, value V) bool

	// RemoveAll removes all values associated with a key
	// Returns the values that were removed, or nil if none
	RemoveAll(key K) []V

	// ContainsKey returns true if this multimap contains at least one key-value mapping with the specified key
	ContainsKey(key K) bool

	// ContainsValue returns true if this multimap contains at least one key-value mapping with the specified value
	ContainsValue(value V) bool

	// ContainsEntry returns true if this multimap contains the specified key-value mapping
	ContainsEntry(key K, value V) bool

	// Get returns all values associated with the specified key
	// Returns nil if no values are associated with the key
	Get(key K) []V

	// Keys returns all distinct keys in this multimap
	Keys() []K

    // Values returns all values in this multimap
    Values() []V

    // Entries returns all key-value pairs in this multimap
    Entries() []common.Entry[K, V]

	// KeySet returns a set view of the distinct keys in this multimap
	KeySet() []K

	// AsMap returns a map view of this multimap, mapping each key to its collection of values
	AsMap() map[K][]V

	// ForEach executes the given function for each key-value pair in this multimap
	ForEach(func(K, V))
}

// Deprecated compatibility alias removed; use common.Entry/common.NewEntry directly.