package set

import (
	"fmt"
	"strings"
)

// Compare 比较两个可比较类型的值
// 注意：此函数仅支持 int 类型，对于其他类型请使用自定义比较器
func Compare[E any](a, b E) int {
	// 尝试转换为 int 类型
	cmpA, okA := any(a).(int)
	cmpB, okB := any(b).(int)
	if okA && okB {
		if cmpA == cmpB {
			return 0
		}
		if cmpA < cmpB {
			return -1
		}
		return 1
	}

	// 尝试转换为 string 类型
	strA, okA := any(a).(string)
	strB, okB := any(b).(string)
	if okA && okB {
		if strA == strB {
			return 0
		}
		if strA < strB {
			return -1
		}
		return 1
	}

	// 对于不支持的类型，返回 0（相等）
	// 这样可以避免 panic，但建议使用自定义比较器
	return 0
}

// TreeSet 是一个基于红黑树的有序Set实现
type TreeSet[E comparable] struct {
	comparator func(a, b E) int // 比较器函数
	root       *treeNode[E]     // 根节点
	size       int              // 元素数量
}

// 红黑树节点颜色
type color bool

const (
	red   color = true
	black color = false
)

// treeNode 是红黑树的节点
type treeNode[E comparable] struct {
	value  E
	left   *treeNode[E]
	right  *treeNode[E]
	parent *treeNode[E]
	color  color
}

// NewTreeSet 创建一个新的TreeSet，使用默认比较器
func NewTreeSet[E comparable]() *TreeSet[E] {
	return &TreeSet[E]{
		comparator: func(a, b E) int { return Compare(a, b) },
		root:       nil,
		size:       0,
	}
}

// NewTreeSetWithComparator 创建一个新的TreeSet，使用指定的比较器
func NewTreeSetWithComparator[E comparable](comparator func(a, b E) int) *TreeSet[E] {
	return &TreeSet[E]{
		comparator: comparator,
		root:       nil,
		size:       0,
	}
}

// NewTreeSetFromSlice 从切片创建一个新的TreeSet
func NewTreeSetFromSlice[E comparable](elements []E) *TreeSet[E] {
	result := NewTreeSet[E]()
	for _, e := range elements {
		result.Add(e)
	}
	return result
}

// isRed 检查节点是否为红色
func isRed[E comparable](node *treeNode[E]) bool {
	if node == nil {
		return false
	}
	return node.color == red
}

