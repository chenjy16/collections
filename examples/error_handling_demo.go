package main

import (
	"errors"
	"fmt"
	"github.com/chenjianyu/collections/container/common"
	"github.com/chenjianyu/collections/container/list"
	"github.com/chenjianyu/collections/container/queue"
	"github.com/chenjianyu/collections/container/stack"
)

func main() {
	fmt.Println("=== 错误处理优化演示 ===\n")

	// 演示 ArrayList 的详细错误信息
	fmt.Println("1. ArrayList 错误处理演示:")
	arrayList := list.New[int]()
	
	// 尝试访问空列表
	_, err := arrayList.Get(0)
	if err != nil {
		fmt.Printf("   访问空列表错误: %v\n", err)
		if errors.Is(err, common.ErrIndexOutOfBounds) {
			fmt.Printf("   ✓ 错误类型匹配: IndexOutOfBounds\n")
		}
	}
	
	// 添加一些元素后尝试无效索引
	arrayList.Add(1)
	arrayList.Add(2)
	_, err = arrayList.Get(5)
	if err != nil {
		fmt.Printf("   无效索引错误: %v\n", err)
	}
	
	// 尝试无效范围
	_, err = arrayList.SubList(2, 1)
	if err != nil {
		fmt.Printf("   无效范围错误: %v\n", err)
	}

	fmt.Println()

	// 演示 LinkedList 的详细错误信息
	fmt.Println("2. LinkedList 错误处理演示:")
	linkedList := list.NewLinkedList[string]()
	
	// 尝试获取空列表的第一个元素
	_, err = linkedList.GetFirst()
	if err != nil {
		fmt.Printf("   获取空列表首元素错误: %v\n", err)
		if errors.Is(err, common.ErrEmptyContainer) {
			fmt.Printf("   ✓ 错误类型匹配: EmptyContainer\n")
		}
	}

	fmt.Println()

	// 演示 ArrayStack 的详细错误信息
	fmt.Println("3. ArrayStack 错误处理演示:")
	
	// 有容量限制的栈
	limitedStack := stack.WithCapacity[int](2)
	limitedStack.Push(1)
	limitedStack.Push(2)
	
	// 尝试向满栈推入元素
	err = limitedStack.Push(3)
	if err != nil {
		fmt.Printf("   满栈推入错误: %v\n", err)
		if errors.Is(err, common.ErrFullContainer) {
			fmt.Printf("   ✓ 错误类型匹配: FullContainer\n")
		}
	}
	
	// 空栈弹出
	emptyStack := stack.New[int]()
	_, err = emptyStack.Pop()
	if err != nil {
		fmt.Printf("   空栈弹出错误: %v\n", err)
	}

	fmt.Println()

	// 演示 PriorityQueue 的详细错误信息
	fmt.Println("4. PriorityQueue 错误处理演示:")
	
	// 空队列移除
	emptyQueue := queue.NewPriorityQueueWithComparator[int](func(a, b int) int {
		return a - b
	})
	_, err = emptyQueue.Remove()
	if err != nil {
		fmt.Printf("   空队列移除错误: %v\n", err)
	}

	fmt.Println()

	// 演示 LinkedListQueue 的详细错误信息
	fmt.Println("5. LinkedListQueue 错误处理演示:")
	
	// 空队列操作
	emptyLLQueue := queue.New[string]()
	_, err = emptyLLQueue.GetFirst()
	if err != nil {
		fmt.Printf("   空队列获取首元素错误: %v\n", err)
	}

	fmt.Println()

	// 演示错误工厂函数的使用
	fmt.Println("6. 错误工厂函数演示:")
	
	// 创建具体的错误实例
	indexErr := common.IndexOutOfBoundsError(10, 5)
	fmt.Printf("   索引越界错误: %v\n", indexErr)
	
	rangeErr := common.InvalidRangeError(5, 2)
	fmt.Printf("   无效范围错误: %v\n", rangeErr)
	
	countErr := common.NegativeCountError(-3)
	fmt.Printf("   负数计数错误: %v\n", countErr)
	
	keyErr := common.KeyNotFoundError("missing_key")
	fmt.Printf("   键未找到错误: %v\n", keyErr)
	
	elementErr := common.ElementNotFoundError(42)
	fmt.Printf("   元素未找到错误: %v\n", elementErr)

	fmt.Println("\n=== 错误处理优化完成 ===")
	fmt.Println("✓ 统一了所有容器的错误定义")
	fmt.Println("✓ 提供了详细的错误信息")
	fmt.Println("✓ 支持错误类型检查")
	fmt.Println("✓ 提供了错误工厂函数")
}