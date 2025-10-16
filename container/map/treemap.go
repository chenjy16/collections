package maps

import (
    "fmt"
    "strings"
    "github.com/chenjianyu/collections/container/common"
)

// TreeMap is an ordered Map implementation based on red-black tree
type TreeMap[K comparable, V any] struct {
	comparator func(a, b K) int // Key comparator function
	root       *mapNode[K, V]   // Root node
	size       int              // Element count
}

// mapNode is a red-black tree node
type mapNode[K comparable, V any] struct {
	key    K
	value  V
	color  color
	left   *mapNode[K, V]
	right  *mapNode[K, V]
	parent *mapNode[K, V]
}

// NewTreeMap creates a new TreeMap using default comparator
func NewTreeMap[K comparable, V any]() *TreeMap[K, V] {
    return &TreeMap[K, V]{
        comparator: func(a, b K) int {
            return common.CompareNatural(a, b)
        },
    }
}

// NewTreeMapWithComparator creates a new TreeMap using specified comparator
func NewTreeMapWithComparator[K comparable, V any](comparator func(a, b K) int) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		comparator: comparator,
	}
}

// isRedMap checks if node is red
func isRedMap[K comparable, V any](node *mapNode[K, V]) bool {
	if node == nil {
		return false
	}
	return node.color == red
}

// rotateLeft left rotation
func (m *TreeMap[K, V]) rotateLeft(h *mapNode[K, V]) *mapNode[K, V] {
	x := h.right
	h.right = x.left
	if x.left != nil {
		x.left.parent = h
	}
	x.parent = h.parent
	if h.parent == nil {
		m.root = x
	} else if h == h.parent.left {
		h.parent.left = x
	} else {
		h.parent.right = x
	}
	x.left = h
	h.parent = x
	x.color = h.color
	h.color = red
	return x
}

// rotateRight right rotation
func (m *TreeMap[K, V]) rotateRight(h *mapNode[K, V]) *mapNode[K, V] {
	x := h.left
	h.left = x.right
	if x.right != nil {
		x.right.parent = h
	}
	x.parent = h.parent
	if h.parent == nil {
		m.root = x
	} else if h == h.parent.right {
		h.parent.right = x
	} else {
		h.parent.left = x
	}
	x.right = h
	h.parent = x
	x.color = h.color
	h.color = red
	return x
}

// flipColors color flip
func (m *TreeMap[K, V]) flipColors(h *mapNode[K, V]) {
	h.color = red
	if h.left != nil {
		h.left.color = black
	}
	if h.right != nil {
		h.right.color = black
	}
}

