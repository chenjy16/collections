package maps

import (
	"fmt"
	"strings"
	"sync"
)

// Use existing color type as node color
// color type is already defined in treeset.go

// Linked list node threshold, convert to red-black tree when list length exceeds this value
const treeifyThreshold = 8

// Red-black tree node threshold, convert to linked list when tree node count is less than this value
const untreeifyThreshold = 6

// Minimum treeify capacity, when hash table capacity is less than this value, prioritize expansion over treeification
const minTreeifyCapacity = 64

// Initial hash table size
const initialCapacity = 16

// Load factor, expand when element count exceeds capacity multiplied by load factor
const loadFactor = 0.75

// LinkedHashMapNode is a linked list/red-black tree node
type LinkedHashMapNode[K comparable, V any] struct {
	key   K
	value V
	hash  uint64

	// Linked list pointers
	next *LinkedHashMapNode[K, V]

	// Red-black tree pointers
	left   *LinkedHashMapNode[K, V]
	right  *LinkedHashMapNode[K, V]
	parent *LinkedHashMapNode[K, V]
	color  color

	// Mark whether node is a tree node
	isTreeNode bool
}

// LinkedHashMap is a Map implementation based on separate chaining and red-black trees
type LinkedHashMap[K comparable, V any] struct {
	table     []*LinkedHashMapNode[K, V] // Hash bucket array
	size      int                        // Element count
	threshold int                        // Resize threshold
	mutex     sync.RWMutex               // Read-write lock for thread safety
}

// NewLinkedHashMap creates a new LinkedHashMap
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	capacity := initialCapacity
	return &LinkedHashMap[K, V]{
		table:     make([]*LinkedHashMapNode[K, V], capacity),
		size:      0,
		threshold: int(float64(capacity) * loadFactor),
	}
}

// NewLinkedHashMapWithCapacity creates a LinkedHashMap with specified initial capacity
func NewLinkedHashMapWithCapacity[K comparable, V any](capacity int) *LinkedHashMap[K, V] {
	if capacity < initialCapacity {
		capacity = initialCapacity
	} else {
		// Ensure capacity is a power of 2
		capacity = tableSizeFor(capacity)
	}

	return &LinkedHashMap[K, V]{
		table:     make([]*LinkedHashMapNode[K, V], capacity),
		size:      0,
		threshold: int(float64(capacity) * loadFactor),
	}
}

// tableSizeFor returns the smallest power of 2 greater than or equal to cap
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

// hash calculates the hash value of the key
func (m *LinkedHashMap[K, V]) hash(key K) uint64 {
	return Hash(key)
}

// Put associates the specified value with the specified key in this map
func (m *LinkedHashMap[K, V]) Put(key K, value V) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var oldValue V
	existed := false

	hashValue := m.hash(key)
	index := int(hashValue % uint64(len(m.table)))

	// If bucket is empty, create new node
	if m.table[index] == nil {
		m.table[index] = &LinkedHashMapNode[K, V]{
			key:   key,
			value: value,
			hash:  hashValue,
		}
		m.size++

		// Check if resize is needed
		m.checkResize()

		return oldValue, existed
	}

	// If it's a tree node, use tree search and insertion
	if m.table[index].isTreeNode {
		return m.putTreeVal(index, key, value, hashValue)
	}

	// Linked list search and insertion
	p := m.table[index]
	var prev *LinkedHashMapNode[K, V]
	count := 0

	// Traverse the linked list
	for p != nil {
		count++

		// If same key is found, update value
		if p.hash == hashValue && Equal(p.key, key) {
			oldValue = p.value
			p.value = value
			return oldValue, true
		}

		prev = p
		p = p.next
	}

	// No same key found, add new node to end of list
	newNode := &LinkedHashMapNode[K, V]{
		key:   key,
		value: value,
		hash:  hashValue,
	}
	prev.next = newNode
	m.size++

	// Check if linked list needs to be converted to red-black tree
	if count >= treeifyThreshold-1 {
		m.treeifyBin(index)
	}

	// Check if resize is needed
	m.checkResize()

	return oldValue, existed
}

