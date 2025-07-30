# Collections

Collections is a Go-based collection library that provides common collection types and utilities found in Java SDK but missing from Go's standard library. This library follows Go language conventions and idioms, providing rich collection operation functionality for Go developers.

## Features

- 🚀 Rich collection types: List, Set, Map, Queue, Stack, etc.
- 🔧 Generic support (based on Go 1.18+)
- ⚡ High-performance implementations
- 🎯 Go-idiomatic API design
- 🧪 Complete documentation and testing
- 🔒 Thread-safe implementations available
- 📊 Mathematical set operations support (union, intersection, difference)
- 🔄 Iterator pattern support

## Project Statistics

- **Total Go files**: 28
- **Source code lines**: 4,731 lines (excluding tests)
- **Test code lines**: 3,915 lines
- **Test files**: 9
- **Average test coverage**: 70%+ (some modules reach 100%)

## Installation

```bash
go get github.com/chenjianyu/collections
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/chenjianyu/collections/container/list"
	"github.com/chenjianyu/collections/container/set"
	"github.com/chenjianyu/collections/container/map"
)

func main() {
	// ArrayList example
	list := list.New[int]()
	list.Add(1, 2, 3)
	fmt.Println("List size:", list.Size())

	// HashSet example
	set := set.New[string]()
	set.Add("apple", "banana")
	fmt.Println("Set size:", set.Size())

	// HashMap example
	hashMap := maps.NewHashMap[string, int]()
	hashMap.Put("one", 1)
	hashMap.Put("two", 2)
	if val, ok := hashMap.Get("one"); ok {
		fmt.Println("Value for 'one':", val)
	}

	// ConcurrentSkipListSet example (thread-safe ordered set)
	skipSet := set.NewConcurrentSkipListSet[int]()
	skipSet.Add(5, 2, 8, 1)
	fmt.Println("SkipSet (sorted):", skipSet) // Output: [1, 2, 5, 8]
}
```

## Supported Collection Types

### 📋 List
- **ArrayList**: List implementation based on dynamic arrays
- **LinkedList**: List implementation based on doubly linked lists

### 🔗 Set
- **HashSet**: Set implementation based on hash tables
- **TreeSet**: Ordered Set implementation based on red-black trees
- **ConcurrentSkipListSet**: Thread-safe ordered Set implementation based on skip lists

### 🗺️ Map
- **HashMap**: Map implementation based on hash tables
- **TreeMap**: Ordered Map implementation based on red-black trees
- **LinkedHashMap**: Map implementation based on chaining and red-black trees, similar to Java's HashMap
- **ConcurrentHashMap**: Thread-safe hash table implementation

`ConcurrentHashMap` is a thread-safe hash table implementation that borrows design ideas from Java's `ConcurrentHashMap`, using segment locking technology to achieve high concurrent performance.

## Core Features

### 1. Thread Safety
- Uses segment locking technology, dividing the hash table into multiple segments
- Each segment locks independently, reducing lock contention
- Supports high concurrent read and write operations

### 2. Segment Lock Design
- Default 16 segments, adjustable through capacity parameters
- Read operations are lock-free, write operations only lock relevant segments
- Operations on different segments can execute in parallel

### 3. Dynamic Resizing
- Supports automatic resizing when load factor exceeds threshold
- Maintains thread safety during resizing process
- Progressive resizing to reduce performance impact

### 4. Rich API
- Implements complete `Map` interface
- Provides atomic operations: `PutIfAbsent`, `Replace`, `ReplaceIf`
- Supports compute operations: `ComputeIfAbsent`, `ComputeIfPresent`
- Provides batch operations and conversion functionality

## Use Cases

### High Concurrent Read/Write
- Cache systems in multi-threaded environments
- Concurrent data processing
- Real-time data statistics

### Read-Heavy Scenarios
- Configuration management systems
- Metadata storage
- Lookup tables and indexes

### Atomic Operation Requirements
- Counters and statistics
- State management
- Conditional updates

## Performance Characteristics

