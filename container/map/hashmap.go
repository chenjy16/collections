package maps

import (
	"fmt"
	"strings"
	"sync"
)

// 使用已有的color类型作为节点颜色
// color类型在treeset.go中已定义

// 链表节点阈值，当链表长度超过此值时转换为红黑树
const treeifyThreshold = 8

// 红黑树节点阈值，当红黑树节点数小于此值时转换为链表
const untreeifyThreshold = 6

// 最小树化容量，当哈希表容量小于此值时，优先扩容而不是树化
const minTreeifyCapacity = 64

// 初始哈希表大小
const initialCapacity = 16

// 负载因子，当哈希表中的元素数量超过容量乘以负载因子时进行扩容
const loadFactor = 0.75

// LinkedHashMapNode 是链表/红黑树节点
type LinkedHashMapNode[K comparable, V any] struct {
	key   K
	value V
	hash  uint64

	// 链表指针
	next *LinkedHashMapNode[K, V]

	// 红黑树指针
	left   *LinkedHashMapNode[K, V]
	right  *LinkedHashMapNode[K, V]
	parent *LinkedHashMapNode[K, V]
	color  color

	// 标记节点是否为树节点
	isTreeNode bool
}

// LinkedHashMap 是一个基于链地址法和红黑树的Map实现
type LinkedHashMap[K comparable, V any] struct {
	table     []*LinkedHashMapNode[K, V] // 哈希桶数组
	size      int                        // 元素数量
	threshold int                        // 扩容阈值
	mutex     sync.RWMutex               // 读写锁，保证线程安全
}

// NewLinkedHashMap 创建一个新的LinkedHashMap
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	capacity := initialCapacity
	return &LinkedHashMap[K, V]{
		table:     make([]*LinkedHashMapNode[K, V], capacity),
		size:      0,
		threshold: int(float64(capacity) * loadFactor),
	}
}

// NewLinkedHashMapWithCapacity 创建一个具有指定初始容量的LinkedHashMap
func NewLinkedHashMapWithCapacity[K comparable, V any](capacity int) *LinkedHashMap[K, V] {
	if capacity < initialCapacity {
		capacity = initialCapacity
	} else {
		// 确保容量是2的幂
		capacity = tableSizeFor(capacity)
	}

	return &LinkedHashMap[K, V]{
		table:     make([]*LinkedHashMapNode[K, V], capacity),
		size:      0,
		threshold: int(float64(capacity) * loadFactor),
	}
}

// tableSizeFor 返回大于等于cap的最小2的幂
func tableSizeFor(cap int) int {
	n := cap - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}

// hash 计算键的哈希值
func (m *LinkedHashMap[K, V]) hash(key K) uint64 {
	return Hash(key)
}

// Put 将指定的值与此映射中的指定键关联
func (m *LinkedHashMap[K, V]) Put(key K, value V) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var oldValue V
	existed := false

	hashValue := m.hash(key)
	index := int(hashValue % uint64(len(m.table)))

	// 如果桶为空，创建新节点
	if m.table[index] == nil {
		m.table[index] = &LinkedHashMapNode[K, V]{
			key:   key,
			value: value,
			hash:  hashValue,
		}
		m.size++

		// 检查是否需要扩容
		m.checkResize()

		return oldValue, existed
	}

	// 如果是树节点，使用树查找和插入
	if m.table[index].isTreeNode {
		return m.putTreeVal(index, key, value, hashValue)
	}

	// 链表查找和插入
	p := m.table[index]
	var prev *LinkedHashMapNode[K, V]
	count := 0

	// 遍历链表
	for p != nil {
		count++

		// 如果找到相同的键，更新值
		if p.hash == hashValue && Equal(p.key, key) {
			oldValue = p.value
			p.value = value
			return oldValue, true
		}

		prev = p
		p = p.next
	}

	// 没有找到相同的键，添加新节点到链表尾部
	newNode := &LinkedHashMapNode[K, V]{
		key:   key,
		value: value,
		hash:  hashValue,
	}
	prev.next = newNode
	m.size++

	// 检查是否需要将链表转换为红黑树
	if count >= treeifyThreshold-1 {
		m.treeifyBin(index)
	}

	// 检查是否需要扩容
	m.checkResize()

	return oldValue, existed
}

