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
  - [Multiset](#-multiset)
  - [Map](#-map)
  - [Multimap](#-multimap)
  - [Queue](#-queue)
  - [Stack](#-stack)
  - [Range](#-range)
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

- **Total Go Files**: 36
- **Source Code Lines**: 6,200+ lines (excluding tests)
- **Test Code Lines**: 4,500+ lines
- **Test Files**: 10
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
    "github.com/chenjianyu/collections/container/multimap"
    "github.com/chenjianyu/collections/container/range"
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
    
    // ArrayListMultimap example
    multiMap := multimap.NewArrayListMultimap[string, int]()
    multiMap.Put("numbers", 1)
    multiMap.Put("numbers", 2)
    multiMap.Put("numbers", 3)
    multiMap.Put("letters", 97) // ASCII for 'a'
    fmt.Println("Values for 'numbers':", multiMap.Get("numbers")) // Output: [1, 2, 3]
    fmt.Println("Multimap size:", multiMap.Size()) // Output: 4
    
    // Range example
    r1 := ranges.ClosedRange(1, 10)     // [1, 10]
    r2 := ranges.OpenRange(5, 15)       // (5, 15)
    fmt.Println("Range contains 5:", r1.Contains(5)) // true
    
    // RangeSet example
    rs := ranges.NewTreeRangeSet[int]()
    rs.Add(r1)
    rs.Add(r2)
    fmt.Println("RangeSet contains 8:", rs.ContainsValue(8)) // true
    
    // RangeMap example
    rm := ranges.NewTreeRangeMap[int, string]()
    rm.Put(r1, "first range")
    rm.Put(r2, "second range")
    if value, ok := rm.Get(8); ok {
        fmt.Println("Value for 8:", value) // "first range"
    }
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

### 🔢 Multiset

Collections that allow duplicate elements with counting functionality (similar to Guava's Multiset).

- **HashMultiset**: Hash table-based multiset
  - O(1) average add/remove/count operations
  - No ordering guarantee
  - Best for fast counting operations

- **TreeMultiset**: Red-black tree-based ordered multiset
  - O(log n) operations
  - Natural ordering maintained
  - Sorted iteration of elements

- **LinkedHashMultiset**: Hash multiset that maintains insertion order
  - O(1) average operations
  - Preserves insertion order
  - Ideal for ordered counting

- **ConcurrentHashMultiset**: Thread-safe hash multiset
  - Segment-based locking for high concurrency
  - Lock-free reads
  - Scalable concurrent counting

- **ImmutableMultiset**: Immutable multiset implementation
  - Copy-on-write semantics
  - Thread-safe by design
  - Functional programming friendly

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

### 🔄 Multimap

Key-value pair collections where each key can map to multiple values (similar to Guava's Multimap).

- **ArrayListMultimap**: Multimap with ArrayList values
  - Values for each key stored in ArrayList
  - Allows duplicate values per key
  - Maintains insertion order of values

- **HashMultimap**: Multimap with HashSet values
  - Values for each key stored in HashSet
  - No duplicate values per key
  - Fast lookup operations

- **LinkedHashMultimap**: Multimap that maintains insertion order
  - Preserves insertion order of keys and values
  - No duplicate values per key
  - Predictable iteration order

- **TreeMultimap**: Ordered multimap implementation
  - Keys maintained in sorted order
  - Values for each key stored in sorted sets
  - Range operations supported

- **ImmutableMultimap**: Immutable multimap implementation
  - Thread-safe by design
  - Copy-on-write semantics
  - Functional programming friendly

- **ImmutableListMultimap**: Immutable multimap with list values
  - Immutable list of values per key
  - Preserves duplicates and order
  - Thread-safe by design

- **ImmutableSetMultimap**: Immutable multimap with set values
  - Immutable set of values per key
  - No duplicates per key
  - Thread-safe by design

### 📤 Queue

FIFO (First-In-First-Out) data structures.

- **LinkedQueue**: Linked list-based queue
- **PriorityQueue**: Heap-based priority queue

### 📚 Stack

LIFO (Last-In-First-Out) data structures.

- **ArrayStack**: Array-based stack implementation
- **LinkedStack**: Linked list-based stack implementation

### 📏 Range

Guava-style range collections for interval data management.

- **Range[T]**: Represents a range of values with configurable bounds
  - Support for open and closed bounds
  - Range operations (intersection, union, contains)
  - Comparable value support

- **RangeSet[T]**: Set of non-overlapping ranges
  - **TreeRangeSet**: Mutable implementation with O(log n) operations
  - **ImmutableRangeSet**: Immutable implementation returning new instances
  - Set operations (union, intersection, difference, complement)
  - Efficient range merging and splitting

- **RangeMap[K,V]**: Mapping from ranges to values
  - **TreeRangeMap**: Mutable implementation with O(log n) operations  
  - **ImmutableRangeMap**: Immutable implementation returning new instances
  - Non-overlapping range keys
  - Efficient range-based lookups

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
├── multiset/   # Multiset implementations (HashMultiset, TreeMultiset, etc.)
├── map/        # Map implementations (HashMap, TreeMap, etc.)
├── multimap/   # Multimap implementations (ArrayListMultimap, HashMultimap, etc.)
├── queue/      # Queue implementations
├── stack/      # Stack implementations
└── range/      # Range implementations (Range, RangeSet, RangeMap)
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

// Multiset operations (counting collections)
multiset := multiset.NewHashMultiset[string]()
multiset.Add("apple")
multiset.Add("banana")
multiset.Add("apple")  // Duplicates allowed
multiset.AddCount("orange", 3)

fmt.Println(multiset.Count("apple"))   // 2
fmt.Println(multiset.TotalSize())      // 6
fmt.Println(multiset.DistinctElements()) // 3

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

#### Multiset Interface
```go
Add(element E)                    // Add single element
AddCount(element E, count int)    // Add element with count
Remove(element E) bool            // Remove one occurrence
RemoveAll(element E) int          // Remove all occurrences
Count(element E) int              // Get element count
SetCount(element E, count int)    // Set element count
TotalSize() int                   // Total number of elements
DistinctElements() int            // Number of unique elements
ElementSet() []E                  // Get unique elements
EntrySet() []Entry[E]             // Get element-count pairs
Union(other Multiset[E]) Multiset[E]        // Multiset union
Intersection(other Multiset[E]) Multiset[E] // Multiset intersection
Difference(other Multiset[E]) Multiset[E]   // Multiset difference
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