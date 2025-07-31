package set

import (
	"fmt"
	"testing"
)

func BenchmarkLinkedHashSet_Add(b *testing.B) {
	set := NewLinkedHashSet[int]()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		set.Add(i)
	}
}

func BenchmarkLinkedHashSet_Contains(b *testing.B) {
	set := NewLinkedHashSet[int]()
	
	// Pre-populate the set
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Contains(i % 1000)
	}
}

func BenchmarkLinkedHashSet_Remove(b *testing.B) {
	set := NewLinkedHashSet[int]()
	
	// Pre-populate the set with enough elements
	for i := 0; i < b.N+1000; i++ {
		set.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Remove(i) // Remove elements sequentially
	}
}

func BenchmarkLinkedHashSet_Iterator(b *testing.B) {
	set := NewLinkedHashSet[int]()
	
	// Pre-populate the set
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator := set.Iterator()
		for iterator.HasNext() {
			iterator.Next()
		}
	}
}

func BenchmarkLinkedHashSet_ToSlice(b *testing.B) {
	set := NewLinkedHashSet[int]()
	
	// Pre-populate the set
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ToSlice()
	}
}

func BenchmarkLinkedHashSet_ForEach(b *testing.B) {
	set := NewLinkedHashSet[int]()
	
	// Pre-populate the set
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ForEach(func(element int) {
			// Do nothing, just iterate
		})
	}
}

// Comparison benchmarks between LinkedHashSet and HashSet
func BenchmarkComparison_Add(b *testing.B) {
	b.Run("LinkedHashSet", func(b *testing.B) {
		set := NewLinkedHashSet[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			set.Add(i)
		}
	})
	
	b.Run("HashSet", func(b *testing.B) {
		set := New[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			set.Add(i)
		}
	})
}

func BenchmarkComparison_Contains(b *testing.B) {
	linkedSet := NewLinkedHashSet[int]()
	hashSet := New[int]()
	
	// Pre-populate both sets
	for i := 0; i < 1000; i++ {
		linkedSet.Add(i)
		hashSet.Add(i)
	}
	
	b.Run("LinkedHashSet", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			linkedSet.Contains(i % 1000)
		}
	})
	
	b.Run("HashSet", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hashSet.Contains(i % 1000)
		}
	})
}

func BenchmarkComparison_Iteration(b *testing.B) {
	linkedSet := NewLinkedHashSet[int]()
	hashSet := New[int]()
	
	// Pre-populate both sets
	for i := 0; i < 1000; i++ {
		linkedSet.Add(i)
		hashSet.Add(i)
	}
	
	b.Run("LinkedHashSet_Iterator", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			iterator := linkedSet.Iterator()
			for iterator.HasNext() {
				iterator.Next()
			}
		}
	})
	
	b.Run("HashSet_Iterator", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			iterator := hashSet.Iterator()
			for iterator.HasNext() {
				iterator.Next()
			}
		}
	})
	
	b.Run("LinkedHashSet_ForEach", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			linkedSet.ForEach(func(element int) {
				// Do nothing
			})
		}
	})
	
	b.Run("HashSet_ForEach", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hashSet.ForEach(func(element int) {
				// Do nothing
			})
		}
	})
}

// PerformanceExample demonstrates performance characteristics
func PerformanceExample() {
	fmt.Println("LinkedHashSet Performance Characteristics:")
	fmt.Println("==========================================")
	
	// Create a LinkedHashSet and measure operations
	set := NewLinkedHashSet[int]()
	
	// Add operation - O(1) average case
	fmt.Println("Add operation: O(1) average case")
	for i := 0; i < 10; i++ {
		set.Add(i)
	}
	
	// Contains operation - O(1) average case
	fmt.Println("Contains operation: O(1) average case")
	fmt.Printf("Contains 5: %t\n", set.Contains(5))
	
	// Remove operation - O(1) average case
	fmt.Println("Remove operation: O(1) average case")
	fmt.Printf("Removed 5: %t\n", set.Remove(5))
	
	// Iteration - O(n) but maintains insertion order
	fmt.Println("Iteration: O(n) with insertion order maintained")
	fmt.Print("Elements in insertion order: ")
	set.ForEach(func(element int) {
		fmt.Printf("%d ", element)
	})
	fmt.Println()
	
	// Memory overhead
	fmt.Println("\nMemory Characteristics:")
	fmt.Println("- Hash table for O(1) lookup")
	fmt.Println("- Doubly linked list for order maintenance")
	fmt.Println("- Node map for O(1) node access")
	fmt.Println("- Higher memory usage than HashSet but provides ordering")
	
	// Use cases
	fmt.Println("\nBest Use Cases:")
	fmt.Println("- When you need both fast lookup AND insertion order")
	fmt.Println("- Cache implementations (LRU cache)")
	fmt.Println("- Maintaining processing order")
	fmt.Println("- Deterministic iteration over sets")
}