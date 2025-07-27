// Package queue 提供队列数据结构的实现
package queue

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptyQueue 表示队列为空错误
	ErrEmptyQueue = errors.New("queue is empty")

	// ErrFullQueue 表示队列已满错误
	ErrFullQueue = errors.New("queue is full")
)

// node 是链表的节点
type node[E comparable] struct {
	value E
	next  *node[E]
	prev  *node[E]
}

// LinkedList 是一个基于双向链表的Deque实现
type LinkedList[E comparable] struct {
	head   *node[E]
	tail   *node[E]
	size   int
	maxCap int // 最大容量，0表示无界
}

// New 创建一个新的无界LinkedList
func New[E comparable]() *LinkedList[E] {
	return &LinkedList[E]{maxCap: 0}
}

// WithCapacity 创建一个具有指定最大容量的LinkedList
func WithCapacity[E comparable](maxCapacity int) *LinkedList[E] {
	return &LinkedList[E]{maxCap: maxCapacity}
}

// FromSlice 从切片创建一个新的LinkedList
func FromSlice[E comparable](elements []E) *LinkedList[E] {
	list := New[E]()
	for _, e := range elements {
		list.Add(e)
	}
	return list
}

// Size 返回队列中的元素数量
func (list *LinkedList[E]) Size() int {
	return list.size
}

// IsEmpty 检查队列是否为空
func (list *LinkedList[E]) IsEmpty() bool {
	return list.size == 0
}

// isFull 检查队列是否已满
func (list *LinkedList[E]) isFull() bool {
	return list.maxCap > 0 && list.size >= list.maxCap
}

// Clear 清空队列
func (list *LinkedList[E]) Clear() {
	list.head = nil
	list.tail = nil
	list.size = 0
}

// Contains 检查队列是否包含指定元素
func (list *LinkedList[E]) Contains(e E) bool {
	if list.IsEmpty() {
		return false
	}

	current := list.head
	for current != nil {
		if current.value == e {
			return true
		}
		current = current.next
	}

	return false
}

// ForEach 对队列中的每个元素执行给定的操作
func (list *LinkedList[E]) ForEach(f func(E)) {
	current := list.head
	for current != nil {
		f(current.value)
		current = current.next
	}
}

// String 返回队列的字符串表示
func (list *LinkedList[E]) String() string {
	if list.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	current := list.head
	for current != nil {
		sb.WriteString(fmt.Sprintf("%v", current.value))
		if current.next != nil {
			sb.WriteString(", ")
		}
		current = current.next
	}

	sb.WriteString("]")
	return sb.String()
}

// Add 将元素添加到队列尾部
func (list *LinkedList[E]) Add(e E) error {
	if list.isFull() {
		return ErrFullQueue
	}

	return list.addLast(e)
}

// Offer 将元素添加到队列尾部
func (list *LinkedList[E]) Offer(e E) bool {
	if list.isFull() {
		return false
	}

	_ = list.addLast(e)
	return true
}

// Remove 移除并返回队列头部的元素
func (list *LinkedList[E]) Remove() (E, error) {
	return list.RemoveFirst()
}

// Poll 移除并返回队列头部的元素
func (list *LinkedList[E]) Poll() (E, bool) {
	return list.PollFirst()
}

// Element 返回队列头部的元素，但不移除
func (list *LinkedList[E]) Element() (E, error) {
	return list.GetFirst()
}

// Peek 返回队列头部的元素，但不移除
func (list *LinkedList[E]) Peek() (E, bool) {
	return list.PeekFirst()
}

// AddFirst 将元素添加到队列头部
func (list *LinkedList[E]) AddFirst(e E) error {
	if list.isFull() {
		return ErrFullQueue
	}

	newNode := &node[E]{value: e}

	if list.IsEmpty() {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.next = list.head
		list.head.prev = newNode
		list.head = newNode
	}

	list.size++
	return nil
}

// AddLast 将元素添加到队列尾部
func (list *LinkedList[E]) AddLast(e E) error {
	if list.isFull() {
		return ErrFullQueue
	}

	return list.addLast(e)
}

// addLast 是内部方法，将元素添加到队列尾部
func (list *LinkedList[E]) addLast(e E) error {
	newNode := &node[E]{value: e}

	if list.IsEmpty() {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.prev = list.tail
		list.tail.next = newNode
		list.tail = newNode
	}

	list.size++
	return nil
}

// OfferFirst 将元素添加到队列头部
func (list *LinkedList[E]) OfferFirst(e E) bool {
	if list.isFull() {
		return false
	}

	_ = list.AddFirst(e)
	return true
}

// OfferLast 将元素添加到队列尾部
func (list *LinkedList[E]) OfferLast(e E) bool {
	if list.isFull() {
		return false
	}

	_ = list.addLast(e)
	return true
}

// RemoveFirst 移除并返回队列头部的元素
func (list *LinkedList[E]) RemoveFirst() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	value := list.head.value

	if list.size == 1 {
		list.head = nil
		list.tail = nil
	} else {
		list.head = list.head.next
		list.head.prev = nil
	}

	list.size--
	return value, nil
}

// RemoveLast 移除并返回队列尾部的元素
func (list *LinkedList[E]) RemoveLast() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	value := list.tail.value

	if list.size == 1 {
		list.head = nil
		list.tail = nil
	} else {
		list.tail = list.tail.prev
		list.tail.next = nil
	}

	list.size--
	return value, nil
}

// PollFirst 移除并返回队列头部的元素
func (list *LinkedList[E]) PollFirst() (E, bool) {
	value, err := list.RemoveFirst()
	return value, err == nil
}

// PollLast 移除并返回队列尾部的元素
func (list *LinkedList[E]) PollLast() (E, bool) {
	value, err := list.RemoveLast()
	return value, err == nil
}

// GetFirst 返回队列头部的元素，但不移除
func (list *LinkedList[E]) GetFirst() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return list.head.value, nil
}

// GetLast 返回队列尾部的元素，但不移除
func (list *LinkedList[E]) GetLast() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return list.tail.value, nil
}

// PeekFirst 返回队列头部的元素，但不移除
func (list *LinkedList[E]) PeekFirst() (E, bool) {
	value, err := list.GetFirst()
	return value, err == nil
}

// PeekLast 返回队列尾部的元素，但不移除
func (list *LinkedList[E]) PeekLast() (E, bool) {
	value, err := list.GetLast()
	return value, err == nil
}

// ToSlice 返回包含队列所有元素的切片
func (list *LinkedList[E]) ToSlice() []E {
	result := make([]E, list.size)
	i := 0
	current := list.head

	for current != nil {
		result[i] = current.value
		current = current.next
		i++
	}

	return result
}
