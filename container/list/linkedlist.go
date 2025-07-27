// Package list 提供列表数据结构的实现
package list

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

var (
	// ErrEmptyList 表示列表为空错误
	ErrEmptyList = errors.New("list is empty")
)

// Node 是双向链表的节点
type Node[E any] struct {
	Value E
	Next  *Node[E]
	Prev  *Node[E]
}

// LinkedList 是基于双向链表的List实现
type LinkedList[E any] struct {
	head *Node[E]
	tail *Node[E]
	size int
}

// NewLinkedList 创建一个新的LinkedList
func NewLinkedList[E any]() *LinkedList[E] {
	return &LinkedList[E]{}
}

// LinkedListFromSlice 从切片创建一个新的LinkedList
func LinkedListFromSlice[E any](elements []E) *LinkedList[E] {
	list := NewLinkedList[E]()
	for _, e := range elements {
		list.Add(e)
	}
	return list
}

// Size 返回列表中的元素数量
func (list *LinkedList[E]) Size() int {
	return list.size
}

// IsEmpty 检查列表是否为空
func (list *LinkedList[E]) IsEmpty() bool {
	return list.size == 0
}

// Clear 清空列表
func (list *LinkedList[E]) Clear() {
	list.head = nil
	list.tail = nil
	list.size = 0
}

// Contains 检查列表是否包含指定元素
func (list *LinkedList[E]) Contains(e E) bool {
	return list.IndexOf(e) != -1
}

// ForEach 对列表中的每个元素执行给定的操作
func (list *LinkedList[E]) ForEach(f func(E)) {
	current := list.head
	for current != nil {
		f(current.Value)
		current = current.Next
	}
}

// String 返回列表的字符串表示
func (list *LinkedList[E]) String() string {
	if list.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	current := list.head
	for current != nil {
		sb.WriteString(fmt.Sprintf("%v", current.Value))
		if current.Next != nil {
			sb.WriteString(", ")
		}
		current = current.Next
	}

	sb.WriteString("]")
	return sb.String()
}

// Add 添加元素到列表末尾
func (list *LinkedList[E]) Add(e E) bool {
	list.AddLast(e)
	return true
}

// Insert 在指定位置插入元素
func (list *LinkedList[E]) Insert(index int, e E) error {
	if index < 0 || index > list.size {
		return ErrIndexOutOfBounds
	}

	if index == 0 {
		list.AddFirst(e)
		return nil
	}

	if index == list.size {
		list.AddLast(e)
		return nil
	}

	newNode := &Node[E]{Value: e}
	current := list.getNodeAt(index)

	newNode.Next = current
	newNode.Prev = current.Prev
	current.Prev.Next = newNode
	current.Prev = newNode

	list.size++
	return nil
}

// Get 获取指定索引的元素
func (list *LinkedList[E]) Get(index int) (E, error) {
	if index < 0 || index >= list.size {
		return *new(E), ErrIndexOutOfBounds
	}

	node := list.getNodeAt(index)
	return node.Value, nil
}

// Set 替换指定索引的元素
func (list *LinkedList[E]) Set(index int, e E) (E, error) {
	if index < 0 || index >= list.size {
		return *new(E), ErrIndexOutOfBounds
	}

	node := list.getNodeAt(index)
	oldValue := node.Value
	node.Value = e
	return oldValue, nil
}

// RemoveAt 移除指定索引的元素
func (list *LinkedList[E]) RemoveAt(index int) (E, error) {
	if index < 0 || index >= list.size {
		return *new(E), ErrIndexOutOfBounds
	}

	if index == 0 {
		return list.RemoveFirst()
	}

	if index == list.size-1 {
		return list.RemoveLast()
	}

	node := list.getNodeAt(index)
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev

	list.size--
	return node.Value, nil
}

// Remove 移除第一个匹配的元素
func (list *LinkedList[E]) Remove(e E) bool {
	index := list.IndexOf(e)
	if index == -1 {
		return false
	}

	_, _ = list.RemoveAt(index)
	return true
}

// IndexOf 返回指定元素在列表中第一次出现的索引
func (list *LinkedList[E]) IndexOf(e E) int {
	current := list.head
	index := 0

	for current != nil {
		if common.Equal(current.Value, e) {
			return index
		}
		current = current.Next
		index++
	}

	return -1
}

// LastIndexOf 返回指定元素在列表中最后一次出现的索引
func (list *LinkedList[E]) LastIndexOf(e E) int {
	current := list.tail
	index := list.size - 1

	for current != nil {
		if common.Equal(current.Value, e) {
			return index
		}
		current = current.Prev
		index--
	}

	return -1
}

