package set

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/chenjianyu/collections/container/common"
)

// skipListNode 跳表节点
type skipListNode[E comparable] struct {
	value E
	next  []*skipListNode[E] // 各层的下一个节点
}

// ConcurrentSkipListSet 是一个基于跳表的线程安全有序集合
// 类似于Java中的ConcurrentSkipListSet
type ConcurrentSkipListSet[E comparable] struct {
	head       *skipListNode[E] // 头节点
	maxLevel   int              // 最大层数
	level      int              // 当前层数
	size       int              // 元素数量
	comparator func(a, b E) int // 比较函数
	mutex      sync.RWMutex     // 读写锁
	random     *rand.Rand       // 随机数生成器
}

const (
	maxLevels = 32   // 最大层数
	p         = 0.25 // 晋升概率
)

// NewConcurrentSkipListSet 创建一个新的ConcurrentSkipListSet
func NewConcurrentSkipListSet[E comparable]() *ConcurrentSkipListSet[E] {
	return NewConcurrentSkipListSetWithComparator(func(a, b E) int { return Compare(a, b) })
}

// NewConcurrentSkipListSetWithComparator 创建一个具有指定比较器的ConcurrentSkipListSet
func NewConcurrentSkipListSetWithComparator[E comparable](comparator func(a, b E) int) *ConcurrentSkipListSet[E] {
	head := &skipListNode[E]{
		next: make([]*skipListNode[E], maxLevels),
	}

	return &ConcurrentSkipListSet[E]{
		head:       head,
		maxLevel:   maxLevels,
		level:      1,
		size:       0,
		comparator: comparator,
		random:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewConcurrentSkipListSetFromSlice 从切片创建一个新的ConcurrentSkipListSet
func NewConcurrentSkipListSetFromSlice[E comparable](elements []E) *ConcurrentSkipListSet[E] {
	set := NewConcurrentSkipListSet[E]()
	for _, e := range elements {
		set.Add(e)
	}
	return set
}

// randomLevel 生成随机层数
func (set *ConcurrentSkipListSet[E]) randomLevel() int {
	level := 1
	for level < set.maxLevel && set.random.Float64() < p {
		level++
	}
	return level
}

// findPredecessors 查找每层的前驱节点
func (set *ConcurrentSkipListSet[E]) findPredecessors(value E) []*skipListNode[E] {
	predecessors := make([]*skipListNode[E], set.maxLevel)
	current := set.head

	// 从最高层开始向下搜索
	for i := set.level - 1; i >= 0; i-- {
		for current.next[i] != nil && set.comparator(current.next[i].value, value) < 0 {
			current = current.next[i]
		}
		predecessors[i] = current
	}

	return predecessors
}

// Add 添加一个元素到集合中
func (set *ConcurrentSkipListSet[E]) Add(e E) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	predecessors := set.findPredecessors(e)
	current := predecessors[0].next[0]

	// 如果元素已存在，返回false
	if current != nil && set.comparator(current.value, e) == 0 {
		return false
	}

	// 生成新节点的层数
	newLevel := set.randomLevel()

	// 如果新层数超过当前层数，更新层数
	if newLevel > set.level {
		for i := set.level; i < newLevel; i++ {
			predecessors[i] = set.head
		}
		set.level = newLevel
	}

	// 创建新节点
	newNode := &skipListNode[E]{
		value: e,
		next:  make([]*skipListNode[E], newLevel),
	}

	// 插入新节点
	for i := 0; i < newLevel; i++ {
		newNode.next[i] = predecessors[i].next[i]
		predecessors[i].next[i] = newNode
	}

	set.size++
	return true
}

// Remove 从集合中移除指定元素
func (set *ConcurrentSkipListSet[E]) Remove(e E) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	predecessors := set.findPredecessors(e)
	current := predecessors[0].next[0]

	// 如果元素不存在，返回false
	if current == nil || set.comparator(current.value, e) != 0 {
		return false
	}

	// 移除节点
	for i := 0; i < len(current.next); i++ {
		predecessors[i].next[i] = current.next[i]
	}

	// 更新层数
	for set.level > 1 && set.head.next[set.level-1] == nil {
		set.level--
	}

	set.size--
	return true
}

// Contains 检查集合是否包含指定元素
func (set *ConcurrentSkipListSet[E]) Contains(e E) bool {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	current := set.head

	// 从最高层开始向下搜索
	for i := set.level - 1; i >= 0; i-- {
		for current.next[i] != nil && set.comparator(current.next[i].value, e) < 0 {
			current = current.next[i]
		}
	}

	current = current.next[0]
	return current != nil && set.comparator(current.value, e) == 0
}