// putTreeVal 在红黑树中插入或更新节点
func (m *LinkedHashMap[K, V]) putTreeVal(index int, key K, value V, hash uint64) (V, bool) {
	var oldValue V
	existed := false

	root := m.table[index]
	p := root

	// 树查找
	for p != nil {
		cmp := 0
		if p.hash > hash {
			cmp = -1
		} else if p.hash < hash {
			cmp = 1
		} else if Equal(key, p.key) {
			// 找到相同的键，更新值
			oldValue = p.value
			p.value = value
			return oldValue, true
		} else {
			// 哈希相同但键不同，使用键的比较
			cmp = Compare(key, p.key)
		}

		// 根据比较结果决定向左还是向右
		if cmp < 0 {
			if p.left == nil {
				// 插入到左子节点
				p.left = &LinkedHashMapNode[K, V]{
					key:        key,
					value:      value,
					hash:       hash,
					isTreeNode: true,
					parent:     p,
					color:      red,
				}
				m.size++
				m.balanceInsertion(root, p.left)
				return oldValue, existed
			}
			p = p.left
		} else {
			if p.right == nil {
				// 插入到右子节点
				p.right = &LinkedHashMapNode[K, V]{
					key:        key,
					value:      value,
					hash:       hash,
					isTreeNode: true,
					parent:     p,
					color:      red,
				}
				m.size++
				m.balanceInsertion(root, p.right)
				return oldValue, existed
			}
			p = p.right
		}
	}

	// 如果树为空，创建根节点
	m.table[index] = &LinkedHashMapNode[K, V]{
		key:        key,
		value:      value,
		hash:       hash,
		isTreeNode: true,
		color:      black, // 根节点为黑色
	}
	m.size++

	return oldValue, existed
}

// balanceInsertion 插入后平衡红黑树
func (m *LinkedHashMap[K, V]) balanceInsertion(root *LinkedHashMapNode[K, V], x *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	// 红黑树平衡调整
	x.color = red

	for x != nil && x != root && x.parent.color == red {
		if parentOf(x) == leftOf(parentOf(parentOf(x))) {
			y := rightOf(parentOf(parentOf(x)))
			if colorOf(y) == red {
				setColor(parentOf(x), black)
				setColor(y, black)
				setColor(parentOf(parentOf(x)), red)
				x = parentOf(parentOf(x))
			} else {
				if x == rightOf(parentOf(x)) {
					x = parentOf(x)
					root = m.rotateLeft(root, x)
				}
				setColor(parentOf(x), black)
				setColor(parentOf(parentOf(x)), red)
				root = m.rotateRight(root, parentOf(parentOf(x)))
			}
		} else {
			y := leftOf(parentOf(parentOf(x)))
			if colorOf(y) == red {
				setColor(parentOf(x), black)
				setColor(y, black)
				setColor(parentOf(parentOf(x)), red)
				x = parentOf(parentOf(x))
			} else {
				if x == leftOf(parentOf(x)) {
					x = parentOf(x)
					root = m.rotateRight(root, x)
				}
				setColor(parentOf(x), black)
				setColor(parentOf(parentOf(x)), red)
				root = m.rotateLeft(root, parentOf(parentOf(x)))
			}
		}
	}

	root.color = black
	return root
}

// rotateLeft 左旋转
func (m *LinkedHashMap[K, V]) rotateLeft(root *LinkedHashMapNode[K, V], p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p != nil {
		r := p.right
		p.right = r.left
		if r.left != nil {
			r.left.parent = p
		}
		r.parent = p.parent

		if p.parent == nil {
			root = r
		} else if p == p.parent.left {
			p.parent.left = r
		} else {
			p.parent.right = r
		}

		r.left = p
		p.parent = r
	}
	return root
}

