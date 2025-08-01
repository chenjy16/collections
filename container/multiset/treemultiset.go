package multiset

import (
	"fmt"
	"strings"
	"sync"

	"github.com/chenjianyu/collections/container/common"
)

// TreeMultiset is a multiset implementation based on a balanced binary search tree
// It maintains elements in sorted order and provides O(log n) time complexity for basic operations
type TreeMultiset[E comparable] struct {
	root *treeNode[E]
	size int
	mu   sync.RWMutex
	cmp  func(E, E) int
}

type treeNode[E comparable] struct {
	element E
	count   int
	left    *treeNode[E]
	right   *treeNode[E]
	height  int
}

// NewTreeMultiset creates a new empty TreeMultiset with default comparison
func NewTreeMultiset[E comparable]() *TreeMultiset[E] {
	return &TreeMultiset[E]{
		cmp: defaultCompare[E],
	}
}

// NewTreeMultisetWithComparator creates a new TreeMultiset with custom comparator
func NewTreeMultisetWithComparator[E comparable](cmp func(E, E) int) *TreeMultiset[E] {
	return &TreeMultiset[E]{
		cmp: cmp,
	}
}

// NewTreeMultisetFromSlice creates a new TreeMultiset from a slice
func NewTreeMultisetFromSlice[E comparable](elements []E) *TreeMultiset[E] {
	ms := NewTreeMultiset[E]()
	for _, element := range elements {
		ms.Add(element)
	}
	return ms
}

