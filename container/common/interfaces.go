// Package common provides common interfaces for the container library
package common

// Container represents a collection of elements
type Container[E any] interface {
	// Size returns the number of elements in the container
	Size() int
	// IsEmpty returns true if the container is empty
	IsEmpty() bool
	// Clear removes all elements from the container
	Clear()
	// Contains returns true if the container contains the specified element
	Contains(element E) bool
	// String returns the string representation of the container
	String() string
}

// Iterable represents a collection that can be iterated
type Iterable[E any] interface {
	// Iterator returns an iterator for the collection
	Iterator() Iterator[E]
	// ForEach executes the given function for each element
	ForEach(func(E))
}

// Iterator represents an iterator for traversing elements
type Iterator[E any] interface {
	// HasNext returns true if there are more elements to iterate
	HasNext() bool
	// Next returns the next element
	Next() (E, bool)
	// Remove removes the current element (optional operation)
	Remove() bool
}

// Comparable represents a type that can be compared
type Comparable interface {
	// CompareTo compares this object with another object
	// Returns:
	//   - negative number: this < other
	//   - zero: this == other
	//   - positive number: this > other
	CompareTo(other interface{}) int
}
