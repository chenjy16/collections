// Package main provides examples of using the Gkit container library
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/chenjianyu/collections/container/list"
	maps "github.com/chenjianyu/collections/container/map"
	"github.com/chenjianyu/collections/container/queue"
	setpkg "github.com/chenjianyu/collections/container/set"
	"github.com/chenjianyu/collections/container/stack"
)

// RunContainerExamples runs all container examples
func RunContainerExamples() {
	fmt.Println("=== Running Gkit Container Library Examples ===")

	// Run various container examples
	ArrayListExample()
	HashSetExample()
	LinkedHashSetExample()
	LinkedListExample()
	PriorityQueueExample()
	ArrayStackExample()
	CopyOnWriteMapExample()
	ConcurrentHashMapExample()
}

// ArrayListExample demonstrates ArrayList usage
func ArrayListExample() {
	fmt.Println("\n=== ArrayList Example ===")

	// Create a new ArrayList
	arrayList := list.New[int]()

	// Add elements
	fmt.Println("Adding elements:")
	for i := 1; i <= 5; i++ {
		arrayList.Add(i)
		fmt.Printf("Added %d, size: %d\n", i, arrayList.Size())
	}

	// Insert element at specific position
	fmt.Println("\nInserting 100 at index 2:")
	arrayList.Insert(2, 100)
	fmt.Printf("After insertion: %s\n", arrayList)

	// Get elements
	fmt.Println("\nGetting elements:")
	val, err := arrayList.Get(2)
	if err == nil {
		fmt.Printf("Element at index 2: %d\n", val)
	}

	// Set elements
	fmt.Println("\nSetting elements:")
	arrayList.Set(2, 200)
	fmt.Printf("ArrayList after setting: %s\n", arrayList)

	// Check containment
	fmt.Println("\nChecking containment:")
	fmt.Printf("ArrayList contains 200: %t\n", arrayList.Contains(200))
	fmt.Printf("ArrayList contains 100: %t\n", arrayList.Contains(100))

	// Find index
	fmt.Println("\nFinding index:")
	fmt.Printf("Index of 200: %d\n", arrayList.IndexOf(200))
	fmt.Printf("Index of 100: %d\n", arrayList.IndexOf(100))

	// Remove elements
	fmt.Println("\nRemoving elements:")
	removed, _ := arrayList.RemoveAt(2)
	fmt.Printf("Removed element at index 2: %d\n", removed)
	fmt.Printf("ArrayList after removal: %s\n", arrayList)

	// Create sublist
	fmt.Println("\nCreating sublist:")
	subList, _ := arrayList.SubList(1, 3)
	fmt.Printf("Sublist[1:3]: %s\n", subList)

	// Convert to slice
	fmt.Println("\nConverting to slice:")
	slice := arrayList.ToSlice()
	fmt.Printf("Slice: %v\n", slice)

	// Use ForEach
	fmt.Println("\nUsing ForEach:")
	fmt.Print("Elements: ")
	arrayList.ForEach(func(e int) {
		fmt.Printf("%d ", e)
	})
	fmt.Println()

	// Clear list
	fmt.Println("\nClearing list:")
	arrayList.Clear()
	fmt.Printf("Size after clearing: %d\n", arrayList.Size())
	fmt.Printf("Is empty: %t\n", arrayList.IsEmpty())
}

