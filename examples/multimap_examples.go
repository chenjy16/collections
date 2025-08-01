package main

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/multimap"
)

// RunMultimapExamples demonstrates the usage of various Multimap implementations
func RunMultimapExamples() {
	fmt.Println("\n=== Running Multimap Examples ===")

	// Run examples for each Multimap implementation
	arrayListMultimapExample()
	hashMultimapExample()
	linkedHashMultimapExample()
	treeMultimapExample()
	immutableMultimapExample()
	immutableListMultimapExample()
	immutableSetMultimapExample()
	wordCountExample()

	fmt.Println("\n=== All Multimap Examples Completed Successfully ===")
}

// arrayListMultimapExample demonstrates the usage of ArrayListMultimap
func arrayListMultimapExample() {
	fmt.Println("\n--- ArrayListMultimap Example ---")

	// Create a new ArrayListMultimap
	multimap := multimap.NewArrayListMultimap[string, string]()

	// Add key-value mappings
	multimap.Put("fruits", "apple")
	multimap.Put("fruits", "banana")
	multimap.Put("fruits", "apple") // ArrayListMultimap allows duplicates
	multimap.Put("vegetables", "carrot")
	multimap.Put("vegetables", "broccoli")

	// Print the multimap
	fmt.Println("Multimap:", multimap)
	fmt.Println("Size:", multimap.Size())

	// Get values for a key
	fruits := multimap.Get("fruits")
	fmt.Println("Fruits:", fruits)

	// Check if the multimap contains a key-value mapping
	fmt.Println("Contains 'fruits' -> 'apple':", multimap.ContainsEntry("fruits", "apple"))
	fmt.Println("Contains 'fruits' -> 'orange':", multimap.ContainsEntry("fruits", "orange"))

	// Replace values for a key
	oldFruits := multimap.ReplaceValues("fruits", []string{"orange", "grape"})
	fmt.Println("Old fruits:", oldFruits)
	fmt.Println("New fruits:", multimap.Get("fruits"))

	// Remove a key-value mapping
	multimap.Remove("vegetables", "carrot")
	fmt.Println("After removing 'vegetables' -> 'carrot':", multimap)

	// Remove all values for a key
	removedVegetables := multimap.RemoveAll("vegetables")
	fmt.Println("Removed vegetables:", removedVegetables)
	fmt.Println("After removing all vegetables:", multimap)

	// Iterate over all key-value pairs
	fmt.Println("Iterating over all key-value pairs:")
	multimap.ForEach(func(key string, value string) {
		fmt.Printf("  %s -> %s\n", key, value)
	})

	// Clear the multimap
	multimap.Clear()
	fmt.Println("After clearing:", multimap)
	fmt.Println("Is empty:", multimap.IsEmpty())
}

// hashMultimapExample demonstrates the usage of HashMultimap
func hashMultimapExample() {
	fmt.Println("\n--- HashMultimap Example ---")

	// Create a new HashMultimap
	multimap := multimap.NewHashMultimap[string, int]()

	// Add key-value mappings
	multimap.Put("A", 1)
	multimap.Put("A", 2)
	multimap.Put("A", 1) // Duplicate value will be ignored
	multimap.Put("B", 3)
	multimap.Put("B", 4)

	// Print the multimap
	fmt.Println("Multimap:", multimap)
	fmt.Println("Size:", multimap.Size())

	// Get values for a key
	aValues := multimap.Get("A")
	fmt.Println("Values for key 'A':", aValues)

	// Get all keys
	keys := multimap.Keys()
	fmt.Println("All keys:", keys)

	// Get all values
	values := multimap.Values()
	fmt.Println("All values:", values)

	// Get all entries
	entries := multimap.Entries()
	fmt.Println("All entries:")
	for _, entry := range entries {
		fmt.Printf("  %v -> %v\n", entry.Key, entry.Value)
	}

	// Get as map
	asMap := multimap.AsMap()
	fmt.Println("As map:", asMap)
}

