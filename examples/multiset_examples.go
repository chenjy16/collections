// Package main provides examples of using the multiset implementations
package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chenjianyu/collections/container/multiset"
)

// RunMultisetExamples runs all multiset examples
func RunMultisetExamples() {
	fmt.Println("=== Running Multiset Examples ===")

	// Run various multiset examples
	HashMultisetExample()
	TreeMultisetExample()
	LinkedHashMultisetExample()
	ConcurrentHashMultisetExample()
	ImmutableMultisetExample()
	MultisetOperationsExample()
	WordFrequencyExample()
	PerformanceComparisonExample()
	FromSliceExample()

	fmt.Println("All multiset examples completed successfully!")
}

// HashMultisetExample demonstrates basic usage of HashMultiset
func HashMultisetExample() {
	fmt.Println("\n=== HashMultiset Example ===")
	
	// Create a new HashMultiset
	ms := multiset.NewHashMultiset[string]()
	
	// Add elements
	ms.Add("apple")
	ms.Add("banana")
	ms.Add("apple")
	ms.AddCount("orange", 3)
	
	fmt.Printf("Total size: %d\n", ms.TotalSize())
	fmt.Printf("Distinct elements: %d\n", ms.DistinctElements())
	fmt.Printf("Count of 'apple': %d\n", ms.Count("apple"))
	fmt.Printf("Count of 'orange': %d\n", ms.Count("orange"))
	
	// Iterate through elements
	fmt.Print("Elements: ")
	ms.ForEach(func(element string) {
		fmt.Printf("%s ", element)
	})
	fmt.Println("")
	
	// Get element set (unique elements)
	elementSet := ms.ElementSet()
	fmt.Printf("Element set: %v\n", elementSet)
	
	// Remove elements
	ms.Remove("apple")
	fmt.Printf("After removing one 'apple': %d\n", ms.Count("apple"))
	
	ms.RemoveAll("orange")
	fmt.Printf("After removing all 'orange': %d\n", ms.Count("orange"))
	
	fmt.Println(ms.String())
}

// TreeMultisetExample demonstrates TreeMultiset with ordering
func TreeMultisetExample() {
	fmt.Println("\n=== TreeMultiset Example ===")
	
	// Create a new TreeMultiset
	ms := multiset.NewTreeMultiset[int]()
	
	// Add elements in random order
	elements := []int{5, 2, 8, 2, 1, 5, 9, 1, 1}
	for _, elem := range elements {
		ms.Add(elem)
	}
	
	fmt.Printf("Total size: %d\n", ms.TotalSize())
	fmt.Printf("Distinct elements: %d\n", ms.DistinctElements())
	
	// TreeMultiset maintains sorted order
	fmt.Print("Sorted elements: ")
	ms.ForEach(func(element int) {
		fmt.Printf("%d ", element)
	})
	fmt.Println("")
	
	// Show entry set with counts
	fmt.Println("Entry set:")
	for _, entry := range ms.EntrySet() {
		fmt.Printf("  %d: %d times\n", entry.Element, entry.Count)
	}
	
	fmt.Println(ms.String())
}

// LinkedHashMultisetExample demonstrates LinkedHashMultiset with insertion order
func LinkedHashMultisetExample() {
	fmt.Println("\n=== LinkedHashMultiset Example ===")
	
	// Create a new LinkedHashMultiset
	ms := multiset.NewLinkedHashMultiset[string]()
	
	// Add elements - insertion order is preserved
	words := []string{"hello", "world", "hello", "go", "world", "programming"}
	for _, word := range words {
		ms.Add(word)
	}
	
	fmt.Printf("Total size: %d\n", ms.TotalSize())
	fmt.Printf("Distinct elements: %d\n", ms.DistinctElements())
	
	// LinkedHashMultiset preserves insertion order
	fmt.Print("Elements in insertion order: ")
	ms.ForEach(func(element string) {
		fmt.Printf("%s ", element)
	})
	fmt.Println("")
	
	// Show entry set with counts in insertion order
	fmt.Println("Entry set (insertion order):")
	for _, entry := range ms.EntrySet() {
		fmt.Printf("  %s: %d times\n", entry.Element, entry.Count)
	}
	
	fmt.Println(ms.String())
}