// rotateRight 右旋转
func (m *LinkedHashMap[K, V]) rotateRight(root *LinkedHashMapNode[K, V], p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p != nil {
		l := p.left
		p.left = l.right
		if l.right != nil {
			l.right.parent = p
		}
		l.parent = p.parent

		if p.parent == nil {
			root = l
		} else if p == p.parent.right {
			p.parent.right = l
		} else {
			p.parent.left = l
		}

		l.right = p
		p.parent = l
	}
	return root
}

// parentOf 获取节点的父节点
func parentOf[K comparable, V any](p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p == nil {
		return nil
	}
	return p.parent
}

// leftOf 获取节点的左子节点
func leftOf[K comparable, V any](p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p == nil {
		return nil
	}
	return p.left
}

// rightOf 获取节点的右子节点
func rightOf[K comparable, V any](p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p == nil {
		return nil
	}
	return p.right
}

// colorOf 获取节点的颜色
func colorOf[K comparable, V any](p *LinkedHashMapNode[K, V]) color {
	if p == nil {
		return black
	}
	return p.color
}

// setColor 设置节点的颜色
func setColor[K comparable, V any](p *LinkedHashMapNode[K, V], c color) {
	if p != nil {
		p.color = c
	}
}

// treeifyBin 将指定索引处的链表转换为红黑树
func (m *LinkedHashMap[K, V]) treeifyBin(index int) {
	// 如果哈希表容量小于最小树化容量，优先扩容
	if len(m.table) < minTreeifyCapacity {
		m.resize()
		return
	}

	// 将链表转换为红黑树
	root := m.table[index]
	if root == nil {
		return
	}

	// 标记所有节点为树节点
	p := root
	for p != nil {
		p.isTreeNode = true
		p = p.next
	}

	// 构建红黑树
	root = m.buildTree(root)
	m.table[index] = root
}

// buildTree 从链表构建红黑树
func (m *LinkedHashMap[K, V]) buildTree(head *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	var root *LinkedHashMapNode[K, V]

	// 遍历链表，将每个节点插入到红黑树中
	p := head
	for p != nil {
		next := p.next
		p.left = nil
		p.right = nil

		if root == nil {
			p.parent = nil
			p.color = black
			root = p
		} else {
			k := p.key
			h := p.hash
			dir := 0

			// 查找插入位置
			cur := root
			for {
				ph := cur.hash
				pk := cur.key
				if h < ph {
					dir = -1
				} else if h > ph {
					dir = 1
				} else if Equal(k, pk) {
					dir = 0
				} else {
					dir = Compare(k, pk)
				}

				if dir < 0 {
					if cur.left == nil {
						cur.left = p
						break
					}
					cur = cur.left
				} else {
					if cur.right == nil {
						cur.right = p
						break
					}
					cur = cur.right
				}
			}

			p.parent = cur
			p.color = red

			// 平衡红黑树
			root = m.balanceInsertion(root, p)
		}

		p = next
	}

	return root
}

// checkResize 检查是否需要扩容
func (m *LinkedHashMap[K, V]) checkResize() {
	if m.size > m.threshold {
		m.resize()
	}
}

// resize 扩容哈希表
func (m *LinkedHashMap[K, V]) resize() {
	oldCap := len(m.table)
	oldTab := m.table

	// 计算新容量
	newCap := oldCap * 2
	if newCap < initialCapacity {
		newCap = initialCapacity
	}

	// 创建新表
	newTab := make([]*LinkedHashMapNode[K, V], newCap)
	m.table = newTab
	m.threshold = int(float64(newCap) * loadFactor)

	// 如果旧表为空，直接返回
	if oldCap == 0 {
		return
	}

	// 将旧表中的元素重新分配到新表中
	for i := 0; i < oldCap; i++ {
		e := oldTab[i]
		if e == nil {
			continue
		}

		// 清空旧表引用
		oldTab[i] = nil

		// 如果是单个节点
		if e.next == nil {
			newIdx := int(e.hash % uint64(newCap))
			newTab[newIdx] = e
			continue
		}

		// 如果是树节点
		if e.isTreeNode {
			m.splitTreeBin(newTab, e, i, oldCap)
			continue
		}

		// 如果是链表，拆分为两个链表
		// 一个链表放在原位置，一个链表放在原位置+oldCap
		var loHead, loTail, hiHead, hiTail *LinkedHashMapNode[K, V]

		// 遍历链表
		for e != nil {
			next := e.next

			// 根据哈希值决定放在哪个链表
			if (e.hash & uint64(oldCap)) == 0 {
				// 放在原位置
				if loTail == nil {
					loHead = e
				} else {
					loTail.next = e
				}
				loTail = e
			} else {
				// 放在原位置+oldCap
				if hiTail == nil {
					hiHead = e
				} else {
					hiTail.next = e
				}
				hiTail = e
			}

			e = next
		}

		// 更新链表引用
		if loTail != nil {
			loTail.next = nil
			newTab[i] = loHead
		}

		if hiTail != nil {
			hiTail.next = nil
			newTab[i+oldCap] = hiHead
		}
	}
}