// putTreeVal inserts or updates node in red-black tree
func (m *LinkedHashMap[K, V]) putTreeVal(index int, key K, value V, hash uint64) (V, bool) {
	var oldValue V
	existed := false

	root := m.table[index]
	p := root

	// Tree search
	for p != nil {
		cmp := 0
		if p.hash > hash {
			cmp = -1
		} else if p.hash < hash {
			cmp = 1
		} else if Equal(key, p.key) {
			// Found same key, update value
			oldValue = p.value
			p.value = value
			return oldValue, true
		} else {
			// Same hash but different key, use key comparison
			cmp = Compare(key, p.key)
		}

		// Decide left or right based on comparison result
		if cmp < 0 {
			if p.left == nil {
				// Insert as left child
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
				// Insert as right child
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

	// If tree is empty, create root node
	m.table[index] = &LinkedHashMapNode[K, V]{
		key:        key,
		value:      value,
		hash:       hash,
		isTreeNode: true,
		color:      black, // Root node is black
	}
	m.size++

	return oldValue, existed
}

// balanceInsertion after insertion
func (m *LinkedHashMap[K, V]) balanceInsertion(root *LinkedHashMapNode[K, V], x *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	// Red-black tree balancing adjustment
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

// rotateLeft left rotation
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

// rotateRight right rotation
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

// parentOf get node's parent
func parentOf[K comparable, V any](p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p == nil {
		return nil
	}
	return p.parent
}

// leftOf get node's left child
func leftOf[K comparable, V any](p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p == nil {
		return nil
	}
	return p.left
}

// rightOf get node's right child
func rightOf[K comparable, V any](p *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	if p == nil {
		return nil
	}
	return p.right
}

// colorOf get node's color
func colorOf[K comparable, V any](p *LinkedHashMapNode[K, V]) color {
	if p == nil {
		return black
	}
	return p.color
}

// setColor set node's color
func setColor[K comparable, V any](p *LinkedHashMapNode[K, V], c color) {
	if p != nil {
		p.color = c
	}
}

// treeifyBin convert specified index's list to red-black tree
func (m *LinkedHashMap[K, V]) treeifyBin(index int) {
	// If hash table capacity is less than minimum treeify capacity, prioritize expansion
	if len(m.table) < minTreeifyCapacity {
		m.resize()
		return
	}

	// Convert list to red-black tree
	root := m.table[index]
	if root == nil {
		return
	}

	// Mark all nodes as tree nodes
	p := root
	for p != nil {
		p.isTreeNode = true
		p = p.next
	}

	// Build red-black tree
	root = m.buildTree(root)
	m.table[index] = root
}

// buildTree build red-black tree from list
func (m *LinkedHashMap[K, V]) buildTree(head *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	var root *LinkedHashMapNode[K, V]

	// Traverse list, insert each node to red-black tree
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

			// Find insert position
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

			// Balance red-black tree
			root = m.balanceInsertion(root, p)
		}

		p = next
	}

	return root
}

// checkResize check if resize is needed
func (m *LinkedHashMap[K, V]) checkResize() {
	if m.size > m.threshold {
		m.resize()
	}
}

// resize resize hash table
func (m *LinkedHashMap[K, V]) resize() {
	oldCap := len(m.table)
	oldTab := m.table

	// Compute new capacity
	newCap := oldCap * 2
	if newCap < initialCapacity {
		newCap = initialCapacity
	}

	// Create new table
	newTab := make([]*LinkedHashMapNode[K, V], newCap)
	m.table = newTab
	m.threshold = int(float64(newCap) * loadFactor)

	// If old table is empty, directly return
	if oldCap == 0 {
		return
	}

	// Reallocate old table elements to new table
	for i := 0; i < oldCap; i++ {
		e := oldTab[i]
		if e == nil {
			continue
		}

		// Clear old table reference
		oldTab[i] = nil

		// If it's a single node
		if e.next == nil {
			newIdx := int(e.hash % uint64(newCap))
			newTab[newIdx] = e
			continue
		}

		// If it's a tree node
		if e.isTreeNode {
			m.splitTreeBin(newTab, e, i, oldCap)
			continue
		}

		// If it's a list, split into two lists
		// One list placed at original position, one list placed at original position+oldCap
		var loHead, loTail, hiHead, hiTail *LinkedHashMapNode[K, V]

		// Traverse list
		for e != nil {
			next := e.next

			// Decide where to put the node
			if (e.hash & uint64(oldCap)) == 0 {
				// Place at original position
				if loTail == nil {
					loHead = e
				} else {
					loTail.next = e
				}
				loTail = e
			} else {
				// Place at original position+oldCap
				if hiTail == nil {
					hiHead = e
				} else {
					hiTail.next = e
				}
				hiTail = e
			}

			e = next
		}

		// Update list reference
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

// splitTreeBin split tree node
func (m *LinkedHashMap[K, V]) splitTreeBin(newTab []*LinkedHashMapNode[K, V], root *LinkedHashMapNode[K, V], index, oldCap int) {
	// Convert tree node back to list
	if len(newTab) <= untreeifyThreshold {
		var head, tail *LinkedHashMapNode[K, V]

		// Traverse tree, build list
		m.treeToList(root, &head, &tail)

		// Split list
		var loHead, loTail, hiHead, hiTail *LinkedHashMapNode[K, V]
		p := head

		for p != nil {
			next := p.next
			p.left = nil
			p.right = nil
			p.parent = nil
			p.isTreeNode = false

			// Decide where to put the node
			if (p.hash & uint64(oldCap)) == 0 {
				// Place at original position
				if loTail == nil {
					loHead = p
				} else {
					loTail.next = p
				}
				loTail = p
			} else {
				// Place at original position+oldCap
				if hiTail == nil {
					hiHead = p
				} else {
					hiTail.next = p
				}
				hiTail = p
			}

			p = next
		}

		// Update list reference
		if loTail != nil {
			loTail.next = nil
			newTab[index] = loHead
		}

		if hiTail != nil {
			hiTail.next = nil
			newTab[index+oldCap] = hiHead
		}
	} else {
		// Split tree
		var loTree, hiTree *LinkedHashMapNode[K, V]

		// Traverse tree, build two new trees
		m.splitTree(root, &loTree, &hiTree, oldCap)

		// Update tree reference
		if loTree != nil {
			newTab[index] = loTree
		}

		if hiTree != nil {
			newTab[index+oldCap] = hiTree
		}
	}
}

// treeToList convert tree to list
func (m *LinkedHashMap[K, V]) treeToList(root *LinkedHashMapNode[K, V], head, tail **LinkedHashMapNode[K, V]) {
	// In-order traverse tree, build list
	if root == nil {
		return
	}

	// Recursive left sub-tree
	m.treeToList(root.left, head, tail)

	// Handle current node
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

	// Recursive right sub-tree
	m.treeToList(root.right, head, tail)
}

// splitTree split tree
func (m *LinkedHashMap[K, V]) splitTree(root *LinkedHashMapNode[K, V], loTree, hiTree **LinkedHashMapNode[K, V], oldCap int) {
	if root == nil {
		return
	}

	// Save sub-node reference
	left := root.left
	right := root.right

	// Clean current node reference
	root.left = nil
	root.right = nil
	root.parent = nil

	// Decide where to put the node
	if (root.hash & uint64(oldCap)) == 0 {
		// Place at original position
		if *loTree == nil {
			*loTree = root
			root.color = black
		} else {
			// Simple insert to loTree, no balancing
			m.insertNodeSimple(loTree, root)
		}
	} else {
		// Place at original position+oldCap
		if *hiTree == nil {
			*hiTree = root
			root.color = black
		} else {
			// Simple insert to hiTree, no balancing
			m.insertNodeSimple(hiTree, root)
		}
	}

	// Recursive handle sub-tree
	m.splitTree(left, loTree, hiTree, oldCap)
	m.splitTree(right, loTree, hiTree, oldCap)
}

// insertNodeSimple simple insert node to tree, no balancing
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

// Get return value mapped to the key
func (m *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	hashValue := m.hash(key)
	index := int(hashValue % uint64(len(m.table)))

	// If bucket is empty, return zero value
	if m.table[index] == nil {
		return *new(V), false
	}

	// If it's a tree node, use tree search
	if m.table[index].isTreeNode {
		return m.getTreeVal(m.table[index], key, hashValue)
	}

	// Linked list search
	p := m.table[index]
	for p != nil {
		if p.hash == hashValue && Equal(p.key, key) {
			return p.value, true
		}
		p = p.next
	}

	return *new(V), false
}

// getTreeVal search node in red-black tree
func (m *LinkedHashMap[K, V]) getTreeVal(root *LinkedHashMapNode[K, V], key K, hash uint64) (V, bool) {
	p := root

	// Tree search
	for p != nil {
		cmp := 0
		if p.hash > hash {
			cmp = -1
		} else if p.hash < hash {
			cmp = 1
		} else if Equal(key, p.key) {
			// Find same key, return value
			return p.value, true
		} else {
			// Hash same but different key, use key comparison
			cmp = Compare(key, p.key)
		}

		// Decide left or right based on comparison result
		if cmp < 0 {
			p = p.left
		} else {
			p = p.right
		}
	}

	return *new(V), false
}

// Remove if exists, removes mapping relationship for the key
func (m *LinkedHashMap[K, V]) Remove(key K) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	hashValue := m.hash(key)
	index := int(hashValue % uint64(len(m.table)))

	// If bucket is empty, return zero value
	if m.table[index] == nil {
		return *new(V), false
	}

	// If it's a tree node, use tree delete
	if m.table[index].isTreeNode {
		return m.removeTreeNode(index, key, hashValue)
	}

	// Linked list delete
	p := m.table[index]
	var prev *LinkedHashMapNode[K, V]

	for p != nil {
		if p.hash == hashValue && Equal(p.key, key) {
			// Find node to delete
			oldValue := p.value

			if prev == nil {
				// Delete list head
				m.table[index] = p.next
			} else {
				// Delete list middle or tail
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

// removeTreeNode remove node from red-black tree
func (m *LinkedHashMap[K, V]) removeTreeNode(index int, key K, hash uint64) (V, bool) {
	root := m.table[index]
	p := root

	// Find node to delete
	for p != nil {
		cmp := 0
		if p.hash > hash {
			cmp = -1
		} else if p.hash < hash {
			cmp = 1
		} else if Equal(key, p.key) {
			// Find node to delete
			break
		} else {
			// Hash same but different key, use key comparison
			cmp = Compare(key, p.key)
		}

		// Decide left or right based on comparison result
		if cmp < 0 {
			p = p.left
		} else {
			p = p.right
		}
	}

	// If node is not found, return zero value
	if p == nil {
		return *new(V), false
	}

	oldValue := p.value

	// Delete node
	if p.left != nil && p.right != nil {
		// If node has two children, find successor
		s := p.right
		for s.left != nil {
			s = s.left
		}

		// Replace current node value with successor value
		p.key = s.key
		p.value = s.value
		p.hash = s.hash

		// Delete successor
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
		// If node has at most one child
		replacement := p.left
		if p.left == nil {
			replacement = p.right
		}

		// Replace node with child node
		if p.parent == nil {
			// If it's a root node
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

	// Balance red-black tree
	if p.color == black {
		m.balanceDeletion(root, p)
	}

	m.size--

	// If tree is too small, convert to list
	if m.size <= untreeifyThreshold {
		m.untreeify(index)
	}

	return oldValue, true
}

// balanceDeletion after deletion
func (m *LinkedHashMap[K, V]) balanceDeletion(root *LinkedHashMapNode[K, V], x *LinkedHashMapNode[K, V]) *LinkedHashMapNode[K, V] {
	// Red-black tree delete balancing adjustment
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

// untreeify convert specified index's red-black tree to list
func (m *LinkedHashMap[K, V]) untreeify(index int) {
	root := m.table[index]
	if root == nil || !root.isTreeNode {
		return
	}

	// Convert tree to list
	var head, tail *LinkedHashMapNode[K, V]
	m.treeToList(root, &head, &tail)

	// Update list reference
	m.table[index] = head
}

// ContainsKey if this mapping contains the key's mapping, returns true
func (m *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, found := m.Get(key)
	return found
}

// ContainsValue if this mapping maps one or more keys to the specified value, returns true
func (m *LinkedHashMap[K, V]) ContainsValue(value V) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.traverseAllWithEarlyExit(func(node *LinkedHashMapNode[K, V]) bool {
		return Equal(node.value, value)
	})
}

// Size returns the number of key-value mapping relationships in this mapping
func (m *LinkedHashMap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.size
}

// IsEmpty if this mapping does not contain key-value mapping relationships, returns true
func (m *LinkedHashMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear removes all mapping relationships from this mapping
func (m *LinkedHashMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear hash table, help GC recycle memory
	for i := range m.table {
		if m.table[i] != nil {
			m.clearNode(m.table[i])
			m.table[i] = nil
		}
	}

	m.size = 0
}

// clearNode recursively clean up nodes to help GC
func (m *LinkedHashMap[K, V]) clearNode(node *LinkedHashMapNode[K, V]) {
	if node == nil {
		return
	}

	// If it's a tree node, recursively clean up child nodes
	if node.isTreeNode {
		m.clearNode(node.left)
		m.clearNode(node.right)
		node.left = nil
		node.right = nil
		node.parent = nil
	}

	// Clean list node
	for node.next != nil {
		next := node.next
		node.next = nil
		node = next
	}
}

// Keys returns the keys contained in this mapping
func (m *LinkedHashMap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]K, 0, m.size)
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		keys = append(keys, node.key)
	})

	return keys
}

// Values returns the values contained in this mapping
func (m *LinkedHashMap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	values := make([]V, 0, m.size)
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		values = append(values, node.value)
	})

	return values
}

// inOrderTraversal in-order traverse red-black tree
func (m *LinkedHashMap[K, V]) inOrderTraversal(root *LinkedHashMapNode[K, V], f func(*LinkedHashMapNode[K, V])) {
	if root == nil {
		return
	}

	m.inOrderTraversal(root.left, f)
	f(root)
	m.inOrderTraversal(root.right, f)
}

// traverseAll traverse all nodes, including linked list and tree nodes
func (m *LinkedHashMap[K, V]) traverseAll(f func(*LinkedHashMapNode[K, V])) {
	// Traverse hash table
	for _, e := range m.table {
		p := e
		for p != nil {
			f(p)
			if p.isTreeNode {
				// If it's a tree node, use in-order traverse other nodes
				m.inOrderTraversal(p, func(node *LinkedHashMapNode[K, V]) {
					if node != p { // Avoid repeat processing root node
						f(node)
					}
				})
				break
			}
			p = p.next
		}
	}
}

// traverseAllWithEarlyExit traverse all nodes, support early exit
func (m *LinkedHashMap[K, V]) traverseAllWithEarlyExit(f func(*LinkedHashMapNode[K, V]) bool) bool {
	// Traverse hash table
	for _, e := range m.table {
		p := e
		for p != nil {
			if f(p) {
				return true
			}
			if p.isTreeNode {
				// If it's a tree node, use in-order traverse other nodes
				found := false
				m.inOrderTraversalWithEarlyExit(p, func(node *LinkedHashMapNode[K, V]) bool {
					if node != p { // Avoid repeat processing root node
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

// inOrderTraversalWithEarlyExit in-order traverse red-black tree, support early exit
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

// ForEach execute the action for each entry
func (m *LinkedHashMap[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		f(node.key, node.value)
	})
}

// PutAll put all mapping relationships from the specified mapping to this mapping
func (m *LinkedHashMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(k K, v V) {
		m.Put(k, v)
	})
}

// String returns the string representation of the mapping
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

// Entries returns the mapping relationships contained in this mapping
func (m *LinkedHashMap[K, V]) Entries() []Pair[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entries := make([]Pair[K, V], 0, m.size)
	m.traverseAll(func(node *LinkedHashMapNode[K, V]) {
		entries = append(entries, NewPair(node.key, node.value))
	})

	return entries
}
