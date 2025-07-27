// Package set 提供集合数据结构的实现
package set

import (
	"github.com/chenjianyu/collections/container/common"
)

// Set 表示不包含重复元素的集合
type Set[E any] interface {
	common.Container[E]
	common.Iterable[E]

	// Add 添加元素到集合
	// 如果集合已包含该元素，则返回false，否则返回true
	Add(e E) bool

	// Remove 从集合中移除指定元素
	// 如果集合包含该元素，则返回true，否则返回false
	Remove(e E) bool

	// ToSlice 返回包含集合所有元素的切片
	ToSlice() []E

	// Union 返回此集合与另一个集合的并集
	Union(other Set[E]) Set[E]

	// Intersection 返回此集合与另一个集合的交集
	Intersection(other Set[E]) Set[E]

	// Difference 返回此集合与另一个集合的差集
	// 即在此集合中但不在另一个集合中的元素
	Difference(other Set[E]) Set[E]

	// IsSubsetOf 检查此集合是否为另一个集合的子集
	IsSubsetOf(other Set[E]) bool

	// IsSupersetOf 检查此集合是否为另一个集合的超集
	IsSupersetOf(other Set[E]) bool
}
