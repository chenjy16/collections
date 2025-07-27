// Package queue 提供队列数据结构的实现
package queue

import (
	"github.com/chenjianyu/collections/container/common"
)

// Queue 表示一个先进先出(FIFO)的队列
type Queue[E any] interface {
	common.Container[E]

	// Add 将元素添加到队列尾部
	// 如果队列已满（对于有界队列），则返回错误
	Add(e E) error

	// Offer 将元素添加到队列尾部
	// 如果队列已满（对于有界队列），则返回false
	Offer(e E) bool

	// Remove 移除并返回队列头部的元素
	// 如果队列为空，则返回错误
	Remove() (E, error)

	// Poll 移除并返回队列头部的元素
	// 如果队列为空，则返回零值和false
	Poll() (E, bool)

	// Element 返回队列头部的元素，但不移除
	// 如果队列为空，则返回错误
	Element() (E, error)

	// Peek 返回队列头部的元素，但不移除
	// 如果队列为空，则返回零值和false
	Peek() (E, bool)
}

// Deque 表示一个双端队列，支持在两端添加和移除元素
type Deque[E any] interface {
	Queue[E]

	// AddFirst 将元素添加到队列头部
	// 如果队列已满（对于有界队列），则返回错误
	AddFirst(e E) error

	// AddLast 将元素添加到队列尾部
	// 如果队列已满（对于有界队列），则返回错误
	AddLast(e E) error

	// OfferFirst 将元素添加到队列头部
	// 如果队列已满（对于有界队列），则返回false
	OfferFirst(e E) bool

	// OfferLast 将元素添加到队列尾部
	// 如果队列已满（对于有界队列），则返回false
	OfferLast(e E) bool

	// RemoveFirst 移除并返回队列头部的元素
	// 如果队列为空，则返回错误
	RemoveFirst() (E, error)

	// RemoveLast 移除并返回队列尾部的元素
	// 如果队列为空，则返回错误
	RemoveLast() (E, error)

	// PollFirst 移除并返回队列头部的元素
	// 如果队列为空，则返回零值和false
	PollFirst() (E, bool)

	// PollLast 移除并返回队列尾部的元素
	// 如果队列为空，则返回零值和false
	PollLast() (E, bool)

	// GetFirst 返回队列头部的元素，但不移除
	// 如果队列为空，则返回错误
	GetFirst() (E, error)

	// GetLast 返回队列尾部的元素，但不移除
	// 如果队列为空，则返回错误
	GetLast() (E, error)

	// PeekFirst 返回队列头部的元素，但不移除
	// 如果队列为空，则返回零值和false
	PeekFirst() (E, bool)

	// PeekLast 返回队列尾部的元素，但不移除
	// 如果队列为空，则返回零值和false
	PeekLast() (E, bool)
}