### Concurrent Performance
- **Read operations**: Lock-free, supports high concurrency
- **Write operations**: Segment locking, reduces contention
- **Mixed operations**: Read and write can execute in parallel

### Benchmark Results
```
BenchmarkConcurrentHashMapGet-10     14684961    79.86 ns/op    159 B/op    10 allocs/op
BenchmarkConcurrentHashMapPut-10     11166087   122.4 ns/op     90 B/op     5 allocs/op
BenchmarkConcurrentHashMapMixed-10    9794205   120.7 ns/op    174 B/op    11 allocs/op
```

### Performance Advantages
- Excellent read operation performance (~80ns)
- Reasonable write operation overhead (~120ns)
- Stable mixed operation performance

## Basic Usage

### Creating Instances
```go
// Using default capacity
chm := maps.NewConcurrentHashMap[string, int]()

// Specifying initial capacity
chm := maps.NewConcurrentHashMapWithCapacity[string, int](100)

// Creating from existing map
sourceMap := map[string]int{"key1": 1, "key2": 2}
chm := maps.ConcurrentHashMapFromMap(sourceMap)
```

### Basic Operations
```go
// Add/update elements
oldValue, existed := chm.Put("key1", 100)

// Get elements
value, exists := chm.Get("key1")

// Remove elements
oldValue, removed := chm.Remove("key1")

// Check containment
hasKey := chm.ContainsKey("key1")
hasValue := chm.ContainsValue(100)

// Get size
size := chm.Size()
isEmpty := chm.IsEmpty()

// Clear
chm.Clear()
```

## API Methods

### Basic Operations
- `Put(key K, value V) (V, bool)` - Add or update key-value pair
- `Get(key K) (V, bool)` - Get value for specified key
- `Remove(key K) (V, bool)` - Remove specified key
- `ContainsKey(key K) bool` - Check if contains specified key
- `ContainsValue(value V) bool` - Check if contains specified value
- `Size() int` - Get number of elements
- `IsEmpty() bool` - Check if empty
- `Clear()` - Clear all elements

### Advanced Operations
- `PutIfAbsent(key K, value V) (V, bool)` - Add only if key doesn't exist
- `Replace(key K, value V) (V, bool)` - Replace only existing keys
- `ReplaceIf(key K, oldValue V, newValue V) bool` - Conditional replacement

### Compute Operations
- `ComputeIfAbsent(key K, mappingFunction func(K) V) V` - Compute missing values
- `ComputeIfPresent(key K, remappingFunction func(K, V) V) (V, bool)` - Recompute existing values

### Batch Operations
- `PutAll(other Map[K, V])` - Batch add all elements from another Map
- `PutAllFromMap(m map[K]V)` - Batch add all elements from Go native map

### Conversion and Traversal
- `Keys() []K` - Get all keys
- `Values() []V` - Get all values
- `ToMap() map[K]V` - Convert to Go native map
- `Snapshot() map[K]V` - Get snapshot of current state
- `ForEach(fn func(K, V))` - Traverse all elements

### Container Operations
- `String() string` - String representation

## Concurrency Safety Mechanism

### Segment Lock Architecture
```go
type ConcurrentHashMap[K comparable, V any] struct {
    segments []*segment[K, V]  // Segment array
    segmentMask int            // Segment mask
    segmentShift int           // Segment shift
}

type segment[K comparable, V any] struct {
    mu      sync.RWMutex       // Read-write lock
    buckets []*bucket[K, V]    // Bucket array
    size    int                // Segment size
    threshold int              // Resize threshold
}
```

### Locking Strategy
- **Read operations**: Use read locks, support concurrent reading
- **Write operations**: Use write locks, ensure data consistency
- **Segment isolation**: Operations on different segments don't interfere

### Hash Distribution
- Uses high-quality hash functions to ensure uniform distribution
- Fast segment and bucket location through bit operations
- Supports dynamic resizing and rehashing

## Memory Usage

### Memory Structure
- Segmented storage reduces memory fragmentation
- Chaining for conflict resolution saves space
- Lazy allocation, resize on demand

### Memory Optimization
- Reasonable load factor (0.75)
- Progressive resizing strategy
- Timely cleanup of unused nodes

