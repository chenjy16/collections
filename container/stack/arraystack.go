// Package stack 提供栈数据结构的实现
package stack

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptyStack 表示栈为空错误
	ErrEmptyStack = errors.New("stack is empty")

	// ErrFullStack 表示栈已满错误
	ErrFullStack = errors.New("stack is full")
)

// ArrayStack 是一个基于切片的栈实现
type ArrayStack[E comparable] struct {
	elements []E
	maxCap   int // 最大容量，0表示无界
}

// New 创建一个新的无界ArrayStack
func New[E comparable]() *ArrayStack[E] {
	return &ArrayStack[E]{
		elements: make([]E, 0),
		maxCap:   0,
	}
}

// WithCapacity 创建一个具有指定最大容量的ArrayStack
func WithCapacity[E comparable](maxCapacity int) *ArrayStack[E] {
	return &ArrayStack[E]{
		elements: make([]E, 0, maxCapacity),
		maxCap:   maxCapacity,
	}
}

// FromSlice 从切片创建一个新的ArrayStack
// 切片中的第一个元素将是栈底元素，最后一个元素将是栈顶元素
func FromSlice[E comparable](elements []E) *ArrayStack[E] {
	stack := New[E]()
	stack.elements = make([]E, len(elements))
	copy(stack.elements, elements)
	return stack
}

// Size 返回栈中的元素数量
func (s *ArrayStack[E]) Size() int {
	return len(s.elements)
}

// IsEmpty 检查栈是否为空
func (s *ArrayStack[E]) IsEmpty() bool {
	return len(s.elements) == 0
}

// isFull 检查栈是否已满
func (s *ArrayStack[E]) isFull() bool {
	return s.maxCap > 0 && len(s.elements) >= s.maxCap
}

// Clear 清空栈
func (s *ArrayStack[E]) Clear() {
	s.elements = s.elements[:0]
}

// Contains 检查栈是否包含指定元素
func (s *ArrayStack[E]) Contains(e E) bool {
	for _, v := range s.elements {
		if v == e {
			return true
		}
	}
	return false
}

// ForEach 对栈中的每个元素执行给定的操作
// 遍历顺序是从栈底到栈顶
func (s *ArrayStack[E]) ForEach(f func(E)) {
	for _, e := range s.elements {
		f(e)
	}
}

// String 返回栈的字符串表示
func (s *ArrayStack[E]) String() string {
	if s.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")

	for i, e := range s.elements {
		sb.WriteString(fmt.Sprintf("%v", e))
		if i < len(s.elements)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]")
	return sb.String()
}

// Push 将元素推入栈顶
func (s *ArrayStack[E]) Push(e E) error {
	if s.isFull() {
		return ErrFullStack
	}

	s.elements = append(s.elements, e)
	return nil
}

// Pop 移除并返回栈顶的元素
func (s *ArrayStack[E]) Pop() (E, error) {
	if s.IsEmpty() {
		return *new(E), ErrEmptyStack
	}

	lastIndex := len(s.elements) - 1
	result := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]

	return result, nil
}

// Peek 返回栈顶的元素，但不移除
func (s *ArrayStack[E]) Peek() (E, error) {
	if s.IsEmpty() {
		return *new(E), ErrEmptyStack
	}

	return s.elements[len(s.elements)-1], nil
}

// Search 搜索元素在栈中的位置
func (s *ArrayStack[E]) Search(e E) int {
	for i := len(s.elements) - 1; i >= 0; i-- {
		if s.elements[i] == e {
			return len(s.elements) - i
		}
	}
	return -1
}

// ToSlice 返回包含栈中所有元素的切片
func (s *ArrayStack[E]) ToSlice() []E {
	result := make([]E, len(s.elements))
	copy(result, s.elements)
	return result
}