// defaultCompare provides default comparison for comparable types
func defaultCompare[E comparable](a, b E) int {
	if any(a) == any(b) {
		return 0
	}
	// For string comparison
	if sa, ok := any(a).(string); ok {
		if sb, ok := any(b).(string); ok {
			if sa < sb {
				return -1
			} else if sa > sb {
				return 1
			}
			return 0
		}
	}
	// For numeric types, we'll use a simple approach
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// Add adds one occurrence of the specified element
func (ms *TreeMultiset[E]) Add(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := 0
	ms.root, prevCount = ms.addNode(ms.root, element)
	ms.size++
	return prevCount
}

func (ms *TreeMultiset[E]) addNode(node *treeNode[E], element E) (*treeNode[E], int) {
	if node == nil {
		return &treeNode[E]{
			element: element,
			count:   1,
			height:  1,
		}, 0
	}
	
	cmp := ms.cmp(element, node.element)
	var prevCount int
	
	if cmp < 0 {
		node.left, prevCount = ms.addNode(node.left, element)
	} else if cmp > 0 {
		node.right, prevCount = ms.addNode(node.right, element)
	} else {
		prevCount = node.count
		node.count++
		return node, prevCount
	}
	
	return ms.balance(node), prevCount
}

// AddCount adds the specified number of occurrences of the element
func (ms *TreeMultiset[E]) AddCount(element E, count int) int {
	if count < 0 {
		panic("count cannot be negative")
	}
	if count == 0 {
		return ms.Count(element)
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := 0
	ms.root, prevCount = ms.addCountNode(ms.root, element, count)
	ms.size += count
	return prevCount
}

func (ms *TreeMultiset[E]) addCountNode(node *treeNode[E], element E, count int) (*treeNode[E], int) {
	if node == nil {
		return &treeNode[E]{
			element: element,
			count:   count,
			height:  1,
		}, 0
	}
	
	cmp := ms.cmp(element, node.element)
	var prevCount int
	
	if cmp < 0 {
		node.left, prevCount = ms.addCountNode(node.left, element, count)
	} else if cmp > 0 {
		node.right, prevCount = ms.addCountNode(node.right, element, count)
	} else {
		prevCount = node.count
		node.count += count
		return node, prevCount
	}
	
	return ms.balance(node), prevCount
}

// Remove removes one occurrence of the specified element
func (ms *TreeMultiset[E]) Remove(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := 0
	ms.root, prevCount = ms.removeNode(ms.root, element, 1)
	if prevCount > 0 {
		ms.size--
	}
	return prevCount
}

func (ms *TreeMultiset[E]) removeNode(node *treeNode[E], element E, count int) (*treeNode[E], int) {
	if node == nil {
		return nil, 0
	}
	
	cmp := ms.cmp(element, node.element)
	var prevCount int
	
	if cmp < 0 {
		node.left, prevCount = ms.removeNode(node.left, element, count)
	} else if cmp > 0 {
		node.right, prevCount = ms.removeNode(node.right, element, count)
	} else {
		prevCount = node.count
		node.count -= count
		if node.count <= 0 {
			return ms.deleteNode(node), prevCount
		}
		return node, prevCount
	}
	
	return ms.balance(node), prevCount
}

// RemoveCount removes the specified number of occurrences of the element
func (ms *TreeMultiset[E]) RemoveCount(element E, count int) int {
	if count < 0 {
		panic("count cannot be negative")
	}
	if count == 0 {
		return ms.Count(element)
	}
	
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := 0
	actualRemoved := 0
	ms.root, prevCount, actualRemoved = ms.removeCountNode(ms.root, element, count)
	ms.size -= actualRemoved
	return prevCount
}

func (ms *TreeMultiset[E]) removeCountNode(node *treeNode[E], element E, count int) (*treeNode[E], int, int) {
	if node == nil {
		return nil, 0, 0
	}
	
	cmp := ms.cmp(element, node.element)
	var prevCount, actualRemoved int
	
	if cmp < 0 {
		node.left, prevCount, actualRemoved = ms.removeCountNode(node.left, element, count)
	} else if cmp > 0 {
		node.right, prevCount, actualRemoved = ms.removeCountNode(node.right, element, count)
	} else {
		prevCount = node.count
		actualRemoved = count
		if actualRemoved > prevCount {
			actualRemoved = prevCount
		}
		
		node.count -= actualRemoved
		if node.count <= 0 {
			return ms.deleteNode(node), prevCount, actualRemoved
		}
		return node, prevCount, actualRemoved
	}
	
	return ms.balance(node), prevCount, actualRemoved
}

// RemoveAll removes all occurrences of the specified element
func (ms *TreeMultiset[E]) RemoveAll(element E) int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	prevCount := 0
	ms.root, prevCount = ms.removeAllNode(ms.root, element)
	ms.size -= prevCount
	return prevCount
}

func (ms *TreeMultiset[E]) removeAllNode(node *treeNode[E], element E) (*treeNode[E], int) {
	if node == nil {
		return nil, 0
	}
	
	cmp := ms.cmp(element, node.element)
	var prevCount int
	
	if cmp < 0 {
		node.left, prevCount = ms.removeAllNode(node.left, element)
	} else if cmp > 0 {
		node.right, prevCount = ms.removeAllNode(node.right, element)
	} else {
		prevCount = node.count
		return ms.deleteNode(node), prevCount
	}
	
	return ms.balance(node), prevCount
}

func (ms *TreeMultiset[E]) deleteNode(node *treeNode[E]) *treeNode[E] {
	if node.left == nil {
		return node.right
	}
	if node.right == nil {
		return node.left
	}
	
	// Find inorder successor
	successor := ms.findMin(node.right)
	node.element = successor.element
	node.count = successor.count
	node.right = ms.deleteMin(node.right)
	
	return ms.balance(node)
}

func (ms *TreeMultiset[E]) findMin(node *treeNode[E]) *treeNode[E] {
	for node.left != nil {
		node = node.left
	}
	return node
}

func (ms *TreeMultiset[E]) deleteMin(node *treeNode[E]) *treeNode[E] {
	if node.left == nil {
		return node.right
	}
	node.left = ms.deleteMin(node.left)
	return ms.balance(node)
}

// Count returns the number of occurrences of the specified element
func (ms *TreeMultiset[E]) Count(element E) int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	node := ms.findNode(ms.root, element)
	if node == nil {
		return 0
	}
	return node.count
}

func (ms *TreeMultiset[E]) findNode(node *treeNode[E], element E) *treeNode[E] {
	if node == nil {
		return nil
	}
	
	cmp := ms.cmp(element, node.element)
	if cmp < 0 {
		return ms.findNode(node.left, element)
	} else if cmp > 0 {
		return ms.findNode(node.right, element)
	}
	return node
}

// SetCount sets the count of the specified element to the given value
func (ms *TreeMultiset[E]) SetCount(element E, count int) int {
	if count < 0 {
		panic("count cannot be negative")
	}
	
	prevCount := ms.Count(element)
	
	if count == 0 {
		ms.RemoveAll(element)
	} else {
		diff := count - prevCount
		if diff > 0 {
			ms.AddCount(element, diff)
		} else if diff < 0 {
			ms.RemoveCount(element, -diff)
		}
	}
	
	return prevCount
}

// Contains checks if the multiset contains the specified element
func (ms *TreeMultiset[E]) Contains(element E) bool {
	return ms.Count(element) > 0
}

// IsEmpty returns true if the multiset contains no elements
func (ms *TreeMultiset[E]) IsEmpty() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.size == 0
}

