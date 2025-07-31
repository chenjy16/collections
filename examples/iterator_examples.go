package main

import (
	"fmt"
	"github.com/chenjianyu/collections/container/list"
	"github.com/chenjianyu/collections/container/set"
)

// IteratorExamples demonstrates the usage and benefits of Iterator pattern
func IteratorExamples() {
	fmt.Println("=== Iterator Pattern Usage Examples ===")

	// 1. ArrayList Iterator Example
	arrayListIteratorExample()

	// 2. LinkedList Iterator Example
	linkedListIteratorExample()

	// 3. HashSet Iterator Example
	hashSetIteratorExample()

	// 4. TreeSet Iterator Example
	treeSetIteratorExample()

	// 5. Iterator vs ForEach Comparison
	iteratorVsForEachComparison()

	// 6. Iterator Safe Removal Operations
	iteratorSafeRemovalExample()

	// 7. Multiple Iterator Concurrent Usage
	multipleIteratorsExample()
}

// arrayListIteratorExample demonstrates ArrayList iterator usage
func arrayListIteratorExample() {
	fmt.Println("\n--- ArrayList Iterator Example ---")

	// Create ArrayList and add elements
	arrayList := list.New[int]()
	for i := 1; i <= 5; i++ {
		arrayList.Add(i * 10)
	}
	fmt.Printf("ArrayList: %s\n", arrayList)

	// Use Iterator to traverse
	fmt.Println("Using Iterator to traverse:")
	iterator := arrayList.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok {
			fmt.Printf("  Element: %d\n", element)
		}
	}

	// Use Iterator for conditional removal
	fmt.Println("\nUsing Iterator to remove elements at even positions:")
	iterator = arrayList.Iterator()
	position := 0
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok {
			if position%2 == 0 { // Remove elements at even positions
				if iterator.Remove() {
					fmt.Printf("  Removed element: %d (position: %d)\n", element, position)
				}
			}
			position++
		}
	}
	fmt.Printf("ArrayList after removal: %s\n", arrayList)
}

// linkedListIteratorExample demonstrates LinkedList iterator usage
func linkedListIteratorExample() {
	fmt.Println("\n--- LinkedList Iterator Example ---")

	// Create LinkedList and add elements
	linkedList := list.NewLinkedList[string]()
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, word := range words {
		linkedList.Add(word)
	}
	fmt.Printf("LinkedList: %s\n", linkedList)

	// Use Iterator to find specific elements
	fmt.Println("Using Iterator to find words containing 'a':")
	iterator := linkedList.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok {
			if containsChar(element, 'a') {
				fmt.Printf("  Found: %s\n", element)
			}
		}
	}

	// Use Iterator to remove words with length less than 5
	fmt.Println("\nUsing Iterator to remove words with length less than 5:")
	iterator = linkedList.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok && len(element) < 5 {
			if iterator.Remove() {
				fmt.Printf("  Removed: %s\n", element)
			}
		}
	}
	fmt.Printf("LinkedList after removal: %s\n", linkedList)
}

// hashSetIteratorExample demonstrates HashSet iterator usage
func hashSetIteratorExample() {
	fmt.Println("\n--- HashSet Iterator Example ---")

	// Create HashSet and add elements
	hashSet := set.New[int]()
	numbers := []int{15, 23, 8, 42, 16, 4, 35}
	for _, num := range numbers {
		hashSet.Add(num)
	}
	fmt.Printf("HashSet: %s\n", hashSet)

	// Use Iterator to count even and odd numbers
	fmt.Println("Using Iterator to count even and odd numbers:")
	iterator := hashSet.Iterator()
	evenCount, oddCount := 0, 0
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok {
			if element%2 == 0 {
				evenCount++
				fmt.Printf("  Even: %d\n", element)
			} else {
				oddCount++
				fmt.Printf("  Odd: %d\n", element)
			}
		}
	}
	fmt.Printf("Even count: %d, Odd count: %d\n", evenCount, oddCount)

	// Use Iterator to remove elements less than 20
	fmt.Println("\nUsing Iterator to remove elements less than 20:")
	iterator = hashSet.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok && element < 20 {
			if iterator.Remove() {
				fmt.Printf("  Removed: %d\n", element)
			}
		}
	}
	fmt.Printf("HashSet after removal: %s\n", hashSet)
}

// treeSetIteratorExample demonstrates TreeSet iterator usage
func treeSetIteratorExample() {
	fmt.Println("\n--- TreeSet Iterator Example ---")

	// Create TreeSet and add elements
	treeSet := set.NewTreeSet[int]()
	numbers := []int{50, 30, 70, 20, 40, 60, 80}
	for _, num := range numbers {
		treeSet.Add(num)
	}
	fmt.Printf("TreeSet (sorted): %s\n", treeSet)

	// Use Iterator to traverse in order
	fmt.Println("Using Iterator to traverse in order:")
	iterator := treeSet.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok {
			fmt.Printf("  %d", element)
		}
	}
	fmt.Println()

	// Use Iterator to find elements within range
	fmt.Println("\nUsing Iterator to find elements in range 30-60:")
	iterator = treeSet.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok && element >= 30 && element <= 60 {
			fmt.Printf("  Element in range: %d\n", element)
		}
	}
}