// splitTreeBin 拆分树节点
func (m *LinkedHashMap[K, V]) splitTreeBin(newTab []*LinkedHashMapNode[K, V], root *LinkedHashMapNode[K, V], index, oldCap int) {
	// 将树节点转换回链表
	if len(newTab) <= untreeifyThreshold {
		var head, tail *LinkedHashMapNode[K, V]

		// 遍历树，构建链表
		m.treeToList(root, &head, &tail)

		// 拆分链表
		var loHead, loTail, hiHead, hiTail *LinkedHashMapNode[K, V]
		p := head

		for p != nil {
			next := p.next
			p.left = nil
			p.right = nil
			p.parent = nil
			p.isTreeNode = false

			// 根据哈希值决定放在哪个链表
			if (p.hash & uint64(oldCap)) == 0 {
				// 放在原位置
				if loTail == nil {
					loHead = p
				} else {
					loTail.next = p
				}
				loTail = p
			} else {
				// 放在原位置+oldCap
				if hiTail == nil {
					hiHead = p
				} else {
					hiTail.next = p
				}
				hiTail = p
			}

			p = next
		}

		// 更新链表引用
		if loTail != nil {
			loTail.next = nil
			newTab[index] = loHead
		}

		if hiTail != nil {
			hiTail.next = nil
			newTab[index+oldCap] = hiHead
		}
	} else {
		// 拆分树
		var loTree, hiTree *LinkedHashMapNode[K, V]

		// 遍历树，构建两棵新树
		m.splitTree(root, &loTree, &hiTree, oldCap)

		// 更新树引用
		if loTree != nil {
			newTab[index] = loTree
		}

		if hiTree != nil {
			newTab[index+oldCap] = hiTree
		}
	}
}

// treeToList 将树转换为链表
func (m *LinkedHashMap[K, V]) treeToList(root *LinkedHashMapNode[K, V], head, tail **LinkedHashMapNode[K, V]) {
	// 中序遍历树，构建链表
	if root == nil {
		return
	}

	// 递归左子树
	m.treeToList(root.left, head, tail)

	// 处理当前节点
	root.left = nil
	root.right = nil
	root.parent = nil
	root.isTreeNode = false

	if *tail == nil {
		*head = root
	} else {
		(*tail).next = root
	}
	*tail = root

	// 递归右子树
	m.treeToList(root.right, head, tail)
}

// splitTree 拆分树
func (m *LinkedHashMap[K, V]) splitTree(root *LinkedHashMapNode[K, V], loTree, hiTree **LinkedHashMapNode[K, V], oldCap int) {
	if root == nil {
		return
	}

	// 先保存子节点引用
	left := root.left
	right := root.right

	// 清理当前节点的引用
	root.left = nil
	root.right = nil
	root.parent = nil

	// 根据哈希值决定放在哪棵树
	if (root.hash & uint64(oldCap)) == 0 {
		// 放在原位置
		if *loTree == nil {
			*loTree = root
			root.color = black
		} else {
			// 简单插入到loTree，不进行平衡
			m.insertNodeSimple(loTree, root)
		}
	} else {
		// 放在原位置+oldCap
		if *hiTree == nil {
			*hiTree = root
			root.color = black
		} else {
			// 简单插入到hiTree，不进行平衡
			m.insertNodeSimple(hiTree, root)
		}
	}

	// 递归处理子树
	m.splitTree(left, loTree, hiTree, oldCap)
	m.splitTree(right, loTree, hiTree, oldCap)
}

