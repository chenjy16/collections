package common

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestCommonErrors tests the common error variables
func TestCommonErrors(t *testing.T) {
	t.Run("ErrorVariables", func(t *testing.T) {
		assert.NotNil(t, ErrIndexOutOfBounds)
		assert.NotNil(t, ErrEmptyContainer)
		assert.NotNil(t, ErrFullContainer)
		assert.NotNil(t, ErrInvalidRange)
		assert.NotNil(t, ErrNegativeCount)
		assert.NotNil(t, ErrImmutableOperation)
		assert.NotNil(t, ErrKeyNotFound)
		assert.NotNil(t, ErrElementNotFound)
		assert.NotNil(t, ErrDuplicateKey)
		assert.NotNil(t, ErrInvalidArgument)
		assert.NotNil(t, ErrConcurrentAccess)
	})

	t.Run("ErrorMessages", func(t *testing.T) {
		assert.Equal(t, "index out of bounds", ErrIndexOutOfBounds.Error())
		assert.Equal(t, "container is empty", ErrEmptyContainer.Error())
		assert.Equal(t, "container is full", ErrFullContainer.Error())
		assert.Equal(t, "invalid range", ErrInvalidRange.Error())
		assert.Equal(t, "count cannot be negative", ErrNegativeCount.Error())
		assert.Equal(t, "operation not allowed on immutable collection", ErrImmutableOperation.Error())
		assert.Equal(t, "key not found", ErrKeyNotFound.Error())
		assert.Equal(t, "element not found", ErrElementNotFound.Error())
		assert.Equal(t, "duplicate key", ErrDuplicateKey.Error())
		assert.Equal(t, "invalid argument", ErrInvalidArgument.Error())
		assert.Equal(t, "concurrent access violation", ErrConcurrentAccess.Error())
	})
}

// TestIndexOutOfBoundsError tests the IndexOutOfBoundsError function
func TestIndexOutOfBoundsError(t *testing.T) {
	t.Run("ValidIndexAndSize", func(t *testing.T) {
		err := IndexOutOfBoundsError(5, 3)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
		assert.Contains(t, err.Error(), "index 5")
		assert.Contains(t, err.Error(), "size 3")
	})

	t.Run("NegativeIndex", func(t *testing.T) {
		err := IndexOutOfBoundsError(-1, 10)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
		assert.Contains(t, err.Error(), "index -1")
	})

	t.Run("ZeroSize", func(t *testing.T) {
		err := IndexOutOfBoundsError(0, 0)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
		assert.Contains(t, err.Error(), "size 0")
	})
}

// TestInvalidRangeError tests the InvalidRangeError function
func TestInvalidRangeError(t *testing.T) {
	t.Run("StartGreaterThanEnd", func(t *testing.T) {
		err := InvalidRangeError(5, 3)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidRange))
		assert.Contains(t, err.Error(), "start 5")
		assert.Contains(t, err.Error(), "end 3")
	})

	t.Run("NegativeRange", func(t *testing.T) {
		err := InvalidRangeError(-5, -1)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidRange))
		assert.Contains(t, err.Error(), "start -5")
		assert.Contains(t, err.Error(), "end -1")
	})

	t.Run("ValidRange", func(t *testing.T) {
		err := InvalidRangeError(1, 5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidRange))
		assert.Contains(t, err.Error(), "start 1")
		assert.Contains(t, err.Error(), "end 5")
	})
}

// TestNegativeCountError tests the NegativeCountError function
func TestNegativeCountError(t *testing.T) {
	t.Run("NegativeCount", func(t *testing.T) {
		err := NegativeCountError(-5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNegativeCount))
		assert.Contains(t, err.Error(), "-5")
	})

	t.Run("ZeroCount", func(t *testing.T) {
		err := NegativeCountError(0)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNegativeCount))
		assert.Contains(t, err.Error(), "0")
	})
}

// TestImmutableOperationError tests the ImmutableOperationError function
func TestImmutableOperationError(t *testing.T) {
	t.Run("AddOperation", func(t *testing.T) {
		err := ImmutableOperationError("Add", "WithAdd()")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrImmutableOperation))
		assert.Contains(t, err.Error(), "Add")
		assert.Contains(t, err.Error(), "WithAdd()")
	})

	t.Run("RemoveOperation", func(t *testing.T) {
		err := ImmutableOperationError("Remove", "WithRemove()")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrImmutableOperation))
		assert.Contains(t, err.Error(), "Remove")
		assert.Contains(t, err.Error(), "WithRemove()")
	})

	t.Run("EmptyStrings", func(t *testing.T) {
		err := ImmutableOperationError("", "")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrImmutableOperation))
	})
}

