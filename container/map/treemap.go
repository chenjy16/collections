package maps

import (
	"fmt"
	"strings"
)

// TreeMap 是一个基于红黑树的有序Map实现
type TreeMap[K comparable, V any] struct {
	comparator func(a, b K) int // 键的比较器函数
	root       *mapNode[K, V]   // 根节点
	size       int              // 元素数量
}

// mapNode 是红黑树的节点
type mapNode[K comparable, V any] struct {
	key    K
	value  V
	left   *mapNode[K, V]
	right  *mapNode[K, V]
	parent *mapNode[K, V]
	color  color
}

// NewTreeMap 创建一个新的TreeMap，使用默认比较器
func NewTreeMap[K comparable, V any]() *TreeMap[K, V] {
	return &TreeMap[K, V]{
		comparator: func(a, b K) int { return Compare(a, b) },
		root:       nil,
		size:       0,
	}
}

// NewTreeMapWithComparator 创建一个新的TreeMap，使用指定的比较器
func NewTreeMapWithComparator[K comparable, V any](comparator func(a, b K) int) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		comparator: comparator,
		root:       nil,
		size:       0,
	}
}

// isRedMap 检查节点是否为红色
func isRedMap[K comparable, V any](node *mapNode[K, V]) bool {
	if node == nil {
		return false
	}
	return node.color == red
}

// rotateLeft 左旋转
func (m *TreeMap[K, V]) rotateLeft(h *mapNode[K, V]) *mapNode[K, V] {
	x := h.right
	h.right = x.left
	if x.left != nil {
		x.left.parent = h
	}
	x.left = h
	x.color = h.color
	h.color = red
	x.parent = h.parent
	h.parent = x
	return x
}

// rotateRight 右旋转
func (m *TreeMap[K, V]) rotateRight(h *mapNode[K, V]) *mapNode[K, V] {
	x := h.left
	h.left = x.right
	if x.right != nil {
		x.right.parent = h
	}
	x.right = h
	x.color = h.color
	h.color = red
	x.parent = h.parent
	h.parent = x
	return x
}

// flipColors 颜色翻转
func (m *TreeMap[K, V]) flipColors(h *mapNode[K, V]) {
	h.color = !h.color
	if h.left != nil {
		h.left.color = !h.left.color
	}
	if h.right != nil {
		h.right.color = !h.right.color
	}
}

// put 插入或更新节点
func (m *TreeMap[K, V]) put(h *mapNode[K, V], key K, value V) (*mapNode[K, V], V, bool) {
	var oldValue V
	existed := false

	if h == nil {
		m.size++
		return &mapNode[K, V]{key: key, value: value, color: red}, oldValue, existed
	}

	cmp := m.comparator(key, h.key)

	if cmp < 0 {
		h.left, oldValue, existed = m.put(h.left, key, value)
		if h.left != nil {
			h.left.parent = h
		}
	} else if cmp > 0 {
		h.right, oldValue, existed = m.put(h.right, key, value)
		if h.right != nil {
			h.right.parent = h
		}
	} else {
		// 键已存在，更新值
		oldValue = h.value
		h.value = value
		existed = true
	}

	// 红黑树平衡调整
	if isRedMap(h.right) && !isRedMap(h.left) {
		h = m.rotateLeft(h)
	}
	if isRedMap(h.left) && isRedMap(h.left.left) {
		h = m.rotateRight(h)
	}
	if isRedMap(h.left) && isRedMap(h.right) {
		m.flipColors(h)
	}

	return h, oldValue, existed
}

// Put 将指定的值与此映射中的指定键关联
func (m *TreeMap[K, V]) Put(key K, value V) (V, bool) {
	var oldValue V
	var existed bool
	m.root, oldValue, existed = m.put(m.root, key, value)
	if m.root != nil {
		m.root.color = black
		m.root.parent = nil
	}
	return oldValue, existed
}

// findMin 查找最小节点
func findMinMap[K comparable, V any](h *mapNode[K, V]) *mapNode[K, V] {
	if h == nil {
		return nil
	}
	if h.left == nil {
		return h
	}
	return findMinMap(h.left)
}