// linkedHashMultimapExample demonstrates the usage of LinkedHashMultimap
func linkedHashMultimapExample() {
	fmt.Println("\n--- LinkedHashMultimap Example ---")

	// Create a new LinkedHashMultimap
	multimap := multimap.NewLinkedHashMultimap[string, string]()

	// Add key-value mappings in a specific order
	multimap.Put("C", "cat")
	multimap.Put("A", "apple")
	multimap.Put("B", "banana")
	multimap.Put("A", "apricot")

	// Print the multimap
	fmt.Println("Multimap:", multimap)

	// Demonstrate insertion order preservation
	fmt.Println("Keys in insertion order:", multimap.Keys())

	// Iterate over entries in insertion order
	fmt.Println("Entries in insertion order:")
	multimap.ForEach(func(key string, value string) {
		fmt.Printf("  %s -> %s\n", key, value)
	})

	// Remove and add to show order preservation
	multimap.Remove("A", "apple")
	multimap.Put("D", "dog")

	// Print the updated multimap
	fmt.Println("Updated multimap:", multimap)
	fmt.Println("Updated keys in insertion order:", multimap.Keys())
}

// ComparableString is a comparable type for TreeMultimap examples
type ComparableString string

// CompareTo compares this ComparableString with another
func (a ComparableString) CompareTo(other interface{}) int {
	if b, ok := other.(ComparableString); ok {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	return 0
}

// ComparableInt is a comparable type for TreeMultimap examples
type ComparableInt int

// CompareTo compares this ComparableInt with another
func (a ComparableInt) CompareTo(other interface{}) int {
	if b, ok := other.(ComparableInt); ok {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}
	return 0
}

// treeMultimapExample demonstrates the usage of TreeMultimap
func treeMultimapExample() {
	fmt.Println("\n--- TreeMultimap Example ---")

	// Create a new TreeMultimap
	multimap := multimap.NewTreeMultimap[ComparableString, ComparableInt]()

	// Add key-value mappings in random order
	multimap.Put("C", 3)
	multimap.Put("A", 1)
	multimap.Put("D", 4)
	multimap.Put("B", 2)
	multimap.Put("A", 5)

	// Print the multimap
	fmt.Println("Multimap:", multimap)

	// Demonstrate natural ordering of keys
	fmt.Println("Keys in natural order:", multimap.Keys())

	// Demonstrate natural ordering of values for a key
	fmt.Println("Values for key 'A' in natural order:", multimap.Get("A"))

	// Iterate over entries in natural key order
	fmt.Println("Entries in natural key order:")
	multimap.ForEach(func(key ComparableString, value ComparableInt) {
		fmt.Printf("  %s -> %d\n", key, value)
	})
}

// immutableMultimapExample demonstrates the usage of ImmutableMultimap
func immutableMultimapExample() {
	fmt.Println("\n--- ImmutableMultimap Example ---")

	// Create a mutable multimap first
	mutable := multimap.NewHashMultimap[string, int]()
	mutable.Put("A", 1)
	mutable.Put("A", 2)
	mutable.Put("B", 3)

	// Create an immutable copy
	immutable := multimap.FromMultimap[string, int](mutable)

	// Print the immutable multimap
	fmt.Println("Immutable multimap:", immutable)
	fmt.Println("Size:", immutable.Size())
	fmt.Println("Values for key 'A':", immutable.Get("A"))

	// Demonstrate immutability
	fmt.Println("Attempting to modify the immutable multimap will cause a panic (commented out):")
	// immutable.Put("C", 4) // This would panic
	// immutable.Remove("A", 1) // This would panic
	// immutable.Clear() // This would panic

	// Create an immutable multimap directly using Of
	immutable2 := multimap.Of[string, int]("X", 10, "Y", 20, "X", 30)
	fmt.Println("Immutable multimap created with Of():", immutable2)
}

// immutableListMultimapExample demonstrates the usage of ImmutableListMultimap
func immutableListMultimapExample() {
	fmt.Println("\n--- ImmutableListMultimap Example ---")

	// Create a mutable multimap first
	mutable := multimap.NewArrayListMultimap[string, string]()
	mutable.Put("colors", "red")
	mutable.Put("colors", "blue")
	mutable.Put("colors", "red") // Duplicate preserved in list
	mutable.Put("shapes", "circle")

	// Create an immutable copy
	immutable := multimap.FromArrayListMultimap(mutable)

	// Print the immutable multimap
	fmt.Println("Immutable list multimap:", immutable)
	fmt.Println("Size:", immutable.Size())
	fmt.Println("Values for key 'colors' (preserves duplicates and order):", immutable.Get("colors"))

	// Create an immutable list multimap directly using ListOf
	immutable2 := multimap.ListOf[string, string](
		"fruits", "apple",
		"fruits", "banana",
		"fruits", "apple", // Duplicate preserved
		"vegetables", "carrot",
	)
	fmt.Println("Immutable list multimap created with ListOf():", immutable2)
	fmt.Println("Values for key 'fruits' (preserves duplicates and order):", immutable2.Get("fruits"))
}

// immutableSetMultimapExample demonstrates the usage of ImmutableSetMultimap
func immutableSetMultimapExample() {
	fmt.Println("\n--- ImmutableSetMultimap Example ---")

	// Create a mutable multimap first
	mutable := multimap.NewHashMultimap[string, string]()
	mutable.Put("colors", "red")
	mutable.Put("colors", "blue")
	mutable.Put("colors", "red") // Duplicate will be ignored in set
	mutable.Put("shapes", "circle")

	// Create an immutable copy
	immutable := multimap.FromHashMultimap(mutable)

	// Print the immutable multimap
	fmt.Println("Immutable set multimap:", immutable)
	fmt.Println("Size:", immutable.Size())
	fmt.Println("Values for key 'colors' (eliminates duplicates):", immutable.Get("colors"))

	// Create an immutable set multimap directly using SetOf
	immutable2 := multimap.SetOf[string, string](
		"fruits", "apple",
		"fruits", "banana",
		"fruits", "apple", // Duplicate will be ignored
		"vegetables", "carrot",
	)
	fmt.Println("Immutable set multimap created with SetOf():", immutable2)
	fmt.Println("Values for key 'fruits' (eliminates duplicates):", immutable2.Get("fruits"))
}

// wordCountExample demonstrates a practical use case for Multimap
func wordCountExample() {
	fmt.Println("\n--- Word Count Example using Multimap ---")

	// Sample text
	text := `
	The quick brown fox jumps over the lazy dog.
	The dog barks, but the fox keeps running.
	Quick thinking helps the fox escape.
	`

	// Create a multimap to store word counts by first letter
	wordCounts := multimap.NewHashMultimap[string, string]()

	// Process the text
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '\'')
	})

	for _, word := range words {
		word = strings.ToLower(word)
		if len(word) > 0 {
			firstLetter := string(word[0])
			wordCounts.Put(firstLetter, word)
		}
	}

	// Print words grouped by first letter
	fmt.Println("Words grouped by first letter:")
	for _, letter := range wordCounts.Keys() {
		words := wordCounts.Get(letter)
		fmt.Printf("  %s: %v\n", letter, words)
	}

	// Count occurrences of each word
	wordOccurrences := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word)
		if len(word) > 0 {
			wordOccurrences[word]++
		}
	}

	// Group words by occurrence count using a multimap
	occurrenceGroups := multimap.NewTreeMultimap[wordCountComparable, string]()
	for word, count := range wordOccurrences {
		occurrenceGroups.Put(wordCountComparable(count), word)
	}

	// Print words grouped by occurrence count
	fmt.Println("\nWords grouped by occurrence count:")
	for _, count := range occurrenceGroups.Keys() {
		words := occurrenceGroups.Get(count)
		fmt.Printf("  %d occurrence(s): %v\n", count, words)
	}
}

// wordCountComparable is a comparable wrapper for int
type wordCountComparable int

func (a wordCountComparable) CompareTo(b wordCountComparable) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}