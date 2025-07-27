// Package maps 提供映射数据结构的实现
package maps

// Map 表示键值对的集合，其中每个键最多只能映射到一个值
type Map[K any, V any] interface {
	// Put 将指定的值与此映射中的指定键关联
	// 如果映射先前包含该键的映射关系，则返回旧值和true，否则返回零值和false
	Put(key K, value V) (V, bool)

	// Get 返回指定键所映射的值
	// 如果此映射不包含该键的映射关系，则返回零值和false
	Get(key K) (V, bool)

	// Remove 如果存在，则从此映射中移除指定键的映射关系
	// 返回先前与指定键关联的值，如果没有该键的映射关系，则返回零值和false
	Remove(key K) (V, bool)

	// ContainsKey 如果此映射包含指定键的映射关系，则返回true
	ContainsKey(key K) bool

	// ContainsValue 如果此映射将一个或多个键映射到指定值，则返回true
	ContainsValue(value V) bool

	// Size 返回此映射中的键值映射关系数
	Size() int

	// IsEmpty 如果此映射未包含键值映射关系，则返回true
	IsEmpty() bool

	// Clear 从此映射中移除所有映射关系
	Clear()

	// Keys 返回此映射中包含的键的集合视图
	Keys() []K

	// Values 返回此映射中包含的值的集合视图
	Values() []V

	// ForEach 对此映射中的每个键值对执行给定的操作
	ForEach(f func(K, V))

	// String 返回此映射的字符串表示形式
	String() string
}
