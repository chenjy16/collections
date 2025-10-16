package maps

import "github.com/chenjianyu/collections/container/common"

// Map is a generic map interface
type Map[K comparable, V any] interface {
	// Put associates the specified value with the specified key in this map
	Put(key K, value V) (V, bool)

	// Get returns the value mapped to the specified key
	Get(key K) (V, bool)

	// Remove if exists, removes mapping relationship for the key from this map
	Remove(key K) (V, bool)

	// ContainsKey if this map contains mapping relationship for the specified key, returns true
	ContainsKey(key K) bool

	// ContainsValue if this map maps one or more keys to the specified value, returns true
	ContainsValue(value V) bool

	// Size returns the number of key-value mapping relationships in this map
	Size() int

	// IsEmpty if this map does not contain key-value mapping relationships, returns true
	IsEmpty() bool

	// Clear removes all mapping relationships from this map
	Clear()

	// Keys returns the keys contained in this map
	Keys() []K

	// Values returns the values contained in this map
	Values() []V

    // Entries returns the mapping relationships contained in this map
    Entries() []common.Entry[K, V]

	// ForEach executes the given operation for each entry in this map
	ForEach(f func(K, V))

	// String returns the string representation of the map
	String() string

	// PutAll copies all mapping relationships from the specified map to this map
	PutAll(other Map[K, V])
}