// insertNodeSimple 简单插入节点到树中，不进行平衡
func (m *LinkedHashMap[K, V]) insertNodeSimple(root **LinkedHashMapNode[K, V], node *LinkedHashMapNode[K, V]) {
	p := *root
	for {
		cmp := 0
		if node.hash < p.hash {
			cmp = -1
		} else if node.hash > p.hash {
			cmp = 1
		} else {
			cmp = Compare(node.key, p.key)
		}

		if cmp < 0 {
			if p.left == nil {
				p.left = node
				node.parent = p
				node.color = red
				break
			}
			p = p.left
		} else {
			if p.right == nil {
				p.right = node
				node.parent = p
				node.color = red
				break
			}
			p = p.right
		}
	}
}

// Get 返回指定键所映射的值
func (m *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	hashValue := m.hash(key)
	index := int(hashValue % uint64(len(m.table)))

	// 如果桶为空，返回零值
	if m.table[index] == nil {
		return *new(V), false
	}

	// 如果是树节点，使用树查找
	if m.table[index].isTreeNode {
		return m.getTreeVal(m.table[index], key, hashValue)
	}

	// 链表查找
	p := m.table[index]
	for p != nil {
		if p.hash == hashValue && Equal(p.key, key) {
			return p.value, true
		}
		p = p.next
	}

	return *new(V), false
}

// getTreeVal 在红黑树中查找节点
func (m *LinkedHashMap[K, V]) getTreeVal(root *LinkedHashMapNode[K, V], key K, hash uint64) (V, bool) {
	p := root

	// 树查找
	for p != nil {
		cmp := 0
		if p.hash > hash {
			cmp = -1
		} else if p.hash < hash {
			cmp = 1
		} else if Equal(key, p.key) {
			// 找到相同的键，返回值
			return p.value, true
		} else {
			// 哈希相同但键不同，使用键的比较
			cmp = Compare(key, p.key)
		}

		// 根据比较结果决定向左还是向右
		if cmp < 0 {
			p = p.left
		} else {
			p = p.right
		}
	}

	return *new(V), false
}

// Remove 如果存在，则从此映射中移除指定键的映射关系
func (m *LinkedHashMap[K, V]) Remove(key K) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	hashValue := m.hash(key)
	index := int(hashValue % uint64(len(m.table)))

	// 如果桶为空，返回零值
	if m.table[index] == nil {
		return *new(V), false
	}

	// 如果是树节点，使用树删除
	if m.table[index].isTreeNode {
		return m.removeTreeNode(index, key, hashValue)
	}

	// 链表删除
	p := m.table[index]
	var prev *LinkedHashMapNode[K, V]

	for p != nil {
		if p.hash == hashValue && Equal(p.key, key) {
			// 找到要删除的节点
			oldValue := p.value

			if prev == nil {
				// 删除链表头部
				m.table[index] = p.next
			} else {
				// 删除链表中间或尾部
				prev.next = p.next
			}

			m.size--
			return oldValue, true
		}

		prev = p
		p = p.next
	}

	return *new(V), false
}