// moveRedLeft 将左子节点或其兄弟节点变红
func (m *TreeMap[K, V]) moveRedLeft(h *mapNode[K, V]) *mapNode[K, V] {
	m.flipColors(h)
	if isRedMap(h.right.left) {
		h.right = m.rotateRight(h.right)
		h = m.rotateLeft(h)
		m.flipColors(h)
	}
	return h
}

// moveRedRight 将右子节点或其兄弟节点变红
func (m *TreeMap[K, V]) moveRedRight(h *mapNode[K, V]) *mapNode[K, V] {
	m.flipColors(h)
	if isRedMap(h.left.left) {
		h = m.rotateRight(h)
		m.flipColors(h)
	}
	return h
}

// balance 平衡节点
func (m *TreeMap[K, V]) balance(h *mapNode[K, V]) *mapNode[K, V] {
	if isRedMap(h.right) {
		h = m.rotateLeft(h)
	}
	if isRedMap(h.left) && isRedMap(h.left.left) {
		h = m.rotateRight(h)
	}
	if isRedMap(h.left) && isRedMap(h.right) {
		m.flipColors(h)
	}
	return h
}

// removeMin 移除最小节点
func (m *TreeMap[K, V]) removeMin(h *mapNode[K, V]) *mapNode[K, V] {
	if h.left == nil {
		// 只有在这是一个独立的删除操作时才减少 size
		// 当从 remove 方法调用时，size 会在 remove 方法中处理
		return nil
	}

	if !isRedMap(h.left) && !isRedMap(h.left.left) {
		h = m.moveRedLeft(h)
	}

	h.left = m.removeMin(h.left)
	if h.left != nil {
		h.left.parent = h
	}

	return m.balance(h)
}

// remove 移除节点
func (m *TreeMap[K, V]) remove(h *mapNode[K, V], key K) (*mapNode[K, V], V, bool) {
	var removedValue V
	removed := false

	if h == nil {
		return nil, removedValue, removed
	}

	cmp := m.comparator(key, h.key)

	if cmp < 0 {
		// 键在左子树中
		if h.left != nil {
			if !isRedMap(h.left) && !isRedMap(h.left.left) {
				h = m.moveRedLeft(h)
			}
			h.left, removedValue, removed = m.remove(h.left, key)
			if h.left != nil {
				h.left.parent = h
			}
		}
	} else if cmp > 0 {
		// 键在右子树中
		if isRedMap(h.left) {
			h = m.rotateRight(h)
			cmp = m.comparator(key, h.key) // 重新计算比较结果
		}
		if cmp > 0 && h.right != nil {
			if !isRedMap(h.right) && !isRedMap(h.right.left) {
				h = m.moveRedRight(h)
			}
			h.right, removedValue, removed = m.remove(h.right, key)
			if h.right != nil {
				h.right.parent = h
			}
		}
	} else {
		// 找到要删除的节点
		removedValue = h.value
		removed = true
		m.size--

		if h.right == nil {
			// 没有右子树，直接返回左子树
			return h.left, removedValue, removed
		} else if h.left == nil {
			// 没有左子树，直接返回右子树
			return h.right, removedValue, removed
		} else {
			// 有两个子树，用右子树的最小节点替换当前节点
			minNode := findMinMap(h.right)
			h.key = minNode.key
			h.value = minNode.value
			// 删除右子树的最小节点，但不减少 size（因为我们只是移动了节点）
			h.right = m.removeMinWithoutSizeChange(h.right)
			if h.right != nil {
				h.right.parent = h
			}
		}
	}

	if h != nil {
		h = m.balance(h)
	}
	return h, removedValue, removed
}

