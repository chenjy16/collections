package maps

// Pair is a simple struct used by MapOf helpers; mirrors common.Entry
// Deprecated: prefer common.Entry/common.NewEntry in new code.
type Pair[K, V any] struct {
    Key   K
    Value V
}

// NewPair constructs a Pair for use with MapOf and helpers.
func NewPair[K, V any](key K, value V) Pair[K, V] {
    return Pair[K, V]{Key: key, Value: value}
}