package multiset

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Benchmark HashMultiset operations
func BenchmarkHashMultiset_Add(b *testing.B) {
	ms := NewHashMultiset[int]()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ms.Add(i % 1000)
	}
}

func BenchmarkHashMultiset_Count(b *testing.B) {
	ms := NewHashMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(i % 1000)
	}
}

func BenchmarkHashMultiset_Remove(b *testing.B) {
	ms := NewHashMultiset[int]()
	for i := 0; i < b.N; i++ {
		ms.AddCount(i%1000, 10)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Remove(i % 1000)
	}
}

func BenchmarkHashMultiset_Iterator(b *testing.B) {
	ms := NewHashMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.AddCount(i, 5)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator := ms.Iterator()
		for iterator.HasNext() {
			iterator.Next()
		}
	}
}

// Benchmark TreeMultiset operations
func BenchmarkTreeMultiset_Add(b *testing.B) {
	ms := NewTreeMultiset[int]()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ms.Add(i % 1000)
	}
}

func BenchmarkTreeMultiset_Count(b *testing.B) {
	ms := NewTreeMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(i % 1000)
	}
}

func BenchmarkTreeMultiset_Remove(b *testing.B) {
	ms := NewTreeMultiset[int]()
	for i := 0; i < b.N; i++ {
		ms.AddCount(i%1000, 10)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Remove(i % 1000)
	}
}

func BenchmarkTreeMultiset_Iterator(b *testing.B) {
	ms := NewTreeMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.AddCount(i, 5)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator := ms.Iterator()
		for iterator.HasNext() {
			iterator.Next()
		}
	}
}

// Benchmark LinkedHashMultiset operations
func BenchmarkLinkedHashMultiset_Add(b *testing.B) {
	ms := NewLinkedHashMultiset[int]()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ms.Add(i % 1000)
	}
}

func BenchmarkLinkedHashMultiset_Count(b *testing.B) {
	ms := NewLinkedHashMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(i % 1000)
	}
}

func BenchmarkLinkedHashMultiset_Remove(b *testing.B) {
	ms := NewLinkedHashMultiset[int]()
	for i := 0; i < b.N; i++ {
		ms.AddCount(i%1000, 10)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Remove(i % 1000)
	}
}

func BenchmarkLinkedHashMultiset_Iterator(b *testing.B) {
	ms := NewLinkedHashMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.AddCount(i, 5)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator := ms.Iterator()
		for iterator.HasNext() {
			iterator.Next()
		}
	}
}

// Benchmark ConcurrentHashMultiset operations
func BenchmarkConcurrentHashMultiset_Add(b *testing.B) {
	ms := NewConcurrentHashMultiset[int]()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ms.Add(i % 1000)
	}
}

func BenchmarkConcurrentHashMultiset_Count(b *testing.B) {
	ms := NewConcurrentHashMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.Add(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(i % 1000)
	}
}

func BenchmarkConcurrentHashMultiset_Remove(b *testing.B) {
	ms := NewConcurrentHashMultiset[int]()
	for i := 0; i < b.N; i++ {
		ms.AddCount(i%1000, 10)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Remove(i % 1000)
	}
}

func BenchmarkConcurrentHashMultiset_Iterator(b *testing.B) {
	ms := NewConcurrentHashMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms.AddCount(i, 5)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator := ms.Iterator()
		for iterator.HasNext() {
			iterator.Next()
		}
	}
}

// Benchmark ImmutableMultiset operations
func BenchmarkImmutableMultiset_WithAdd(b *testing.B) {
	ms := NewImmutableMultiset[int]()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ms = ms.WithAdd(i % 1000)
	}
}

func BenchmarkImmutableMultiset_Count(b *testing.B) {
	ms := NewImmutableMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms = ms.WithAdd(i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(i % 1000)
	}
}

func BenchmarkImmutableMultiset_WithRemove(b *testing.B) {
	ms := NewImmutableMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms = ms.WithAddCount(i, 10)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms = ms.WithRemove(i % 1000)
	}
}

func BenchmarkImmutableMultiset_Iterator(b *testing.B) {
	ms := NewImmutableMultiset[int]()
	for i := 0; i < 1000; i++ {
		ms = ms.WithAddCount(i, 5)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator := ms.Iterator()
		for iterator.HasNext() {
			iterator.Next()
		}
	}
}

