// Package common provides common error types for the container library
package common

import (
	"errors"
	"fmt"
)

// Common error variables
var (
	ErrIndexOutOfBounds   = errors.New("index out of bounds")
	ErrEmptyContainer     = errors.New("container is empty")
	ErrFullContainer      = errors.New("container is full")
	ErrInvalidRange       = errors.New("invalid range")
	ErrNegativeCount      = errors.New("count cannot be negative")
	ErrImmutableOperation = errors.New("operation not allowed on immutable collection")
	ErrKeyNotFound        = errors.New("key not found")
	ErrElementNotFound    = errors.New("element not found")
	ErrDuplicateKey       = errors.New("duplicate key")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrConcurrentAccess   = errors.New("concurrent access violation")
)

// Error factory functions for specific error scenarios

// IndexOutOfBoundsError creates a specific index out of bounds error
func IndexOutOfBoundsError(index, size int) error {
	return fmt.Errorf("%w: index %d, size %d", ErrIndexOutOfBounds, index, size)
}

// InvalidRangeError creates a specific invalid range error
func InvalidRangeError(start, end int) error {
	return fmt.Errorf("%w: start %d, end %d", ErrInvalidRange, start, end)
}

// NegativeCountError creates a specific negative count error
func NegativeCountError(count int) error {
	return fmt.Errorf("%w: %d", ErrNegativeCount, count)
}

// ImmutableOperationError creates a specific immutable operation error
func ImmutableOperationError(operation, suggestion string) error {
	return fmt.Errorf("%w: %s - use %s instead", ErrImmutableOperation, operation, suggestion)
}

// KeyNotFoundError creates a specific key not found error
func KeyNotFoundError(key interface{}) error {
	return fmt.Errorf("%w: %v", ErrKeyNotFound, key)
}

// ElementNotFoundError creates a specific element not found error
func ElementNotFoundError(element interface{}) error {
	return fmt.Errorf("%w: %v", ErrElementNotFound, element)
}

// DuplicateKeyError creates a specific duplicate key error
func DuplicateKeyError(key interface{}) error {
	return fmt.Errorf("%w: %v", ErrDuplicateKey, key)
}

// InvalidArgumentError creates a specific invalid argument error
func InvalidArgumentError(argument string, reason string) error {
	return fmt.Errorf("%w: %s - %s", ErrInvalidArgument, argument, reason)
}

// ConcurrentAccessError creates a specific concurrent access error
func ConcurrentAccessError(operation string) error {
	return fmt.Errorf("%w: %s", ErrConcurrentAccess, operation)
}

// Container-specific error factory functions

// EmptyContainerError creates a specific empty container error with container type
func EmptyContainerError(containerType string) error {
	return fmt.Errorf("%w: %s", ErrEmptyContainer, containerType)
}

// FullContainerError creates a specific full container error with container type and capacity
func FullContainerError(containerType string, capacity int) error {
	return fmt.Errorf("%w: %s (capacity: %d)", ErrFullContainer, containerType, capacity)
}