// SubList 返回列表中指定范围的视图
func (list *LinkedList[E]) SubList(fromIndex, toIndex int) (List[E], error) {
	if fromIndex < 0 || toIndex > list.size || fromIndex > toIndex {
		return nil, ErrIndexOutOfBounds
	}

	subList := NewLinkedList[E]()
	current := list.getNodeAt(fromIndex)

	for i := fromIndex; i < toIndex && current != nil; i++ {
		subList.Add(current.Value)
		current = current.Next
	}

	return subList, nil
}

// ToSlice 返回包含列表所有元素的切片
func (list *LinkedList[E]) ToSlice() []E {
	result := make([]E, list.size)
	current := list.head
	index := 0

	for current != nil {
		result[index] = current.Value
		current = current.Next
		index++
	}

	return result
}

// AddFirst 将元素添加到列表头部
func (list *LinkedList[E]) AddFirst(e E) {
	newNode := &Node[E]{Value: e}

	if list.IsEmpty() {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.Next = list.head
		list.head.Prev = newNode
		list.head = newNode
	}

	list.size++
}

// AddLast 将元素添加到列表尾部
func (list *LinkedList[E]) AddLast(e E) {
	newNode := &Node[E]{Value: e}

	if list.IsEmpty() {
		list.head = newNode
		list.tail = newNode
	} else {
		newNode.Prev = list.tail
		list.tail.Next = newNode
		list.tail = newNode
	}

	list.size++
}

// RemoveFirst 移除并返回列表头部的元素
func (list *LinkedList[E]) RemoveFirst() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyList
	}

	value := list.head.Value

	if list.size == 1 {
		list.head = nil
		list.tail = nil
	} else {
		list.head = list.head.Next
		list.head.Prev = nil
	}

	list.size--
	return value, nil
}

// RemoveLast 移除并返回列表尾部的元素
func (list *LinkedList[E]) RemoveLast() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyList
	}

	value := list.tail.Value

	if list.size == 1 {
		list.head = nil
		list.tail = nil
	} else {
		list.tail = list.tail.Prev
		list.tail.Next = nil
	}

	list.size--
	return value, nil
}

// GetFirst 返回列表头部的元素，但不移除
func (list *LinkedList[E]) GetFirst() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyList
	}

	return list.head.Value, nil
}

// GetLast 返回列表尾部的元素，但不移除
func (list *LinkedList[E]) GetLast() (E, error) {
	if list.IsEmpty() {
		return *new(E), ErrEmptyList
	}

	return list.tail.Value, nil
}

// Head 返回头节点
func (list *LinkedList[E]) Head() *Node[E] {
	return list.head
}

// Tail 返回尾节点
func (list *LinkedList[E]) Tail() *Node[E] {
	return list.tail
}

// getNodeAt 获取指定索引的节点（内部方法）
func (list *LinkedList[E]) getNodeAt(index int) *Node[E] {
	if index < list.size/2 {
		// 从头部开始搜索
		current := list.head
		for i := 0; i < index; i++ {
			current = current.Next
		}
		return current
	} else {
		// 从尾部开始搜索
		current := list.tail
		for i := list.size - 1; i > index; i-- {
			current = current.Prev
		}
		return current
	}
}

// Iterator 返回一个迭代器，用于遍历列表中的元素
func (list *LinkedList[E]) Iterator() common.Iterator[E] {
	return &linkedListIterator[E]{list: list, current: list.head, lastReturned: nil}
}

// linkedListIterator 是LinkedList的迭代器实现
type linkedListIterator[E any] struct {
	list         *LinkedList[E]
	current      *Node[E]
	lastReturned *Node[E]
}

// HasNext 检查迭代器是否还有下一个元素
func (it *linkedListIterator[E]) HasNext() bool {
	return it.current != nil
}

// Next 返回迭代器中的下一个元素
func (it *linkedListIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		return *new(E), false
	}

	it.lastReturned = it.current
	value := it.current.Value
	it.current = it.current.Next
	return value, true
}

// Remove 移除迭代器最后返回的元素
func (it *linkedListIterator[E]) Remove() bool {
	if it.lastReturned == nil {
		return false
	}

	nodeToRemove := it.lastReturned
	it.lastReturned = nil

	// 更新链表结构
	if nodeToRemove.Prev != nil {
		nodeToRemove.Prev.Next = nodeToRemove.Next
	} else {
		it.list.head = nodeToRemove.Next
	}

	if nodeToRemove.Next != nil {
		nodeToRemove.Next.Prev = nodeToRemove.Prev
	} else {
		it.list.tail = nodeToRemove.Prev
	}

	it.list.size--
	return true
}