// ConcurrentHashMultisetExample demonstrates ConcurrentHashMultiset with concurrent access
func ConcurrentHashMultisetExample() {
	fmt.Println("\n=== ConcurrentHashMultiset Example ===")
	
	// Create a new ConcurrentHashMultiset
	ms := multiset.NewConcurrentHashMultiset[int]()
	
	// Simulate concurrent access
	var wg sync.WaitGroup
	numGoroutines := 10
	elementsPerGoroutine := 100
	
	// Start multiple goroutines adding elements
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < elementsPerGoroutine; j++ {
				ms.Add(start + j%10) // Add numbers 0-9 from each goroutine
			}
		}(i * 10)
	}
	
	wg.Wait()
	
	fmt.Printf("Total size after concurrent operations: %d\n", ms.TotalSize())
	fmt.Printf("Distinct elements: %d\n", ms.DistinctElements())
	
	// Show some counts
	for i := 0; i < 10; i++ {
		fmt.Printf("Count of %d: %d\n", i, ms.Count(i))
	}
}

// ImmutableMultisetExample demonstrates ImmutableMultiset
func ImmutableMultisetExample() {
	fmt.Println("\n=== ImmutableMultiset Example ===")
	
	// Create a new ImmutableMultiset
	ms1 := multiset.NewImmutableMultiset[string]()
	
	// Add elements (returns new instances)
	ms2 := ms1.WithAdd("apple")
	ms3 := ms2.WithAdd("banana").WithAdd("apple")
	ms4 := ms3.WithAddCount("orange", 3)
	
	fmt.Printf("ms1 size: %d\n", ms1.TotalSize())
	fmt.Printf("ms2 size: %d\n", ms2.TotalSize())
	fmt.Printf("ms3 size: %d\n", ms3.TotalSize())
	fmt.Printf("ms4 size: %d\n", ms4.TotalSize())
	
	// Original instances remain unchanged
	fmt.Printf("ms1 is empty: %t\n", ms1.IsEmpty())
	fmt.Printf("ms4 count of 'apple': %d\n", ms4.Count("apple"))
	fmt.Printf("ms4 count of 'orange': %d\n", ms4.Count("orange"))
	
	// Remove elements (returns new instances)
	ms5 := ms4.WithRemove("apple")
	fmt.Printf("ms5 count of 'apple': %d\n", ms5.Count("apple"))
	fmt.Printf("ms4 count of 'apple' (unchanged): %d\n", ms4.Count("apple"))
	
	fmt.Println(ms4.String())
}

// MultisetOperationsExample demonstrates set operations
func MultisetOperationsExample() {
	fmt.Println("\n=== Multiset Set Operations Example ===")
	
	// Create two multisets
	ms1 := multiset.NewHashMultiset[string]()
	ms1.AddCount("a", 2)
	ms1.AddCount("b", 3)
	ms1.AddCount("c", 1)
	
	ms2 := multiset.NewHashMultiset[string]()
	ms2.AddCount("b", 1)
	ms2.AddCount("c", 2)
	ms2.AddCount("d", 4)
	
	fmt.Printf("ms1: %s\n", ms1.String())
	fmt.Printf("ms2: %s\n", ms2.String())
	
	// Union
	union := ms1.Union(ms2)
	fmt.Printf("Union: %s\n", union.String())
	
	// Intersection
	intersection := ms1.Intersection(ms2)
	fmt.Printf("Intersection: %s\n", intersection.String())
	
	// Difference
	difference := ms1.Difference(ms2)
	fmt.Printf("Difference (ms1 - ms2): %s\n", difference.String())
	
	// Subset/Superset checks
	ms3 := multiset.NewHashMultiset[string]()
	ms3.Add("b")
	ms3.Add("c")
	
	fmt.Printf("ms3: %s\n", ms3.String())
	fmt.Printf("ms3 is subset of ms1: %t\n", ms3.IsSubsetOf(ms1))
	fmt.Printf("ms1 is superset of ms3: %t\n", ms1.IsSupersetOf(ms3))
}

