// Package stack 提供栈数据结构的实现
package stack

import (
	"github.com/chenjianyu/collections/container/common"
)

// Stack 表示一个后进先出(LIFO)的栈
type Stack[E any] interface {
	common.Container[E]

	// Push 将元素推入栈顶
	// 如果栈已满（对于有界栈），则返回错误
	Push(e E) error

	// Pop 移除并返回栈顶的元素
	// 如果栈为空，则返回错误
	Pop() (E, error)

	// Peek 返回栈顶的元素，但不移除
	// 如果栈为空，则返回错误
	Peek() (E, error)

	// Search 搜索元素在栈中的位置
	// 如果找到，返回从1开始的位置索引（1表示栈顶）
	// 如果未找到，返回-1
	Search(e E) int

	// ToSlice 返回包含栈中所有元素的切片
	// 切片中的第一个元素是栈底元素，最后一个元素是栈顶元素
	ToSlice() []E
}
