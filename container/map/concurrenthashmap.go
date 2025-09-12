package maps

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/chenjianyu/collections/container/common"
)

// ConcurrentHashMap is a thread-safe hash map implementation
// using segmented locking strategy, similar to Java's ConcurrentHashMap
type ConcurrentHashMap[K comparable, V any] struct {
	segments     []*segment[K, V]
	segmentMask  uint32
	segmentCount int
	hashStrategy common.HashStrategy[K]
}

// segment represents a segment of ConcurrentHashMap
type segment[K comparable, V any] struct {
	buckets []bucket[K, V]
	mutex   sync.RWMutex
	size    int
}

// bucket represents a hash bucket, using linked list to resolve conflicts
type bucket[K, V any] struct {
	key   K
	value V
	next  *bucket[K, V]
}

// Default number of segments, must be a power of 2
const defaultSegments = 16

// Default initial bucket count per segment
const defaultBuckets = 16

// Load factor threshold for concurrent hash map
const concurrentLoadFactor = 0.75

// NewConcurrentHashMap creates a new ConcurrentHashMap with default hash strategy
func NewConcurrentHashMap[K comparable, V any]() *ConcurrentHashMap[K, V] {
	return NewConcurrentHashMapWithHashStrategy[K, V](common.NewComparableHashStrategy[K]())
}

// NewConcurrentHashMapWithHashStrategy creates a new ConcurrentHashMap with custom hash strategy
func NewConcurrentHashMapWithHashStrategy[K comparable, V any](hashStrategy common.HashStrategy[K]) *ConcurrentHashMap[K, V] {
	return NewConcurrentHashMapWithCapacityAndHashStrategy[K, V](defaultSegments*defaultBuckets, hashStrategy)
}

// NewConcurrentHashMapWithCapacity creates a ConcurrentHashMap with specified initial capacity and default hash strategy
func NewConcurrentHashMapWithCapacity[K comparable, V any](capacity int) *ConcurrentHashMap[K, V] {
	return NewConcurrentHashMapWithCapacityAndHashStrategy[K, V](capacity, common.NewComparableHashStrategy[K]())
}

// NewConcurrentHashMapWithCapacityAndHashStrategy creates a ConcurrentHashMap with specified initial capacity and custom hash strategy
func NewConcurrentHashMapWithCapacityAndHashStrategy[K comparable, V any](capacity int, hashStrategy common.HashStrategy[K]) *ConcurrentHashMap[K, V] {
	if capacity <= 0 {
		capacity = defaultSegments * defaultBuckets
	}

	segmentCount := defaultSegments
	bucketsPerSegment := capacity / segmentCount
	if bucketsPerSegment < defaultBuckets {
		bucketsPerSegment = defaultBuckets
	}

	segments := make([]*segment[K, V], segmentCount)
	for i := range segments {
		segments[i] = &segment[K, V]{
			buckets: make([]bucket[K, V], bucketsPerSegment),
		}
	}

	return &ConcurrentHashMap[K, V]{
		segments:     segments,
		segmentCount: segmentCount,
		segmentMask:  uint32(segmentCount - 1),
		hashStrategy: hashStrategy,
	}
}

// ConcurrentHashMapFromMap creates ConcurrentHashMap from existing map with default hash strategy
func ConcurrentHashMapFromMap[K comparable, V any](m map[K]V) *ConcurrentHashMap[K, V] {
	return ConcurrentHashMapFromMapWithHashStrategy[K, V](m, common.NewComparableHashStrategy[K]())
}

// ConcurrentHashMapFromMapWithHashStrategy creates ConcurrentHashMap from existing map with custom hash strategy
func ConcurrentHashMapFromMapWithHashStrategy[K comparable, V any](m map[K]V, hashStrategy common.HashStrategy[K]) *ConcurrentHashMap[K, V] {
	chm := NewConcurrentHashMapWithCapacityAndHashStrategy[K, V](len(m), hashStrategy)
	for k, v := range m {
		chm.Put(k, v)
	}
	return chm
}

// Put associates the specified value with the specified key in this map
func (chm *ConcurrentHashMap[K, V]) Put(key K, value V) (V, bool) {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucketPtr := &segment.buckets[bucketIndex]

	// Check if key already exists
	for current := bucketPtr.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			oldValue := current.value
			current.value = value
			return oldValue, true
		}
	}

	// Add new node
	newNode := &bucket[K, V]{
		key:   key,
		value: value,
		next:  bucketPtr.next,
	}
	bucketPtr.next = newNode
	segment.size++

	// Check if resize is needed
	if float64(segment.size) > float64(len(segment.buckets))*concurrentLoadFactor {
		chm.resizeSegment(segment)
	}

	var zero V
	return zero, false
}

