// Package queue 提供队列数据结构的实现
package queue

import (
	"github.com/chenjianyu/collections/container/list"
)

// LinkedList 是一个基于双向链表的Deque实现，复用list包中的LinkedList
type LinkedList[E comparable] struct {
	list   *list.LinkedList[E]
	maxCap int // 最大容量，0表示无界
}

// New 创建一个新的无界LinkedList
func New[E comparable]() *LinkedList[E] {
	return &LinkedList[E]{
		list:   list.NewLinkedList[E](),
		maxCap: 0,
	}
}

// WithCapacity 创建一个具有指定最大容量的LinkedList
func WithCapacity[E comparable](maxCapacity int) *LinkedList[E] {
	return &LinkedList[E]{
		list:   list.NewLinkedList[E](),
		maxCap: maxCapacity,
	}
}

// FromSlice 从切片创建一个新的LinkedList
func FromSlice[E comparable](elements []E) *LinkedList[E] {
	queue := New[E]()
	for _, e := range elements {
		queue.Add(e)
	}
	return queue
}

// Size 返回队列中的元素数量
func (queue *LinkedList[E]) Size() int {
	return queue.list.Size()
}

// IsEmpty 检查队列是否为空
func (queue *LinkedList[E]) IsEmpty() bool {
	return queue.list.IsEmpty()
}

// isFull 检查队列是否已满
func (queue *LinkedList[E]) isFull() bool {
	return queue.maxCap > 0 && queue.list.Size() >= queue.maxCap
}

// Clear 清空队列
func (queue *LinkedList[E]) Clear() {
	queue.list.Clear()
}

// Contains 检查队列是否包含指定元素
func (queue *LinkedList[E]) Contains(e E) bool {
	return queue.list.Contains(e)
}

// ForEach 对队列中的每个元素执行给定的操作
func (queue *LinkedList[E]) ForEach(f func(E)) {
	queue.list.ForEach(f)
}

// String 返回队列的字符串表示
func (queue *LinkedList[E]) String() string {
	return queue.list.String()
}

// Add 将元素添加到队列尾部
func (queue *LinkedList[E]) Add(e E) error {
	if queue.isFull() {
		return ErrFullQueue
	}

	queue.list.AddLast(e)
	return nil
}

// Offer 将元素添加到队列尾部
func (queue *LinkedList[E]) Offer(e E) bool {
	if queue.isFull() {
		return false
	}

	queue.list.AddLast(e)
	return true
}

// Remove 移除并返回队列头部的元素
func (queue *LinkedList[E]) Remove() (E, error) {
	return queue.RemoveFirst()
}

// Poll 移除并返回队列头部的元素
func (queue *LinkedList[E]) Poll() (E, bool) {
	return queue.PollFirst()
}

// Element 返回队列头部的元素，但不移除
func (queue *LinkedList[E]) Element() (E, error) {
	return queue.GetFirst()
}

// Peek 返回队列头部的元素，但不移除
func (queue *LinkedList[E]) Peek() (E, bool) {
	return queue.PeekFirst()
}

// AddFirst 将元素添加到队列头部
func (queue *LinkedList[E]) AddFirst(e E) error {
	if queue.isFull() {
		return ErrFullQueue
	}

	queue.list.AddFirst(e)
	return nil
}

// AddLast 将元素添加到队列尾部
func (queue *LinkedList[E]) AddLast(e E) error {
	if queue.isFull() {
		return ErrFullQueue
	}

	queue.list.AddLast(e)
	return nil
}

// OfferFirst 将元素添加到队列头部
func (queue *LinkedList[E]) OfferFirst(e E) bool {
	if queue.isFull() {
		return false
	}

	queue.list.AddFirst(e)
	return true
}

// OfferLast 将元素添加到队列尾部
func (queue *LinkedList[E]) OfferLast(e E) bool {
	if queue.isFull() {
		return false
	}

	queue.list.AddLast(e)
	return true
}

// RemoveFirst 移除并返回队列头部的元素
func (queue *LinkedList[E]) RemoveFirst() (E, error) {
	if queue.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return queue.list.RemoveFirst()
}

// RemoveLast 移除并返回队列尾部的元素
func (queue *LinkedList[E]) RemoveLast() (E, error) {
	if queue.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return queue.list.RemoveLast()
}

// PollFirst 移除并返回队列头部的元素
func (queue *LinkedList[E]) PollFirst() (E, bool) {
	value, err := queue.RemoveFirst()
	return value, err == nil
}

// PollLast 移除并返回队列尾部的元素
func (queue *LinkedList[E]) PollLast() (E, bool) {
	value, err := queue.RemoveLast()
	return value, err == nil
}

// GetFirst 返回队列头部的元素，但不移除
func (queue *LinkedList[E]) GetFirst() (E, error) {
	if queue.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return queue.list.GetFirst()
}

// GetLast 返回队列尾部的元素，但不移除
func (queue *LinkedList[E]) GetLast() (E, error) {
	if queue.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return queue.list.GetLast()
}

// PeekFirst 返回队列头部的元素，但不移除
func (queue *LinkedList[E]) PeekFirst() (E, bool) {
	value, err := queue.GetFirst()
	return value, err == nil
}

// PeekLast 返回队列尾部的元素，但不移除
func (queue *LinkedList[E]) PeekLast() (E, bool) {
	value, err := queue.GetLast()
	return value, err == nil
}

// ToSlice 返回包含队列所有元素的切片
func (queue *LinkedList[E]) ToSlice() []E {
	return queue.list.ToSlice()
}