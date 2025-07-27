// Package queue 提供队列数据结构的实现
package queue

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)
// PriorityQueue 是一个基于二叉堆的优先队列实现
// 默认为最小堆，可以通过提供自定义比较器来创建最大堆
type PriorityQueue[E any] struct {
	elements []E
	compare  func(E, E) int
	maxCap   int // 最大容量，0表示无界
}

// NewPriorityQueue 创建一个新的无界优先队列，使用默认比较器
// 默认比较器要求元素类型实现common.Comparable接口
func NewPriorityQueue[E common.Comparable]() *PriorityQueue[E] {
	return &PriorityQueue[E]{
		elements: make([]E, 0),
		compare: func(a, b E) int {
			return a.CompareTo(b)
		},
		maxCap: 0,
	}
}

// NewPriorityQueueWithComparator 创建一个新的无界优先队列，使用自定义比较器
func NewPriorityQueueWithComparator[E any](compare func(E, E) int) *PriorityQueue[E] {
	return &PriorityQueue[E]{
		elements: make([]E, 0),
		compare:  compare,
		maxCap:   0,
	}
}

// WithCapacity 创建一个具有指定最大容量的优先队列，使用默认比较器
func NewPriorityQueueWithCapacity[E common.Comparable](maxCapacity int) *PriorityQueue[E] {
	return &PriorityQueue[E]{
		elements: make([]E, 0, maxCapacity),
		compare: func(a, b E) int {
			return a.CompareTo(b)
		},
		maxCap: maxCapacity,
	}
}

// WithCapacityAndComparator 创建一个具有指定最大容量的优先队列，使用自定义比较器
func NewPriorityQueueWithCapacityAndComparator[E any](maxCapacity int, compare func(E, E) int) *PriorityQueue[E] {
	return &PriorityQueue[E]{
		elements: make([]E, 0, maxCapacity),
		compare:  compare,
		maxCap:   maxCapacity,
	}
}

// NewFromSlice 从切片创建一个新的优先队列，使用默认比较器
func NewFromSlice[E common.Comparable](elements []E) *PriorityQueue[E] {
	pq := NewPriorityQueue[E]()
	for _, e := range elements {
		pq.Add(e)
	}
	return pq
}

// NewFromSliceWithComparator 从切片创建一个新的优先队列，使用自定义比较器
func NewFromSliceWithComparator[E any](elements []E, compare func(E, E) int) *PriorityQueue[E] {
	pq := NewPriorityQueueWithComparator(compare)
	for _, e := range elements {
		pq.Add(e)
	}
	return pq
}

// Size 返回队列中的元素数量
func (pq *PriorityQueue[E]) Size() int {
	return len(pq.elements)
}

// IsEmpty 检查队列是否为空
func (pq *PriorityQueue[E]) IsEmpty() bool {
	return len(pq.elements) == 0
}

// isFull 检查队列是否已满
func (pq *PriorityQueue[E]) isFull() bool {
	return pq.maxCap > 0 && len(pq.elements) >= pq.maxCap
}

// Clear 清空队列
func (pq *PriorityQueue[E]) Clear() {
	pq.elements = pq.elements[:0]
}

// Contains 检查队列是否包含指定元素
func (pq *PriorityQueue[E]) Contains(e E) bool {
	for _, v := range pq.elements {
		if pq.compare(v, e) == 0 {
			return true
		}
	}
	return false
}

// ForEach 对队列中的每个元素执行给定的操作
// 注意：遍历顺序不保证是优先级顺序
func (pq *PriorityQueue[E]) ForEach(f func(E)) {
	for _, e := range pq.elements {
		f(e)
	}
}

// String 返回队列的字符串表示
func (pq *PriorityQueue[E]) String() string {
	if pq.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	for i, e := range pq.elements {
		sb.WriteString(fmt.Sprintf("%v", e))
		if i < len(pq.elements)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]")
	return sb.String()
}

// Add 将元素添加到队列中
func (pq *PriorityQueue[E]) Add(e E) error {
	if pq.isFull() {
		return ErrFullQueue
	}

	pq.elements = append(pq.elements, e)
	pq.siftUp(len(pq.elements) - 1)
	return nil
}

// Offer 将元素添加到队列中
func (pq *PriorityQueue[E]) Offer(e E) bool {
	if pq.isFull() {
		return false
	}

	_ = pq.Add(e)
	return true
}

// Remove 移除并返回队列中优先级最高的元素
func (pq *PriorityQueue[E]) Remove() (E, error) {
	if pq.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	result := pq.elements[0]
	lastIdx := len(pq.elements) - 1
	pq.elements[0] = pq.elements[lastIdx]
	pq.elements = pq.elements[:lastIdx]

	if !pq.IsEmpty() {
		pq.siftDown(0)
	}

	return result, nil
}

// Poll 移除并返回队列中优先级最高的元素
func (pq *PriorityQueue[E]) Poll() (E, bool) {
	result, err := pq.Remove()
	return result, err == nil
}

// Element 返回队列中优先级最高的元素，但不移除
func (pq *PriorityQueue[E]) Element() (E, error) {
	if pq.IsEmpty() {
		return *new(E), ErrEmptyQueue
	}

	return pq.elements[0], nil
}

// Peek 返回队列中优先级最高的元素，但不移除
func (pq *PriorityQueue[E]) Peek() (E, bool) {
	result, err := pq.Element()
	return result, err == nil
}

// ToSlice 返回包含队列所有元素的切片
// 注意：返回的切片不保证是按优先级排序的
func (pq *PriorityQueue[E]) ToSlice() []E {
	result := make([]E, len(pq.elements))
	copy(result, pq.elements)
	return result
}

// ToSortedSlice 返回按优先级排序的队列元素切片
func (pq *PriorityQueue[E]) ToSortedSlice() []E {
	result := make([]E, 0, len(pq.elements))
	tempQueue := &PriorityQueue[E]{
		elements: make([]E, len(pq.elements)),
		compare:  pq.compare,
	}
	copy(tempQueue.elements, pq.elements)

	for !tempQueue.IsEmpty() {
		val, _ := tempQueue.Remove()
		result = append(result, val)
	}

	return result
}

// 内部方法：上移操作，用于维护堆属性
func (pq *PriorityQueue[E]) siftUp(index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if pq.compare(pq.elements[index], pq.elements[parentIndex]) >= 0 {
			break
		}
		pq.elements[index], pq.elements[parentIndex] = pq.elements[parentIndex], pq.elements[index]
		index = parentIndex
	}
}

// 内部方法：下移操作，用于维护堆属性
func (pq *PriorityQueue[E]) siftDown(index int) {
	lastIndex := len(pq.elements) - 1
	for {
		smallestIndex := index
		leftChildIndex := 2*index + 1
		rightChildIndex := 2*index + 2

		if leftChildIndex <= lastIndex && pq.compare(pq.elements[leftChildIndex], pq.elements[smallestIndex]) < 0 {
			smallestIndex = leftChildIndex
		}

		if rightChildIndex <= lastIndex && pq.compare(pq.elements[rightChildIndex], pq.elements[smallestIndex]) < 0 {
			smallestIndex = rightChildIndex
		}

		if smallestIndex == index {
			break
		}

		pq.elements[index], pq.elements[smallestIndex] = pq.elements[smallestIndex], pq.elements[index]
		index = smallestIndex
	}
}
