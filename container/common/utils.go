// Package common 提供容器库的通用工具和接口
package common

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"strings"
)

// Equal 比较两个值是否相等
// 使用reflect.DeepEqual进行深度比较
func Equal[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}

// Hash 计算值的哈希码
// 使用FNV-1a算法对值的字符串表示进行哈希
func Hash[T any](value T) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", value)))
	return int(h.Sum32())
}

// Compare 比较两个值
// 返回:
//   - 负数: a < b
//   - 0: a == b
//   - 正数: a > b
func Compare[T any](a, b T) int {
	// 处理nil值
	aIsNil := reflect.ValueOf(a).IsZero()
	bIsNil := reflect.ValueOf(b).IsZero()

	if aIsNil && bIsNil {
		return 0
	}
	if aIsNil {
		return -1
	}
	if bIsNil {
		return 1
	}

	// 尝试类型断言为可比较类型
	switch v1 := any(a).(type) {
	case int:
		return v1 - any(b).(int)
	case int8:
		return int(v1) - int(any(b).(int8))
	case int16:
		return int(v1) - int(any(b).(int16))
	case int32:
		return int(v1) - int(any(b).(int32))
	case int64:
		v2 := any(b).(int64)
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	case uint:
		v2 := any(b).(uint)
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	case uint8:
		return int(v1) - int(any(b).(uint8))
	case uint16:
		return int(v1) - int(any(b).(uint16))
	case uint32:
		v2 := any(b).(uint32)
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	case uint64:
		v2 := any(b).(uint64)
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	case float32:
		v2 := any(b).(float32)
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	case float64:
		v2 := any(b).(float64)
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	case string:
		return strings.Compare(v1, any(b).(string))
	case fmt.Stringer:
		return strings.Compare(v1.String(), any(b).(fmt.Stringer).String())
	default:
		// 对于不可直接比较的类型，使用哈希值比较
		return Hash(a) - Hash(b)
	}
}

// Pair 表示键值对
type Pair[K any, V any] struct {
	Key   K
	Value V
}

// NewPair 创建一个新的键值对
func NewPair[K any, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// String 返回Pair的字符串表示
func (p Pair[K, V]) String() string {
	return fmt.Sprintf("%v=%v", p.Key, p.Value)
}
