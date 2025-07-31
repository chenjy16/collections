# Go Collections Library

A comprehensive Go-based collection library that provides common collection types and utilities found in Java SDK but missing from Go's standard library. This library follows Go language conventions and idioms, providing rich collection operation functionality for Go developers.

## Table of Contents

- [Features](#features)
- [Project Statistics](#project-statistics)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Collection Types](#collection-types)
  - [List](#-list)
  - [Set](#-set)
  - [Map](#-map)
  - [Queue](#-queue)
  - [Stack](#-stack)
- [Architecture Design](#architecture-design)
- [Core Features](#core-features)
- [Iterator Pattern](#iterator-pattern)
- [Performance](#performance)
- [Thread Safety](#thread-safety)
- [Examples](#examples)
- [API Reference](#api-reference)
- [Contributing](#contributing)
- [License](#license)

## Features

- 🚀 **Rich Collection Types**: List, Set, Map, Queue, Stack implementations
- 🔧 **Generic Support**: Full Go 1.18+ generics support with type safety
- ⚡ **High Performance**: Optimized algorithms and efficient memory management
- 🎯 **Go-Idiomatic**: Clean API design following Go conventions
- 🧪 **Well Tested**: Comprehensive test coverage with 3,915+ lines of test code
- 🔒 **Thread Safety**: Concurrent implementations available where needed
- 📊 **Set Operations**: Mathematical operations (union, intersection, difference)
- 🔄 **Iterator Pattern**: Unified traversal interface across all collections

## Project Statistics

- **Total Go Files**: 28
- **Source Code Lines**: 4,731 lines (excluding tests)
- **Test Code Lines**: 3,915 lines
- **Test Files**: 9
- **Average Test Coverage**: 70%+ (some modules reach 100%)

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
    arrayList := list.New[int]()
    arrayList.Add(1, 2, 3)
    fmt.Println("List size:", arrayList.Size())

    // HashSet example
    hashSet := set.New[string]()
    hashSet.Add("apple", "banana")
    fmt.Println("Set size:", hashSet.Size())

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

## Collection Types

### 📋 List

Dynamic array and linked list implementations with rich operations.

- **ArrayList**: Dynamic array-based implementation
  - O(1) access by index
  - O(1) amortized append
  - Automatic capacity management
  
- **LinkedList**: Doubly linked list implementation
  - O(1) insertion/deletion at any position
  - O(n) access by index
  - Memory efficient for frequent modifications

### 🔗 Set

Unique element collections with various ordering guarantees.

- **HashSet**: Hash table-based set
  - O(1) average add/remove/contains
  - No ordering guarantee
  - Best for fast lookups

- **LinkedHashSet**: Hash set that maintains insertion order
  - O(1) average operations
  - Preserves insertion order
  - Ideal for ordered iteration

- **TreeSet**: Red-black tree-based ordered set
  - O(log n) operations
  - Natural ordering maintained
  - Range operations supported

- **ConcurrentSkipListSet**: Thread-safe ordered set
  - O(log n) operations
  - Lock-free concurrent access
  - Scalable for high concurrency

### 🗺️ Map

Key-value pair collections with different characteristics.

- **HashMap**: Hash table-based map
  - O(1) average operations
  - No ordering guarantee
  - General-purpose mapping

- **TreeMap**: Red-black tree-based ordered map
  - O(log n) operations
  - Key ordering maintained
  - Range queries supported

- **LinkedHashMap**: Hash map with insertion order
  - O(1) average operations
  - Preserves insertion order
  - Hybrid performance benefits

- **ConcurrentHashMap**: Thread-safe hash map
  - Segment-based locking
  - High concurrent performance
  - Lock-free reads

### 📤 Queue

FIFO (First-In-First-Out) data structures.

- **LinkedQueue**: Linked list-based queue
- **PriorityQueue**: Heap-based priority queue

### 📚 Stack

LIFO (Last-In-First-Out) data structures.

- **ArrayStack**: Array-based stack implementation
- **LinkedStack**: Linked list-based stack implementation

## Architecture Design

### 🏗️ Core Interfaces

```go
// Base container interface
type Container interface {
    Size() int
    IsEmpty() bool
    Clear()
    String() string
}

// Iteration support
type Iterable[E any] interface {
    ForEach(func(E))
}

// Iterator pattern
type Iterator[E any] interface {
    HasNext() bool
    Next() (E, bool)
    Remove() bool
}

// Comparison interface
type Comparable[T any] interface {
    CompareTo(T) int
}
```

### 📁 Module Structure

```
container/
├── common/     # Common interfaces and utilities
├── list/       # List implementations (ArrayList, LinkedList)
├── set/        # Set implementations (HashSet, TreeSet, etc.)
├── map/        # Map implementations (HashMap, TreeMap, etc.)
├── queue/      # Queue implementations
└── stack/      # Stack implementations
```

## Core Features

### Generic Type Safety

- Complete Go 1.18+ generics support
- Compile-time type checking
- Zero runtime type conversion overhead
- Clean and intuitive API design

### Memory Management

- Efficient memory allocation strategies
- Automatic capacity management
- Lazy initialization where appropriate
- Memory-conscious data structures

### Algorithm Optimization

- Optimized algorithms for common operations
- Benchmark-driven performance tuning
- Efficient iteration patterns
- Cache-friendly data layouts

## Iterator Pattern

The library provides a unified iterator interface across all collection types:

### Iterator Interface

```go
type Iterator[E any] interface {
    HasNext() bool     // Check if more elements exist
    Next() (E, bool)   // Get next element
    Remove() bool      // Remove current element (optional)
}
```

### Key Benefits

1. **Unified Traversal**: Same interface for all collection types
2. **Safe Modification**: Remove elements safely during iteration
3. **Flexible Control**: Pause, resume, or break iteration as needed
4. **Concurrent Support**: Multiple independent iterators

### Usage Example

```go
list := list.New[int]()
list.Add(1, 2, 3, 4, 5)

// Safe removal during iteration
iterator := list.Iterator()
for iterator.HasNext() {
    value, _ := iterator.Next()
    if value%2 == 0 {
        iterator.Remove() // Safe removal
    }
}
```

## Performance

### Benchmark Results

Collection performance characteristics (operations per second):

```
BenchmarkArrayList_Add-10           50000000    25.2 ns/op
BenchmarkArrayList_Get-10          100000000    10.1 ns/op
BenchmarkLinkedList_Add-10          30000000    45.3 ns/op
BenchmarkHashSet_Add-10             20000000    65.4 ns/op
BenchmarkHashSet_Contains-10       100000000    12.8 ns/op
BenchmarkTreeSet_Add-10             10000000   125.7 ns/op
BenchmarkConcurrentHashMap_Get-10   15000000    79.9 ns/op
BenchmarkConcurrentHashMap_Put-10   11000000   122.4 ns/op
```

### Performance Guidelines

- **ArrayList**: Best for index-based access and append operations
- **LinkedList**: Best for frequent insertions/deletions
- **HashSet**: Best for fast membership testing
- **TreeSet**: Best when ordering is required
- **ConcurrentHashMap**: Best for high-concurrency scenarios

## Thread Safety

### Thread-Safe Collections

- **ConcurrentHashMap**: Segment-based locking for high concurrency
- **ConcurrentSkipListSet**: Lock-free skip list implementation
- **CopyOnWriteMap**: Copy-on-write for read-heavy scenarios

### Concurrency Features

- **Lock-Free Reads**: Where possible, reads don't require locks
- **Fine-Grained Locking**: Minimal lock contention
- **Atomic Operations**: Built-in atomic operations support
- **Memory Consistency**: Proper memory barriers and synchronization

## Examples

### Basic Operations

```go
// List operations
list := list.New[string]()
list.Add("apple", "banana", "cherry")
list.Insert(1, "orange")
fmt.Println(list.Get(1)) // "orange"

// Set operations
set1 := set.New[int]()
set1.Add(1, 2, 3)
set2 := set.New[int]()
set2.Add(3, 4, 5)

union := set1.Union(set2)        // {1, 2, 3, 4, 5}
intersection := set1.Intersection(set2) // {3}
difference := set1.Difference(set2)     // {1, 2}

// Map operations
m := maps.NewHashMap[string, int]()
m.Put("one", 1)
m.Put("two", 2)
m.PutIfAbsent("three", 3)

keys := m.Keys()     // ["one", "two", "three"]
values := m.Values() // [1, 2, 3]
```

### Advanced Usage

```go
// Custom comparison for TreeSet
type Person struct {
    Name string
    Age  int
}

func (p Person) CompareTo(other Person) int {
    if p.Age != other.Age {
        return p.Age - other.Age
    }
    return strings.Compare(p.Name, other.Name)
}

treeSet := set.NewTreeSet[Person]()
treeSet.Add(Person{"Alice", 30})
treeSet.Add(Person{"Bob", 25})
// Automatically sorted by age, then name
```

## API Reference

### Common Operations

All collections implement the base `Container` interface:

```go
Size() int          // Get number of elements
IsEmpty() bool      // Check if empty
Clear()            // Remove all elements
String() string    // String representation
```

### Collection-Specific APIs

#### List Interface
```go
Add(elements ...E)              // Add elements to end
Insert(index int, element E)    // Insert at specific position
Get(index int) (E, bool)       // Get element by index
Set(index int, element E) bool // Set element at index
RemoveAt(index int) (E, bool)  // Remove by index
IndexOf(element E) int         // Find index of element
SubList(start, end int) List[E] // Get sublist
```

#### Set Interface
```go
Add(elements ...E) bool        // Add elements
Remove(element E) bool         // Remove element
Contains(element E) bool       // Check membership
Union(other Set[E]) Set[E]     // Set union
Intersection(other Set[E]) Set[E] // Set intersection
Difference(other Set[E]) Set[E]   // Set difference
```

#### Map Interface
```go
Put(key K, value V) (V, bool)     // Add/update entry
Get(key K) (V, bool)              // Get value by key
Remove(key K) (V, bool)           // Remove entry
ContainsKey(key K) bool           // Check key existence
ContainsValue(value V) bool       // Check value existence
Keys() []K                        // Get all keys
Values() []V                      // Get all values
```

## Contributing

We welcome contributions! Please feel free to submit issues and pull requests.

### Development Guidelines

1. Follow Go conventions and idioms
2. Add comprehensive tests for new features
3. Update documentation for API changes
4. Run benchmarks for performance-critical changes
5. Ensure thread safety where applicable

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Note**: This library requires Go 1.18+ for generics support.