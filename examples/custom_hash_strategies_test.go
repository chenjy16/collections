package examples

import (
	"strings"
	"testing"

	"github.com/chenjianyu/collections/container/common"
	"github.com/chenjianyu/collections/container/map"
	"github.com/chenjianyu/collections/container/set"
)

// Example: Case-insensitive string hash strategy
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

// Example: Custom struct with specialized hash strategy
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

// Example: Case-insensitive person hash strategy
type CaseInsensitivePersonHashStrategy struct{}

func (s *CaseInsensitivePersonHashStrategy) Hash(p Person) uint64 {
	lower := strings.ToLower(p.Name)
	hash := uint64(0)
	for i := 0; i < len(lower); i++ {
		hash = hash*31 + uint64(lower[i])
	}
	return hash
}

func (s *CaseInsensitivePersonHashStrategy) Equals(a, b Person) bool {
	return strings.ToLower(a.Name) == strings.ToLower(b.Name)
}

func TestCaseInsensitiveStringHashSet(t *testing.T) {
	strategy := &CaseInsensitiveHashStrategy{}
	hashSet := set.NewWithHashStrategy[string](strategy)

	// Add some strings with different cases
	hashSet.Add("Hello")
	hashSet.Add("HELLO")
	hashSet.Add("hello")
	hashSet.Add("World")
	hashSet.Add("WORLD")

	// Should only contain 2 unique elements (case-insensitive)
	if hashSet.Size() != 2 {
		t.Errorf("Expected size 2, got %d", hashSet.Size())
	}

	// Should contain elements regardless of case
	if !hashSet.Contains("hello") {
		t.Error("Expected to contain 'hello'")
	}
	if !hashSet.Contains("HELLO") {
		t.Error("Expected to contain 'HELLO'")
	}
	if !hashSet.Contains("Hello") {
		t.Error("Expected to contain 'Hello'")
	}
	if !hashSet.Contains("world") {
		t.Error("Expected to contain 'world'")
	}

	t.Logf("HashSet contents: %s", hashSet.String())
}

func TestCaseInsensitiveStringLinkedHashSet(t *testing.T) {
	strategy := &CaseInsensitiveHashStrategy{}
	linkedHashSet := set.NewLinkedHashSetWithHashStrategy[string](strategy)

	// Add some strings with different cases
	linkedHashSet.Add("Hello")
	linkedHashSet.Add("HELLO") // Should not be added (duplicate)
	linkedHashSet.Add("World")
	linkedHashSet.Add("hello") // Should not be added (duplicate)

	// Should only contain 2 unique elements (case-insensitive)
	if linkedHashSet.Size() != 2 {
		t.Errorf("Expected size 2, got %d", linkedHashSet.Size())
	}

	// Should maintain insertion order
	elements := linkedHashSet.ToSlice()
	if len(elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(elements))
	}
	if elements[0] != "Hello" {
		t.Errorf("Expected first element to be 'Hello', got '%s'", elements[0])
	}
	if elements[1] != "World" {
		t.Errorf("Expected second element to be 'World', got '%s'", elements[1])
	}

	t.Logf("LinkedHashSet contents: %s", linkedHashSet.String())
}

func TestCaseInsensitiveStringConcurrentHashMap(t *testing.T) {
	strategy := &CaseInsensitiveHashStrategy{}
	hashMap := maps.NewConcurrentHashMapWithHashStrategy[string, int](strategy)

	// Add some key-value pairs with different cases
	hashMap.Put("Hello", 1)
	hashMap.Put("HELLO", 2) // Should replace the previous value
	hashMap.Put("World", 3)
	hashMap.Put("hello", 4) // Should replace the previous value

	// Should only contain 2 unique keys (case-insensitive)
	if hashMap.Size() != 2 {
		t.Errorf("Expected size 2, got %d", hashMap.Size())
	}

	// Should retrieve values regardless of case
	if val, exists := hashMap.Get("hello"); !exists || val != 4 {
		t.Errorf("Expected to get value 4 for 'hello', got %d (exists: %t)", val, exists)
	}
	if val, exists := hashMap.Get("HELLO"); !exists || val != 4 {
		t.Errorf("Expected to get value 4 for 'HELLO', got %d (exists: %t)", val, exists)
	}
	if val, exists := hashMap.Get("Hello"); !exists || val != 4 {
		t.Errorf("Expected to get value 4 for 'Hello', got %d (exists: %t)", val, exists)
	}
	if val, exists := hashMap.Get("world"); !exists || val != 3 {
		t.Errorf("Expected to get value 3 for 'world', got %d (exists: %t)", val, exists)
	}

	t.Logf("ConcurrentHashMap contents: %s", hashMap.String())
}

func TestPersonByNameHashSet(t *testing.T) {
	strategy := &PersonByNameHashStrategy{}
	hashSet := set.NewWithHashStrategy[Person](strategy)

	// Add persons with same name but different ages
	p1 := Person{Name: "John", Age: 25}
	p2 := Person{Name: "John", Age: 30} // Same name, different age
	p3 := Person{Name: "Jane", Age: 25}

	hashSet.Add(p1)
	hashSet.Add(p2) // Should not be added (same name)
	hashSet.Add(p3)

	// Should only contain 2 unique persons (by name)
	if hashSet.Size() != 2 {
		t.Errorf("Expected size 2, got %d", hashSet.Size())
	}

	// Should contain persons regardless of age when name matches
	if !hashSet.Contains(Person{Name: "John", Age: 40}) {
		t.Error("Expected to contain person with name 'John' regardless of age")
	}
	if !hashSet.Contains(Person{Name: "Jane", Age: 50}) {
		t.Error("Expected to contain person with name 'Jane' regardless of age")
	}

	t.Logf("HashSet contents: %s", hashSet.String())
}