// Size returns the number of distinct elements in the multiset
func (ms *TreeMultiset[E]) Size() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.countNodes(ms.root)
}

func (ms *TreeMultiset[E]) countNodes(node *treeNode[E]) int {
	if node == nil {
		return 0
	}
	return 1 + ms.countNodes(node.left) + ms.countNodes(node.right)
}

// TotalSize returns the total number of elements (including duplicates)
func (ms *TreeMultiset[E]) TotalSize() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.size
}

// DistinctElements returns the number of distinct elements
func (ms *TreeMultiset[E]) DistinctElements() int {
	return ms.Size()
}

// Clear removes all elements from the multiset
func (ms *TreeMultiset[E]) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.root = nil
	ms.size = 0
}

// ElementSet returns a slice of distinct elements in sorted order
func (ms *TreeMultiset[E]) ElementSet() []E {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	var elements []E
	ms.inorderElements(ms.root, &elements)
	return elements
}

func (ms *TreeMultiset[E]) inorderElements(node *treeNode[E], elements *[]E) {
	if node != nil {
		ms.inorderElements(node.left, elements)
		*elements = append(*elements, node.element)
		ms.inorderElements(node.right, elements)
	}
}

// EntrySet returns a slice of entries (element-count pairs) in sorted order
func (ms *TreeMultiset[E]) EntrySet() []Entry[E] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	var entries []Entry[E]
	ms.inorderEntries(ms.root, &entries)
	return entries
}

func (ms *TreeMultiset[E]) inorderEntries(node *treeNode[E], entries *[]Entry[E]) {
	if node != nil {
		ms.inorderEntries(node.left, entries)
		*entries = append(*entries, Entry[E]{Element: node.element, Count: node.count})
		ms.inorderEntries(node.right, entries)
	}
}

// ToSlice returns a slice containing all elements (including duplicates) in sorted order
func (ms *TreeMultiset[E]) ToSlice() []E {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	var result []E
	ms.inorderSlice(ms.root, &result)
	return result
}

func (ms *TreeMultiset[E]) inorderSlice(node *treeNode[E], result *[]E) {
	if node != nil {
		ms.inorderSlice(node.left, result)
		for i := 0; i < node.count; i++ {
			*result = append(*result, node.element)
		}
		ms.inorderSlice(node.right, result)
	}
}

// Iterator returns an iterator over the multiset elements in sorted order
func (ms *TreeMultiset[E]) Iterator() common.Iterator[E] {
	return &treeMultisetIterator[E]{
		multiset: ms,
		entries:  ms.EntrySet(),
		index:    0,
		current:  0,
	}
}

// ForEach executes the given function for each element in the multiset in sorted order
func (ms *TreeMultiset[E]) ForEach(fn func(E)) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	ms.forEachInorder(ms.root, fn)
}

func (ms *TreeMultiset[E]) forEachInorder(node *treeNode[E], fn func(E)) {
	if node != nil {
		ms.forEachInorder(node.left, fn)
		for i := 0; i < node.count; i++ {
			fn(node.element)
		}
		ms.forEachInorder(node.right, fn)
	}
}

// Union returns a new multiset containing the union of this and another multiset
func (ms *TreeMultiset[E]) Union(other Multiset[E]) Multiset[E] {
	result := NewTreeMultisetWithComparator(ms.cmp)
	
	// Add all elements from this multiset
	for _, entry := range ms.EntrySet() {
		result.AddCount(entry.Element, entry.Count)
	}
	
	// Add elements from other multiset, taking maximum count
	for _, entry := range other.EntrySet() {
		currentCount := result.Count(entry.Element)
		if entry.Count > currentCount {
			result.SetCount(entry.Element, entry.Count)
		}
	}
	
	return result
}

// Intersection returns a new multiset containing the intersection
func (ms *TreeMultiset[E]) Intersection(other Multiset[E]) Multiset[E] {
	result := NewTreeMultisetWithComparator(ms.cmp)
	
	for _, entry := range ms.EntrySet() {
		otherCount := other.Count(entry.Element)
		if otherCount > 0 {
			minCount := entry.Count
			if otherCount < minCount {
				minCount = otherCount
			}
			result.AddCount(entry.Element, minCount)
		}
	}
	
	return result
}

// Difference returns a new multiset containing elements in this but not in other
func (ms *TreeMultiset[E]) Difference(other Multiset[E]) Multiset[E] {
	result := NewTreeMultisetWithComparator(ms.cmp)
	
	for _, entry := range ms.EntrySet() {
		otherCount := other.Count(entry.Element)
		if entry.Count > otherCount {
			result.AddCount(entry.Element, entry.Count-otherCount)
		}
	}
	
	return result
}