// Size 返回集合中的元素数量
func (set *ConcurrentSkipListSet[E]) Size() int {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	return set.size
}

// IsEmpty 如果集合不包含元素，返回true
func (set *ConcurrentSkipListSet[E]) IsEmpty() bool {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	return set.size == 0
}

// Clear 移除集合中的所有元素
func (set *ConcurrentSkipListSet[E]) Clear() {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	// 重新初始化头节点
	set.head = &skipListNode[E]{
		next: make([]*skipListNode[E], set.maxLevel),
	}
	set.level = 1
	set.size = 0
}

// ToSlice 返回包含集合中所有元素的切片（按排序顺序）
func (set *ConcurrentSkipListSet[E]) ToSlice() []E {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	result := make([]E, 0, set.size)
	current := set.head.next[0]

	for current != nil {
		result = append(result, current.value)
		current = current.next[0]
	}

	return result
}

// ForEach 对集合中的每个元素执行给定的操作（按排序顺序）
func (set *ConcurrentSkipListSet[E]) ForEach(f func(E)) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	current := set.head.next[0]
	for current != nil {
		f(current.value)
		current = current.next[0]
	}
}

// AddAll 将指定集合中的所有元素添加到此集合
func (set *ConcurrentSkipListSet[E]) AddAll(s Set[E]) bool {
	changed := false
	s.ForEach(func(e E) {
		if set.Add(e) {
			changed = true
		}
	})
	return changed
}

// RemoveAll 移除此集合中也包含在指定集合中的所有元素
func (set *ConcurrentSkipListSet[E]) RemoveAll(s Set[E]) bool {
	changed := false
	s.ForEach(func(e E) {
		if set.Remove(e) {
			changed = true
		}
	})
	return changed
}

// RetainAll 仅保留此集合中包含在指定集合中的元素
func (set *ConcurrentSkipListSet[E]) RetainAll(s Set[E]) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	toRemove := make([]E, 0)
	current := set.head.next[0]

	for current != nil {
		if !s.Contains(current.value) {
			toRemove = append(toRemove, current.value)
		}
		current = current.next[0]
	}

	changed := false
	for _, e := range toRemove {
		if set.removeUnsafe(e) {
			changed = true
		}
	}

	return changed
}

// removeUnsafe 在不加锁的情况下移除元素（内部使用）
func (set *ConcurrentSkipListSet[E]) removeUnsafe(e E) bool {
	predecessors := set.findPredecessors(e)
	current := predecessors[0].next[0]

	if current == nil || set.comparator(current.value, e) != 0 {
		return false
	}

	for i := 0; i < len(current.next); i++ {
		predecessors[i].next[i] = current.next[i]
	}

	for set.level > 1 && set.head.next[set.level-1] == nil {
		set.level--
	}

	set.size--
	return true
}

// ContainsAll 如果此集合包含指定集合中的所有元素，返回true
func (set *ConcurrentSkipListSet[E]) ContainsAll(s Set[E]) bool {
	result := true
	s.ForEach(func(e E) {
		if !set.Contains(e) {
			result = false
		}
	})
	return result
}

// First 返回集合中的第一个（最小）元素
func (set *ConcurrentSkipListSet[E]) First() (E, bool) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	if set.size == 0 {
		return *new(E), false
	}

	return set.head.next[0].value, true
}

// Last 返回集合中的最后一个（最大）元素
func (set *ConcurrentSkipListSet[E]) Last() (E, bool) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	if set.size == 0 {
		return *new(E), false
	}

	current := set.head
	for i := set.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			current = current.next[i]
		}
	}

	return current.value, true
}

// Lower 返回集合中严格小于给定元素的最大元素
func (set *ConcurrentSkipListSet[E]) Lower(e E) (E, bool) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	predecessors := set.findPredecessors(e)
	predecessor := predecessors[0]

	if predecessor == set.head {
		return *new(E), false
	}

	return predecessor.value, true
}

// Higher 返回集合中严格大于给定元素的最小元素
func (set *ConcurrentSkipListSet[E]) Higher(e E) (E, bool) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	current := set.head

	for i := set.level - 1; i >= 0; i-- {
		for current.next[i] != nil && set.comparator(current.next[i].value, e) <= 0 {
			current = current.next[i]
		}
	}

	if current.next[0] != nil {
		return current.next[0].value, true
	}

	return *new(E), false
}

