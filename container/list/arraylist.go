// Package list 提供列表数据结构的实现
package list

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

var (
	// ErrIndexOutOfBounds 表示索引超出范围错误
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrInvalidRange 表示无效的范围错误
	ErrInvalidRange = errors.New("invalid range: fromIndex > toIndex")
)

// ArrayList 是一个基于动态数组的List实现
type ArrayList[E any] struct {
	elements []E
}

// New 创建一个新的ArrayList
func New[E any]() *ArrayList[E] {
	return &ArrayList[E]{elements: make([]E, 0)}
}

// WithCapacity 创建一个具有指定初始容量的ArrayList
func WithCapacity[E any](initialCapacity int) *ArrayList[E] {
	if initialCapacity < 0 {
		initialCapacity = 0
	}
	return &ArrayList[E]{elements: make([]E, 0, initialCapacity)}
}

// FromSlice 从切片创建一个新的ArrayList
func FromSlice[E any](elements []E) *ArrayList[E] {
	result := New[E]()
	result.elements = append(result.elements, elements...)
	return result
}

// Add 添加一个元素到列表末尾
func (list *ArrayList[E]) Add(e E) bool {
	list.elements = append(list.elements, e)
	return true
}

// Insert 在指定位置插入元素
func (list *ArrayList[E]) Insert(index int, e E) error {
	if index < 0 || index > len(list.elements) {
		return ErrIndexOutOfBounds
	}

	// 在末尾添加
	if index == len(list.elements) {
		list.Add(e)
		return nil
	}

	// 在中间插入
	list.elements = append(list.elements, *new(E)) // 扩展切片
	copy(list.elements[index+1:], list.elements[index:])
	list.elements[index] = e
	return nil
}

// Get 获取指定索引的元素
func (list *ArrayList[E]) Get(index int) (E, error) {
	if index < 0 || index >= len(list.elements) {
		return *new(E), ErrIndexOutOfBounds
	}
	return list.elements[index], nil
}

// Set 替换指定索引的元素
func (list *ArrayList[E]) Set(index int, e E) (E, error) {
	if index < 0 || index >= len(list.elements) {
		return *new(E), ErrIndexOutOfBounds
	}

	old := list.elements[index]
	list.elements[index] = e
	return old, nil
}

// RemoveAt 移除指定索引的元素
func (list *ArrayList[E]) RemoveAt(index int) (E, error) {
	if index < 0 || index >= len(list.elements) {
		return *new(E), ErrIndexOutOfBounds
	}

	removed := list.elements[index]
	// 移动元素以填补空缺
	copy(list.elements[index:], list.elements[index+1:])
	// 缩小切片
	list.elements = list.elements[:len(list.elements)-1]

	return removed, nil
}

// Remove 移除第一个匹配的元素
func (list *ArrayList[E]) Remove(e E) bool {
	index := list.IndexOf(e)
	if index == -1 {
		return false
	}

	_, err := list.RemoveAt(index)
	return err == nil
}

// Contains 检查列表是否包含指定元素
func (list *ArrayList[E]) Contains(e E) bool {
	return list.IndexOf(e) != -1
}

// IndexOf 返回指定元素在列表中第一次出现的索引
func (list *ArrayList[E]) IndexOf(e E) int {
	for i, v := range list.elements {
		if common.Equal(v, e) {
			return i
		}
	}
	return -1
}

// LastIndexOf 返回指定元素在列表中最后一次出现的索引
func (list *ArrayList[E]) LastIndexOf(e E) int {
	for i := len(list.elements) - 1; i >= 0; i-- {
		if common.Equal(list.elements[i], e) {
			return i
		}
	}
	return -1
}

// Size 返回列表中的元素数量
func (list *ArrayList[E]) Size() int {
	return len(list.elements)
}

// IsEmpty 检查列表是否为空
func (list *ArrayList[E]) IsEmpty() bool {
	return len(list.elements) == 0
}

// Clear 清空列表
func (list *ArrayList[E]) Clear() {
	list.elements = make([]E, 0)
}

// ToSlice 返回包含列表所有元素的切片
func (list *ArrayList[E]) ToSlice() []E {
	result := make([]E, len(list.elements))
	copy(result, list.elements)
	return result
}

// SubList 返回列表中指定范围的视图
func (list *ArrayList[E]) SubList(fromIndex, toIndex int) (List[E], error) {
	if fromIndex < 0 || toIndex > len(list.elements) {
		return nil, ErrIndexOutOfBounds
	}
	if fromIndex > toIndex {
		return nil, ErrInvalidRange
	}

	result := New[E]()
	result.elements = append(result.elements, list.elements[fromIndex:toIndex]...)
	return result, nil
}

// ForEach 对列表中的每个元素执行给定的操作
func (list *ArrayList[E]) ForEach(f func(E)) {
	for _, e := range list.elements {
		f(e)
	}
}

// String 返回列表的字符串表示
func (list *ArrayList[E]) String() string {
	if list.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	for i, e := range list.elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", e))
	}

	sb.WriteString("]")
	return sb.String()
}

// Iterator 返回一个迭代器，用于遍历列表中的元素
func (list *ArrayList[E]) Iterator() common.Iterator[E] {
	return &arrayListIterator[E]{list: list, cursor: 0, lastRet: -1}
}

// arrayListIterator 是ArrayList的迭代器实现
type arrayListIterator[E any] struct {
	list    *ArrayList[E]
	cursor  int // 下一个元素的索引
	lastRet int // 最后返回的元素的索引，如果没有则为-1
}

// HasNext 检查迭代器是否还有下一个元素
func (it *arrayListIterator[E]) HasNext() bool {
	return it.cursor < it.list.Size()
}

// Next 返回迭代器中的下一个元素
func (it *arrayListIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		return *new(E), false
	}

	it.lastRet = it.cursor
	it.cursor++
	elem, _ := it.list.Get(it.lastRet)
	return elem, true
}

// Remove 移除迭代器最后返回的元素
func (it *arrayListIterator[E]) Remove() bool {
	if it.lastRet < 0 {
		return false
	}

	_, err := it.list.RemoveAt(it.lastRet)
	if err != nil {
		return false
	}

	it.cursor = it.lastRet
	it.lastRet = -1
	return true
}