// HashSetExample demonstrates HashSet usage
func HashSetExample() {
	fmt.Println("\n=== HashSet Example ===")

	// Create a new HashSet
	set := setpkg.New[string]()

	// Add elements
	fmt.Println("Adding elements:")
	fruits := []string{"apple", "banana", "orange", "apple", "grape"}
	for _, fruit := range fruits {
		set.Add(fruit)
		fmt.Printf("Added '%s', size: %d\n", fruit, set.Size())
	}

	// Check containment
	fmt.Println("\nChecking containment:")
	fmt.Printf("Contains 'banana': %t\n", set.Contains("banana"))
	fmt.Printf("Contains 'watermelon': %t\n", set.Contains("watermelon"))

	// Remove elements
	fmt.Println("\nRemoving elements:")
	set.Remove("banana")
	fmt.Printf("After removing 'banana': %s\n", set)

	// Create another set
	fmt.Println("\nCreating another set:")
	set2 := setpkg.FromSlice([]string{"grape", "watermelon", "lemon"})
	fmt.Printf("Set2: %s\n", set2)

	// Set operations
	fmt.Println("\nSet operations:")
	unionSet := set.Union(set2)
	fmt.Printf("Union: %s\n", unionSet)

	intersectionSet := set.Intersection(set2)
	fmt.Printf("Intersection: %s\n", intersectionSet)

	differenceSet := set.Difference(set2)
	fmt.Printf("Difference: %s\n", differenceSet)

	// Subset and superset
	fmt.Println("\nSubset and superset:")
	subset := setpkg.FromSlice([]string{"apple"})
	fmt.Printf("Subset: %s\n", subset)
	fmt.Printf("Is subset: %t\n", subset.IsSubsetOf(set))
	fmt.Printf("Is superset: %t\n", set.IsSupersetOf(subset))

	// Convert to slice
	fmt.Println("\nConverting to slice:")
	slice := set.ToSlice()
	fmt.Printf("Slice: %v\n", slice)

	// Use ForEach
	fmt.Println("\nUsing ForEach:")
	fmt.Print("Elements: ")
	set.ForEach(func(e string) {
		fmt.Printf("%s ", e)
	})
	fmt.Println()

	// Clear set
	fmt.Println("\nClearing set:")
	set.Clear()
	fmt.Printf("Size after clearing: %d\n", set.Size())
	fmt.Printf("Is empty: %t\n", set.IsEmpty())
}

// LinkedListExample demonstrates LinkedList usage
func LinkedListExample() {
	fmt.Println("\n=== LinkedList Example ===")

	// Create a new LinkedList
	list := queue.New[int]()

	// Add elements
	fmt.Println("Adding elements:")
	for i := 0; i < 5; i++ {
		_ = list.Add(i)
		fmt.Printf("Added %d, size: %d\n", i, list.Size())
	}

	// Use as queue
	fmt.Println("\nUsing as queue:")
	val, _ := list.Peek()
	fmt.Printf("Queue head element: %d\n", val)

	val, _ = list.Remove()
	fmt.Printf("Removed queue head element: %d\n", val)
	fmt.Printf("Queue after removal: %s\n", list)

	// Use as deque
	fmt.Println("\nUsing as deque:")
	_ = list.AddFirst(10)
	_ = list.AddLast(20)
	fmt.Printf("Deque after adding elements: %s\n", list)

	val, _ = list.GetFirst()
	fmt.Printf("Deque head element: %d\n", val)

	val, _ = list.GetLast()
	fmt.Printf("Deque tail element: %d\n", val)

	val, _ = list.RemoveFirst()
	fmt.Printf("Removed deque head element: %d\n", val)

	val, _ = list.RemoveLast()
	fmt.Printf("Removed deque tail element: %d\n", val)

	fmt.Printf("Deque after removal: %s\n", list)

	// Convert to slice
	fmt.Println("\nConverting to slice:")
	slice := list.ToSlice()
	fmt.Printf("Slice: %v\n", slice)

	// Use ForEach
	fmt.Println("\nUsing ForEach:")
	fmt.Print("Elements: ")
	list.ForEach(func(e int) {
		fmt.Printf("%d ", e)
	})
	fmt.Println()

	// Clear list
	fmt.Println("\nClearing list:")
	list.Clear()
	fmt.Printf("Size after clearing: %d\n", list.Size())
	fmt.Printf("Is empty: %t\n", list.IsEmpty())
}

