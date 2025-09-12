// Package common provides strategy interfaces for customizable hashing and comparison
package common

import (
	"hash/fnv"
	"reflect"
	"strings"
)

// HashStrategy defines the interface for custom hash functions
type HashStrategy[T any] interface {
	// Hash computes the hash value for the given element
	Hash(element T) uint64

	// Equals checks if two elements are equal
	// This is important for hash collision resolution
	Equals(a, b T) bool
}

// ComparatorStrategy defines the interface for custom comparison functions
type ComparatorStrategy[T any] interface {
	// Compare compares two elements
	// Returns: negative if a < b, zero if a == b, positive if a > b
	Compare(a, b T) int
}

// DefaultHashStrategy provides a default hash strategy using the existing Hash function
type DefaultHashStrategy[T any] struct{}

// Hash implements HashStrategy.Hash using the common Hash function
func (dhs *DefaultHashStrategy[T]) Hash(element T) uint64 {
	return Hash(element)
}

// Equals implements HashStrategy.Equals using reflect.DeepEqual
func (dhs *DefaultHashStrategy[T]) Equals(a, b T) bool {
	return reflect.DeepEqual(a, b)
}

// NewDefaultHashStrategy creates a new default hash strategy
func NewDefaultHashStrategy[T any]() HashStrategy[T] {
	return &DefaultHashStrategy[T]{}
}

// ComparableHashStrategy provides hash strategy for comparable types
type ComparableHashStrategy[T comparable] struct{}

// Hash implements HashStrategy.Hash for comparable types
func (chs *ComparableHashStrategy[T]) Hash(element T) uint64 {
	return Hash(element)
}

// Equals implements HashStrategy.Equals for comparable types using ==
func (chs *ComparableHashStrategy[T]) Equals(a, b T) bool {
	return a == b
}

// NewComparableHashStrategy creates a new hash strategy for comparable types
func NewComparableHashStrategy[T comparable]() HashStrategy[T] {
	return &ComparableHashStrategy[T]{}
}

// FunctionalHashStrategy allows using function-based hash strategies
type FunctionalHashStrategy[T any] struct {
	hashFunc   func(T) uint64
	equalsFunc func(T, T) bool
}

// Hash implements HashStrategy.Hash using the provided function
func (fhs *FunctionalHashStrategy[T]) Hash(element T) uint64 {
	return fhs.hashFunc(element)
}

// Equals implements HashStrategy.Equals using the provided function
func (fhs *FunctionalHashStrategy[T]) Equals(a, b T) bool {
	return fhs.equalsFunc(a, b)
}

// NewFunctionalHashStrategy creates a new functional hash strategy
func NewFunctionalHashStrategy[T any](hashFunc func(T) uint64, equalsFunc func(T, T) bool) HashStrategy[T] {
	return &FunctionalHashStrategy[T]{
		hashFunc:   hashFunc,
		equalsFunc: equalsFunc,
	}
}

// DefaultComparatorStrategy provides a default comparator using the existing CompareGeneric function
type DefaultComparatorStrategy[T comparable] struct{}

// Compare implements ComparatorStrategy.Compare
func (dcs *DefaultComparatorStrategy[T]) Compare(a, b T) int {
	return CompareGeneric(a, b)
}

// NewDefaultComparatorStrategy creates a new default comparator strategy
func NewDefaultComparatorStrategy[T comparable]() ComparatorStrategy[T] {
	return &DefaultComparatorStrategy[T]{}
}

// FunctionalComparatorStrategy allows using function-based comparison strategies
type FunctionalComparatorStrategy[T any] struct {
	compareFunc func(T, T) int
}

// Compare implements ComparatorStrategy.Compare using the provided function
func (fcs *FunctionalComparatorStrategy[T]) Compare(a, b T) int {
	return fcs.compareFunc(a, b)
}

// NewFunctionalComparatorStrategy creates a new functional comparator strategy
func NewFunctionalComparatorStrategy[T any](compareFunc func(T, T) int) ComparatorStrategy[T] {
	return &FunctionalComparatorStrategy[T]{
		compareFunc: compareFunc,
	}
}

// StringLengthHashStrategy is an example custom hash strategy for strings based on length
type StringLengthHashStrategy struct{}

// Hash computes hash based on string length (example custom strategy)
func (slhs *StringLengthHashStrategy) Hash(element string) uint64 {
	h := fnv.New64a()
	// Hash based on string length and first character
	if len(element) > 0 {
		h.Write([]byte{byte(len(element)), element[0]})
	} else {
		h.Write([]byte{0})
	}
	return h.Sum64()
}

// Equals compares strings normally
func (slhs *StringLengthHashStrategy) Equals(a, b string) bool {
	return a == b
}

// NewStringLengthHashStrategy creates a new string length-based hash strategy
func NewStringLengthHashStrategy() HashStrategy[string] {
	return &StringLengthHashStrategy{}
}

// CaseInsensitiveStringHashStrategy provides case-insensitive hashing for strings
type CaseInsensitiveStringHashStrategy struct{}

// Hash computes hash for lowercase version of string
func (cishs *CaseInsensitiveStringHashStrategy) Hash(element string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(strings.ToLower(element)))
	return h.Sum64()
}

// Equals compares strings case-insensitively
func (cishs *CaseInsensitiveStringHashStrategy) Equals(a, b string) bool {
	return strings.EqualFold(a, b)
}

// NewCaseInsensitiveStringHashStrategy creates a new case-insensitive string hash strategy
func NewCaseInsensitiveStringHashStrategy() HashStrategy[string] {
	return &CaseInsensitiveStringHashStrategy{}
}