// iteratorVsForEachComparison compares Iterator and ForEach approaches
func iteratorVsForEachComparison() {
	fmt.Println("\n--- Iterator vs ForEach Comparison ---")

	arrayList := list.New[int]()
	for i := 1; i <= 5; i++ {
		arrayList.Add(i)
	}

	// Using ForEach (simple traversal)
	fmt.Println("Using ForEach to traverse:")
	arrayList.ForEach(func(element int) {
		fmt.Printf("  %d", element)
	})
	fmt.Println()

	// Using Iterator (controllable traversal process)
	fmt.Println("Using Iterator to traverse (can be interrupted):")
	iterator := arrayList.Iterator()
	count := 0
	for iterator.HasNext() && count < 3 { // Only traverse first 3 elements
		element, ok := iterator.Next()
		if ok {
			fmt.Printf("  %d", element)
			count++
		}
	}
	fmt.Println()

	fmt.Println("\nIterator advantages:")
	fmt.Println("  1. Can control traversal process (pause, interrupt)")
	fmt.Println("  2. Can safely remove elements")
	fmt.Println("  3. Can get more information during traversal")
	fmt.Println("  4. Supports multiple concurrent iterators")
}

// iteratorSafeRemovalExample demonstrates safe element removal using iterator
func iteratorSafeRemovalExample() {
	fmt.Println("\n--- Iterator Safe Removal Example ---")

	arrayList := list.New[int]()
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, num := range numbers {
		arrayList.Add(num)
	}
	fmt.Printf("Original list: %s\n", arrayList)

	// Wrong way to remove (causes index issues)
	fmt.Println("\n❌ Wrong way to remove (direct removal in loop):")
	fmt.Println("  // This approach causes index confusion")
	fmt.Println("  for i := 0; i < list.Size(); i++ {")
	fmt.Println("      if element % 2 == 0 {")
	fmt.Println("          list.RemoveAt(i) // Dangerous!")
	fmt.Println("      }")
	fmt.Println("  }")

	// Correct way to remove (using Iterator)
	fmt.Println("\n✅ Correct way to remove (using Iterator):")
	iterator := arrayList.Iterator()
	for iterator.HasNext() {
		element, ok := iterator.Next()
		if ok && element%2 == 0 { // Remove even numbers
			if iterator.Remove() {
				fmt.Printf("  Safely removed even number: %d\n", element)
			}
		}
	}
	fmt.Printf("List after removing even numbers: %s\n", arrayList)
}

// multipleIteratorsExample demonstrates using multiple iterators concurrently
func multipleIteratorsExample() {
	fmt.Println("\n--- Multiple Iterator Concurrent Usage Example ---")

	arrayList := list.New[int]()
	for i := 1; i <= 10; i++ {
		arrayList.Add(i)
	}
	fmt.Printf("List: %s\n", arrayList)

	// Create two independent iterators
	iterator1 := arrayList.Iterator()
	iterator2 := arrayList.Iterator()

	fmt.Println("\nUsing two independent Iterators:")
	fmt.Println("Iterator1 traverses first half:")
	count := 0
	for iterator1.HasNext() && count < 5 {
		element, ok := iterator1.Next()
		if ok {
			fmt.Printf("  Iterator1: %d\n", element)
			count++
		}
	}

	fmt.Println("Iterator2 traverses all elements:")
	for iterator2.HasNext() {
		element, ok := iterator2.Next()
		if ok {
			fmt.Printf("  Iterator2: %d\n", element)
		}
	}

	fmt.Println("Iterator1 continues traversing remaining part:")
	for iterator1.HasNext() {
		element, ok := iterator1.Next()
		if ok {
			fmt.Printf("  Iterator1: %d\n", element)
		}
	}
}

// containsChar checks if a string contains a specific character
func containsChar(s string, c rune) bool {
	for _, char := range s {
		if char == c {
			return true
		}
	}
	return false
}

// Summary of Iterator pattern's core functions and advantages
func printIteratorSummary() {
	fmt.Println("\n=== Iterator Pattern Functions and Advantages Summary ===")
	fmt.Println()
	fmt.Println("🎯 Iterator's core functions:")
	fmt.Println("  1. Unified traversal interface - Provides unified traversal for different data structures")
	fmt.Println("  2. Encapsulates traversal logic - Hides implementation details, provides clean traversal interface")
	fmt.Println("  3. Supports multiple traversals - Can have multiple independent iterators simultaneously")
	fmt.Println("  4. Safe modification operations - Safely remove elements during traversal")
	fmt.Println()
	fmt.Println("✨ Iterator's main advantages:")
	fmt.Println("  1. Decoupling - Separates traversal logic from data structure implementation")
	fmt.Println("  2. Consistency - All collection types use the same traversal interface")
	fmt.Println("  3. Flexibility - Can control traversal process, supports pause and resume")
	fmt.Println("  4. Safety - Avoids problems caused by direct collection modification during traversal")
	fmt.Println("  5. Concurrency - Supports multiple iterators working simultaneously")
	fmt.Println()
	fmt.Println("🔧 Iterator interface methods:")
	fmt.Println("  - HasNext() bool     : Check if there are more elements")
	fmt.Println("  - Next() (E, bool)   : Get next element")
	fmt.Println("  - Remove() bool      : Remove current element (optional operation)")
	fmt.Println()
	fmt.Println("📚 Applicable scenarios:")
	fmt.Println("  1. Need to traverse all elements in a collection")
	fmt.Println("  2. Need to remove elements during traversal")
	fmt.Println("  3. Need to control traversal process (pause, skip, etc.)")
	fmt.Println("  4. Need multiple independent traversal processes")
	fmt.Println("  5. Need unified handling of different collection types")
}

func RunIteratorExamples() {
	IteratorExamples()
	printIteratorSummary()
}