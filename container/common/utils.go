// Package common provides common utilities and interfaces for the container library
package common

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"strings"
)

// Equal compares two values for equality
// Uses reflect.DeepEqual for deep comparison
func Equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Hash calculates the hash code of a value
// Uses FNV-1a algorithm to hash the string representation of the value
func Hash(v interface{}) uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	return h.Sum64()
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
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case float64:
		if vb, ok := b.(float64); ok {
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

// Pair represents a key-value pair
type Pair[K, V any] struct {
	Key   K
	Value V
}

// NewPair creates a new key-value pair
func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// String returns the string representation of the Pair
func (p Pair[K, V]) String() string {
	return fmt.Sprintf("(%v, %v)", p.Key, p.Value)
}
