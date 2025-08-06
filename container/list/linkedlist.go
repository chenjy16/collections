// Package list provides implementations of list data structures
package list

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// Node is a node in the doubly linked list
type Node[E any] struct {
	data E
	prev *Node[E]
	next *Node[E]
}

// LinkedList is a List implementation based on doubly linked list
type LinkedList[E any] struct {
	head *Node[E]
	tail *Node[E]
	size int
}

// NewLinkedList creates a new LinkedList
func NewLinkedList[E any]() *LinkedList[E] {
	return &LinkedList[E]{}
}

// LinkedListFromSlice creates a new LinkedList from a slice
func LinkedListFromSlice[E any](slice []E) *LinkedList[E] {
	list := NewLinkedList[E]()
	for _, element := range slice {
		list.Add(element)
	}
	return list
}

// Size returns the number of elements in the list
func (list *LinkedList[E]) Size() int {
	return list.size
}

// IsEmpty checks if the list is empty
func (list *LinkedList[E]) IsEmpty() bool {
	return list.size == 0
}

// Clear empties the list
func (list *LinkedList[E]) Clear() {
	list.head = nil
	list.tail = nil
	list.size = 0
}

// Contains checks if the list contains the specified element
func (list *LinkedList[E]) Contains(element E) bool {
	return list.IndexOf(element) != -1
}

// ForEach executes the given operation on each element in the list
func (list *LinkedList[E]) ForEach(f func(E)) {
	current := list.head
	for current != nil {
		f(current.data)
		current = current.next
	}
}

