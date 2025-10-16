// Package common provides common utilities and interfaces for the container library
package common

import (
    "fmt"
    "hash"
    "hash/fnv"
    "math"
    "reflect"
    "sort"
    "strings"
)

// Equal compares two values for equality
// Uses reflect.DeepEqual for deep comparison
func Equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Hash calculates the hash code of a value
// Uses a unified approach with proper handling of all types including floats
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
		writeInt64ToHash(h, v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		writeUint64ToHash(h, v.Uint())
	case reflect.Float32:
		// Handle float32 properly by converting to bits
		f32 := float32(v.Float())
		bits := math.Float32bits(f32)
		writeUint64ToHash(h, uint64(bits))
	case reflect.Float64:
		// Handle float64 properly by converting to bits
		f64 := v.Float()
		bits := math.Float64bits(f64)
		writeUint64ToHash(h, bits)
	case reflect.String:
		h.Write([]byte(v.String()))
    case reflect.Slice, reflect.Array:
        // Hash length first to distinguish between different length arrays
        writeInt64ToHash(h, int64(v.Len()))
        for i := 0; i < v.Len(); i++ {
            hashValue(v.Index(i), h)
        }
    case reflect.Map:
        // Deterministic hashing for maps: hash length and entries in sorted key order
        keys := v.MapKeys()
        writeInt64ToHash(h, int64(len(keys)))
        // Sort keys using the generic comparator for stability across runs
        sort.Slice(keys, func(i, j int) bool {
            return Compare(keys[i].Interface(), keys[j].Interface()) < 0
        })
        for _, k := range keys {
            // Hash key then value for each entry
            hashValue(k, h)
            val := v.MapIndex(k)
            hashValue(val, h)
        }
    case reflect.Struct:
        for i := 0; i < v.NumField(); i++ {
            hashValue(v.Field(i), h)
        }
	case reflect.Ptr:
		if !v.IsNil() {
			hashValue(v.Elem(), h)
		} else {
			// Hash nil pointer as a special value
			h.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF})
		}
	case reflect.Interface:
		if !v.IsNil() {
			hashValue(v.Elem(), h)
		} else {
			// Hash nil interface as a special value
			h.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF})
		}
	default:
		// For other types, use their string representation
		h.Write([]byte(fmt.Sprintf("%v", v.Interface())))
	}
}

// writeInt64ToHash writes int64 value to hash
func writeInt64ToHash(h hash.Hash64, val int64) {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(val >> (8 * i))
	}
	h.Write(bytes)
}

// writeUint64ToHash writes uint64 value to hash
func writeUint64ToHash(h hash.Hash64, val uint64) {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(val >> (8 * i))
	}
	h.Write(bytes)
}

// Compare compares two values
// Returns:
//   - negative number: a < b
//   - zero: a == b
//   - positive number: a > b
func Compare(a, b interface{}) int {
	// Handle nil values
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try type assertion for comparable types
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case int8:
		if vb, ok := b.(int8); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case int16:
		if vb, ok := b.(int16); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case int32:
		if vb, ok := b.(int32); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case int64:
		if vb, ok := b.(int64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case uint:
		if vb, ok := b.(uint); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case uint8:
		if vb, ok := b.(uint8); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case uint16:
		if vb, ok := b.(uint16); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case uint32:
		if vb, ok := b.(uint32); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case uint64:
		if vb, ok := b.(uint64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case float32:
		if vb, ok := b.(float32); ok {
			// Handle NaN explicitly for deterministic ordering
			if math.IsNaN(float64(va)) {
				if math.IsNaN(float64(vb)) {
					return 0
				}
				return -1
			}
			if math.IsNaN(float64(vb)) {
				return 1
			}
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case float64:
		if vb, ok := b.(float64); ok {
			// Handle NaN explicitly for deterministic ordering
			if math.IsNaN(va) {
				if math.IsNaN(vb) {
					return 0
				}
				return -1
			}
			if math.IsNaN(vb) {
				return 1
			}
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case string:
		if vb, ok := b.(string); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	}

	// Check if both implement fmt.Stringer
	if sa, ok := a.(fmt.Stringer); ok {
		if sb, ok := b.(fmt.Stringer); ok {
			return strings.Compare(sa.String(), sb.String())
		}
	}

	// For types that cannot be directly compared, use hash values for comparison
	hashA := Hash(a)
	hashB := Hash(b)
	if hashA < hashB {
		return -1
	} else if hashA > hashB {
		return 1
	}
	return 0
}

// CompareGeneric compares two comparable values using generics
// Returns:
//   - negative number: a < b
//   - zero: a == b
//   - positive number: a > b
func CompareGeneric[T comparable](a, b T) int {
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
		// Handle NaN explicitly for deterministic ordering
		if math.IsNaN(float64(va)) {
			if math.IsNaN(float64(vb)) {
				return 0
			}
			return -1
		}
		if math.IsNaN(float64(vb)) {
			return 1
		}
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	case float64:
		vb := any(b).(float64)
		// Handle NaN explicitly for deterministic ordering
		if math.IsNaN(va) {
			if math.IsNaN(vb) {
				return 0
			}
			return -1
		}
		if math.IsNaN(vb) {
			return 1
		}
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

// CompareNatural first uses Comparable.CompareTo when available, otherwise falls back to CompareGeneric
// This provides consistent "natural" ordering across custom comparable types and built-ins.
func CompareNatural[T comparable](a, b T) int {
    if ca, ok := any(a).(Comparable); ok {
        return ca.CompareTo(any(b))
    }
    return CompareGeneric[T](a, b)
}

// Entry represents a key-value pair used across Map/Multimap APIs
type Entry[K, V any] struct {
    Key   K
    Value V
}

// NewEntry creates a new key-value pair
func NewEntry[K, V any](key K, value V) Entry[K, V] {
    return Entry[K, V]{Key: key, Value: value}
}

// String returns the string representation of the Entry
func (e Entry[K, V]) String() string {
    return fmt.Sprintf("(%v, %v)", e.Key, e.Value)
}

// Pair is a compatibility struct mirroring Entry for older Go versions
// Deprecated: prefer Entry/NewEntry directly.
type Pair[K, V any] struct {
    Key   K
    Value V
}

// String returns the string representation of the Pair
func (p Pair[K, V]) String() string {
    return fmt.Sprintf("(%v, %v)", p.Key, p.Value)
}

// NewPair is a compatibility constructor that forwards to NewEntry.
// Deprecated: use NewEntry instead.
func NewPair[K, V any](key K, value V) Pair[K, V] {
    return Pair[K, V]{Key: key, Value: value}
}

// ZeroValue returns the zero value of type T
// This is a cleaner alternative to *new(T)
func ZeroValue[T any]() T {
	var zero T
	return zero
}
