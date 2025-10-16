package set

import (
    "fmt"
    "strings"

    "github.com/chenjianyu/collections/container/common"
)

// defaultIntComparator provides a default comparator only for int types.
// For non-int element types, please use NewTreeSetWithComparator or NewTreeSetWithComparatorStrategy.
func defaultIntComparator[E comparable](a, b E) int {
    cmpA, okA := any(a).(int)
    cmpB, okB := any(b).(int)
    if okA && okB {
        if cmpA == cmpB {
            return 0
        } else if cmpA < cmpB {
            return -1
        }
        return 1
    }
    // Non-int types are not supported by the default comparator
    // Fall back to equality to avoid undefined ordering; pass a custom comparator for proper ordering
    return 0
}

// TreeSet is a set implementation based on red-black tree
// Elements are stored in sorted order
type TreeSet[E comparable] struct {
	root       *treeNode[E]
	size       int
	comparator func(a, b E) int
}

// treeNode represents a node in the red-black tree
type treeNode[E comparable] struct {
	value  E
	color  bool // true for red, false for black
	left   *treeNode[E]
	right  *treeNode[E]
	parent *treeNode[E]
}

// NewTreeSet creates a new TreeSet using default comparator
// Default comparator prefers Comparable.CompareTo, otherwise falls back to natural generic ordering
func NewTreeSet[E comparable]() *TreeSet[E] {
    return &TreeSet[E]{
        root:       nil,
        size:       0,
        comparator: common.CompareNatural[E],
    }
}

// NewTreeSetWithComparator creates a new TreeSet using custom comparator
func NewTreeSetWithComparator[E comparable](comparator func(a, b E) int) *TreeSet[E] {
    return &TreeSet[E]{
        root:       nil,
        size:       0,
        comparator: comparator,
    }
}

// NewTreeSetWithComparatorStrategy creates a new TreeSet using a ComparatorStrategy from common package
func NewTreeSetWithComparatorStrategy[E comparable](strategy common.ComparatorStrategy[E]) *TreeSet[E] {
    return &TreeSet[E]{
        root: nil,
        size: 0,
        comparator: func(a, b E) int {
            return strategy.Compare(a, b)
        },
    }
}

// Size returns the number of elements in the set
func (ts *TreeSet[E]) Size() int {
	return ts.size
}

// IsEmpty checks if the set is empty
func (ts *TreeSet[E]) IsEmpty() bool {
	return ts.size == 0
}

// Clear clears all elements from the set
func (ts *TreeSet[E]) Clear() {
	ts.root = nil
	ts.size = 0
}

// Contains checks if the set contains the specified element
func (ts *TreeSet[E]) Contains(element E) bool {
	return ts.findNode(element) != nil
}

// Add adds an element to the set
// Returns false if the set already contains the element, otherwise returns true
func (ts *TreeSet[E]) Add(element E) bool {
	if ts.root == nil {
		ts.root = &treeNode[E]{
			value: element,
			color: false, // root is always black
		}
		ts.size++
		return true
	}

	node := ts.root
	for {
		cmp := ts.comparator(element, node.value)
		if cmp == 0 {
			return false // element already exists
		} else if cmp < 0 {
			if node.left == nil {
				newNode := &treeNode[E]{
					value:  element,
					color:  true, // new node is red
					parent: node,
				}
				node.left = newNode
				ts.insertFixup(newNode)
				ts.size++
				return true
			}
			node = node.left
		} else {
			if node.right == nil {
				newNode := &treeNode[E]{
					value:  element,
					color:  true, // new node is red
					parent: node,
				}
				node.right = newNode
				ts.insertFixup(newNode)
				ts.size++
				return true
			}
			node = node.right
		}
	}
}

// Remove removes the specified element from the set
// Returns true if the set contained the element, otherwise returns false
func (ts *TreeSet[E]) Remove(element E) bool {
	node := ts.findNode(element)
	if node == nil {
		return false
	}

	ts.deleteNode(node)
	ts.size--
	return true
}

