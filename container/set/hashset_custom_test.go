package set

import (
	"testing"
	"github.com/chenjianyu/collections/container/common"
)

func TestHashSetWithDefaultStrategy(t *testing.T) {
	set := New[string]()
	
	// Test basic operations
	if !set.Add("hello") {
		t.Error("Should be able to add 'hello'")
	}
	
	if set.Add("hello") {
		t.Error("Should not be able to add duplicate 'hello'")
	}
	
	if !set.Contains("hello") {
		t.Error("Set should contain 'hello'")
	}
	
	if set.Size() != 1 {
		t.Errorf("Expected size 1, got %d", set.Size())
	}
}

func TestHashSetWithCustomHashStrategy(t *testing.T) {
	// Create a hash set with case-insensitive string strategy
	strategy := common.NewCaseInsensitiveStringHashStrategy()
	set := NewWithHashStrategy[string](strategy)
	
	// Test case-insensitive behavior
	if !set.Add("Hello") {
		t.Error("Should be able to add 'Hello'")
	}
	
	if set.Add("HELLO") {
		t.Error("Should not be able to add 'HELLO' due to case-insensitive strategy")
	}
	
	if set.Add("hello") {
		t.Error("Should not be able to add 'hello' due to case-insensitive strategy")
	}
	
	if !set.Contains("Hello") {
		t.Error("Set should contain 'Hello'")
	}
	
	if !set.Contains("HELLO") {
		t.Error("Set should contain 'HELLO' (case-insensitive)")
	}
	
	if !set.Contains("hello") {
		t.Error("Set should contain 'hello' (case-insensitive)")
	}
	
	if set.Size() != 1 {
		t.Errorf("Expected size 1, got %d", set.Size())
	}
}

func TestHashSetWithStringLengthStrategy(t *testing.T) {
	// Create a hash set with string length-based strategy
	strategy := common.NewStringLengthHashStrategy()
	set := NewWithHashStrategy[string](strategy)
	
	// Add strings of different lengths
	if !set.Add("a") {
		t.Error("Should be able to add 'a'")
	}
	
	if !set.Add("bb") {
		t.Error("Should be able to add 'bb'")
	}
	
	if !set.Add("ccc") {
		t.Error("Should be able to add 'ccc'")
	}
	
	// Test that all strings are present
	if !set.Contains("a") {
		t.Error("Set should contain 'a'")
	}
	
	if !set.Contains("bb") {
		t.Error("Set should contain 'bb'")
	}
	
	if !set.Contains("ccc") {
		t.Error("Set should contain 'ccc'")
	}
	
	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}
}

func TestHashSetWithFunctionalStrategy(t *testing.T) {
	// Create a custom functional strategy that only considers first character
	strategy := common.NewFunctionalHashStrategy[string](
		func(s string) uint64 {
			if len(s) == 0 {
				return 0
			}
			return uint64(s[0])
		},
		func(a, b string) bool {
			if len(a) == 0 && len(b) == 0 {
				return true
			}
			if len(a) == 0 || len(b) == 0 {
				return false
			}
			return a[0] == b[0]
		},
	)
	
	set := NewWithHashStrategy[string](strategy)
	
	// Add strings starting with 'a'
	if !set.Add("apple") {
		t.Error("Should be able to add 'apple'")
	}
	
	if set.Add("apricot") {
		t.Error("Should not be able to add 'apricot' (starts with same letter)")
	}
	
	// Add string starting with different letter
	if !set.Add("banana") {
		t.Error("Should be able to add 'banana'")
	}
	
	// Test contains with first character matching
	if !set.Contains("amazing") {
		t.Error("Set should contain 'amazing' (starts with 'a')")
	}
	
	if !set.Contains("brilliant") {
		t.Error("Set should contain 'brilliant' (starts with 'b')")
	}
	
	if set.Contains("chocolate") {
		t.Error("Set should not contain 'chocolate' (starts with 'c')")
	}
	
	if set.Size() != 2 {
		t.Errorf("Expected size 2, got %d", set.Size())
	}
}

