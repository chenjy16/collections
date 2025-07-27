// Package list 提供列表数据结构的实现
package list

import (
	"github.com/chenjianyu/collections/container/common"
)

// List 表示有序的元素集合，允许重复元素
type List[E any] interface {
	common.Container[E]
	common.Iterable[E]

	// Add 添加元素到列表末尾
	// 返回是否成功添加
	Add(e E) bool

	// Insert 在指定位置插入元素
	// 如果索引无效，返回错误
	Insert(index int, e E) error

	// Get 获取指定索引的元素
	// 如果索引无效，返回零值和错误
	Get(index int) (E, error)

	// Set 替换指定索引的元素
	// 返回被替换的元素和是否成功
	Set(index int, e E) (E, error)

	// RemoveAt 移除指定索引的元素
	// 返回被移除的元素和是否成功
	RemoveAt(index int) (E, error)

	// Remove 移除第一个匹配的元素
	// 返回是否成功移除
	Remove(e E) bool

	// IndexOf 返回指定元素在列表中第一次出现的索引
	// 如果不存在，返回-1
	IndexOf(e E) int

	// LastIndexOf 返回指定元素在列表中最后一次出现的索引
	// 如果不存在，返回-1
	LastIndexOf(e E) int

	// SubList 返回列表中指定范围的视图
	// 如果索引无效，返回错误
	SubList(fromIndex, toIndex int) (List[E], error)

	// ToSlice 返回包含列表所有元素的切片
	ToSlice() []E
}
