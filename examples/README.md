# Custom Hash Strategies Guide

This guide demonstrates how to use custom hash strategies with the Go collections library to customize how elements are hashed and compared.

## Overview

The collections library supports custom hash strategies for hash-based collections:
- `HashSet`
- `LinkedHashSet` 
- `ConcurrentHashMap`

Custom hash strategies allow you to:
- Implement case-insensitive string comparisons
- Compare objects by specific fields only
- Create custom equality semantics for complex types

## Basic Usage

### HashStrategy Interface

All hash strategies must implement the `HashStrategy[T]` interface:

```go
type HashStrategy[T any] interface {
    Hash(element T) uint64
    Equals(a, b T) bool
}
```

### Creating Collections with Custom Strategies

```go
// HashSet with custom strategy
strategy := &CaseInsensitiveHashStrategy{}
hashSet := set.NewWithHashStrategy[string](strategy)

// LinkedHashSet with custom strategy  
linkedHashSet := set.NewLinkedHashSetWithHashStrategy[string](strategy)

// ConcurrentHashMap with custom strategy
hashMap := maps.NewConcurrentHashMapWithHashStrategy[string, int](strategy)
```

## Examples

### 1. Case-Insensitive String Strategy

```go
type CaseInsensitiveHashStrategy struct{}

func (s *CaseInsensitiveHashStrategy) Hash(element string) uint64 {
    lower := strings.ToLower(element)
    hash := uint64(0)
    for i := 0; i < len(lower); i++ {
        hash = hash*31 + uint64(lower[i])
    }
    return hash
}

func (s *CaseInsensitiveHashStrategy) Equals(a, b string) bool {
    return strings.ToLower(a) == strings.ToLower(b)
}

// Usage
strategy := &CaseInsensitiveHashStrategy{}
hashSet := set.NewWithHashStrategy[string](strategy)

hashSet.Add("Hello")
hashSet.Add("HELLO") // Won't be added (duplicate)
hashSet.Add("hello") // Won't be added (duplicate)

fmt.Println(hashSet.Size()) // Output: 1
fmt.Println(hashSet.Contains("HeLLo")) // Output: true
```

### 2. Custom Struct Strategy (Compare by Field)

```go
type Person struct {
    Name string
    Age  int
}

type PersonByNameHashStrategy struct{}

func (s *PersonByNameHashStrategy) Hash(p Person) uint64 {
    hash := uint64(0)
    for i := 0; i < len(p.Name); i++ {
        hash = hash*31 + uint64(p.Name[i])
    }
    return hash
}

func (s *PersonByNameHashStrategy) Equals(a, b Person) bool {
    return a.Name == b.Name // Only compare by name, ignore age
}

// Usage
strategy := &PersonByNameHashStrategy{}
hashSet := set.NewWithHashStrategy[Person](strategy)

p1 := Person{Name: "John", Age: 25}
p2 := Person{Name: "John", Age: 30} // Same name, different age

hashSet.Add(p1)
hashSet.Add(p2) // Won't be added (same name)

fmt.Println(hashSet.Size()) // Output: 1
fmt.Println(hashSet.Contains(Person{Name: "John", Age: 40})) // Output: true
```

### 3. Functional Hash Strategy

You can also create strategies using functional constructors:

```go
hashStrategy := common.NewFunctionalHashStrategy(
    func(s string) uint64 {
        lower := strings.ToLower(s)
        hash := uint64(0)
        for i := 0; i < len(lower); i++ {
            hash = hash*31 + uint64(lower[i])
        }
        return hash
    },
    func(a, b string) bool {
        return strings.ToLower(a) == strings.ToLower(b)
    },
)

hashSet := set.NewWithHashStrategy[string](hashStrategy)
```

## Advanced Examples

### Set Operations with Custom Strategies

```go
strategy := &CaseInsensitiveHashStrategy{}

set1 := set.NewWithHashStrategy[string](strategy)
set1.Add("Hello")
set1.Add("World")

set2 := set.NewWithHashStrategy[string](strategy)
set2.Add("HELLO") // Matches "Hello" from set1
set2.Add("Go")

union := set1.Union(set2)        // ["Hello", "World", "Go"]
intersection := set1.Intersection(set2) // ["Hello"]
difference := set1.Difference(set2)     // ["World"]
```

### ConcurrentHashMap with Custom Strategy

```go
strategy := &CaseInsensitiveHashStrategy{}
hashMap := maps.NewConcurrentHashMapWithHashStrategy[string, int](strategy)

hashMap.Put("Hello", 1)
hashMap.Put("HELLO", 2) // Replaces the previous value

value, exists := hashMap.Get("hello") // value=2, exists=true
```

## Performance Considerations

- Custom hash strategies may have different performance characteristics
- Simple strategies (like default comparable) are generally faster
- Complex strategies with string operations may be slower but provide more functionality
- Use benchmarks to measure performance impact for your use case

## Default Strategy

If no custom strategy is provided, collections use the default `ComparableHashStrategy[T]` which:
- Uses Go's built-in hash function for the type
- Uses `==` operator for equality comparison
- Works with any `comparable` type

```go
// These are equivalent
hashSet1 := set.New[string]()
hashSet2 := set.NewWithHashStrategy[string](common.NewComparableHashStrategy[string]())
```

## Best Practices

1. **Consistency**: Ensure `Hash(a) == Hash(b)` when `Equals(a, b) == true`
2. **Performance**: Keep hash functions simple and fast
3. **Distribution**: Aim for good hash distribution to avoid collisions
4. **Immutability**: Don't modify objects used as keys after adding them to collections
5. **Testing**: Thoroughly test custom strategies with edge cases

## Thread Safety

- `HashStrategy` implementations should be stateless and thread-safe
- Multiple goroutines may call the same strategy methods concurrently
- Avoid mutable state in strategy implementations