func TestCaseInsensitivePersonHashSet(t *testing.T) {
	strategy := &CaseInsensitivePersonHashStrategy{}
	hashSet := set.NewWithHashStrategy[Person](strategy)

	// Add persons with same name in different cases
	p1 := Person{Name: "John", Age: 25}
	p2 := Person{Name: "JOHN", Age: 30} // Same name different case
	p3 := Person{Name: "john", Age: 35} // Same name different case
	p4 := Person{Name: "Jane", Age: 25}

	hashSet.Add(p1)
	hashSet.Add(p2) // Should not be added (same name, case-insensitive)
	hashSet.Add(p3) // Should not be added (same name, case-insensitive)
	hashSet.Add(p4)

	// Should only contain 2 unique persons (by name, case-insensitive)
	if hashSet.Size() != 2 {
		t.Errorf("Expected size 2, got %d", hashSet.Size())
	}

	// Should contain persons regardless of name case
	if !hashSet.Contains(Person{Name: "JOHN", Age: 40}) {
		t.Error("Expected to contain person with name 'JOHN' (case-insensitive)")
	}
	if !hashSet.Contains(Person{Name: "jane", Age: 50}) {
		t.Error("Expected to contain person with name 'jane' (case-insensitive)")
	}

	t.Logf("HashSet contents: %s", hashSet.String())
}

func TestHashSetUnionWithCustomStrategy(t *testing.T) {
	strategy := &CaseInsensitiveHashStrategy{}
	
	set1 := set.NewWithHashStrategy[string](strategy)
	set1.Add("Hello")
	set1.Add("World")

	set2 := set.NewWithHashStrategy[string](strategy)
	set2.Add("HELLO") // Should be considered duplicate
	set2.Add("Go")

	union := set1.Union(set2)
	
	// Should contain 3 unique elements (case-insensitive)
	if union.Size() != 3 {
		t.Errorf("Expected union size 3, got %d", union.Size())
	}

	// Should contain all elements regardless of case
	if !union.Contains("hello") {
		t.Error("Expected union to contain 'hello'")
	}
	if !union.Contains("world") {
		t.Error("Expected union to contain 'world'")
	}
	if !union.Contains("go") {
		t.Error("Expected union to contain 'go'")
	}

	t.Logf("Union result: %s", union.String())
}

func TestHashSetIntersectionWithCustomStrategy(t *testing.T) {
	strategy := &CaseInsensitiveHashStrategy{}
	
	set1 := set.NewWithHashStrategy[string](strategy)
	set1.Add("Hello")
	set1.Add("World")

	set2 := set.NewWithHashStrategy[string](strategy)
	set2.Add("HELLO") // Should match with "Hello" from set1
	set2.Add("Go")

	intersection := set1.Intersection(set2)
	
	// Should contain 1 element (case-insensitive intersection)
	if intersection.Size() != 1 {
		t.Errorf("Expected intersection size 1, got %d", intersection.Size())
	}

	// Should contain the intersecting element
	if !intersection.Contains("hello") {
		t.Error("Expected intersection to contain 'hello'")
	}

	t.Logf("Intersection result: %s", intersection.String())
}

func TestFunctionalHashStrategy(t *testing.T) {
	// Create a functional hash strategy for case-insensitive strings
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
	hashSet.Add("Hello")
	hashSet.Add("HELLO") // Should not be added (duplicate)
	hashSet.Add("World")

	if hashSet.Size() != 2 {
		t.Errorf("Expected size 2, got %d", hashSet.Size())
	}

	if !hashSet.Contains("hello") {
		t.Error("Expected to contain 'hello'")
	}

	t.Logf("Functional strategy HashSet: %s", hashSet.String())
}

func TestDefaultComparableHashStrategy(t *testing.T) {
	// Test default strategy with regular strings
	hashSet := set.New[string]() // Uses default strategy
	hashSet.Add("Hello")
	hashSet.Add("hello") // Should be added (case-sensitive by default)
	hashSet.Add("World")

	if hashSet.Size() != 3 {
		t.Errorf("Expected size 3, got %d", hashSet.Size())
	}

	if !hashSet.Contains("Hello") {
		t.Error("Expected to contain 'Hello'")
	}
	if !hashSet.Contains("hello") {
		t.Error("Expected to contain 'hello'")
	}

	t.Logf("Default strategy HashSet: %s", hashSet.String())
}

// Benchmark tests to compare performance with and without custom strategies
func BenchmarkHashSetWithDefaultStrategy(b *testing.B) {
	hashSet := set.New[string]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashSet.Add("test" + string(rune(i%1000)))
	}
}

func BenchmarkHashSetWithCustomStrategy(b *testing.B) {
	strategy := &CaseInsensitiveHashStrategy{}
	hashSet := set.NewWithHashStrategy[string](strategy)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashSet.Add("test" + string(rune(i%1000)))
	}
}

func BenchmarkConcurrentHashMapWithDefaultStrategy(b *testing.B) {
	hashMap := maps.NewConcurrentHashMap[string, int]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashMap.Put("key"+string(rune(i%1000)), i)
	}
}

func BenchmarkConcurrentHashMapWithCustomStrategy(b *testing.B) {
	strategy := &CaseInsensitiveHashStrategy{}
	hashMap := maps.NewConcurrentHashMapWithHashStrategy[string, int](strategy)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashMap.Put("key"+string(rune(i%1000)), i)
	}
}