// removeMinWithoutSizeChange 移除最小节点但不改变 size
func (m *TreeMap[K, V]) removeMinWithoutSizeChange(h *mapNode[K, V]) *mapNode[K, V] {
	if h == nil {
		return nil
	}

	if h.left == nil {
		return h.right
	}

	if !isRedMap(h.left) && (h.left.left == nil || !isRedMap(h.left.left)) {
		h = m.moveRedLeft(h)
	}

	h.left = m.removeMinWithoutSizeChange(h.left)
	if h.left != nil {
		h.left.parent = h
	}

	return m.balance(h)
}

// Remove 如果存在，则从此映射中移除指定键的映射关系
func (m *TreeMap[K, V]) Remove(key K) (V, bool) {
	if m.root == nil {
		return *new(V), false
	}

	if !isRedMap(m.root.left) && !isRedMap(m.root.right) {
		m.root.color = red
	}

	var removedValue V
	var removed bool
	m.root, removedValue, removed = m.remove(m.root, key)
	if m.root != nil {
		m.root.color = black
	}

	return removedValue, removed
}

// find 查找节点
func (m *TreeMap[K, V]) find(h *mapNode[K, V], key K) *mapNode[K, V] {
	if h == nil {
		return nil
	}

	cmp := m.comparator(key, h.key)
	if cmp < 0 {
		return m.find(h.left, key)
	} else if cmp > 0 {
		return m.find(h.right, key)
	} else {
		return h
	}
}

// Get 返回指定键所映射的值
func (m *TreeMap[K, V]) Get(key K) (V, bool) {
	node := m.find(m.root, key)
	if node == nil {
		return *new(V), false
	}
	return node.value, true
}

// ContainsKey 如果此映射包含指定键的映射关系，则返回true
func (m *TreeMap[K, V]) ContainsKey(key K) bool {
	return m.find(m.root, key) != nil
}

// Size 返回此映射中的键-值映射关系数
func (m *TreeMap[K, V]) Size() int {
	return m.size
}

// IsEmpty 如果此映射不包含键-值映射关系，则返回true
func (m *TreeMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

// Clear 从此映射中移除所有映射关系
func (m *TreeMap[K, V]) Clear() {
	m.root = nil
	m.size = 0
}

// inOrderTraversalMap 中序遍历
func inOrderTraversalMap[K comparable, V any](node *mapNode[K, V], keys *[]K, values *[]V, entries *[]Pair[K, V]) {
	if node == nil {
		return
	}
	inOrderTraversalMap(node.left, keys, values, entries)
	if keys != nil {
		*keys = append(*keys, node.key)
	}
	if values != nil {
		*values = append(*values, node.value)
	}
	if entries != nil {
		*entries = append(*entries, NewPair(node.key, node.value))
	}
	inOrderTraversalMap(node.right, keys, values, entries)
}

// Keys 返回此映射中包含的键的集合视图（按顺序）
func (m *TreeMap[K, V]) Keys() []K {
	result := make([]K, 0, m.size)
	inOrderTraversalMap(m.root, &result, nil, nil)
	return result
}

// Values 返回此映射中包含的值的集合视图（按键的顺序）
func (m *TreeMap[K, V]) Values() []V {
	result := make([]V, 0, m.size)
	inOrderTraversalMap(m.root, nil, &result, nil)
	return result
}

// Entries 返回此映射中包含的映射关系的集合视图（按键的顺序）
func (m *TreeMap[K, V]) Entries() []Pair[K, V] {
	result := make([]Pair[K, V], 0, m.size)
	inOrderTraversalMap(m.root, nil, nil, &result)
	return result
}

// ForEach 对此映射中的每个条目执行给定的操作（按键的顺序）
func (m *TreeMap[K, V]) ForEach(f func(K, V)) {
	entries := m.Entries()
	for _, entry := range entries {
		f(entry.Key, entry.Value)
	}
}

// String 返回映射的字符串表示
func (m *TreeMap[K, V]) String() string {
	if m.IsEmpty() {
		return "{}"
	}

	entries := m.Entries()
	var sb strings.Builder
	sb.WriteString("{")

	for i, entry := range entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v=%v", entry.Key, entry.Value))
	}

	sb.WriteString("}")
	return sb.String()
}
