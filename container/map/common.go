package maps

import (
	"fmt"
	"hash"
	"hash/fnv"
	"reflect"
)

// color represents the color of red-black tree nodes
type color bool

const (
	red   color = true
	black color = false
)

// Pair represents a key-value pair
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// NewPair creates a new key-value pair
func NewPair[K comparable, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// String returns the string representation of the pair
func (p Pair[K, V]) String() string {
	return fmt.Sprintf("(%v, %v)", p.Key, p.Value)
}

// Hash calculates the hash value of any type
func Hash(v interface{}) uint64 {
	h := fnv.New64a()
	hashValue(reflect.ValueOf(v), h)
	return h.Sum64()
}

// hashValue recursively calculates hash value
func hashValue(v reflect.Value, h hash.Hash64) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			h.Write([]byte{1})
		} else {
			h.Write([]byte{0})
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		writeInt64(h, v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		writeUint64(h, v.Uint())
	case reflect.Float32, reflect.Float64:
		writeUint64(h, v.Uint())
	case reflect.String:
		h.Write([]byte(v.String()))
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			hashValue(v.Index(i), h)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			hashValue(v.Field(i), h)
		}
	case reflect.Ptr:
		if !v.IsNil() {
			hashValue(v.Elem(), h)
		}
	case reflect.Interface:
		if !v.IsNil() {
			hashValue(v.Elem(), h)
		}
	default:
		// For other types, use their string representation
		h.Write([]byte(fmt.Sprintf("%v", v.Interface())))
	}
}

// writeInt64 writes int64 value to hash
func writeInt64(h hash.Hash64, val int64) {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(val >> (8 * i))
	}
	h.Write(bytes)
}

// writeUint64 writes uint64 value to hash
func writeUint64(h hash.Hash64, val uint64) {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(val >> (8 * i))
	}
	h.Write(bytes)
}

// Equal checks if two values are equal
func Equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Compare compares two comparable values
func Compare[T comparable](a, b T) int {
	// Use type assertion to handle specific types
	switch va := any(a).(type) {
	case int:
		vb := any(b).(int)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case int8:
		vb := any(b).(int8)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case int16:
		vb := any(b).(int16)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case int32:
		vb := any(b).(int32)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case int64:
		vb := any(b).(int64)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case uint:
		vb := any(b).(uint)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case uint8:
		vb := any(b).(uint8)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case uint16:
		vb := any(b).(uint16)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case uint32:
		vb := any(b).(uint32)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case uint64:
		vb := any(b).(uint64)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case float32:
		vb := any(b).(float32)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case float64:
		vb := any(b).(float64)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case string:
		vb := any(b).(string)
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	default:
		// For other comparable types, check equality first
		if a == b {
			return 0
		}
		// For non-comparable types or when equality fails, use hash values for comparison
		hashA := Hash(a)
		hashB := Hash(b)
		if hashA < hashB {
			return -1
		} else if hashA > hashB {
			return 1
		}
		return 0
	}
}
