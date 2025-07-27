package common

import (
	"testing"
)

func TestEqual(t *testing.T) {
	// 测试基本类型
	if !Equal(1, 1) {
		t.Error("Equal(1, 1) should return true")
	}

	if Equal(1, 2) {
		t.Error("Equal(1, 2) should return false")
	}

	// 测试字符串
	if !Equal("hello", "hello") {
		t.Error("Equal(\"hello\", \"hello\") should return true")
	}

	if Equal("hello", "world") {
		t.Error("Equal(\"hello\", \"world\") should return false")
	}

	// 测试切片
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	slice3 := []int{1, 2, 4}

	if !Equal(slice1, slice2) {
		t.Error("Equal should return true for equal slices")
	}

	if Equal(slice1, slice3) {
		t.Error("Equal should return false for different slices")
	}

	// 测试结构体
	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	p3 := Person{Name: "Bob", Age: 25}

	if !Equal(p1, p2) {
		t.Error("Equal should return true for equal structs")
	}

	if Equal(p1, p3) {
		t.Error("Equal should return false for different structs")
	}
}

func TestHash(t *testing.T) {
	// 测试相同值产生相同哈希
	if Hash(42) != Hash(42) {
		t.Error("Hash should return same value for same input")
	}

	if Hash("hello") != Hash("hello") {
		t.Error("Hash should return same value for same string")
	}

	// 测试不同值产生不同哈希（大概率）
	if Hash(1) == Hash(2) {
		t.Error("Hash should return different values for different inputs (probabilistic)")
	}

	if Hash("hello") == Hash("world") {
		t.Error("Hash should return different values for different strings (probabilistic)")
	}

	// 测试复杂类型
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
	// 测试整数比较
	if Compare(1, 2) >= 0 {
		t.Error("Compare(1, 2) should return negative value")
	}

	if Compare(2, 1) <= 0 {
		t.Error("Compare(2, 1) should return positive value")
	}

	if Compare(1, 1) != 0 {
		t.Error("Compare(1, 1) should return 0")
	}

	// 测试字符串比较
	if Compare("apple", "banana") >= 0 {
		t.Error("Compare(\"apple\", \"banana\") should return negative value")
	}

	if Compare("banana", "apple") <= 0 {
		t.Error("Compare(\"banana\", \"apple\") should return positive value")
	}

	if Compare("apple", "apple") != 0 {
		t.Error("Compare(\"apple\", \"apple\") should return 0")
	}

	// 测试浮点数比较
	if Compare(1.5, 2.5) >= 0 {
		t.Error("Compare(1.5, 2.5) should return negative value")
	}

	if Compare(2.5, 1.5) <= 0 {
		t.Error("Compare(2.5, 1.5) should return positive value")
	}

	if Compare(1.5, 1.5) != 0 {
		t.Error("Compare(1.5, 1.5) should return 0")
	}

	// 测试不同类型的整数
	if Compare(int8(1), int8(2)) >= 0 {
		t.Error("Compare(int8(1), int8(2)) should return negative value")
	}

	if Compare(int16(2), int16(1)) <= 0 {
		t.Error("Compare(int16(2), int16(1)) should return positive value")
	}

	if Compare(int32(1), int32(1)) != 0 {
		t.Error("Compare(int32(1), int32(1)) should return 0")
	}

	if Compare(int64(1), int64(2)) >= 0 {
		t.Error("Compare(int64(1), int64(2)) should return negative value")
	}

	// 测试无符号整数
	if Compare(uint(1), uint(2)) >= 0 {
		t.Error("Compare(uint(1), uint(2)) should return negative value")
	}

	if Compare(uint8(2), uint8(1)) <= 0 {
		t.Error("Compare(uint8(2), uint8(1)) should return positive value")
	}

	if Compare(uint16(1), uint16(1)) != 0 {
		t.Error("Compare(uint16(1), uint16(1)) should return 0")
	}

	if Compare(uint32(1), uint32(2)) >= 0 {
		t.Error("Compare(uint32(1), uint32(2)) should return negative value")
	}

	if Compare(uint64(2), uint64(1)) <= 0 {
		t.Error("Compare(uint64(2), uint64(1)) should return positive value")
	}

	// 测试float32
	if Compare(float32(1.5), float32(2.5)) >= 0 {
		t.Error("Compare(float32(1.5), float32(2.5)) should return negative value")
	}
}

func TestCompareWithStringer(t *testing.T) {
	// 测试实现了Stringer接口的类型
	type StringerType struct {
		value string
	}

	s1 := StringerType{value: "apple"}
	s2 := StringerType{value: "banana"}
	s3 := StringerType{value: "apple"}

	if Compare(s1, s2) >= 0 {
		t.Error("Compare should return negative value for StringerType")
	}

	if Compare(s2, s1) <= 0 {
		t.Error("Compare should return positive value for StringerType")
	}

	if Compare(s1, s3) != 0 {
		t.Error("Compare should return 0 for equal StringerType")
	}
}

func TestCompareWithNonComparableTypes(t *testing.T) {
	// 测试不可直接比较的类型（使用哈希值比较）
	type CustomType struct {
		value int
	}

	c1 := CustomType{value: 1}
	c2 := CustomType{value: 2}
	c3 := CustomType{value: 1}

	// 对于不可比较的类型，Compare使用哈希值比较
	// 相同的值应该返回0
	if Compare(c1, c3) != 0 {
		t.Error("Compare should return 0 for equal custom types")
	}

	// 不同的值可能返回不同的结果（基于哈希值）
	result := Compare(c1, c2)
	if result == 0 && c1.value != c2.value {
		// 这种情况很少见，但可能发生哈希冲突
		t.Logf("Hash collision detected for different values")
	}
}

func TestPair(t *testing.T) {
	// 测试NewPair
	pair := NewPair("key", 42)

	if pair.Key != "key" {
		t.Errorf("Expected key 'key', got '%v'", pair.Key)
	}

	if pair.Value != 42 {
		t.Errorf("Expected value 42, got %v", pair.Value)
	}

	// 测试String方法
	expected := "key=42"
	if pair.String() != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, pair.String())
	}

	// 测试不同类型的Pair
	pairFloat := NewPair(1.5, "value")
	expectedFloat := "1.5=value"
	if pairFloat.String() != expectedFloat {
		t.Errorf("Expected string '%s', got '%s'", expectedFloat, pairFloat.String())
	}

	// 测试复杂类型的Pair
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
	// 测试Pair的相等性
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
	// 测试哈希值的一致性
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