// WordFrequencyExample demonstrates word frequency counting
func WordFrequencyExample() {
	fmt.Println("\n=== Word Frequency Example ===")
	
	text := "the quick brown fox jumps over the lazy dog the fox is quick"
	words := strings.Fields(strings.ToLower(text))
	
	// Count word frequencies using HashMultiset
	wordCount := multiset.NewHashMultiset[string]()
	for _, word := range words {
		wordCount.Add(word)
	}
	
	fmt.Printf("Total words: %d\n", wordCount.TotalSize())
	fmt.Printf("Unique words: %d\n", wordCount.DistinctElements())
	
	// Get word frequencies sorted by count
	type wordFreq struct {
		word  string
		count int
	}
	
	var frequencies []wordFreq
	for _, entry := range wordCount.EntrySet() {
		frequencies = append(frequencies, wordFreq{entry.Element, entry.Count})
	}
	
	// Sort by frequency (descending)
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].count > frequencies[j].count
	})
	
	fmt.Println("Word frequencies:")
	for _, freq := range frequencies {
		fmt.Printf("  %s: %d\n", freq.word, freq.count)
	}
}

// PerformanceComparisonExample demonstrates performance comparison
func PerformanceComparisonExample() {
	fmt.Println("\n=== Performance Comparison Example ===")
	
	const numElements = 10000
	
	// Test HashMultiset
	start := time.Now()
	hashMs := multiset.NewHashMultiset[int]()
	for i := 0; i < numElements; i++ {
		hashMs.Add(i % 1000)
	}
	hashTime := time.Since(start)
	
	// Test TreeMultiset
	start = time.Now()
	treeMs := multiset.NewTreeMultiset[int]()
	for i := 0; i < numElements; i++ {
		treeMs.Add(i % 1000)
	}
	treeTime := time.Since(start)
	
	// Test LinkedHashMultiset
	start = time.Now()
	linkedMs := multiset.NewLinkedHashMultiset[int]()
	for i := 0; i < numElements; i++ {
		linkedMs.Add(i % 1000)
	}
	linkedTime := time.Since(start)
	
	fmt.Printf("Adding %d elements:\n", numElements)
	fmt.Printf("  HashMultiset: %v\n", hashTime)
	fmt.Printf("  TreeMultiset: %v\n", treeTime)
	fmt.Printf("  LinkedHashMultiset: %v\n", linkedTime)
	
	// Test lookup performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		hashMs.Count(i)
	}
	hashLookupTime := time.Since(start)
	
	start = time.Now()
	for i := 0; i < 1000; i++ {
		treeMs.Count(i)
	}
	treeLookupTime := time.Since(start)
	
	fmt.Printf("1000 lookups:\n")
	fmt.Printf("  HashMultiset: %v\n", hashLookupTime)
	fmt.Printf("  TreeMultiset: %v\n", treeLookupTime)
}

// FromSliceExample demonstrates creating multisets from slices
func FromSliceExample() {
	fmt.Println("\n=== Creating Multisets from Slices ===")
	
	// Sample data
	numbers := []int{1, 2, 3, 2, 1, 4, 3, 2, 5, 1}
	
	// Create different types of multisets from slice
	hashMs := multiset.NewHashMultisetFromSlice(numbers)
	treeMs := multiset.NewTreeMultisetFromSlice(numbers)
	linkedMs := multiset.NewLinkedHashMultisetFromSlice(numbers)
	immutableMs := multiset.NewImmutableMultisetFromSlice(numbers)
	
	fmt.Printf("Original slice: %v\n", numbers)
	fmt.Printf("HashMultiset: %s\n", hashMs.String())
	fmt.Printf("TreeMultiset: %s\n", treeMs.String())
	fmt.Printf("LinkedHashMultiset: %s\n", linkedMs.String())
	fmt.Printf("ImmutableMultiset: %s\n", immutableMs.String())
	
	// Convert back to slice
	hashSlice := hashMs.ToSlice()
	fmt.Printf("Back to slice (HashMultiset): %v\n", hashSlice)
}