// String returns the string representation of the list
func (list *LinkedList[E]) String() string {
	if list.IsEmpty() {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	current := list.head
	first := true
	for current != nil {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v", current.data))
		first = false
		current = current.next
	}
	builder.WriteString("]")
	return builder.String()
}

// Add adds an element to the end of the list
func (list *LinkedList[E]) Add(element E) bool {
	list.AddLast(element)
	return true
}

// Insert inserts an element at the specified position
func (list *LinkedList[E]) Insert(index int, element E) error {
	if index < 0 || index > list.size {
		return common.IndexOutOfBoundsError(index, list.size)
	}

	if index == 0 {
		list.AddFirst(element)
		return nil
	}

	if index == list.size {
		list.AddLast(element)
		return nil
	}

	node := list.getNodeAt(index)
	newNode := &Node[E]{data: element}

	newNode.prev = node.prev
	newNode.next = node
	node.prev.next = newNode
	node.prev = newNode

	list.size++
	return nil
}

// Get retrieves the element at the specified index
func (list *LinkedList[E]) Get(index int) (E, error) {
	if index < 0 || index >= list.size {
		return *new(E), common.IndexOutOfBoundsError(index, list.size)
	}

	node := list.getNodeAt(index)
	return node.data, nil
}

// Set replaces the element at the specified index
func (list *LinkedList[E]) Set(index int, element E) (E, bool) {
	if index < 0 || index >= list.size {
		return *new(E), false
	}

	node := list.getNodeAt(index)
	oldData := node.data
	node.data = element
	return oldData, true
}

// RemoveAt removes the element at the specified index
func (list *LinkedList[E]) RemoveAt(index int) (E, bool) {
	if index < 0 || index >= list.size {
		return *new(E), false
	}

	if index == 0 {
		return list.RemoveFirst()
	}

	if index == list.size-1 {
		return list.RemoveLast()
	}

	node := list.getNodeAt(index)
	data := node.data

	node.prev.next = node.next
	node.next.prev = node.prev

	list.size--
	return data, true
}

// Remove removes the first occurrence of the specified element
func (list *LinkedList[E]) Remove(element E) bool {
	current := list.head
	for current != nil {
		if common.Equal(current.data, element) {
			if current.prev != nil {
				current.prev.next = current.next
			} else {
				list.head = current.next
			}

			if current.next != nil {
				current.next.prev = current.prev
			} else {
				list.tail = current.prev
			}

			list.size--
			return true
		}
		current = current.next
	}
	return false
}

// IndexOf returns the index of the first occurrence of the specified element in the list
func (list *LinkedList[E]) IndexOf(element E) int {
	current := list.head
	index := 0
	for current != nil {
		if common.Equal(current.data, element) {
			return index
		}
		current = current.next
		index++
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of the specified element in the list
func (list *LinkedList[E]) LastIndexOf(element E) int {
	current := list.tail
	index := list.size - 1
	for current != nil {
		if common.Equal(current.data, element) {
			return index
		}
		current = current.prev
		index--
	}
	return -1
}

// SubList returns a view of the specified range in the list
func (list *LinkedList[E]) SubList(fromIndex, toIndex int) (List[E], error) {
	if fromIndex < 0 || toIndex > list.size || fromIndex > toIndex {
		return nil, common.InvalidRangeError(fromIndex, toIndex)
	}

	subList := NewLinkedList[E]()
	current := list.getNodeAt(fromIndex)
	for i := fromIndex; i < toIndex; i++ {
		subList.Add(current.data)
		current = current.next
	}
	return subList, nil
}

// ToSlice returns a slice containing all elements in the list
func (list *LinkedList[E]) ToSlice() []E {
	result := make([]E, list.size)
	current := list.head
	index := 0
	for current != nil {
		result[index] = current.data
		current = current.next
		index++
	}
	return result
}

// AddFirst adds an element to the beginning of the list
func (list *LinkedList[E]) AddFirst(element E) {
	newNode := &Node[E]{data: element}

	if list.head == nil {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.next = list.head
		list.head.prev = newNode
		list.head = newNode
	}

	list.size++
}

// AddLast adds an element to the end of the list
func (list *LinkedList[E]) AddLast(element E) {
	newNode := &Node[E]{data: element}

	if list.tail == nil {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.prev = list.tail
		list.tail.next = newNode
		list.tail = newNode
	}

	list.size++
}

// RemoveFirst removes and returns the first element of the list
func (list *LinkedList[E]) RemoveFirst() (E, bool) {
	if list.head == nil {
		return *new(E), false
	}

	data := list.head.data

	if list.head == list.tail {
		list.head = nil
		list.tail = nil
	} else {
		list.head = list.head.next
		list.head.prev = nil
	}

	list.size--
	return data, true
}

// RemoveLast removes and returns the last element of the list
func (list *LinkedList[E]) RemoveLast() (E, bool) {
	if list.tail == nil {
		return *new(E), false
	}

	data := list.tail.data

	if list.head == list.tail {
		list.head = nil
		list.tail = nil
	} else {
		list.tail = list.tail.prev
		list.tail.next = nil
	}

	list.size--
	return data, true
}

// GetFirst returns the first element of the list without removing it
func (list *LinkedList[E]) GetFirst() (E, error) {
	if list.head == nil {
		return *new(E), common.EmptyContainerError("LinkedList")
	}
	return list.head.data, nil
}

// GetLast returns the last element of the list without removing it
func (list *LinkedList[E]) GetLast() (E, error) {
	if list.tail == nil {
		return *new(E), common.EmptyContainerError("LinkedList")
	}
	return list.tail.data, nil
}

// Head returns the head node
func (list *LinkedList[E]) Head() *Node[E] {
	return list.head
}

// Tail returns the tail node
func (list *LinkedList[E]) Tail() *Node[E] {
	return list.tail
}

// getNodeAt retrieves the node at the specified index (internal method)
func (list *LinkedList[E]) getNodeAt(index int) *Node[E] {
	// Search from the head
	if index < list.size/2 {
		current := list.head
		for i := 0; i < index; i++ {
			current = current.next
		}
		return current
	} else {
		// Search from the tail
		current := list.tail
		for i := list.size - 1; i > index; i-- {
			current = current.prev
		}
		return current
	}
}

// Iterator returns an iterator for traversing the elements in the list
func (list *LinkedList[E]) Iterator() common.Iterator[E] {
	return &linkedListIterator[E]{list: list, current: list.head, lastReturned: nil}
}

// linkedListIterator is the iterator implementation for LinkedList
type linkedListIterator[E any] struct {
	list         *LinkedList[E]
	current      *Node[E]
	lastReturned *Node[E]
}

// HasNext checks if the iterator has a next element
func (it *linkedListIterator[E]) HasNext() bool {
	return it.current != nil
}

// Next returns the next element in the iterator
func (it *linkedListIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		return *new(E), false
	}

	data := it.current.data
	it.lastReturned = it.current
	it.current = it.current.next
	return data, true
}

// Remove removes the last element returned by the iterator
func (it *linkedListIterator[E]) Remove() bool {
	if it.lastReturned == nil {
		return false
	}

	// Update the linked list structure
	if it.lastReturned.prev != nil {
		it.lastReturned.prev.next = it.lastReturned.next
	} else {
		it.list.head = it.lastReturned.next
	}

	if it.lastReturned.next != nil {
		it.lastReturned.next.prev = it.lastReturned.prev
	} else {
		it.list.tail = it.lastReturned.prev
	}

	it.list.size--
	it.lastReturned = nil
	return true
}
