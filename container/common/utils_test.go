package common

import (
	"testing"
)

type StringerType struct {
	value string
}

func (s StringerType) String() string {
	return s.value
}

func TestEqual(t *testing.T) {
	// Test basic types
	if !Equal(42, 42) {
		t.Error("Equal integers should return true")
	}

	if Equal(42, 43) {
		t.Error("Different integers should return false")
	}

	// Test strings
	if !Equal("hello", "hello") {
		t.Error("Equal strings should return true")
	}

	if Equal("hello", "world") {
		t.Error("Different strings should return false")
	}

	// Test slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	slice3 := []int{1, 2, 4}

	if !Equal(slice1, slice2) {
		t.Error("Equal slices should return true")
	}

	if Equal(slice1, slice3) {
		t.Error("Different slices should return false")
	}

	// Test structs
	type Person struct {
		Name string
		Age  int
	}

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Alice", Age: 30}
	person3 := Person{Name: "Bob", Age: 25}

	if !Equal(person1, person2) {
		t.Error("Equal structs should return true")
	}

	if Equal(person1, person3) {
		t.Error("Different structs should return false")
	}
}

func TestHash(t *testing.T) {
	// Test basic types
	if Hash(42) == 0 {
		t.Error("Hash should not return 0 for non-zero value")
	}

	if Hash(0) == 0 {
		t.Error("Hash should not return 0 for zero value")
	}

	// Test strings
	if Hash("hello") == 0 {
		t.Error("Hash should not return 0 for non-empty string")
	}

	if Hash("") == 0 {
		t.Error("Hash should not return 0 for empty string")
	}

	// Test slices
	slice := []int{1, 2, 3}
	if Hash(slice) == 0 {
		t.Error("Hash should not return 0 for non-empty slice")
	}

	emptySlice := []int{}
	if Hash(emptySlice) == 0 {
		t.Error("Hash should not return 0 for empty slice")
	}

	// Test structs
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "Alice", Age: 30}
	if Hash(person) == 0 {
		t.Error("Hash should not return 0 for struct")
	}

	// Test that same values produce same hash
	if Hash(42) != Hash(42) {
		t.Error("Same values should produce same hash")
	}

	if Hash("hello") != Hash("hello") {
		t.Error("Same strings should produce same hash")
	}

	// Test that different values produce different hash (probably)
	if Hash(42) == Hash(43) {
		t.Log("Different values produced same hash (hash collision)")
	}

	if Hash("hello") == Hash("world") {
		t.Log("Different strings produced same hash (hash collision)")
	}

	// Test complex types
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	slice3 := []int{3, 2, 1}

	if Hash(slice1) != Hash(slice2) {
		t.Error("Hash should return same value for equal slices")
	}

	if Hash(slice1) == Hash(slice3) {
		t.Error("Hash should return different values for different slices (probabilistic)")
	}
}

func TestCompare(t *testing.T) {
	// Test integer comparison
	if CompareGeneric(1, 2) >= 0 {
		t.Error("CompareGeneric(1, 2) should return negative value")
	}

	if CompareGeneric(2, 1) <= 0 {
		t.Error("CompareGeneric(2, 1) should return positive value")
	}

	if CompareGeneric(1, 1) != 0 {
		t.Error("CompareGeneric(1, 1) should return 0")
	}

	// Test string comparison
	if CompareGeneric("a", "b") >= 0 {
		t.Error("CompareGeneric(\"a\", \"b\") should return negative value")
	}

	if CompareGeneric("b", "a") <= 0 {
		t.Error("CompareGeneric(\"b\", \"a\") should return positive value")
	}

	if CompareGeneric("a", "a") != 0 {
		t.Error("CompareGeneric(\"a\", \"a\") should return 0")
	}

	// Test float comparison
	if CompareGeneric(1.5, 2.5) >= 0 {
		t.Error("CompareGeneric(1.5, 2.5) should return negative value")
	}

	if CompareGeneric(2.5, 1.5) <= 0 {
		t.Error("CompareGeneric(2.5, 1.5) should return positive value")
	}

	if CompareGeneric(1.5, 1.5) != 0 {
		t.Error("CompareGeneric(1.5, 1.5) should return 0")
	}

	// Test uint comparison
	if CompareGeneric(uint(1), uint(2)) >= 0 {
		t.Error("CompareGeneric(uint(1), uint(2)) should return negative value")
	}

	if CompareGeneric(uint8(2), uint8(1)) <= 0 {
		t.Error("CompareGeneric(uint8(2), uint8(1)) should return positive value")
	}

	if CompareGeneric(uint16(1), uint16(1)) != 0 {
		t.Error("CompareGeneric(uint16(1), uint16(1)) should return 0")
	}

	if CompareGeneric(uint32(1), uint32(2)) >= 0 {
		t.Error("CompareGeneric(uint32(1), uint32(2)) should return negative value")
	}

	if CompareGeneric(uint64(2), uint64(1)) <= 0 {
		t.Error("CompareGeneric(uint64(2), uint64(1)) should return positive value")
	}

	// Test float32
	if CompareGeneric(float32(1.5), float32(2.5)) >= 0 {
		t.Error("1.5 should be less than 2.5")
	}

	s1 := StringerType{value: "a"}
	s2 := StringerType{value: "b"}

	// Use non-generic Compare for non-comparable types
	if Compare(s1, s2) >= 0 {
		t.Error("StringerType a should be less than b")
	}

	// Test non-comparable types (using hash values for comparison)
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}

	// For non-comparable types, use the non-generic Compare function
	// Same values should return 0
	slice3 := []int{1, 2, 3}
	if Compare(slice1, slice3) != 0 {
		t.Error("Same slices should return 0")
	}

	// Different values may return different results (based on hash values)
	result := Compare(slice1, slice2)
	// This case is rare, but hash collisions may occur
	if result == 0 {
		t.Log("Hash collision occurred for different slices")
	}
}

