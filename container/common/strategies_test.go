package common

import (
	"testing"
)

func TestDefaultHashStrategy(t *testing.T) {
	strategy := NewDefaultHashStrategy[string]()

	// Test hash consistency
	hash1 := strategy.Hash("test")
	hash2 := strategy.Hash("test")
	if hash1 != hash2 {
		t.Error("Hash should be consistent for same input")
	}

	// Test different inputs produce different hashes (probabilistic)
	hash3 := strategy.Hash("different")
	if hash1 == hash3 {
		t.Log("Hash collision occurred (this is rare but possible)")
	}

	// Test equals
	if !strategy.Equals("test", "test") {
		t.Error("Equal strings should return true")
	}

	if strategy.Equals("test", "different") {
		t.Error("Different strings should return false")
	}
}

func TestComparableHashStrategy(t *testing.T) {
	strategy := NewComparableHashStrategy[int]()

	// Test hash
	hash1 := strategy.Hash(42)
	hash2 := strategy.Hash(42)
	if hash1 != hash2 {
		t.Error("Hash should be consistent for same input")
	}

	// Test equals
	if !strategy.Equals(42, 42) {
		t.Error("Equal integers should return true")
	}

	if strategy.Equals(42, 43) {
		t.Error("Different integers should return false")
	}
}

func TestFunctionalHashStrategy(t *testing.T) {
	// Create a custom hash strategy that hashes only the first character
	strategy := NewFunctionalHashStrategy[string](
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

	// Test custom hash function
	hash1 := strategy.Hash("apple")
	hash2 := strategy.Hash("apricot")
	if hash1 != hash2 {
		t.Error("Both strings start with 'a', should have same hash")
	}

	hash3 := strategy.Hash("banana")
	if hash1 == hash3 {
		t.Error("Different first characters should produce different hashes")
	}

	// Test equals function
	if !strategy.Equals("test", "test") {
		t.Error("Equal strings should return true")
	}
}

func TestDefaultComparatorStrategy(t *testing.T) {
	strategy := NewDefaultComparatorStrategy[int]()

	// Test comparison
	if strategy.Compare(1, 2) >= 0 {
		t.Error("1 should be less than 2")
	}

	if strategy.Compare(2, 1) <= 0 {
		t.Error("2 should be greater than 1")
	}

	if strategy.Compare(1, 1) != 0 {
		t.Error("1 should be equal to 1")
	}
}

func TestFunctionalComparatorStrategy(t *testing.T) {
	// Create a reverse comparator
	strategy := NewFunctionalComparatorStrategy[int](
		func(a, b int) int {
			if a < b {
				return 1 // Reverse: a < b returns positive
			} else if a > b {
				return -1 // Reverse: a > b returns negative
			}
			return 0
		},
	)

	// Test reverse comparison
	if strategy.Compare(1, 2) <= 0 {
		t.Error("With reverse comparator, 1 should be 'greater' than 2")
	}

	if strategy.Compare(2, 1) >= 0 {
		t.Error("With reverse comparator, 2 should be 'less' than 1")
	}

	if strategy.Compare(1, 1) != 0 {
		t.Error("Equal values should still return 0")
	}
}

func TestStringLengthHashStrategy(t *testing.T) {
	strategy := NewStringLengthHashStrategy()

	// Test hash based on length and first character
	hash1 := strategy.Hash("abc")
	hash2 := strategy.Hash("abd") // Same length, different second char

	// Should be different due to different first chars if both are 'a'
	if hash1 == hash2 {
		t.Log("Hash collision occurred (this is possible)")
	}

	// Test equals
	if !strategy.Equals("test", "test") {
		t.Error("Equal strings should return true")
	}

	if strategy.Equals("test", "different") {
		t.Error("Different strings should return false")
	}
}

func TestCaseInsensitiveStringHashStrategy(t *testing.T) {
	strategy := NewCaseInsensitiveStringHashStrategy()

	// Test case-insensitive hash
	hash1 := strategy.Hash("Test")
	hash2 := strategy.Hash("TEST")
	hash3 := strategy.Hash("test")

	if hash1 != hash2 || hash2 != hash3 {
		t.Error("Case variations should produce same hash")
	}

	// Test case-insensitive equals
	if !strategy.Equals("Test", "TEST") {
		t.Error("Case variations should be equal")
	}

	if !strategy.Equals("Hello", "hello") {
		t.Error("Case variations should be equal")
	}

	if strategy.Equals("Hello", "World") {
		t.Error("Different strings should not be equal")
	}
}

// Benchmark tests
func BenchmarkDefaultHashStrategy(b *testing.B) {
	strategy := NewDefaultHashStrategy[string]()
	testString := "benchmark test string"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.Hash(testString)
	}
}

func BenchmarkComparableHashStrategy(b *testing.B) {
	strategy := NewComparableHashStrategy[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.Hash(i)
	}
}

func BenchmarkFunctionalHashStrategy(b *testing.B) {
	strategy := NewFunctionalHashStrategy[string](
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
	testString := "benchmark test string"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.Hash(testString)
	}
}

func BenchmarkStringLengthHashStrategy(b *testing.B) {
	strategy := NewStringLengthHashStrategy()
	testString := "benchmark test string"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.Hash(testString)
	}
}

func BenchmarkCaseInsensitiveStringHashStrategy(b *testing.B) {
	strategy := NewCaseInsensitiveStringHashStrategy()
	testString := "BenchMark Test String"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.Hash(testString)
	}
}