// Get returns the value to which the specified key is mapped
func (chm *ConcurrentHashMap[K, V]) Get(key K) (V, bool) {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.RLock()
	defer segment.mutex.RUnlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucketPtr := &segment.buckets[bucketIndex]

	for current := bucketPtr.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			return current.value, true
		}
	}

	var zero V
	return zero, false
}

// Remove removes the mapping for the specified key from this map if present
func (chm *ConcurrentHashMap[K, V]) Remove(key K) (V, bool) {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucketPtr := &segment.buckets[bucketIndex]

	// If head node is the one to delete
	if bucketPtr.next != nil && chm.hashStrategy.Equals(bucketPtr.next.key, key) {
		removed := bucketPtr.next.value
		bucketPtr.next = bucketPtr.next.next
		segment.size--
		return removed, true
	}

	// Find the node to delete
	for current := bucketPtr.next; current != nil && current.next != nil; current = current.next {
		if chm.hashStrategy.Equals(current.next.key, key) {
			removed := current.next.value
			current.next = current.next.next
			segment.size--
			return removed, true
		}
	}

	var zero V
	return zero, false
}

// ContainsKey returns true if this map contains a mapping for the specified key
func (chm *ConcurrentHashMap[K, V]) ContainsKey(key K) bool {
	_, exists := chm.Get(key)
	return exists
}

// ContainsValue returns true if this map maps one or more keys to the specified value
func (chm *ConcurrentHashMap[K, V]) ContainsValue(value V) bool {
	for _, segment := range chm.segments {
		segment.mutex.RLock()
		for _, bkt := range segment.buckets {
			for current := bkt.next; current != nil; current = current.next {
				if chm.valueEquals(current.value, value) {
					segment.mutex.RUnlock()
					return true
				}
			}
		}
		segment.mutex.RUnlock()
	}
	return false
}

// Size returns the number of key-value mappings in this map
func (chm *ConcurrentHashMap[K, V]) Size() int {
	totalSize := 0
	for _, segment := range chm.segments {
		segment.mutex.RLock()
		totalSize += segment.size
		segment.mutex.RUnlock()
	}
	return totalSize
}

// IsEmpty returns true if this map contains no key-value mappings
func (chm *ConcurrentHashMap[K, V]) IsEmpty() bool {
	return chm.Size() == 0
}

// Clear removes all mappings from this map
func (chm *ConcurrentHashMap[K, V]) Clear() {
	for _, segment := range chm.segments {
		segment.mutex.Lock()
		for i := range segment.buckets {
			segment.buckets[i] = bucket[K, V]{}
		}
		segment.size = 0
		segment.mutex.Unlock()
	}
}

// Keys returns a collection view of the keys contained in this map
func (chm *ConcurrentHashMap[K, V]) Keys() []K {
	var keys []K
	for _, segment := range chm.segments {
		segment.mutex.RLock()
		for _, bkt := range segment.buckets {
			for current := bkt.next; current != nil; current = current.next {
				keys = append(keys, current.key)
			}
		}
		segment.mutex.RUnlock()
	}
	return keys
}

// Values returns a collection view of the values contained in this map
func (chm *ConcurrentHashMap[K, V]) Values() []V {
	var values []V
	for _, segment := range chm.segments {
		segment.mutex.RLock()
		for _, bkt := range segment.buckets {
			for current := bkt.next; current != nil; current = current.next {
				values = append(values, current.value)
			}
		}
		segment.mutex.RUnlock()
	}
	return values
}

// Entries returns a collection view of the key-value pairs contained in this map
func (chm *ConcurrentHashMap[K, V]) Entries() []Pair[K, V] {
	var entries []Pair[K, V]
	for _, segment := range chm.segments {
		segment.mutex.RLock()
		for _, bkt := range segment.buckets {
			for current := bkt.next; current != nil; current = current.next {
				entries = append(entries, Pair[K, V]{Key: current.key, Value: current.value})
			}
		}
		segment.mutex.RUnlock()
	}
	return entries
}

// ForEach performs the given action for each key-value pair in this map
func (chm *ConcurrentHashMap[K, V]) ForEach(action func(K, V)) {
	for _, segment := range chm.segments {
		segment.mutex.RLock()
		for _, bkt := range segment.buckets {
			for current := bkt.next; current != nil; current = current.next {
				action(current.key, current.value)
			}
		}
		segment.mutex.RUnlock()
	}
}

// String returns a string representation of this map
func (chm *ConcurrentHashMap[K, V]) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	first := true
	chm.ForEach(func(key K, value V) {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v=%v", key, value))
		first = false
	})
	builder.WriteString("}")
	return builder.String()
}

// Advanced operation methods