// TestKeyNotFoundError tests the KeyNotFoundError function
func TestKeyNotFoundError(t *testing.T) {
	t.Run("StringKey", func(t *testing.T) {
		err := KeyNotFoundError("testKey")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrKeyNotFound))
		assert.Contains(t, err.Error(), "testKey")
	})

	t.Run("IntKey", func(t *testing.T) {
		err := KeyNotFoundError(42)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrKeyNotFound))
		assert.Contains(t, err.Error(), "42")
	})

	t.Run("NilKey", func(t *testing.T) {
		err := KeyNotFoundError(nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrKeyNotFound))
		assert.Contains(t, err.Error(), "<nil>")
	})

	t.Run("StructKey", func(t *testing.T) {
		type TestStruct struct {
			ID   int
			Name string
		}
		key := TestStruct{ID: 1, Name: "test"}
		err := KeyNotFoundError(key)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrKeyNotFound))
	})
}

// TestElementNotFoundError tests the ElementNotFoundError function
func TestElementNotFoundError(t *testing.T) {
	t.Run("StringElement", func(t *testing.T) {
		err := ElementNotFoundError("testElement")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrElementNotFound))
		assert.Contains(t, err.Error(), "testElement")
	})

	t.Run("IntElement", func(t *testing.T) {
		err := ElementNotFoundError(123)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrElementNotFound))
		assert.Contains(t, err.Error(), "123")
	})

	t.Run("NilElement", func(t *testing.T) {
		err := ElementNotFoundError(nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrElementNotFound))
		assert.Contains(t, err.Error(), "<nil>")
	})
}

// TestDuplicateKeyError tests the DuplicateKeyError function
func TestDuplicateKeyError(t *testing.T) {
	t.Run("StringKey", func(t *testing.T) {
		err := DuplicateKeyError("duplicateKey")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrDuplicateKey))
		assert.Contains(t, err.Error(), "duplicateKey")
	})

	t.Run("IntKey", func(t *testing.T) {
		err := DuplicateKeyError(999)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrDuplicateKey))
		assert.Contains(t, err.Error(), "999")
	})
}

// TestInvalidArgumentError tests the InvalidArgumentError function
func TestInvalidArgumentError(t *testing.T) {
	t.Run("ValidArguments", func(t *testing.T) {
		err := InvalidArgumentError("capacity", "must be positive")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidArgument))
		assert.Contains(t, err.Error(), "capacity")
		assert.Contains(t, err.Error(), "must be positive")
	})

	t.Run("EmptyArgument", func(t *testing.T) {
		err := InvalidArgumentError("", "cannot be empty")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidArgument))
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("EmptyReason", func(t *testing.T) {
		err := InvalidArgumentError("value", "")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidArgument))
		assert.Contains(t, err.Error(), "value")
	})
}

// TestConcurrentAccessError tests the ConcurrentAccessError function
func TestConcurrentAccessError(t *testing.T) {
	t.Run("ReadOperation", func(t *testing.T) {
		err := ConcurrentAccessError("read")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrConcurrentAccess))
		assert.Contains(t, err.Error(), "read")
	})

	t.Run("WriteOperation", func(t *testing.T) {
		err := ConcurrentAccessError("write")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrConcurrentAccess))
		assert.Contains(t, err.Error(), "write")
	})

	t.Run("EmptyOperation", func(t *testing.T) {
		err := ConcurrentAccessError("")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrConcurrentAccess))
	})
}

// TestEmptyContainerError tests the EmptyContainerError function
func TestEmptyContainerError(t *testing.T) {
	t.Run("ListContainer", func(t *testing.T) {
		err := EmptyContainerError("ArrayList")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrEmptyContainer))
		assert.Contains(t, err.Error(), "ArrayList")
	})

	t.Run("QueueContainer", func(t *testing.T) {
		err := EmptyContainerError("Queue")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrEmptyContainer))
		assert.Contains(t, err.Error(), "Queue")
	})

	t.Run("EmptyContainerType", func(t *testing.T) {
		err := EmptyContainerError("")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrEmptyContainer))
	})
}