// ToSlice returns a slice containing all elements in the set
func (ts *TreeSet[E]) ToSlice() []E {
	result := make([]E, 0, ts.size)
	ts.inorderTraversal(ts.root, func(value E) {
		result = append(result, value)
	})
	return result
}

// Union returns a new set containing all elements from this set and the other set
func (ts *TreeSet[E]) Union(other Set[E]) Set[E] {
	result := NewTreeSetWithComparator(ts.comparator)
	ts.ForEach(func(element E) {
		result.Add(element)
	})
	other.ForEach(func(element E) {
		result.Add(element)
	})
	return result
}

// Intersection returns a new set containing elements that exist in both sets
func (ts *TreeSet[E]) Intersection(other Set[E]) Set[E] {
	result := NewTreeSetWithComparator(ts.comparator)
	ts.ForEach(func(element E) {
		if other.Contains(element) {
			result.Add(element)
		}
	})
	return result
}

// Difference returns a new set containing elements that exist in this set but not in the other set
func (ts *TreeSet[E]) Difference(other Set[E]) Set[E] {
	result := NewTreeSetWithComparator(ts.comparator)
	ts.ForEach(func(element E) {
		if !other.Contains(element) {
			result.Add(element)
		}
	})
	return result
}

// IsSubsetOf checks if this set is a subset of the other set
func (ts *TreeSet[E]) IsSubsetOf(other Set[E]) bool {
	if ts.Size() > other.Size() {
		return false
	}

	isSubset := true
	ts.ForEach(func(element E) {
		if !other.Contains(element) {
			isSubset = false
		}
	})
	return isSubset
}

// IsSupersetOf checks if this set is a superset of the other set
func (ts *TreeSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(ts)
}

// ForEach executes the given operation for each element in the set
func (ts *TreeSet[E]) ForEach(fn func(E)) {
	ts.inorderTraversal(ts.root, fn)
}

// Iterator returns an iterator for the set
func (ts *TreeSet[E]) Iterator() common.Iterator[E] {
	elements := ts.ToSlice()
	return &treeSetIterator[E]{
		elements: elements,
		index:    0,
	}
}

// treeSetIterator implements Iterator for TreeSet
type treeSetIterator[E comparable] struct {
	elements []E
	index    int
}

// HasNext returns true if there are more elements to iterate
func (it *treeSetIterator[E]) HasNext() bool {
	return it.index < len(it.elements)
}

// Next returns the next element
func (it *treeSetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		var zero E
		return zero, false
	}
	element := it.elements[it.index]
	it.index++
	return element, true
}

// Remove removes the current element (not supported)
func (it *treeSetIterator[E]) Remove() bool {
	return false // Not supported
}