## Notes

### Usage Recommendations
1. **Capacity Planning**: Set reasonable initial capacity based on expected data volume
2. **Concurrency Control**: Although thread-safe, compound operations still need additional synchronization
3. **Performance Monitoring**: Monitor performance metrics in high concurrency scenarios
4. **Memory Management**: Regularly clean up unnecessary data

### Limitations
1. **Iteration Consistency**: Modifications during traversal may not be immediately reflected
2. **Memory Overhead**: Additional lock and segment overhead compared to regular maps
3. **Complexity**: Complex implementation, relatively difficult to debug

### Best Practices
1. Prioritize use in high concurrency scenarios
2. Set reasonable initial capacity to avoid frequent resizing
3. Use atomic operations instead of compound operations
4. Regularly monitor performance and memory usage

## Complete Examples

For complete usage examples, please refer to the `ConcurrentHashMapExample()` function in `examples/container_examples.go`.

## Technical Implementation

### Core Algorithms
- **Hash Algorithm**: Uses high-quality hash functions to ensure uniform distribution
- **Segmentation Strategy**: Segment allocation based on high bits of hash value
- **Conflict Resolution**: Uses chaining to handle hash conflicts
- **Resize Algorithm**: Progressive resizing maintains stable performance

### Performance Optimization
- **Lock-free Reading**: Read operations don't require locks
- **Fine-grained Locking**: Write operations only lock relevant segments
- **Memory Locality**: Optimized data structure layout
- **Batch Operations**: Reduce lock acquisition frequency

### Generic Support
- Full Go generics support
- Type-safe key-value operations
- Compile-time type checking
- Zero runtime type conversion overhead

### 📤 Queue
- **LinkedQueue**: Queue implementation based on linked lists
- **PriorityQueue**: Priority queue implementation

### 📚 Stack
- **ArrayStack**: Stack implementation based on arrays
- **LinkedStack**: Stack implementation based on linked lists

## Architecture Design

### 🏗️ Core Interfaces
- **Container**: Base container interface providing `Size()`, `IsEmpty()`, `Clear()`, `Contains()`, `String()` methods
- **Iterable**: Provides iteration capability with `ForEach()` method
- **Iterator**: Standard iterator pattern with `HasNext()`, `Next()`, `Remove()` methods
- **Comparable**: Generic comparison interface for custom types

### 📁 Module Structure
```
container/
├── common/     # Common interfaces and utilities
├── list/       # List implementations (ArrayList, LinkedList)
├── set/        # Set implementations (HashSet, TreeSet, ConcurrentSkipListSet)
├── map/        # Map implementations (HashMap, TreeMap, LinkedHashMap)
├── queue/      # Queue implementations
└── stack/      # Stack implementations
```

## Implementation Highlights

### 🔧 Generic Support
- Complete generic type support based on Go 1.18+
- Type-safe operations without runtime type assertions
- Clean and intuitive API design

### 🚀 Advanced Data Structures
- **ConcurrentSkipListSet**: Lock-free concurrent skip list with O(log n) operation complexity
- **TreeMap/TreeSet**: Red-black tree implementation ensuring O(log n) complexity
- **LinkedHashMap**: Hybrid approach using linked list and tree conversion for optimal performance

### ⚡ High-Performance Features
- Efficient memory management
- Optimized algorithms for common operations
- Benchmark tests ensuring performance standards

## Code Quality

### ✅ Strengths
1. **Comprehensive Test Coverage**: 9 test files with 3,915 lines of test code
2. **Clear Architecture**: Well-defined interfaces and modular design
3. **Performance Optimized**: Includes benchmark tests for key operations
4. **Thread Safety**: Provides concurrent implementations where needed
5. **Well Documented**: Rich examples and clear API documentation

### 📊 Test Coverage
- `container/stack`: 100.0% coverage
- `container/queue`: High coverage for core operations
- `container/set`: 54.8% coverage
- `container/map`: Good coverage for most implementations
- `container/list`: Comprehensive test suite

## Contributing

Issues and Pull Requests are welcome!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.