// Test NewPair
func TestCompareWithNonComparableTypes(t *testing.T) {
	// Test non-comparable types (using hash values for comparison)
	type CustomType struct {
		value int
	}

	c1 := CustomType{value: 1}
	c2 := CustomType{value: 2}
	c3 := CustomType{value: 1}

	// For non-comparable types, CompareGeneric uses hash values for comparison
	// Same values should return 0
	if CompareGeneric(c1, c3) != 0 {
		t.Error("CompareGeneric should return 0 for equal custom types")
	}

	// Different values may return different results (based on hash values)
	result := CompareGeneric(c1, c2)
	if result == 0 && c1.value != c2.value {
		// This case is rare, but hash collisions may occur
		t.Logf("Hash collision detected for different values")
	}
}

func TestPair(t *testing.T) {
	// Test NewPair
	pair := NewPair("key", 42)

	if pair.Key != "key" {
		t.Errorf("Expected key 'key', got '%v'", pair.Key)
	}

	if pair.Value != 42 {
		t.Errorf("Expected value 42, got %v", pair.Value)
	}

	// Test String method
	expected := "(key, 42)"
	if pair.String() != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, pair.String())
	}

	// Test different types of Pair
	pairFloat := NewPair(1.5, "value")
	expectedFloat := "(1.5, value)"
	if pairFloat.String() != expectedFloat {
		t.Errorf("Expected string '%s', got '%s'", expectedFloat, pairFloat.String())
	}

	// Test complex types of Pair
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "Alice", Age: 30}
	pairPerson := NewPair("person", person)

	if pairPerson.Key != "person" {
		t.Errorf("Expected key 'person', got '%v'", pairPerson.Key)
	}

	if pairPerson.Value != person {
		t.Errorf("Expected value %v, got %v", person, pairPerson.Value)
	}
}

func TestPairEqual(t *testing.T) {
	// Test Pair equality
	pair1 := NewPair("key", 42)
	pair2 := NewPair("key", 42)
	pair3 := NewPair("key", 43)
	pair4 := NewPair("other", 42)

	if !Equal(pair1, pair2) {
		t.Error("Equal pairs should be equal")
	}

	if Equal(pair1, pair3) {
		t.Error("Pairs with different values should not be equal")
	}

	if Equal(pair1, pair4) {
		t.Error("Pairs with different keys should not be equal")
	}
}

func TestHashConsistency(t *testing.T) {
	// Test hash value consistency
	values := []interface{}{
		1, 2, 3,
		"hello", "world",
		1.5, 2.5,
		[]int{1, 2, 3},
		map[string]int{"a": 1, "b": 2},
	}

	for _, v := range values {
		hash1 := Hash(v)
		hash2 := Hash(v)

		if hash1 != hash2 {
			t.Errorf("Hash should be consistent for value %v", v)
		}
	}
}

func BenchmarkHash(b *testing.B) {
	values := []interface{}{
		42,
		"hello world",
		[]int{1, 2, 3, 4, 5},
		map[string]int{"a": 1, "b": 2, "c": 3},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			Hash(v)
		}
	}
}

func BenchmarkCompare(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compare(i, i+1)
		Compare("hello", "world")
		Compare(1.5, 2.5)
	}
}

func BenchmarkEqual(b *testing.B) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{1, 2, 3, 4, 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(slice1, slice2)
		Equal("hello", "hello")
		Equal(42, 42)
	}
}
