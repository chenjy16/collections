package maps

import (
	"github.com/chenjianyu/collections/container/common"
)

// color 表示红黑树节点的颜色
type color bool

const (
	red   color = true
	black color = false
)

// Hash 计算键的哈希值
func Hash[K any](key K) uint64 {
	return uint64(common.Hash(key))
}

// Equal 比较两个键是否相等
func Equal[K any](a, b K) bool {
	return common.Equal(a, b)
}

// Compare 比较两个键
func Compare[K any](a, b K) int {
	return common.Compare(a, b)
}

// Pair 键值对类型
type Pair[K, V any] struct {
	Key   K
	Value V
}

// NewPair 创建新的键值对
func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}