func TestHashSetUnionWithCustomStrategy(t *testing.T) {
	strategy := common.NewCaseInsensitiveStringHashStrategy()
	
	set1 := NewWithHashStrategy[string](strategy)
	set1.Add("Hello")
	set1.Add("World")
	
	set2 := NewWithHashStrategy[string](strategy)
	set2.Add("HELLO") // Should be considered same as "Hello"
	set2.Add("Go")
	
	union := set1.Union(set2).(*HashSet[string])
	
	// Union should contain Hello (case variations considered same), World, and Go
	expectedSize := 3
	if union.Size() != expectedSize {
		t.Errorf("Expected union size %d, got %d", expectedSize, union.Size())
	}
	
	if !union.Contains("hello") {
		t.Error("Union should contain 'hello' (case-insensitive)")
	}
	
	if !union.Contains("WORLD") {
		t.Error("Union should contain 'WORLD' (case-insensitive)")
	}
	
	if !union.Contains("go") {
		t.Error("Union should contain 'go' (case-insensitive)")
	}
}

func TestHashSetIntersectionWithCustomStrategy(t *testing.T) {
	strategy := common.NewCaseInsensitiveStringHashStrategy()
	
	set1 := NewWithHashStrategy[string](strategy)
	set1.Add("Hello")
	set1.Add("World")
	
	set2 := NewWithHashStrategy[string](strategy)
	set2.Add("HELLO") // Should be considered same as "Hello"
	set2.Add("Go")
	
	intersection := set1.Intersection(set2).(*HashSet[string])
	
	// Intersection should contain only Hello (case variations considered same)
	expectedSize := 1
	if intersection.Size() != expectedSize {
		t.Errorf("Expected intersection size %d, got %d", expectedSize, intersection.Size())
	}
	
	if !intersection.Contains("hello") {
		t.Error("Intersection should contain 'hello' (case-insensitive)")
	}
	
	if intersection.Contains("World") {
		t.Error("Intersection should not contain 'World'")
	}
	
	if intersection.Contains("Go") {
		t.Error("Intersection should not contain 'Go'")
	}
}

func TestHashSetFromSliceWithCustomStrategy(t *testing.T) {
	strategy := common.NewCaseInsensitiveStringHashStrategy()
	slice := []string{"Hello", "HELLO", "hello", "World", "WORLD"}
	
	set := FromSliceWithHashStrategy(slice, strategy)
	
	// Due to case-insensitive strategy, should only have 2 unique elements
	expectedSize := 2
	if set.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, set.Size())
	}
	
	if !set.Contains("hello") {
		t.Error("Set should contain 'hello' (case-insensitive)")
	}
	
	if !set.Contains("WORLD") {
		t.Error("Set should contain 'WORLD' (case-insensitive)")
	}
}

// Benchmark tests
func BenchmarkHashSetDefaultStrategy(b *testing.B) {
	set := New[string]()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Add("test")
		set.Contains("test")
		set.Remove("test")
	}
}

func BenchmarkHashSetCaseInsensitiveStrategy(b *testing.B) {
	strategy := common.NewCaseInsensitiveStringHashStrategy()
	set := NewWithHashStrategy[string](strategy)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Add("Test")
		set.Contains("TEST")
		set.Remove("test")
	}
}

func BenchmarkHashSetFunctionalStrategy(b *testing.B) {
	strategy := common.NewFunctionalHashStrategy[string](
		func(s string) uint64 {
			if len(s) == 0 {
				return 0
			}
			return uint64(s[0])
		},
		func(a, b string) bool {
			return a == b
		},
	)
	set := NewWithHashStrategy[string](strategy)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Add("test")
		set.Contains("test")
		set.Remove("test")
	}
}