// Floor 返回集合中小于或等于给定元素的最大元素
func (set *ConcurrentSkipListSet[E]) Floor(e E) (E, bool) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	predecessors := set.findPredecessors(e)
	current := predecessors[0].next[0]

	// 如果找到相等的元素，返回它
	if current != nil && set.comparator(current.value, e) == 0 {
		return current.value, true
	}

	// 否则返回前驱
	predecessor := predecessors[0]
	if predecessor == set.head {
		return *new(E), false
	}

	return predecessor.value, true
}

// Ceiling 返回集合中大于或等于给定元素的最小元素
func (set *ConcurrentSkipListSet[E]) Ceiling(e E) (E, bool) {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	current := set.head

	for i := set.level - 1; i >= 0; i-- {
		for current.next[i] != nil && set.comparator(current.next[i].value, e) < 0 {
			current = current.next[i]
		}
	}

	if current.next[0] != nil {
		return current.next[0].value, true
	}

	return *new(E), false
}

// PollFirst 获取并移除第一个（最小）元素
func (set *ConcurrentSkipListSet[E]) PollFirst() (E, bool) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.size == 0 {
		return *new(E), false
	}

	first := set.head.next[0].value
	set.removeUnsafe(first)
	return first, true
}

// PollLast 获取并移除最后一个（最大）元素
func (set *ConcurrentSkipListSet[E]) PollLast() (E, bool) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.size == 0 {
		return *new(E), false
	}

	// 找到最后一个元素
	current := set.head
	for i := set.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			current = current.next[i]
		}
	}

	last := current.value
	set.removeUnsafe(last)
	return last, true
}

// String 返回集合的字符串表示
func (set *ConcurrentSkipListSet[E]) String() string {
	set.mutex.RLock()
	defer set.mutex.RUnlock()

	if set.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	current := set.head.next[0]
	first := true
	for current != nil {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", current.value))
		current = current.next[0]
		first = false
	}

	sb.WriteString("]")
	return sb.String()
}

// Difference 返回此集合与另一个集合的差集
func (set *ConcurrentSkipListSet[E]) Difference(other Set[E]) Set[E] {
	result := NewConcurrentSkipListSet[E]()
	set.ForEach(func(e E) {
		if !other.Contains(e) {
			result.Add(e)
		}
	})
	return result
}

// Union 返回此集合与另一个集合的并集
func (set *ConcurrentSkipListSet[E]) Union(other Set[E]) Set[E] {
	result := NewConcurrentSkipListSetWithComparator(set.comparator)

	// 添加当前集合的所有元素
	set.ForEach(func(e E) {
		result.Add(e)
	})

	// 添加另一个集合的所有元素
	other.ForEach(func(e E) {
		result.Add(e)
	})

	return result
}

// Intersection 返回此集合与另一个集合的交集
func (set *ConcurrentSkipListSet[E]) Intersection(other Set[E]) Set[E] {
	result := NewConcurrentSkipListSet[E]()
	set.ForEach(func(e E) {
		if other.Contains(e) {
			result.Add(e)
		}
	})
	return result
}

// IsSubsetOf 检查此集合是否为另一个集合的子集
func (set *ConcurrentSkipListSet[E]) IsSubsetOf(other Set[E]) bool {
	result := true
	set.ForEach(func(e E) {
		if !other.Contains(e) {
			result = false
		}
	})
	return result
}

// IsSupersetOf 检查此集合是否为另一个集合的超集
func (set *ConcurrentSkipListSet[E]) IsSupersetOf(other Set[E]) bool {
	return other.IsSubsetOf(set)
}

// skipListIterator 跳表迭代器
type skipListIterator[E comparable] struct {
	current *skipListNode[E]
}

func (it *skipListIterator[E]) HasNext() bool {
	return it.current != nil
}

func (it *skipListIterator[E]) Next() (E, bool) {
	if it.current == nil {
		var zero E
		return zero, false
	}
	val := it.current.value
	it.current = it.current.next[0]
	return val, true
}

func (it *skipListIterator[E]) Remove() bool {
	// 跳表不支持迭代器移除
	return false
}

// Iterator 返回一个迭代器，用于遍历集合中的元素
func (set *ConcurrentSkipListSet[E]) Iterator() common.Iterator[E] {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	return &skipListIterator[E]{current: set.head.next[0]}
}