// put insert or update node
func (m *TreeMap[K, V]) put(h *mapNode[K, V], key K, value V) (*mapNode[K, V], V, bool) {
	var oldValue V
	existed := false

	if h == nil {
		m.size++
		return &mapNode[K, V]{
			key:   key,
			value: value,
			color: red,
		}, oldValue, existed
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
		// Key already exists, update value
		oldValue = h.value
		h.value = value
		existed = true
		return h, oldValue, existed
	}

	// Red-black tree balancing adjustment
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

// Put associates the specified value with the specified key in this map
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

// findMin find minimum node
func (m *TreeMap[K, V]) findMin(h *mapNode[K, V]) *mapNode[K, V] {
	if h == nil {
		return nil
	}
	for h.left != nil {
		h = h.left
	}
	return h
}

// moveRedLeft make left child or its sibling red
func (m *TreeMap[K, V]) moveRedLeft(h *mapNode[K, V]) *mapNode[K, V] {
	m.flipColors(h)
	if h.right != nil && isRedMap(h.right.left) {
		h.right = m.rotateRight(h.right)
		h = m.rotateLeft(h)
		m.flipColors(h)
	}
	return h
}

// moveRedRight make right child or its sibling red
func (m *TreeMap[K, V]) moveRedRight(h *mapNode[K, V]) *mapNode[K, V] {
	m.flipColors(h)
	if h.left != nil && isRedMap(h.left.left) {
		h = m.rotateRight(h)
		m.flipColors(h)
	}
	return h
}

// balance balance node
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

// removeMin remove minimum node
func (m *TreeMap[K, V]) removeMin(h *mapNode[K, V], decreaseSize bool) *mapNode[K, V] {
	if h.left == nil {
		// Only decrease size when this is a standalone delete operation
		// When called from remove method, size will be handled in remove method
		if decreaseSize {
			m.size--
		}
		return nil
	}

	if !isRedMap(h.left) && !isRedMap(h.left.left) {
		h = m.moveRedLeft(h)
	}

	h.left = m.removeMin(h.left, decreaseSize)
	if h.left != nil {
		h.left.parent = h
	}

	return m.balance(h)
}

// remove remove node
func (m *TreeMap[K, V]) remove(h *mapNode[K, V], key K) (*mapNode[K, V], V, bool) {
	var oldValue V
	found := false

	if h == nil {
		return nil, oldValue, found
	}

	cmp := m.comparator(key, h.key)
	if cmp < 0 {
		if h.left != nil {
			if !isRedMap(h.left) && !isRedMap(h.left.left) {
				h = m.moveRedLeft(h)
			}
			// Key is in left subtree
			h.left, oldValue, found = m.remove(h.left, key)
			if h.left != nil {
				h.left.parent = h
			}
		}
	} else {
		if isRedMap(h.left) {
			h = m.rotateRight(h)
		}
		if cmp == 0 && h.right == nil {
			// Key is in right subtree
			oldValue = h.value
			found = true
			m.size--
			return nil, oldValue, found
		}
		if h.right != nil {
			if !isRedMap(h.right) && !isRedMap(h.right.left) {
				h = m.moveRedRight(h)
			}
			cmp = m.comparator(key, h.key) // Recalculate comparison result
			if cmp == 0 {
				// Found node to delete
				oldValue = h.value
				found = true
				m.size--

				min := m.findMin(h.right)
				h.key = min.key
				h.value = min.value

				// Delete minimum node from right subtree, but don't decrease size (since we just moved the node)
				h.right = m.removeMin(h.right, false)
				if h.right != nil {
					h.right.parent = h
				}
			} else {
				h.right, oldValue, found = m.remove(h.right, key)
				if h.right != nil {
					h.right.parent = h
				}
			}
		}
	}

	if h == nil {
		return nil, oldValue, found
	}

	if h.left == nil && h.right == nil {
		// No right subtree, directly return left subtree
		return nil, oldValue, found
	}

	if h.left == nil {
		// No left subtree, directly return right subtree
		return h.right, oldValue, found
	}

	// Has two subtrees, replace current node with minimum node from right subtree
	return m.balance(h), oldValue, found
}

// removeMinWithoutSizeChange remove minimum node but don't change size
func (m *TreeMap[K, V]) removeMinWithoutSizeChange(h *mapNode[K, V]) *mapNode[K, V] {
	if h.left == nil {
		return nil
	}

	if !isRedMap(h.left) && !isRedMap(h.left.left) {
		h = m.moveRedLeft(h)
	}

	h.left = m.removeMinWithoutSizeChange(h.left)
	if h.left != nil {
		h.left.parent = h
	}

	return m.balance(h)
}

// Remove if exists, removes mapping relationship for the key from this map
func (m *TreeMap[K, V]) Remove(key K) (V, bool) {
	var oldValue V
	found := false

	if m.root == nil {
		return oldValue, found
	}

	m.root, oldValue, found = m.remove(m.root, key)
	if m.root != nil {
		m.root.color = black
		m.root.parent = nil
	}

	return oldValue, found
}

// find find node
func (m *TreeMap[K, V]) find(h *mapNode[K, V], key K) *mapNode[K, V] {
	for h != nil {
		cmp := m.comparator(key, h.key)
		if cmp < 0 {
			h = h.left
		} else if cmp > 0 {
			h = h.right
		} else {
			return h
		}
	}
	return nil
}

// Get returns the value mapped to the specified key
func (m *TreeMap[K, V]) Get(key K) (V, bool) {
	node := m.find(m.root, key)
	if node != nil {
		return node.value, true
	}
	return *new(V), false
}

// ContainsKey if this map contains mapping relationship for the specified key, returns true
func (m *TreeMap[K, V]) ContainsKey(key K) bool {
	_, found := m.Get(key)
	return found
}

// Size returns the number of key-value mapping relationships in this map
func (m *TreeMap[K, V]) Size() int {
	return m.size
}

// IsEmpty if this map does not contain key-value mapping relationships, returns true
func (m *TreeMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

// Clear removes all mapping relationships from this map
func (m *TreeMap[K, V]) Clear() {
	m.root = nil
	m.size = 0
}

// inOrderTraversalMap in-order traversal
func (m *TreeMap[K, V]) inOrderTraversalMap(node *mapNode[K, V], visit func(K, V)) {
	if node == nil {
		return
	}

	m.inOrderTraversalMap(node.left, visit)
	visit(node.key, node.value)
	m.inOrderTraversalMap(node.right, visit)
}

func (m *TreeMap[K, V]) inOrderTraversalMapKeys(node *mapNode[K, V], keys *[]K) {
	if node == nil {
		return
	}

	m.inOrderTraversalMapKeys(node.left, keys)
	*keys = append(*keys, node.key)
	m.inOrderTraversalMapKeys(node.right, keys)
}

func (m *TreeMap[K, V]) inOrderTraversalMapValues(node *mapNode[K, V], values *[]V) {
	if node == nil {
		return
	}

	m.inOrderTraversalMapValues(node.left, values)
	*values = append(*values, node.value)
	m.inOrderTraversalMapValues(node.right, values)
}

func (m *TreeMap[K, V]) inOrderTraversalMapEntries(node *mapNode[K, V], entries *[]common.Entry[K, V]) {
    if node == nil {
        return
    }

    m.inOrderTraversalMapEntries(node.left, entries)
    *entries = append(*entries, common.NewEntry(node.key, node.value))
    m.inOrderTraversalMapEntries(node.right, entries)
}

// Keys returns the keys contained in this map (in order)
func (m *TreeMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.size)
	m.inOrderTraversalMapKeys(m.root, &keys)
	return keys
}

// Values returns the values contained in this map (in key order)
func (m *TreeMap[K, V]) Values() []V {
	values := make([]V, 0, m.size)
	m.inOrderTraversalMapValues(m.root, &values)
	return values
}

// Entries returns the mapping relationships contained in this map (in key order)
func (m *TreeMap[K, V]) Entries() []common.Entry[K, V] {
    entries := make([]common.Entry[K, V], 0, m.size)
    m.inOrderTraversalMapEntries(m.root, &entries)
    return entries
}

// ForEach executes the given operation for each entry in this map (in key order)
func (m *TreeMap[K, V]) ForEach(f func(K, V)) {
	m.inOrderTraversalMap(m.root, f)
}

// ContainsValue if this map maps one or more keys to the specified value, returns true
func (m *TreeMap[K, V]) ContainsValue(value V) bool {
    found := false
    m.inOrderTraversalMap(m.root, func(k K, v V) {
        if !found && common.Equal(v, value) {
            found = true
        }
    })
    return found
}

// String returns the string representation of the map
func (m *TreeMap[K, V]) String() string {
	if m.IsEmpty() {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")
	first := true
	m.inOrderTraversalMap(m.root, func(k K, v V) {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%v=%v", k, v))
		first = false
	})
	builder.WriteString("}")
	return builder.String()
}

// PutAll copies all mapping relationships from the specified map to this map
func (m *TreeMap[K, V]) PutAll(other Map[K, V]) {
	other.ForEach(func(k K, v V) {
		m.Put(k, v)
	})
}