// IsSubsetOf checks if this multiset is a subset of another
func (ms *TreeMultiset[E]) IsSubsetOf(other Multiset[E]) bool {
	for _, entry := range ms.EntrySet() {
		if other.Count(entry.Element) < entry.Count {
			return false
		}
	}
	return true
}

// IsSupersetOf checks if this multiset is a superset of another
func (ms *TreeMultiset[E]) IsSupersetOf(other Multiset[E]) bool {
	return other.IsSubsetOf(ms)
}

// String returns a string representation of the multiset
func (ms *TreeMultiset[E]) String() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	if ms.size == 0 {
		return "TreeMultiset[]"
	}
	
	var builder strings.Builder
	builder.WriteString("TreeMultiset[")
	
	entries := ms.EntrySet()
	for i, entry := range entries {
		if i > 0 {
			builder.WriteString(", ")
		}
		if entry.Count == 1 {
			builder.WriteString(fmt.Sprintf("%v", entry.Element))
		} else {
			builder.WriteString(fmt.Sprintf("%v x %d", entry.Element, entry.Count))
		}
	}
	
	builder.WriteString("]")
	return builder.String()
}

// AVL tree balancing methods
func (ms *TreeMultiset[E]) height(node *treeNode[E]) int {
	if node == nil {
		return 0
	}
	return node.height
}

func (ms *TreeMultiset[E]) updateHeight(node *treeNode[E]) {
	leftHeight := ms.height(node.left)
	rightHeight := ms.height(node.right)
	if leftHeight > rightHeight {
		node.height = leftHeight + 1
	} else {
		node.height = rightHeight + 1
	}
}

func (ms *TreeMultiset[E]) balanceFactor(node *treeNode[E]) int {
	if node == nil {
		return 0
	}
	return ms.height(node.left) - ms.height(node.right)
}

func (ms *TreeMultiset[E]) rotateRight(y *treeNode[E]) *treeNode[E] {
	x := y.left
	y.left = x.right
	x.right = y
	ms.updateHeight(y)
	ms.updateHeight(x)
	return x
}

func (ms *TreeMultiset[E]) rotateLeft(x *treeNode[E]) *treeNode[E] {
	y := x.right
	x.right = y.left
	y.left = x
	ms.updateHeight(x)
	ms.updateHeight(y)
	return y
}

func (ms *TreeMultiset[E]) balance(node *treeNode[E]) *treeNode[E] {
	ms.updateHeight(node)
	
	balance := ms.balanceFactor(node)
	
	// Left heavy
	if balance > 1 {
		if ms.balanceFactor(node.left) < 0 {
			node.left = ms.rotateLeft(node.left)
		}
		return ms.rotateRight(node)
	}
	
	// Right heavy
	if balance < -1 {
		if ms.balanceFactor(node.right) > 0 {
			node.right = ms.rotateRight(node.right)
		}
		return ms.rotateLeft(node)
	}
	
	return node
}

// treeMultisetIterator implements Iterator for TreeMultiset
type treeMultisetIterator[E comparable] struct {
	multiset *TreeMultiset[E]
	entries  []Entry[E]
	index    int
	current  int
}

func (it *treeMultisetIterator[E]) HasNext() bool {
	return it.index < len(it.entries) && (it.current < it.entries[it.index].Count || it.index+1 < len(it.entries))
}

func (it *treeMultisetIterator[E]) Next() (E, bool) {
	if !it.HasNext() {
		var zero E
		return zero, false
	}
	
	if it.current >= it.entries[it.index].Count {
		it.index++
		it.current = 0
	}
	
	element := it.entries[it.index].Element
	it.current++
	return element, true
}

func (it *treeMultisetIterator[E]) Reset() {
	it.entries = it.multiset.EntrySet()
	it.index = 0
	it.current = 0
}

func (it *treeMultisetIterator[E]) Remove() bool {
	if it.index >= len(it.entries) || it.current == 0 {
		return false
	}
	
	element := it.entries[it.index].Element
	it.multiset.Remove(element)
	
	// Refresh entries after removal
	it.entries = it.multiset.EntrySet()
	if it.index >= len(it.entries) {
		it.index = len(it.entries)
		it.current = 0
	} else if it.current > it.entries[it.index].Count {
		it.current = it.entries[it.index].Count
	}
	
	return true
}