// removeTreeNode 从红黑树中删除节点
func (m *LinkedHashMap[K, V]) removeTreeNode(index int, key K, hash uint64) (V, bool) {
	root := m.table[index]
	p := root

	// 查找要删除的节点
	for p != nil {
		cmp := 0
		if p.hash > hash {
			cmp = -1
		} else if p.hash < hash {
			cmp = 1
		} else if Equal(key, p.key) {
			// 找到要删除的节点
			break
		} else {
			// 哈希相同但键不同，使用键的比较
			cmp = Compare(key, p.key)
		}

		// 根据比较结果决定向左还是向右
		if cmp < 0 {
			p = p.left
		} else {
			p = p.right
		}
	}

	// 如果没有找到节点，返回零值
	if p == nil {
		return *new(V), false
	}

	oldValue := p.value

	// 删除节点
	if p.left != nil && p.right != nil {
		// 如果有两个子节点，找到后继节点
		s := p.right
		for s.left != nil {
			s = s.left
		}

		// 用后继节点的值替换当前节点的值
		p.key = s.key
		p.value = s.value
		p.hash = s.hash

		// 删除后继节点
		if s.parent == p {
			p.right = s.right
			if s.right != nil {
				s.right.parent = p
			}
		} else {
			s.parent.left = s.right
			if s.right != nil {
				s.right.parent = s.parent
			}
		}
	} else {
		// 如果最多有一个子节点
		replacement := p.left
		if p.left == nil {
			replacement = p.right
		}

		// 用子节点替换当前节点
		if p.parent == nil {
			// 如果是根节点
			m.table[index] = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		if replacement != nil {
			replacement.parent = p.parent
		}
	}

	// 平衡红黑树
	if p.color == black {
		m.balanceDeletion(root, p)
	}

	m.size--

	// 如果树太小，转换回链表
	if m.size <= untreeifyThreshold {
		m.untreeify(index)
	}

	return oldValue, true
}

// balanceDeletion 删除后平衡红黑树
func (m *LinkedHashMap[K, V]) balanceDeletion(root *LinkedHashMapNode[K, V], x *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	// 红黑树删除平衡调整
	for x != root && colorOf(x) == black {
		if x == leftOf(parentOf(x)) {
			sib := rightOf(parentOf(x))

			if colorOf(sib) == red {
				setColor(sib, black)
				setColor(parentOf(x), red)
				root = m.rotateLeft(root, parentOf(x))
				sib = rightOf(parentOf(x))
			}

			if colorOf(leftOf(sib)) == black && colorOf(rightOf(sib)) == black {
				setColor(sib, red)
				x = parentOf(x)
			} else {
				if colorOf(rightOf(sib)) == black {
					setColor(leftOf(sib), black)
					setColor(sib, red)
					root = m.rotateRight(root, sib)
					sib = rightOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), black)
				setColor(rightOf(sib), black)
				root = m.rotateLeft(root, parentOf(x))
				x = root
			}
		} else {
			sib := leftOf(parentOf(x))

			if colorOf(sib) == red {
				setColor(sib, black)
				setColor(parentOf(x), red)
				root = m.rotateRight(root, parentOf(x))
				sib = leftOf(parentOf(x))
			}

			if colorOf(rightOf(sib)) == black && colorOf(leftOf(sib)) == black {
				setColor(sib, red)
				x = parentOf(x)
			} else {
				if colorOf(leftOf(sib)) == black {
					setColor(rightOf(sib), black)
					setColor(sib, red)
					root = m.rotateLeft(root, sib)
					sib = leftOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), black)
				setColor(leftOf(sib), black)
				root = m.rotateRight(root, parentOf(x))
				x = root
			}
		}
	}

	setColor(x, black)
	return root
}

// untreeify 将指定索引处的红黑树转换为链表
func (m *LinkedHashMap[K, V]) untreeify(index int) {
	root := m.table[index]
	if root == nil || !root.isTreeNode {
		return
	}

	// 将树转换为链表
	var head, tail *LinkedHashMapNode[K, V]
	m.treeToList(root, &head, &tail)

	// 更新链表引用
	m.table[index] = head
}

// ContainsKey 如果此映射包含指定键的映射关系，则返回true
func (m *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, found := m.Get(key)
	return found
}

// ContainsValue 如果此映射将一个或多个键映射到指定值，则返回true
func (m *LinkedHashMap[K, V]) ContainsValue(value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.traverseAllWithEarlyExit(func(node *LinkedHashMapNode[K, V]) bool {
		return Equal(node.value, value)
	})
}

// Size 返回此映射中的键-值映射关系数
func (m *LinkedHashMap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.size
}