// Comparative benchmarks
func BenchmarkMultisetComparison_Add(b *testing.B) {
	implementations := map[string]func() Multiset[int]{
		"HashMultiset":           func() Multiset[int] { return NewHashMultiset[int]() },
		"TreeMultiset":           func() Multiset[int] { return NewTreeMultiset[int]() },
		"LinkedHashMultiset":     func() Multiset[int] { return NewLinkedHashMultiset[int]() },
		"ConcurrentHashMultiset": func() Multiset[int] { return NewConcurrentHashMultiset[int]() },
	}
	
	for name, factory := range implementations {
		b.Run(name, func(b *testing.B) {
			ms := factory()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				ms.Add(i % 1000)
			}
		})
	}
}

func BenchmarkMultisetComparison_Count(b *testing.B) {
	implementations := map[string]func() Multiset[int]{
		"HashMultiset":           func() Multiset[int] { return NewHashMultiset[int]() },
		"TreeMultiset":           func() Multiset[int] { return NewTreeMultiset[int]() },
		"LinkedHashMultiset":     func() Multiset[int] { return NewLinkedHashMultiset[int]() },
		"ConcurrentHashMultiset": func() Multiset[int] { return NewConcurrentHashMultiset[int]() },
	}
	
	for name, factory := range implementations {
		b.Run(name, func(b *testing.B) {
			ms := factory()
			for i := 0; i < 1000; i++ {
				ms.Add(i)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ms.Count(i % 1000)
			}
		})
	}
}

func BenchmarkMultisetComparison_Remove(b *testing.B) {
	implementations := map[string]func() Multiset[int]{
		"HashMultiset":           func() Multiset[int] { return NewHashMultiset[int]() },
		"TreeMultiset":           func() Multiset[int] { return NewTreeMultiset[int]() },
		"LinkedHashMultiset":     func() Multiset[int] { return NewLinkedHashMultiset[int]() },
		"ConcurrentHashMultiset": func() Multiset[int] { return NewConcurrentHashMultiset[int]() },
	}
	
	for name, factory := range implementations {
		b.Run(name, func(b *testing.B) {
			ms := factory()
			for i := 0; i < b.N; i++ {
				ms.AddCount(i%1000, 10)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ms.Remove(i % 1000)
			}
		})
	}
}

// Concurrent benchmark for ConcurrentHashMultiset
func BenchmarkConcurrentHashMultiset_ConcurrentAdd(b *testing.B) {
	ms := NewConcurrentHashMultiset[int]()
	
	b.RunParallel(func(pb *testing.PB) {
		rand.Seed(time.Now().UnixNano())
		for pb.Next() {
			ms.Add(rand.Intn(1000))
		}
	})
}

func BenchmarkConcurrentHashMultiset_ConcurrentMixed(b *testing.B) {
	ms := NewConcurrentHashMultiset[int]()
	
	// Pre-populate
	for i := 0; i < 1000; i++ {
		ms.AddCount(i, 5)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		rand.Seed(time.Now().UnixNano())
		for pb.Next() {
			switch rand.Intn(3) {
			case 0:
				ms.Add(rand.Intn(1000))
			case 1:
				ms.Count(rand.Intn(1000))
			case 2:
				ms.Remove(rand.Intn(1000))
			}
		}
	})
}

// Memory allocation benchmarks
func BenchmarkMultisetMemory_HashMultiset(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ms := NewHashMultiset[int]()
		for j := 0; j < 100; j++ {
			ms.Add(j)
		}
	}
}

func BenchmarkMultisetMemory_TreeMultiset(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ms := NewTreeMultiset[int]()
		for j := 0; j < 100; j++ {
			ms.Add(j)
		}
	}
}

func BenchmarkMultisetMemory_LinkedHashMultiset(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ms := NewLinkedHashMultiset[int]()
		for j := 0; j < 100; j++ {
			ms.Add(j)
		}
	}
}

func BenchmarkMultisetMemory_ImmutableMultiset(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ms := NewImmutableMultiset[int]()
		for j := 0; j < 100; j++ {
			ms = ms.WithAdd(j)
		}
	}
}

// Large dataset benchmarks
func BenchmarkMultisetLarge_HashMultiset(b *testing.B) {
	ms := NewHashMultiset[int]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Add(i)
	}
}

func BenchmarkMultisetLarge_TreeMultiset(b *testing.B) {
	ms := NewTreeMultiset[int]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Add(i)
	}
}

// String operations benchmark
func BenchmarkMultisetString_HashMultiset(b *testing.B) {
	ms := NewHashMultiset[string]()
	for i := 0; i < 1000; i++ {
		ms.Add(fmt.Sprintf("element_%d", i))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(fmt.Sprintf("element_%d", i%1000))
	}
}

func BenchmarkMultisetString_TreeMultiset(b *testing.B) {
	ms := NewTreeMultiset[string]()
	for i := 0; i < 1000; i++ {
		ms.Add(fmt.Sprintf("element_%d", i))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Count(fmt.Sprintf("element_%d", i%1000))
	}
}