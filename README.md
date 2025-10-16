# Go Collections Library

This library follows Go language conventions and idioms, providing rich collection operation functionality for Go developers.

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
  - [Immutable Collections](#-immutable-collections)
  - [Graph](#-graph)
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

- 🚀 **Rich Collection Types**: List, Set, Map, Graph, Queue, Stack implementations
- 🔧 **Generic Support**: Full Go 1.18+ generics support with type safety
- ⚡ **High Performance**: Optimized algorithms and efficient memory management
- 🎯 **Go-Idiomatic**: Clean API design following Go conventions
- 🧪 **Well Tested**: Comprehensive test coverage with 7,757 lines of test code across 19 test files
- 🔒 **Thread Safety**: Concurrent implementations available where needed
- 📊 **Set Operations**: Mathematical operations (union, intersection, difference)
- 🔄 **Iterator Pattern**: Unified traversal interface across all collections



## Installation

```bash
go get github.com/chenjianyu/collections
```

## Quick Start

```go
package main

import (
    "errors"
    "fmt"

    "github.com/chenjianyu/collections/container/common"
    "github.com/chenjianyu/collections/container/graph"
    "github.com/chenjianyu/collections/container/list"
    maps "github.com/chenjianyu/collections/container/map"
    "github.com/chenjianyu/collections/container/multimap"
    "github.com/chenjianyu/collections/container/queue"
    "github.com/chenjianyu/collections/container/range"
    "github.com/chenjianyu/collections/container/set"
    "github.com/chenjianyu/collections/container/stack"
)

func main() {
    // ===== List Example =====
    arrayList := list.New[int]()
    arrayList.Add(1)
    arrayList.Add(2)
    arrayList.Add(3)
    fmt.Println("List size:", arrayList.Size())

    // ===== Set Example =====
    hashSet := set.New[string]()
    hashSet.Add("apple")
    hashSet.Add("banana")
    fmt.Println("Set size:", hashSet.Size())

    // ===== Map Example =====
    hashMap := maps.NewHashMap[string, int]()
    hashMap.Put("one", 1)
    hashMap.Put("two", 2)
    if val, ok := hashMap.Get("one"); ok {
        fmt.Println("Value for 'one':", val)
    }

    // ===== Concurrent Skip List Set =====
    skipSet := set.NewConcurrentSkipListSet[int]()
    skipSet.Add(5)
    skipSet.Add(2)
    skipSet.Add(8)
    skipSet.Add(1)
    fmt.Println("SkipSet (sorted):", skipSet)

    // ===== Multimap Example =====
    multiMap := multimap.NewArrayListMultimap[string, int]()
    multiMap.Put("numbers", 1)
    multiMap.Put("numbers", 2)
    multiMap.Put("numbers", 3)
    multiMap.Put("letters", 97)
    fmt.Println("Values for 'numbers':", multiMap.Get("numbers"))
    fmt.Println("Multimap size:", multiMap.Size())

    // ===== Range Examples =====
    r1 := ranges.ClosedRange(1, 10)
    r2 := ranges.OpenRange(5, 15)
    fmt.Println("Range contains 5:", r1.Contains(5))

    rs := ranges.NewTreeRangeSet[int]()
    rs.Add(r1, r2)
    fmt.Println("RangeSet contains 8:", rs.ContainsValue(8))

    rm := ranges.NewTreeRangeMap[int, string]()
    rm.Put(r1, "first range")
    rm.Put(r2, "second range")
    if value, ok := rm.Get(8); ok {
        fmt.Println("Value for 8:", value)
    }

    // ===== Error Handling for ArrayList =====
    fmt.Println("\nArrayList Error Handling:")
    arrayList = list.New[int]()
    _, err := arrayList.Get(0)
    if err != nil && errors.Is(err, common.ErrIndexOutOfBounds) {
        fmt.Println("✓ IndexOutOfBounds on empty list")
    }
    arrayList.Add(1)
    arrayList.Add(2)
    _, err = arrayList.Get(5)
    if err != nil {
        fmt.Println("Invalid index error:", err)
    }
    _, err = arrayList.SubList(2, 1)
    if err != nil {
        fmt.Println("Invalid range error:", err)
    }

    // ===== Error Handling for LinkedList =====
    fmt.Println("\nLinkedList Error Handling:")
    linkedList := list.NewLinkedList[string]()
    _, err = linkedList.GetFirst()
    if err != nil && errors.Is(err, common.ErrEmptyContainer) {
        fmt.Println("✓ EmptyContainer on GetFirst")
    }

    // ===== Stack with Capacity Limit =====
    fmt.Println("\nStack Error Handling:")
    limitedStack := stack.WithCapacity[int](2)
    limitedStack.Push(1, 2)
    err = limitedStack.Push(3)
    if err != nil && errors.Is(err, common.ErrFullContainer) {
        fmt.Println("✓ FullContainer on Push")
    }
    emptyStack := stack.New[int]()
    _, err = emptyStack.Pop()
    if err != nil {
        fmt.Println("Error popping from empty stack:", err)
    }

    // ===== Priority Queue Error =====
    fmt.Println("\nPriorityQueue Error Handling:")
    emptyQueue := queue.NewPriorityQueueWithComparator[int](func(a, b int) int { return a - b })
    _, err = emptyQueue.Remove()
    if err != nil {
        fmt.Println("Error removing from empty queue:", err)
    }

    // ===== LinkedList Queue Error =====
    fmt.Println("\nLinkedListQueue Error Handling:")
    emptyLLQueue := queue.New[string]()
    _, err = emptyLLQueue.GetFirst()
    if err != nil {
        fmt.Println("Error getting first element from empty queue:", err)
    }

    // ===== Graph Examples =====
    fmt.Println("\nBasic Graph Example:")
    g := graph.UndirectedGraph[string]()
    g.AddNode("A", "B", "C")
    g.PutEdge("A", "B")
    g.PutEdge("B", "C")
    fmt.Println("A connected to B:", g.HasEdgeConnecting("A", "B"))
    fmt.Println("A's neighbors:", g.AdjacentNodes("A"))

    fmt.Println("\nDirected Graph:")
    dg := graph.DirectedGraph[int]()
    dg.AddNode(1, 2, 3)
    dg.PutEdge(1, 2, 2, 3, 3, 1)
    fmt.Println("Node 2 in-degree:", dg.InDegree(2))

    fmt.Println("\nValueGraph:")
    vg := graph.UndirectedValueGraph[string, int]()
    vg.AddNode("X", "Y")
    vg.PutEdgeValue("X", "Y", 100)
    if val, ok := vg.EdgeValue("X", "Y"); ok {
        fmt.Println("Edge X->Y value:", val)
    }

    fmt.Println("\nNetwork Example:")
    net := graph.UndirectedNetwork[string, string]()
    net.AddNode("N1", "N2")
    net.AddEdge("E1", "N1", "N2")
    fmt.Println("Edge endpoints:", net.IncidentNodes("E1"))

    // ===== Immutable Collections =====
    fmt.Println("\nImmutable Collections:")
    immutableList := list.Of(1, 2, 3)
    newList := immutableList.WithAdd(4)
    fmt.Println("Original:", immutableList)
    fmt.Println("New:", newList)

    immutableSet := set.SetOf("a", "b")
    newSet := immutableSet.WithAdd("c")
    fmt.Println("Original:", immutableSet)
    fmt.Println("New:", newSet)

    immutableMap := maps.MapOf(
        maps.Pair[string, int]{Key: "x", Value: 1},
        maps.Pair[string, int]{Key: "y", Value: 2},
    )
    newMap := immutableMap.WithPut("z", 3)
    fmt.Println("Original map size:", immutableMap.Size())
    fmt.Println("New map size:", newMap.Size())
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
  - Fine-grained locking (sync.RWMutex)
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
  - Read-optimized with minimal locking
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
  - Read-optimized with segmented locks

### 🔒 Immutable Collections

immutable collections that provide thread-safe, copy-on-write semantics.

- **ImmutableList**: Immutable list implementation
  - Thread-safe by design
  - Copy-on-write semantics for modifications
  - Returns new instances for all modification operations
  - Supports all standard list operations (get, indexOf, subList)
  - Creation methods: `NewImmutableList()`, `NewImmutableListFromSlice()`, `Of()`

- **ImmutableSet**: Immutable set implementation
  - Thread-safe by design
  - Copy-on-write semantics for modifications
  - Returns new instances for all modification operations
  - Supports set operations (union, intersection, difference)
  - Creation methods: `NewImmutableSet()`, `NewImmutableSetFromSlice()`, `SetOf()`

- **ImmutableMap**: Immutable map implementation
  - Thread-safe by design
  - Copy-on-write semantics for modifications
  - Returns new instances for all modification operations
  - Supports all standard map operations (get, keys, values, entries)
  - Creation methods: `NewImmutableMap()`, `NewImmutableMapFromMap()`, `MapOf()`

#### Semantics

- All direct mutation methods on immutable collections are no-ops and return failure indicators
  - Set: `Add`/`Remove` return `false`
  - Map: `Put`/`Remove` return `(zeroValue, false)`
  - Multiset: mutation methods do not modify the original and return previous counts
- Use `WithXxx` methods to create modified copies:
  - Set: `WithAdd`, `WithRemove`, `WithClear`
  - Map: `WithPut`, `WithRemove`, `WithClear`, `WithPutAll`
  - Multiset: `WithAdd`, `WithAddCount`, `WithRemove`, `WithRemoveCount`, `WithClear`

### 🕸️ Graph

graph data structures for modeling relationships between nodes.

- **Graph[N]**: Basic graph interface for node relationships
  - **MutableGraph**: Mutable graph implementation
  - Supports both directed and undirected graphs
  - Self-loops configurable
  - Node and edge management operations
  - Graph traversal methods (successors, predecessors, adjacent nodes)

- **ValueGraph[N,V]**: Graph with values on edges
  - **MutableValueGraph**: Mutable value graph implementation
  - Associates values with edges
  - All Graph operations plus edge value management
  - Efficient edge value lookup and modification
  - AsGraph() view for basic graph operations

- **Network[N,E]**: Graph with explicit edge objects
  - **MutableNetwork**: Mutable network implementation
  - Explicit edge objects with unique identities
  - Supports parallel edges (configurable)
  - Edge-centric operations (incident nodes, adjacent edges)
  - AsGraph() view for basic graph operations

#### Graph Properties

- **Directedness**: Directed or undirected graphs
- **Self-loops**: Allow or disallow self-loops
- **Parallel edges**: Allow or disallow parallel edges (Network only)
- **Node ordering**: Configurable node iteration order
- **Edge ordering**: Configurable edge iteration order (Network only)

#### Core Implementation

- **Graph Interface Series**： Implemented three core interfaces: Graph, ValueGraph, and Network.
- **Mutable Implementations**： Created complete implementations for MutableGraph, MutableValueGraph, and MutableNetwork.
- **Builder Pattern**： Implemented GraphBuilder, ValueGraphBuilder, and NetworkBuilder, supporting fluent configuration.
- **Convenience Functions**：Provided convenience constructors such as DirectedGraph, UndirectedGraph, etc.

#### Functional Features

- **Directionality Support**： Supports directed and undirected graphs.
- **Self-Loop Support**： Configurable self-loop handling.
- **Parallel Edge Support**： The Network interface supports multiple edges.
- **Edge Value Support**： ValueGraph supports edge weights/values.
- **Element Ordering**： Supports unordered, insertion order, and natural order.
- **View Adapters**： ValueGraph and Network can be converted to a basic Graph view.



#### Builder Pattern

```go
// Graph builders for flexible configuration
graph := NewGraphBuilder[string]().
    Directed().
    AllowSelfLoops().
    Build()

valueGraph := NewValueGraphBuilder[string, int]().
    Undirected().
    Build()

network := NewNetworkBuilder[string, string]().
    Directed().
    AllowParallelEdges().
    Build()
```

#### Convenience Factory Functions

```go
// Quick creation methods
g1 := UndirectedGraph[string]()
g2 := DirectedGraph[int]()
vg1 := UndirectedValueGraph[string, float64]()
vg2 := DirectedValueGraph[int, string]()
n1 := UndirectedNetwork[string, string]()
n2 := DirectedNetwork[int, int]()
```

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
├── graph/      # Graph implementations (Graph, ValueGraph, Network)
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

### Comparator Strategies

- Customize ordering using `common.ComparatorStrategy[T]` across ordered structures
- Supported constructors:
  - `set.NewTreeSetWithComparatorStrategy[T](strategy)` and `set.NewTreeSetWithComparator[T](cmp)`
  - `set.NewConcurrentSkipListSetWithComparatorStrategy[T](strategy)` and `set.NewConcurrentSkipListSetWithComparator[T](cmp)`
  - `ranges.NewTreeRangeSetWithComparatorStrategy[T](strategy)` and `ranges.NewTreeRangeSetWithComparator[T](cmp)`
  - `ranges.NewTreeRangeMapWithComparatorStrategy[K, V](strategy)` and `ranges.NewTreeRangeMapWithComparator[K, V](cmp)`
- Default comparator expectations:
  - `TreeSet`/`ConcurrentSkipListSet`: natural ordering via `common.CompareNatural` (prefers `Comparable.CompareTo`, falls back to generic)
  - `TreeMap`/`TreeMultimap` (keys): natural ordering via `common.CompareNatural`; TreeMultimap values use `TreeSet` with the same comparator
  - `TreeMultiset`: natural ordering via `common.CompareNatural`
  - `PriorityQueue`: default expects `common.Comparable`; provide a comparator for non-Comparable types
  - Range structures: default comparator prefers `Comparable.CompareTo` then built-in natural order; strategies supported via `ComparatorFromStrategy`

Example (reverse order with strategy):

```go
strat := common.NewFunctionalComparatorStrategy[int](func(a, b int) int { return b - a })
ts := set.NewTreeSetWithComparatorStrategy[int](strat)
rs := ranges.NewTreeRangeSetWithComparatorStrategy[int](strat)
rm := ranges.NewTreeRangeMapWithComparatorStrategy[int, string](strat)
```

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
list.Add(1)
list.Add(2)
list.Add(3)
list.Add(4)
list.Add(5)

// Safe removal during iteration
iterator := list.Iterator()
for iterator.HasNext() {
    value, _ := iterator.Next()
    if value%2 == 0 {
        iterator.Remove() // Safe removal
    }
}
```


## Thread Safety

### Thread-Safe Collections

- **ConcurrentHashMap**: Segment-based locking for high concurrency
- **ConcurrentSkipListSet**: Fine-grained locking (sync.RWMutex)
- **ConcurrentHashMultiset**: Segment-based locking with atomic size tracking
- **CopyOnWriteMap**: Copy-on-write for read-heavy scenarios

### Concurrency Features

- **Read-Optimized Reads**: Minimize locking overhead for common read paths
- **Fine-Grained Locking**: Minimal lock contention
- **Atomic Operations**: Built-in atomic operations support
- **Memory Consistency**: Proper memory barriers and synchronization

## Examples

### Basic Operations

```go
// List operations
list := list.New[string]()
list.Add("apple")
list.Add("banana")
list.Add("cherry")
list.Insert(1, "orange")
fmt.Println(list.Get(1)) // "orange"

// Set operations
set1 := set.New[int]()
set1.Add(1)
set1.Add(2)
set1.Add(3)
set2 := set.New[int]()
set2.Add(3)
set2.Add(4)
set2.Add(5)

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
// Use a concurrent map to demonstrate PutIfAbsent (implementation-specific)
m := maps.NewConcurrentHashMap[string, int]()
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

cmp := func(a, b Person) int {
    if a.Age != b.Age {
        return a.Age - b.Age
    }
    return strings.Compare(a.Name, b.Name)
}

treeSet := set.NewTreeSetWithComparator[Person](cmp)
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
Add(element E) bool                              // Add element to end
Insert(index int, element E) error               // Insert at specific position
Get(index int) (E, error)                        // Get element by index
Set(index int, element E) (E, bool)              // Replace element at index
RemoveAt(index int) (E, bool)                    // Remove by index
Remove(element E) bool                           // Remove first occurrence
IndexOf(element E) int                           // Find index of element
LastIndexOf(element E) int                       // Find last index of element
SubList(fromIndex, toIndex int) (List[E], error) // Get sublist view
ToSlice() []E                                    // Copy all elements to slice
```

#### Set Interface
```go
Add(element E) bool            // Add element
Remove(element E) bool         // Remove element
Contains(element E) bool       // Check membership
Union(other Set[E]) Set[E]     // Set union
Intersection(other Set[E]) Set[E] // Set intersection
Difference(other Set[E]) Set[E]   // Set difference
IsSubsetOf(other Set[E]) bool  // Subset check
IsSupersetOf(other Set[E]) bool // Superset check
ToSlice() []E                  // Copy all elements to slice
```

#### Multiset Interface
```go
Add(element E)                                     // Add single element (returns previous count)
AddCount(element E, count int) (int, error)        // Add with count
Remove(element E) int                              // Remove one occurrence (returns previous count)
RemoveCount(element E, count int) (int, error)     // Remove with count
RemoveAll(element E) int                           // Remove all occurrences (returns previous count)
Count(element E) int                               // Get element count
SetCount(element E, count int) (int, error)        // Set element count
ElementSet() []E                                   // Get unique elements
EntrySet() []Entry[E]                              // Get element-count pairs
TotalSize() int                                    // Total number of elements
DistinctElements() int                             // Number of unique elements
ToSlice() []E                                      // All elements including duplicates
Union(other Multiset[E]) Multiset[E]               // Multiset union (max counts)
Intersection(other Multiset[E]) Multiset[E]        // Multiset intersection (min counts)
Difference(other Multiset[E]) Multiset[E]          // Multiset difference
IsSubsetOf(other Multiset[E]) bool                 // Subset by counts
IsSupersetOf(other Multiset[E]) bool               // Superset by counts
```

#### Map Interface
```go
Put(key K, value V) (V, bool)     // Add/update entry
Get(key K) (V, bool)              // Get value by key
Remove(key K) (V, bool)           // Remove entry
ContainsKey(key K) bool           // Check key existence
ContainsValue(value V) bool       // Check value existence
Size() int                        // Number of entries
IsEmpty() bool                    // Whether empty
Clear()                           // Remove all entries
Keys() []K                        // Get all keys
Values() []V                      // Get all values
Entries() []common.Entry[K, V]    // Get all key-value pairs (unified Entry type)
ForEach(func(K, V))               // Iterate over entries
String() string                   // String representation
PutAll(other Map[K, V])           // Copy all entries from another Map
```

#### Multimap Interface
```go
Put(key K, value V) bool                 // Add a key-value mapping
PutAll(multimap Multimap[K, V]) bool     // Add all mappings from another multimap
ReplaceValues(key K, values []V) []V     // Replace values for a key
Remove(key K, value V) bool              // Remove a specific key-value mapping
RemoveAll(key K) []V                     // Remove all values for a key
ContainsKey(key K) bool                  // Whether key exists
ContainsValue(value V) bool              // Whether any mapping to value exists
ContainsEntry(key K, value V) bool       // Whether specific key-value mapping exists
Get(key K) []V                           // All values for a key
Keys() []K                               // All distinct keys
Values() []V                             // All values across keys
Entries() []common.Entry[K, V]           // All key-value pairs
KeySet() []K                             // Distinct keys view
AsMap() map[K][]V                        // Map view (key -> values)
ForEach(func(K, V))                      // Iterate key-value pairs
```

### Concurrent Containers Behavior

- CopyOnWriteMap
  - Uses `sync.RWMutex`; writes acquire `Lock` and replace the underlying map with a fresh copy (copy-on-write).
  - Reads and iteration (`Keys`, `Values`, `Entries`, `ForEach`) acquire `RLock` and operate on the current snapshot.
  - `Snapshot()`/`ToMap()` return a copy for stable iteration; `PutIfAbsent`, `Replace`, `ReplaceIf` are atomic writes.

- ConcurrentHashMap
  - Segmented `sync.RWMutex` locks; each operation locks only the relevant segment for better concurrency.
  - `Keys`/`Values`/`Entries` iterate segments under `RLock` per segment (weakly consistent view; may reflect concurrent updates).
  - Provides advanced ops: `PutIfAbsent`, `Replace`, `ReplaceIf`, `ComputeIfAbsent`, `ComputeIfPresent`, plus `Snapshot`/`ToMap`.

- ConcurrentSkipListSet
  - Uses `sync.RWMutex`; `Add`/`Remove` under `Lock`, reads under `RLock`.
  - `Iterator()` returns a snapshot built via `ToSlice()`; iterator `Remove()` is not supported for concurrent set.
  - Maintains sorted order via skip list; `Union`/`Intersection`/`Difference` create new sets from snapshots.

- ConcurrentHashMultiset
  - Segment-based `sync.RWMutex`; counts stored per segment; total size maintained via atomic ops.
  - `EntrySet()` builds a snapshot by reading each segment under `RLock`; `Iterator` iterates snapshot entries and refreshes after `Remove()`.
  - Supports `AddCount`/`RemoveCount`/`SetCount` with atomic size adjustments; `ForEach` iterates over snapshot entries.

> Note: Methods like `PutIfAbsent` and `Replace*` are implementation-specific to concurrent/copy-on-write maps. The base `Map` interface does not declare them.

#### Graph Interface
```go
// Basic Graph Operations
Nodes() set.Set[N]                // Get all nodes
Edges() set.Set[EndpointPair[N]]  // Get all edges
IsDirected() bool                 // Check if directed
AllowsSelfLoops() bool            // Check if self-loops allowed
NodeOrder() ElementOrder          // Get node ordering
Degree(node N) int                // Get node degree
InDegree(node N) int              // Get in-degree (directed)
OutDegree(node N) int             // Get out-degree (directed)

// Node Relationships
Adjacents(node N) set.Set[N]      // Get adjacent nodes
Predecessors(node N) set.Set[N]   // Get predecessor nodes
Successors(node N) set.Set[N]     // Get successor nodes
IncidentEdges(node N) set.Set[EndpointPair[N]] // Get incident edges

// Connectivity
HasEdgeConnecting(nodeU, nodeV N) bool // Check edge existence
```

#### ValueGraph Interface
```go
// Inherits all Graph methods plus:
EdgeValue(nodeU, nodeV N) (V, bool)           // Get edge value
EdgeValueOrDefault(nodeU, nodeV N, defaultValue V) V // Get edge value with default
AsGraph() Graph[N]                            // View as basic graph
```

#### Network Interface
```go
// Inherits all Graph methods plus:
AllowsParallelEdges() bool                    // Check if parallel edges allowed
EdgeOrder() ElementOrder                      // Get edge ordering
EdgesConnecting(nodeU, nodeV N) set.Set[E]   // Get connecting edges
InEdges(node N) set.Set[E]                   // Get incoming edges
OutEdges(node N) set.Set[E]                  // Get outgoing edges
IncidentNodes(edge E) EndpointPair[N]        // Get edge endpoints
AsGraph() Graph[N]                           // View as basic graph
```

## Recent Improvements


### Quality Assurance

- All tests pass with zero failures
- Fixed string format inconsistencies in range tests
- Enhanced error handling with proper `errors.Is()` usage
- Improved test reliability and maintainability

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