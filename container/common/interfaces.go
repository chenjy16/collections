// Package common 提供容器库的通用工具和接口
package common

// 模块路径: github.com/chenjianyu/collections/container/common

// Container 是所有容器类型的基本接口
type Container[E any] interface {
	// Size 返回容器中的元素数量
	Size() int

	// IsEmpty 检查容器是否为空
	IsEmpty() bool

	// Clear 清空容器中的所有元素
	Clear()

	// Contains 检查容器是否包含指定元素
	Contains(e E) bool

	// ForEach 对容器中的每个元素执行给定的操作
	ForEach(f func(E))

	// String 返回容器的字符串表示
	String() string
}

// Iterable 表示可迭代的容器
type Iterable[E any] interface {
	// Iterator 返回一个迭代器，用于遍历容器中的元素
	Iterator() Iterator[E]
}

// Iterator 表示容器的迭代器
type Iterator[E any] interface {
	// HasNext 检查迭代器是否还有下一个元素
	HasNext() bool

	// Next 返回迭代器中的下一个元素
	// 如果没有下一个元素，第二个返回值为false
	Next() (E, bool)

	// Remove 移除迭代器最后返回的元素
	// 如果在调用Next()之前调用Remove()，或者在同一个Next()之后多次调用Remove()，则返回false
	Remove() bool
}

// Comparable 表示可比较的类型
type Comparable interface {
	// CompareTo 将此对象与指定对象进行比较
	// 返回负整数、零或正整数，分别表示此对象小于、等于或大于指定对象
	CompareTo(other interface{}) int
}