// PriorityQueueExample demonstrates PriorityQueue usage
func PriorityQueueExample() {
	fmt.Println("\n=== PriorityQueue Example ===")

	// Create a min-heap priority queue
	minHeap := queue.NewPriorityQueueWithComparator(func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// Add elements (unordered)
	fmt.Println("Adding elements:")
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = minHeap.Add(e)
		fmt.Printf("Added %d, size: %d\n", e, minHeap.Size())
	}

	// View minimum element
	fmt.Println("\nViewing minimum element:")
	val, _ := minHeap.Peek()
	fmt.Printf("Minimum element: %d\n", val)

	// Remove and get minimum elements
	fmt.Println("\nRemoving and getting minimum elements:")
	for !minHeap.IsEmpty() {
		val, _ := minHeap.Remove()
		fmt.Printf("Removed: %d\n", val)
	}

	// Create a max-heap priority queue
	fmt.Println("\nCreating max heap:")
	maxHeap := queue.NewPriorityQueueWithComparator(func(a, b int) int {
		if a > b {
			return -1 // Reverse comparison result
		} else if a < b {
			return 1
		}
		return 0
	})

	// Add elements
	for _, e := range elements {
		_ = maxHeap.Add(e)
	}

	// Get sorted slice
	fmt.Println("\nGetting sorted slice:")
	sortedSlice := maxHeap.ToSortedSlice()
	fmt.Printf("Sorted slice (descending): %v\n", sortedSlice)

	// Remove and get maximum elements
	fmt.Println("\nRemoving and getting maximum elements:")
	for !maxHeap.IsEmpty() {
		val, _ := maxHeap.Remove()
		fmt.Printf("Removed: %d\n", val)
	}
}

// ArrayStackExample demonstrates ArrayStack usage
func ArrayStackExample() {
	fmt.Println("\n=== ArrayStack Example ===")

	// Create a new ArrayStack
	stack := stack.New[int]()

	// Push elements
	fmt.Println("Pushing elements:")
	for i := 0; i < 5; i++ {
		stack.Push(i)
		fmt.Printf("Pushed %d, size: %d\n", i, stack.Size())
	}

	// Peek at top element
	fmt.Println("\nPeeking at top element:")
	val, _ := stack.Peek()
	fmt.Printf("Top element: %d\n", val)

	// Pop elements
	fmt.Println("\nPopping elements:")
	val, _ = stack.Pop()
	fmt.Printf("Popped top element: %d\n", val)
	fmt.Printf("Stack after popping: %s\n", stack)

	// Search for elements
	fmt.Println("\nSearching for elements:")
	pos := stack.Search(2)
	fmt.Printf("Position of element 2: %d\n", pos)

	pos = stack.Search(10)
	fmt.Printf("Position of element 10: %d\n", pos)

	// Convert to slice
	fmt.Println("\nConverting to slice:")
	slice := stack.ToSlice()
	fmt.Printf("Slice: %v\n", slice)

	// Use ForEach
	fmt.Println("\nUsing ForEach:")
	fmt.Print("Elements: ")
	stack.ForEach(func(e int) {
		fmt.Printf("%d ", e)
	})
	fmt.Println()

	// Clear stack
	fmt.Println("\nClearing stack:")
	stack.Clear()
	fmt.Printf("Size after clearing: %d\n", stack.Size())
	fmt.Printf("Is empty: %t\n", stack.IsEmpty())
}