// rotateLeft 左旋转
func (set *TreeSet[E]) rotateLeft(h *treeNode[E]) *treeNode[E] {
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
func (set *TreeSet[E]) rotateRight(h *treeNode[E]) *treeNode[E] {
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
func (set *TreeSet[E]) flipColors(h *treeNode[E]) {
	h.color = !h.color
	if h.left != nil {
		h.left.color = !h.left.color
	}
	if h.right != nil {
		h.right.color = !h.right.color
	}
}

// put 插入节点
func (set *TreeSet[E]) put(h *treeNode[E], value E) (*treeNode[E], bool) {
	if h == nil {
		set.size++
		return &treeNode[E]{value: value, color: red}, true
	}

	added := false
	cmp := set.comparator(value, h.value)

	if cmp < 0 {
		h.left, added = set.put(h.left, value)
		if h.left != nil {
			h.left.parent = h
		}
	} else if cmp > 0 {
		h.right, added = set.put(h.right, value)
		if h.right != nil {
			h.right.parent = h
		}
	} else {
		// 值已存在，不添加
		return h, false
	}

	// 红黑树平衡调整
	if isRed(h.right) && !isRed(h.left) {
		h = set.rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = set.rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		set.flipColors(h)
	}

	return h, added
}

// Add 添加一个元素到集合中
func (set *TreeSet[E]) Add(e E) bool {
	var added bool
	set.root, added = set.put(set.root, e)
	if set.root != nil {
		set.root.color = black
		set.root.parent = nil
	}
	return added
}

// findMin 查找最小节点
func findMin[E comparable](h *treeNode[E]) *treeNode[E] {
	if h == nil {
		return nil
	}
	if h.left == nil {
		return h
	}
	return findMin(h.left)
}

// moveRedLeft 将左子节点或其兄弟节点变红
func (set *TreeSet[E]) moveRedLeft(h *treeNode[E]) *treeNode[E] {
	set.flipColors(h)
	if isRed(h.right.left) {
		h.right = set.rotateRight(h.right)
		h = set.rotateLeft(h)
		set.flipColors(h)
	}
	return h
}

// moveRedRight 将右子节点或其兄弟节点变红
func (set *TreeSet[E]) moveRedRight(h *treeNode[E]) *treeNode[E] {
	set.flipColors(h)
	if isRed(h.left.left) {
		h = set.rotateRight(h)
		set.flipColors(h)
	}
	return h
}

// balance 平衡节点
func (set *TreeSet[E]) balance(h *treeNode[E]) *treeNode[E] {
	if isRed(h.right) {
		h = set.rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = set.rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		set.flipColors(h)
	}
	return h
}

// removeMin 移除最小节点
func (set *TreeSet[E]) removeMin(h *treeNode[E]) *treeNode[E] {
	if h.left == nil {
		return nil
	}

	if !isRed(h.left) && !isRed(h.left.left) {
		h = set.moveRedLeft(h)
	}

	h.left = set.removeMin(h.left)
	if h.left != nil {
		h.left.parent = h
	}

	return set.balance(h)
}

// remove 移除节点
func (set *TreeSet[E]) remove(h *treeNode[E], value E) (*treeNode[E], bool) {
	if h == nil {
		return nil, false
	}

	removed := false
	cmp := set.comparator(value, h.value)

	if cmp < 0 {
		if h.left == nil {
			return h, false
		}
		if !isRed(h.left) && !isRed(h.left.left) {
			h = set.moveRedLeft(h)
		}
		h.left, removed = set.remove(h.left, value)
		if h.left != nil {
			h.left.parent = h
		}
	} else {
		if isRed(h.left) {
			h = set.rotateRight(h)
		}
		if cmp == 0 && h.right == nil {
			set.size--
			return nil, true
		}
		if !isRed(h.right) && h.right != nil && !isRed(h.right.left) {
			h = set.moveRedRight(h)
		}
		if cmp == 0 {
			minNode := findMin(h.right)
			h.value = minNode.value
			h.right = set.removeMin(h.right)
			if h.right != nil {
				h.right.parent = h
			}
			removed = true
			set.size--
		} else {
			h.right, removed = set.remove(h.right, value)
			if h.right != nil {
				h.right.parent = h
			}
		}
	}

	return set.balance(h), removed
}

// Remove 从集合中移除指定元素
func (set *TreeSet[E]) Remove(e E) bool {
	if !set.Contains(e) {
		return false
	}

	if !isRed(set.root.left) && !isRed(set.root.right) {
		set.root.color = red
	}

	var removed bool
	set.root, removed = set.remove(set.root, e)
	if set.root != nil {
		set.root.color = black
	}

	return removed
}

// find 查找节点
func (set *TreeSet[E]) find(h *treeNode[E], value E) *treeNode[E] {
	if h == nil {
		return nil
	}

	cmp := set.comparator(value, h.value)
	if cmp < 0 {
		return set.find(h.left, value)
	} else if cmp > 0 {
		return set.find(h.right, value)
	} else {
		return h
	}
}

// Contains 检查集合是否包含指定元素
func (set *TreeSet[E]) Contains(e E) bool {
	return set.find(set.root, e) != nil
}

// Size 返回集合中的元素数量
func (set *TreeSet[E]) Size() int {
	return set.size
}

// IsEmpty 检查集合是否为空
func (set *TreeSet[E]) IsEmpty() bool {
	return set.size == 0
}

// Clear 清空集合
func (set *TreeSet[E]) Clear() {
	set.root = nil
	set.size = 0
}

// inOrderTraversal 中序遍历
func inOrderTraversal[E comparable](node *treeNode[E], result *[]E) {
	if node == nil {
		return
	}
	inOrderTraversal(node.left, result)
	*result = append(*result, node.value)
	inOrderTraversal(node.right, result)
}

// ToSlice 返回包含集合所有元素的有序切片
func (set *TreeSet[E]) ToSlice() []E {
	result := make([]E, 0, set.size)
	inOrderTraversal(set.root, &result)
	return result
}

// ForEach 对集合中的每个元素执行给定的操作（按顺序）
func (set *TreeSet[E]) ForEach(f func(E)) {
	elements := set.ToSlice()
	for _, e := range elements {
		f(e)
	}
}

// AddAll 将指定集合中的所有元素添加到此集合
func (set *TreeSet[E]) AddAll(other Set[E]) bool {
	modified := false
	other.ForEach(func(e E) {
		if set.Add(e) {
			modified = true
		}
	})
	return modified
}

// RemoveAll 移除此集合中也包含在指定集合中的所有元素
func (set *TreeSet[E]) RemoveAll(other Set[E]) bool {
	modified := false
	other.ForEach(func(e E) {
		if set.Remove(e) {
			modified = true
		}
	})
	return modified
}

// RetainAll 仅保留此集合中包含在指定集合中的元素
func (set *TreeSet[E]) RetainAll(other Set[E]) bool {
	modified := false
	toRemove := make([]E, 0)

	set.ForEach(func(e E) {
		if !other.Contains(e) {
			toRemove = append(toRemove, e)
			modified = true
		}
	})

	for _, e := range toRemove {
		set.Remove(e)
	}

	return modified
}

// ContainsAll 如果此集合包含指定集合中的所有元素，返回true
func (set *TreeSet[E]) ContainsAll(other Set[E]) bool {
	containsAll := true
	other.ForEach(func(e E) {
		if !set.Contains(e) {
			containsAll = false
		}
	})
	return containsAll
}

// String 返回集合的字符串表示
func (set *TreeSet[E]) String() string {
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