// PutIfAbsent adds key-value pair only if key doesn't exist
// It returns the old value and whether the key was inserted.
// If key exists: returns (old_value, false)
// If key doesn't exist: returns (zero_value, true)
func (chm *ConcurrentHashMap[K, V]) PutIfAbsent(key K, value V) (V, bool) {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucketPtr := &segment.buckets[bucketIndex]

	// Check if key already exists
	for current := bucketPtr.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			return current.value, false // Key exists, not inserted
		}
	}

	// Add new node
	newNode := &bucket[K, V]{
		key:   key,
		value: value,
		next:  bucketPtr.next,
	}
	bucketPtr.next = newNode
	segment.size++

	// Check if resize is needed
	if float64(segment.size) > float64(len(segment.buckets))*concurrentLoadFactor {
		chm.resizeSegment(segment)
	}

	var zero V
	return zero, true // Key was inserted
}

// Replace replaces the value for existing key only
func (chm *ConcurrentHashMap[K, V]) Replace(key K, value V) (V, bool) {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucket := &segment.buckets[bucketIndex]

	for current := bucket.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			oldValue := current.value
			current.value = value
			return oldValue, true
		}
	}

	var zero V
	return zero, false
}

// ReplaceIf replaces the value only if current value equals expected value
func (chm *ConcurrentHashMap[K, V]) ReplaceIf(key K, oldValue, newValue V) bool {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucket := &segment.buckets[bucketIndex]

	for current := bucket.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			if chm.valueEquals(current.value, oldValue) {
				current.value = newValue
				return true
			}
			return false
		}
	}

	return false
}

// PutAll copies all elements from another Map
func (chm *ConcurrentHashMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(key K, value V) {
		chm.Put(key, value)
	})
}

// PutAllFromMap copies all elements from Go native map
func (chm *ConcurrentHashMap[K, V]) PutAllFromMap(m map[K]V) {
	for key, value := range m {
		chm.Put(key, value)
	}
}

// ToMap converts to Go native map
func (chm *ConcurrentHashMap[K, V]) ToMap() map[K]V {
	result := make(map[K]V)
	chm.ForEach(func(key K, value V) {
		result[key] = value
	})
	return result
}

// Snapshot gets a snapshot of current data
func (chm *ConcurrentHashMap[K, V]) Snapshot() map[K]V {
	return chm.ToMap()
}

// ComputeIfAbsent computes and adds value if key doesn't exist
func (chm *ConcurrentHashMap[K, V]) ComputeIfAbsent(key K, mappingFunction func(K) V) V {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucketPtr := &segment.buckets[bucketIndex]

	// Check again (double-checked locking pattern)
	for current := bucketPtr.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			return current.value
		}
	}

	// Compute new value and add
	newValue := mappingFunction(key)
	newNode := &bucket[K, V]{
		key:   key,
		value: newValue,
		next:  bucketPtr.next,
	}
	bucketPtr.next = newNode
	segment.size++

	// Check if resize is needed
	if float64(segment.size) > float64(len(segment.buckets))*concurrentLoadFactor {
		chm.resizeSegment(segment)
	}

	return newValue
}

// ComputeIfPresent recomputes value if key exists
func (chm *ConcurrentHashMap[K, V]) ComputeIfPresent(key K, remappingFunction func(K, V) V) (V, bool) {
	hash := chm.hash(key)
	segmentIndex := hash & chm.segmentMask
	segment := chm.segments[segmentIndex]

	segment.mutex.Lock()
	defer segment.mutex.Unlock()

	bucketIndex := hash % uint32(len(segment.buckets))
	bucket := &segment.buckets[bucketIndex]

	for current := bucket.next; current != nil; current = current.next {
		if chm.hashStrategy.Equals(current.key, key) {
			newValue := remappingFunction(key, current.value)
			current.value = newValue
			return newValue, true
		}
	}

	var zero V
	return zero, false
}

// Internal helper methods

// hash computes hash value for key using the hash strategy
func (chm *ConcurrentHashMap[K, V]) hash(key K) uint32 {
	return uint32(chm.hashStrategy.Hash(key))
}

// valueEquals compares two values for equality
func (chm *ConcurrentHashMap[K, V]) valueEquals(a, b V) bool {
	return reflect.DeepEqual(a, b)
}

// resizeSegment resizes the specified segment
func (chm *ConcurrentHashMap[K, V]) resizeSegment(segment *segment[K, V]) {
	oldBuckets := segment.buckets
	newSize := len(oldBuckets) * 2
	newBuckets := make([]bucket[K, V], newSize)

	// Rehash all elements from old buckets to new buckets
	for _, oldBucket := range oldBuckets {
		for current := oldBucket.next; current != nil; current = current.next {
			hash := chm.hash(current.key)
			bucketIndex := hash % uint32(newSize)
			newNode := &bucket[K, V]{
				key:   current.key,
				value: current.value,
				next:  newBuckets[bucketIndex].next,
			}
			newBuckets[bucketIndex].next = newNode
		}
	}

	// Replace old buckets with new buckets
	segment.buckets = newBuckets
}

// nextPowerOfTwo returns the smallest power of 2 greater than or equal to n
func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}