// IsEmpty 如果此映射不包含键-值映射关系，则返回true
func (m *LinkedHashMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear 从此映射中移除所有映射关系
func (m *LinkedHashMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 清空哈希表，并帮助GC回收内存
	for i := range m.table {
		if m.table[i] != nil {
			m.clearNode(m.table[i])
			m.table[i] = nil
		}
	}

	m.size = 0
}

// clearNode 递归清理节点以帮助GC
func (m *LinkedHashMap[K, V]) clearNode(node *LinkedHashMapNode[K, V]) {
	if node == nil {
		return
	}

	// 如果是树节点，递归清理子节点
	if node.isTreeNode {
		m.clearNode(node.left)
		m.clearNode(node.right)
		node.left = nil
		node.right = nil
		node.parent = nil
	}

	// 清理链表节点
	for node.next != nil {
		next := node.next
		node.next = nil
		node = next
	}
}

// Keys 返回此映射中包含的键的集合视图
func (m *LinkedHashMap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]K, 0, m.size)
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		keys = append(keys, node.key)
	})

	return keys
}

// Values 返回此映射中包含的值的集合视图
func (m *LinkedHashMap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, m.size)
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		values = append(values, node.value)
	})

	return values
}

// inOrderTraversal 中序遍历红黑树
func (m *LinkedHashMap[K, V]) inOrderTraversal(root *LinkedHashMapNode[K, V], f func(*LinkedHashMapNode[K, V])) {
	if root == nil {
		return
	}

	m.inOrderTraversal(root.left, f)
	f(root)
	m.inOrderTraversal(root.right, f)
}

// traverseAll 遍历所有节点，包括链表和树节点
func (m *LinkedHashMap[K, V]) traverseAll(f func(*LinkedHashMapNode[K, V])) {
	// 遍历哈希表
	for _, e := range m.table {
		p := e
		for p != nil {
			f(p)
			if p.isTreeNode {
				// 如果是树节点，使用中序遍历其他节点
				m.inOrderTraversal(p, func(node *LinkedHashMapNode[K, V]) {
					if node != p { // 避免重复处理根节点
						f(node)
					}
				})
				break
			}
			p = p.next
		}
	}
}

// traverseAllWithEarlyExit 遍历所有节点，支持提前退出
func (m *LinkedHashMap[K, V]) traverseAllWithEarlyExit(f func(*LinkedHashMapNode[K, V]) bool) bool {
	// 遍历哈希表
	for _, e := range m.table {
		p := e
		for p != nil {
			if f(p) {
				return true
			}
			if p.isTreeNode {
				// 如果是树节点，使用中序遍历其他节点
				found := false
				m.inOrderTraversalWithEarlyExit(p, func(node *LinkedHashMapNode[K, V]) bool {
					if node != p { // 避免重复处理根节点
						return f(node)
					}
					return false
				}, &found)
				if found {
					return true
				}
				break
			}
			p = p.next
		}
	}
	return false
}

// inOrderTraversalWithEarlyExit 中序遍历红黑树，支持提前退出
func (m *LinkedHashMap[K, V]) inOrderTraversalWithEarlyExit(root *LinkedHashMapNode[K, V], f func(*LinkedHashMapNode[K, V]) bool, found *bool) {
	if root == nil || *found {
		return
	}

	m.inOrderTraversalWithEarlyExit(root.left, f, found)
	if !*found && f(root) {
		*found = true
		return
	}
	m.inOrderTraversalWithEarlyExit(root.right, f, found)
}

// ForEach 对此映射中的每个条目执行给定的操作
func (m *LinkedHashMap[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		f(node.key, node.value)
	})
}

// PutAll 将指定映射中的所有映射关系复制到此映射
func (m *LinkedHashMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(k K, v V) {
		m.Put(k, v)
	})
}

// String 返回映射的字符串表示
func (m *LinkedHashMap[K, V]) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.size == 0 {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{")

	first := true
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v=%v", node.key, node.value))
		first = false
	})

	sb.WriteString("}")
	return sb.String()
}

// Entries 返回此映射中包含的映射关系的集合视图
func (m *LinkedHashMap[K, V]) Entries() []Pair[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entries := make([]Pair[K, V], 0, m.size)
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		entries = append(entries, NewPair(node.key, node.value))
	})

	return entries
}