// String returns the string representation of the set
func (ts *TreeSet[E]) String() string {
	if ts.IsEmpty() {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{")

	first := true
	ts.ForEach(func(element E) {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", element))
		first = false
	})

	sb.WriteString("}")
	return sb.String()
}

// Internal method: find node with specified value
func (ts *TreeSet[E]) findNode(element E) *treeNode[E] {
	node := ts.root
	for node != nil {
		cmp := ts.comparator(element, node.value)
		if cmp == 0 {
			return node
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}
	return nil
}

// Internal method: in-order traversal
func (ts *TreeSet[E]) inorderTraversal(node *treeNode[E], fn func(E)) {
	if node != nil {
		ts.inorderTraversal(node.left, fn)
		fn(node.value)
		ts.inorderTraversal(node.right, fn)
	}
}

// Internal method: left rotation
func (ts *TreeSet[E]) rotateLeft(node *treeNode[E]) {
	right := node.right
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.parent = node.parent
	if node.parent == nil {
		ts.root = right
	} else if node == node.parent.left {
		node.parent.left = right
	} else {
		node.parent.right = right
	}
	right.left = node
	node.parent = right
}

// Internal method: right rotation
func (ts *TreeSet[E]) rotateRight(node *treeNode[E]) {
	left := node.left
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.parent = node.parent
	if node.parent == nil {
		ts.root = left
	} else if node == node.parent.right {
		node.parent.right = left
	} else {
		node.parent.left = left
	}
	left.right = node
	node.parent = left
}

// Internal method: fix red-black tree properties after insertion
func (ts *TreeSet[E]) insertFixup(node *treeNode[E]) {
	for node.parent != nil && node.parent.color {
		if node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle != nil && uncle.color {
				node.parent.color = false
				uncle.color = false
				node.parent.parent.color = true
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					ts.rotateLeft(node)
				}
				node.parent.color = false
				node.parent.parent.color = true
				ts.rotateRight(node.parent.parent)
			}
		} else {
			uncle := node.parent.parent.left
			if uncle != nil && uncle.color {
				node.parent.color = false
				uncle.color = false
				node.parent.parent.color = true
				node = node.parent.parent
			} else {
				if node == node.parent.left {
					node = node.parent
					ts.rotateRight(node)
				}
				node.parent.color = false
				node.parent.parent.color = true
				ts.rotateLeft(node.parent.parent)
			}
		}
	}
	ts.root.color = false
}

// Internal method: delete node
func (ts *TreeSet[E]) deleteNode(node *treeNode[E]) {
	var y *treeNode[E]
	var x *treeNode[E]

	if node.left == nil || node.right == nil {
		y = node
	} else {
		y = ts.successor(node)
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		ts.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != node {
		node.value = y.value
	}

	if !y.color && x != nil {
		ts.deleteFixup(x)
	}
}

// Internal method: find successor node
func (ts *TreeSet[E]) successor(node *treeNode[E]) *treeNode[E] {
	if node.right != nil {
		node = node.right
		for node.left != nil {
			node = node.left
		}
		return node
	}

	parent := node.parent
	for parent != nil && node == parent.right {
		node = parent
		parent = parent.parent
	}
	return parent
}

// Internal method: fix red-black tree properties after deletion
func (ts *TreeSet[E]) deleteFixup(node *treeNode[E]) {
	for node != ts.root && !node.color {
		if node == node.parent.left {
			sibling := node.parent.right
			if sibling.color {
				sibling.color = false
				node.parent.color = true
				ts.rotateLeft(node.parent)
				sibling = node.parent.right
			}
			if (sibling.left == nil || !sibling.left.color) &&
				(sibling.right == nil || !sibling.right.color) {
				sibling.color = true
				node = node.parent
			} else {
				if sibling.right == nil || !sibling.right.color {
					if sibling.left != nil {
						sibling.left.color = false
					}
					sibling.color = true
					ts.rotateRight(sibling)
					sibling = node.parent.right
				}
				sibling.color = node.parent.color
				node.parent.color = false
				if sibling.right != nil {
					sibling.right.color = false
				}
				ts.rotateLeft(node.parent)
				node = ts.root
			}
		} else {
			sibling := node.parent.left
			if sibling.color {
				sibling.color = false
				node.parent.color = true
				ts.rotateRight(node.parent)
				sibling = node.parent.left
			}
			if (sibling.right == nil || !sibling.right.color) &&
				(sibling.left == nil || !sibling.left.color) {
				sibling.color = true
				node = node.parent
			} else {
				if sibling.left == nil || !sibling.left.color {
					if sibling.right != nil {
						sibling.right.color = false
					}
					sibling.color = true
					ts.rotateLeft(sibling)
					sibling = node.parent.left
				}
				sibling.color = node.parent.color
				node.parent.color = false
				if sibling.left != nil {
					sibling.left.color = false
				}
				ts.rotateRight(node.parent)
				node = ts.root
			}
		}
	}
	node.color = false
}