// CopyOnWriteMapExample demonstrates CopyOnWriteMap usage
func CopyOnWriteMapExample() {
	fmt.Println("\n=== CopyOnWriteMap Example ===")

	// Create a new CopyOnWriteMap
	cowMap := maps.NewCopyOnWriteMap[string, int]()

	// Basic operations
	fmt.Println("Basic operations:")
	cowMap.Put("apple", 10)
	cowMap.Put("banana", 20)
	cowMap.Put("cherry", 30)
	fmt.Printf("Map size: %d\n", cowMap.Size())
	fmt.Printf("Map content: %s\n", cowMap.String())

	// Get elements
	fmt.Println("\nGetting elements:")
	if value, exists := cowMap.Get("apple"); exists {
		fmt.Printf("apple: %d\n", value)
	}

	// Check containment
	fmt.Println("\nChecking containment:")
	fmt.Printf("Contains 'banana': %t\n", cowMap.ContainsKey("banana"))
	fmt.Printf("Contains value 20: %t\n", cowMap.ContainsValue(20))

	// Advanced operations
	fmt.Println("\nAdvanced operations:")

	// PutIfAbsent - add only if key doesn't exist
	if _, existed := cowMap.PutIfAbsent("grape", 40); !existed {
		fmt.Println("Successfully added grape")
	}

	if _, existed := cowMap.PutIfAbsent("apple", 50); existed {
		fmt.Println("apple already exists, not inserted")
	}

	// Replace - replace only existing keys
	if oldValue, replaced := cowMap.Replace("apple", 15); replaced {
		fmt.Printf("Replaced apple: %d -> 15\n", oldValue)
	}

	// ReplaceIf - conditional replacement
	if cowMap.ReplaceIf("banana", 20, 25) {
		fmt.Println("Conditionally replaced banana: 20 -> 25")
	}

	// Iterate elements
	fmt.Println("\nIterating elements:")
	cowMap.ForEach(func(key string, value int) {
		fmt.Printf("%s: %d\n", key, value)
	})

	// Get snapshot
	fmt.Println("\nSnapshot functionality:")
	snapshot1 := cowMap.Snapshot()
	fmt.Printf("Snapshot1 size: %d\n", len(snapshot1))

	// Modify original map
	cowMap.Put("date", 50)
	snapshot2 := cowMap.Snapshot()
	fmt.Printf("Snapshot2 size: %d\n", len(snapshot2))
	fmt.Printf("Snapshot independence: %t\n", len(snapshot1) != len(snapshot2))

	// Concurrent safety demonstration
	fmt.Println("\nConcurrent safety demonstration:")
	demonstrateConcurrency(cowMap)

	// Conversion operations
	fmt.Println("\nConversion operations:")
	goMap := cowMap.ToMap()
	fmt.Printf("Converted to Go map: %v\n", goMap)

	keys := cowMap.Keys()
	values := cowMap.Values()
	fmt.Printf("All keys: %v\n", keys)
	fmt.Printf("All values: %v\n", values)

	// Clear map
	fmt.Println("\nClearing map:")
	cowMap.Clear()
	fmt.Printf("Size after clearing: %d\n", cowMap.Size())
	fmt.Printf("Is empty: %t\n", cowMap.IsEmpty())
}

// demonstrateConcurrency demonstrates CopyOnWriteMap's concurrent safety
func demonstrateConcurrency(cowMap *maps.CopyOnWriteMap[string, int]) {
	// Pre-populate some data
	for i := 0; i < 10; i++ {
		cowMap.Put(fmt.Sprintf("key%d", i), i)
	}

	var wg sync.WaitGroup

	// Start multiple readers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", j%10)
				value, _ := cowMap.Get(key)
				_ = value // Use value to avoid compiler optimization
			}
			fmt.Printf("  Reader%d completed\n", readerID)
		}(i)
	}

	// Start a writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 10; i < 20; i++ {
			cowMap.Put(fmt.Sprintf("key%d", i), i)
			time.Sleep(1 * time.Millisecond) // Simulate write operation interval
		}
		fmt.Println("  Writer completed")
	}()

	wg.Wait()
	fmt.Printf("  Concurrent operations completed, final map size: %d\n", cowMap.Size())
}