// TestFullContainerError tests the FullContainerError function
func TestFullContainerError(t *testing.T) {
	t.Run("FixedSizeContainer", func(t *testing.T) {
		err := FullContainerError("FixedSizeQueue", 10)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrFullContainer))
		assert.Contains(t, err.Error(), "FixedSizeQueue")
		assert.Contains(t, err.Error(), "10")
	})

	t.Run("ZeroCapacity", func(t *testing.T) {
		err := FullContainerError("EmptyContainer", 0)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrFullContainer))
		assert.Contains(t, err.Error(), "EmptyContainer")
		assert.Contains(t, err.Error(), "0")
	})

	t.Run("LargeCapacity", func(t *testing.T) {
		err := FullContainerError("LargeContainer", 1000000)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrFullContainer))
		assert.Contains(t, err.Error(), "LargeContainer")
		assert.Contains(t, err.Error(), "1000000")
	})

	t.Run("EmptyContainerType", func(t *testing.T) {
		err := FullContainerError("", 5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrFullContainer))
		assert.Contains(t, err.Error(), "5")
	})
}

// TestErrorWrapping tests error wrapping functionality
func TestErrorWrapping(t *testing.T) {
	t.Run("IndexOutOfBoundsWrapping", func(t *testing.T) {
		err := IndexOutOfBoundsError(10, 5)
		assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
		
		// Test unwrapping
		unwrapped := errors.Unwrap(err)
		assert.Equal(t, ErrIndexOutOfBounds, unwrapped)
	})

	t.Run("InvalidRangeWrapping", func(t *testing.T) {
		err := InvalidRangeError(10, 5)
		assert.True(t, errors.Is(err, ErrInvalidRange))
		
		unwrapped := errors.Unwrap(err)
		assert.Equal(t, ErrInvalidRange, unwrapped)
	})

	t.Run("ImmutableOperationWrapping", func(t *testing.T) {
		err := ImmutableOperationError("Add", "WithAdd")
		assert.True(t, errors.Is(err, ErrImmutableOperation))
		
		unwrapped := errors.Unwrap(err)
		assert.Equal(t, ErrImmutableOperation, unwrapped)
	})
}

// TestErrorComparison tests error comparison functionality
func TestErrorComparison(t *testing.T) {
	t.Run("SameErrorTypes", func(t *testing.T) {
		err1 := IndexOutOfBoundsError(5, 10)
		err2 := IndexOutOfBoundsError(3, 8)
		
		assert.True(t, errors.Is(err1, ErrIndexOutOfBounds))
		assert.True(t, errors.Is(err2, ErrIndexOutOfBounds))
		assert.NotEqual(t, err1.Error(), err2.Error()) // Different messages
	})

	t.Run("DifferentErrorTypes", func(t *testing.T) {
		err1 := IndexOutOfBoundsError(5, 10)
		err2 := InvalidRangeError(5, 3)
		
		assert.True(t, errors.Is(err1, ErrIndexOutOfBounds))
		assert.True(t, errors.Is(err2, ErrInvalidRange))
		assert.False(t, errors.Is(err1, ErrInvalidRange))
		assert.False(t, errors.Is(err2, ErrIndexOutOfBounds))
	})
}

// TestErrorStringRepresentation tests error string representations
func TestErrorStringRepresentation(t *testing.T) {
	t.Run("ComplexErrorMessages", func(t *testing.T) {
		err1 := IndexOutOfBoundsError(100, 50)
		assert.Contains(t, err1.Error(), "index out of bounds")
		assert.Contains(t, err1.Error(), "index 100")
		assert.Contains(t, err1.Error(), "size 50")

		err2 := FullContainerError("PriorityQueue", 1000)
		assert.Contains(t, err2.Error(), "container is full")
		assert.Contains(t, err2.Error(), "PriorityQueue")
		assert.Contains(t, err2.Error(), "capacity: 1000")

		err3 := ImmutableOperationError("Clear", "WithClear()")
		assert.Contains(t, err3.Error(), "operation not allowed on immutable collection")
		assert.Contains(t, err3.Error(), "Clear")
		assert.Contains(t, err3.Error(), "use WithClear() instead")
	})
}