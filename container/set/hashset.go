// Package set 提供集合数据结构的实现
package set

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
)

// HashSet 是一个基于哈希表的Set实现
type HashSet[E any] struct {
	elements map[any]E
}

// New 创建一个新的HashSet
func New[E any]() *HashSet[E] {
	return &HashSet[E]{elements: make(map[any]E)}
}

// FromSlice 从切片创建一个新的HashSet
func FromSlice[E any](elements []E) *HashSet[E] {
	set := New[E]()
	for _, e := range elements {
		set.Add(e)
	}
	return set
}

// Add 添加元素到集合
func (set *HashSet[E]) Add(e E) bool {
	hashCode := common.Hash(e)
	_, exists := set.elements[hashCode]
	if exists {
		// 检查是否真的是相同元素（处理哈希冲突）
		if common.Equal(set.elements[hashCode], e) {
			return false
		}
	}
	set.elements[hashCode] = e
	return true
}

// Remove 从集合中移除指定元素
func (set *HashSet[E]) Remove(e E) bool {
	hashCode := common.Hash(e)
	_, exists := set.elements[hashCode]
	if !exists {
		return false
	}

	// 检查是否真的是相同元素（处理哈希冲突）
	if !common.Equal(set.elements[hashCode], e) {
		return false
	}

	delete(set.elements, hashCode)
	return true
}

// Contains 检查集合是否包含指定元素
func (set *HashSet[E]) Contains(e E) bool {
	hashCode := common.Hash(e)
	val, exists := set.elements[hashCode]
	if !exists {
		return false
	}

	// 检查是否真的是相同元素（处理哈希冲突）
	return common.Equal(val, e)
}

// Size 返回集合中的元素数量
func (set *HashSet[E]) Size() int {
	return len(set.elements)
}

// IsEmpty 检查集合是否为空
func (set *HashSet[E]) IsEmpty() bool {
	return len(set.elements) == 0
}

// Clear 清空集合
func (set *HashSet[E]) Clear() {
	set.elements = make(map[any]E)
}

// ToSlice 返回包含集合所有元素的切片
func (set *HashSet[E]) ToSlice() []E {
	result := make([]E, 0, len(set.elements))
	for _, v := range set.elements {
		result = append(result, v)
	}
	return result
}

// ForEach 对集合中的每个元素执行给定的操作
func (set *HashSet[E]) ForEach(f func(E)) {
	for _, v := range set.elements {
		f(v)
	}
}

// Union 返回此集合与另一个集合的并集
func (set *HashSet[E]) Union(other Set[E]) Set[E] {
	result := New[E]()

	// 添加此集合的所有元素
	for _, v := range set.elements {
		result.Add(v)
	}

	// 添加另一个集合的所有元素
	other.ForEach(func(e E) {
		result.Add(e)
	})

	return result
}

// Intersection 返回此集合与另一个集合的交集
func (set *HashSet[E]) Intersection(other Set[E]) Set[E] {
	result := New[E]()

	// 添加同时存在于两个集合中的元素
	set.ForEach(func(e E) {
		if other.Contains(e) {
			result.Add(e)
		}
	})

	return result
}

// Difference 返回此集合与另一个集合的差集
func (set *HashSet[E]) Difference(other Set[E]) Set[E] {
	result := New[E]()

	// 添加在此集合中但不在另一个集合中的元素
	set.ForEach(func(e E) {
		if !other.Contains(e) {
			result.Add(e)
		}
	})

	return result
}

// IsSubsetOf 检查此集合是否为另一个集合的子集
func (set *HashSet[E]) IsSubsetOf(other Set[E]) bool {
	// 空集是任何集合的子集
	if set.IsEmpty() {
		return true
	}

	// 如果此集合的大小大于另一个集合，则它不可能是子集
	if set.Size() > other.Size() {
		return false
	}

	// 检查此集合的每个元素是否都在另一个集合中
	for _, v := range set.elements {
		if !other.Contains(v) {
			return false
		}
	}

	return true
}

// IsSupersetOf 检查此集合是否为另一个集合的超集
func (set *HashSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(set)
}

// String 返回集合的字符串表示
func (set *HashSet[E]) String() string {
	if set.IsEmpty() {
		return "[]"
	}

	elements := set.ToSlice()
	var sb strings.Builder
	sb.WriteString("[")

	for i, e := range elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", e))
	}

	sb.WriteString("]")
	return sb.String()
}

// Iterator 返回一个迭代器，用于遍历集合中的元素
func (set *HashSet[E]) Iterator() common.Iterator[E] {
	return &hashSetIterator[E]{set: set, elements: set.ToSlice(), cursor: 0, lastRet: -1}
}

// hashSetIterator 是HashSet的迭代器实现
type hashSetIterator[E any] struct {
	set      *HashSet[E]
	elements []E
	cursor   int // 下一个元素的索引
	lastRet  int // 最后返回的元素的索引，如果没有则为-1
}

// HasNext 检查迭代器是否还有下一个元素
func (it *hashSetIterator[E]) HasNext() bool {
	return it.cursor < len(it.elements)
}

// Next 返回迭代器中的下一个元素
func (it *hashSetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		return *new(E), false
	}

	it.lastRet = it.cursor
	it.cursor++
	return it.elements[it.lastRet], true
}

// Remove 移除迭代器最后返回的元素
func (it *hashSetIterator[E]) Remove() bool {
	if it.lastRet < 0 {
		return false
	}

	removed := it.set.Remove(it.elements[it.lastRet])
	if removed {
		// 更新迭代器的元素列表
		it.elements = it.set.ToSlice()
		it.cursor = 0 // 重置游标，因为元素列表已更改
		it.lastRet = -1
	}

	return removed
}