// ConcurrentHashMapExample demonstrates ConcurrentHashMap usage
func ConcurrentHashMapExample() {
	fmt.Println("\n=== ConcurrentHashMap Example ===")

	// Create a new ConcurrentHashMap
	chm := maps.NewConcurrentHashMap[string, int]()

	// Basic operations
	fmt.Println("Basic operations:")
	chm.Put("apple", 10)
	chm.Put("banana", 20)
	chm.Put("cherry", 30)
	fmt.Printf("Map size: %d\n", chm.Size())
	fmt.Printf("Map content: %s\n", chm.String())

	// Get elements
	fmt.Println("\nGetting elements:")
	if value, exists := chm.Get("apple"); exists {
		fmt.Printf("apple: %d\n", value)
	}

	// Check containment
	fmt.Println("\nChecking containment:")
	fmt.Printf("Contains 'banana': %t\n", chm.ContainsKey("banana"))
	fmt.Printf("Contains value 20: %t\n", chm.ContainsValue(20))

	// Advanced operations
	fmt.Println("\nAdvanced operations:")

	// PutIfAbsent - add only if key doesn't exist
	if _, existed := chm.PutIfAbsent("grape", 40); !existed {
		fmt.Println("Successfully added grape")
	}

	if _, existed := chm.PutIfAbsent("apple", 50); existed {
		fmt.Println("apple already exists, not inserted")
	}

	// Replace - replace only existing keys
	if oldValue, replaced := chm.Replace("apple", 15); replaced {
		fmt.Printf("Replaced apple: %d -> 15\n", oldValue)
	}

	// ReplaceIf - conditional replacement
	if chm.ReplaceIf("banana", 20, 25) {
		fmt.Println("Conditionally replaced banana: 20 -> 25")
	}

	// ComputeIfAbsent - compute missing values
	fmt.Println("\nCompute operations:")
	value := chm.ComputeIfAbsent("kiwi", func(key string) int {
		return len(key) * 5 // Calculate value based on key length
	})
	fmt.Printf("Computed value for kiwi: %d\n", value)

	// ComputeIfPresent - compute existing values
	if newValue, exists := chm.ComputeIfPresent("apple", func(key string, oldValue int) int {
		return oldValue * 2
	}); exists {
		fmt.Printf("Recomputed value for apple: %d\n", newValue)
	}

	// Iterate elements
	fmt.Println("\nIterating elements:")
	chm.ForEach(func(key string, value int) {
		fmt.Printf("%s: %d\n", key, value)
	})

	// Batch operations
	fmt.Println("\nBatch operations:")
	batchData := map[string]int{
		"mango":  60,
		"orange": 70,
		"peach":  80,
	}
	chm.PutAllFromMap(batchData)
	fmt.Printf("Size after batch addition: %d\n", chm.Size())

	// Get snapshot
	fmt.Println("\nSnapshot functionality:")
	snapshot1 := chm.Snapshot()
	fmt.Printf("Snapshot1 size: %d\n", len(snapshot1))

	// Modify original map
	chm.Put("date", 90)
	snapshot2 := chm.Snapshot()
	fmt.Printf("Snapshot2 size: %d\n", len(snapshot2))
	fmt.Printf("Snapshot independence: %t\n", len(snapshot1) != len(snapshot2))

	// Concurrent safety demonstration
	fmt.Println("\nConcurrent safety demonstration:")
	demonstrateConcurrentHashMapConcurrency(chm)

	// Conversion operations
	fmt.Println("\nConversion operations:")
	goMap := chm.ToMap()
	fmt.Printf("Converted to Go map size: %d\n", len(goMap))

	keys := chm.Keys()
	values := chm.Values()
	fmt.Printf("Number of keys: %d\n", len(keys))
	fmt.Printf("Number of values: %d\n", len(values))

	// Remove operations
	fmt.Println("\nRemove operations:")
	if oldValue, removed := chm.Remove("banana"); removed {
		fmt.Printf("Removed banana, old value: %d\n", oldValue)
	}
	fmt.Printf("Size after removal: %d\n", chm.Size())

	// Clear map
	fmt.Println("\nClearing map:")
	chm.Clear()
	fmt.Printf("Size after clearing: %d\n", chm.Size())
	fmt.Printf("Is empty: %t\n", chm.IsEmpty())
}

// demonstrateConcurrentHashMapConcurrency demonstrates ConcurrentHashMap's concurrent safety
func demonstrateConcurrentHashMapConcurrency(chm *maps.ConcurrentHashMap[string, int]) {
	// Pre-populate some data
	for i := 0; i < 10; i++ {
		chm.Put(fmt.Sprintf("key%d", i), i)
	}

	var wg sync.WaitGroup

	// Start multiple readers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", j%10)
				value, _ := chm.Get(key)
				_ = value // Use value to avoid compiler optimization
			}
			fmt.Printf("  Reader%d completed\n", readerID)
		}(i)
	}

	// Start multiple writers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("writer%d_key%d", writerID, j)
				chm.Put(key, j)
			}
			fmt.Printf("  Writer%d completed\n", writerID)
		}(i)
	}

	// Start a computer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("key%d", i)
			chm.ComputeIfPresent(key, func(k string, v int) int {
				return v + 100
			})
		}
		fmt.Println("  Computer completed")
	}()

	wg.Wait()
	fmt.Printf("  Concurrent operations completed, final map size: %d\n", chm.Size())
}
