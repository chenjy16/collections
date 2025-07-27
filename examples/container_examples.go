// Package main 提供使用Gkit容器库的示例
package main

import (
	"fmt"

	"github.com/chenjianyu/collections/container/list"
	"github.com/chenjianyu/collections/container/queue"
	setpkg "github.com/chenjianyu/collections/container/set"
	"github.com/chenjianyu/collections/container/stack"
)

// RunContainerExamples 运行所有容器示例
func RunContainerExamples() {
	fmt.Println("=== 运行Gkit容器库示例 ===")

	ArrayListExample()
	HashSetExample()
	LinkedListExample()
	PriorityQueueExample()
	ArrayStackExample()
}

// ArrayListExample 展示ArrayList的用法
func ArrayListExample() {
	fmt.Println("\n=== ArrayList示例 ===")

	// 创建一个新的ArrayList
	arrayList := list.New[int]()

	// 添加元素
	fmt.Println("添加元素:")
	for i := 0; i < 5; i++ {
		arrayList.Add(i)
		fmt.Printf("添加 %d, 大小: %d\n", i, arrayList.Size())
	}

	// 在指定位置插入元素
	fmt.Println("\n在索引2处插入100:")
	arrayList.Insert(2, 100)
	fmt.Printf("ArrayList: %s\n", arrayList)

	// 获取元素
	fmt.Println("\n获取元素:")
	val, err := arrayList.Get(2)
	if err == nil {
		fmt.Printf("索引2处的元素: %d\n", val)
	}

	// 设置元素
	fmt.Println("\n设置元素:")
	arrayList.Set(2, 200)
	fmt.Printf("设置后的ArrayList: %s\n", arrayList)

	// 检查包含
	fmt.Println("\n检查包含:")
	fmt.Printf("ArrayList包含200: %t\n", arrayList.Contains(200))
	fmt.Printf("ArrayList包含100: %t\n", arrayList.Contains(100))

	// 查找索引
	fmt.Println("\n查找索引:")
	fmt.Printf("200的索引: %d\n", arrayList.IndexOf(200))
	fmt.Printf("100的索引: %d\n", arrayList.IndexOf(100))

	// 移除元素
	fmt.Println("\n移除元素:")
	removed, _ := arrayList.RemoveAt(2)
	fmt.Printf("移除索引2处的元素: %d\n", removed)
	fmt.Printf("移除后的ArrayList: %s\n", arrayList)

	// 创建子列表
	fmt.Println("\n创建子列表:")
	subList, _ := arrayList.SubList(1, 3)
	fmt.Printf("子列表[1:3]: %s\n", subList)

	// 转换为切片
	fmt.Println("\n转换为切片:")
	slice := arrayList.ToSlice()
	fmt.Printf("切片: %v\n", slice)

	// 使用ForEach
	fmt.Println("\n使用ForEach:")
	fmt.Print("元素: ")
	arrayList.ForEach(func(e int) {
		fmt.Printf("%d ", e)
	})
	fmt.Println()

	// 清空列表
	fmt.Println("\n清空列表:")
	arrayList.Clear()
	fmt.Printf("清空后的大小: %d\n", arrayList.Size())
	fmt.Printf("是否为空: %t\n", arrayList.IsEmpty())
}

// HashSetExample 展示HashSet的用法
func HashSetExample() {
	fmt.Println("\n=== HashSet示例 ===")

	// 创建一个新的HashSet
	set := setpkg.New[string]()

	// 添加元素
	fmt.Println("添加元素:")
	fruits := []string{"苹果", "香蕉", "橙子", "苹果", "葡萄"}
	for _, fruit := range fruits {
		set.Add(fruit)
		fmt.Printf("添加 '%s', 大小: %d\n", fruit, set.Size())
	}

	// 检查包含
	fmt.Println("\n检查包含:")
	fmt.Printf("包含'香蕉': %t\n", set.Contains("香蕉"))
	fmt.Printf("包含'西瓜': %t\n", set.Contains("西瓜"))

	// 移除元素
	fmt.Println("\n移除元素:")
	set.Remove("香蕉")
	fmt.Printf("移除'香蕉'后: %s\n", set)

	// 创建另一个集合
	fmt.Println("\n创建另一个集合:")
	set2 := setpkg.FromSlice([]string{"葡萄", "西瓜", "柠檬"})
	fmt.Printf("集合2: %s\n", set2)

	// 集合操作
	fmt.Println("\n集合操作:")
	unionSet := set.Union(set2)
	fmt.Printf("并集: %s\n", unionSet)

	intersectionSet := set.Intersection(set2)
	fmt.Printf("交集: %s\n", intersectionSet)

	differenceSet := set.Difference(set2)
	fmt.Printf("差集: %s\n", differenceSet)

	// 子集和超集
	fmt.Println("\n子集和超集:")
	subset := setpkg.FromSlice([]string{"苹果"})
	fmt.Printf("子集: %s\n", subset)
	fmt.Printf("是子集: %t\n", subset.IsSubsetOf(set))
	fmt.Printf("是超集: %t\n", set.IsSupersetOf(subset))

	// 转换为切片
	fmt.Println("\n转换为切片:")
	slice := set.ToSlice()
	fmt.Printf("切片: %v\n", slice)

	// 使用ForEach
	fmt.Println("\n使用ForEach:")
	fmt.Print("元素: ")
	set.ForEach(func(e string) {
		fmt.Printf("%s ", e)
	})
	fmt.Println()

	// 清空集合
	fmt.Println("\n清空集合:")
	set.Clear()
	fmt.Printf("清空后的大小: %d\n", set.Size())
	fmt.Printf("是否为空: %t\n", set.IsEmpty())
}

// LinkedListExample 展示LinkedList的用法
func LinkedListExample() {
	fmt.Println("\n=== LinkedList示例 ===")

	// 创建一个新的LinkedList
	list := queue.New[int]()

	// 添加元素
	fmt.Println("添加元素:")
	for i := 0; i < 5; i++ {
		_ = list.Add(i)
		fmt.Printf("添加 %d, 大小: %d\n", i, list.Size())
	}

	// 作为队列使用
	fmt.Println("\n作为队列使用:")
	val, _ := list.Peek()
	fmt.Printf("队列头部元素: %d\n", val)

	val, _ = list.Remove()
	fmt.Printf("移除队列头部元素: %d\n", val)
	fmt.Printf("移除后的队列: %s\n", list)

	// 作为双端队列使用
	fmt.Println("\n作为双端队列使用:")
	_ = list.AddFirst(10)
	_ = list.AddLast(20)
	fmt.Printf("添加元素后的双端队列: %s\n", list)

	val, _ = list.GetFirst()
	fmt.Printf("双端队列头部元素: %d\n", val)

	val, _ = list.GetLast()
	fmt.Printf("双端队列尾部元素: %d\n", val)

	val, _ = list.RemoveFirst()
	fmt.Printf("移除双端队列头部元素: %d\n", val)

	val, _ = list.RemoveLast()
	fmt.Printf("移除双端队列尾部元素: %d\n", val)

	fmt.Printf("移除后的双端队列: %s\n", list)

	// 转换为切片
	fmt.Println("\n转换为切片:")
	slice := list.ToSlice()
	fmt.Printf("切片: %v\n", slice)

	// 使用ForEach
	fmt.Println("\n使用ForEach:")
	fmt.Print("元素: ")
	list.ForEach(func(e int) {
		fmt.Printf("%d ", e)
	})
	fmt.Println()

	// 清空列表
	fmt.Println("\n清空列表:")
	list.Clear()
	fmt.Printf("清空后的大小: %d\n", list.Size())
	fmt.Printf("是否为空: %t\n", list.IsEmpty())
}

// PriorityQueueExample 展示PriorityQueue的用法
func PriorityQueueExample() {
	fmt.Println("\n=== PriorityQueue示例 ===")

	// 创建一个最小堆优先队列
	minHeap := queue.NewPriorityQueueWithComparator(func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 添加元素（无序）
	fmt.Println("添加元素:")
	elements := []int{5, 3, 1, 4, 2}
	for _, e := range elements {
		_ = minHeap.Add(e)
		fmt.Printf("添加 %d, 大小: %d\n", e, minHeap.Size())
	}

	// 查看最小元素
	fmt.Println("\n查看最小元素:")
	val, _ := minHeap.Peek()
	fmt.Printf("最小元素: %d\n", val)

	// 移除并获取最小元素
	fmt.Println("\n移除并获取最小元素:")
	for !minHeap.IsEmpty() {
		val, _ := minHeap.Remove()
		fmt.Printf("移除: %d\n", val)
	}

	// 创建一个最大堆优先队列
	fmt.Println("\n创建最大堆:")
	maxHeap := queue.NewPriorityQueueWithComparator(func(a, b int) int {
		if a > b {
			return -1 // 反转比较结果
		} else if a < b {
			return 1
		}
		return 0
	})

	// 添加元素
	for _, e := range elements {
		_ = maxHeap.Add(e)
	}

	// 获取排序后的切片
	fmt.Println("\n获取排序后的切片:")
	sortedSlice := maxHeap.ToSortedSlice()
	fmt.Printf("排序后的切片（降序）: %v\n", sortedSlice)

	// 移除并获取最大元素
	fmt.Println("\n移除并获取最大元素:")
	for !maxHeap.IsEmpty() {
		val, _ := maxHeap.Remove()
		fmt.Printf("移除: %d\n", val)
	}
}

// ArrayStackExample 展示ArrayStack的用法
func ArrayStackExample() {
	fmt.Println("\n=== ArrayStack示例 ===")

	// 创建一个新的ArrayStack
	stack := stack.New[int]()

	// 压入元素
	fmt.Println("压入元素:")
	for i := 0; i < 5; i++ {
		_ = stack.Push(i)
		fmt.Printf("压入 %d, 大小: %d\n", i, stack.Size())
	}

	// 查看栈顶元素
	fmt.Println("\n查看栈顶元素:")
	val, _ := stack.Peek()
	fmt.Printf("栈顶元素: %d\n", val)

	// 弹出元素
	fmt.Println("\n弹出元素:")
	val, _ = stack.Pop()
	fmt.Printf("弹出栈顶元素: %d\n", val)
	fmt.Printf("弹出后的栈: %s\n", stack)

	// 搜索元素
	fmt.Println("\n搜索元素:")
	pos := stack.Search(2)
	fmt.Printf("元素2的位置: %d\n", pos)

	pos = stack.Search(10)
	fmt.Printf("元素10的位置: %d\n", pos)

	// 转换为切片
	fmt.Println("\n转换为切片:")
	slice := stack.ToSlice()
	fmt.Printf("切片: %v\n", slice)

	// 使用ForEach
	fmt.Println("\n使用ForEach:")
	fmt.Print("元素: ")
	stack.ForEach(func(e int) {
		fmt.Printf("%d ", e)
	})
	fmt.Println()

	// 清空栈
	fmt.Println("\n清空栈:")
	stack.Clear()
	fmt.Printf("清空后的大小: %d\n", stack.Size())
	fmt.Printf("是否为空: %t\n", stack.IsEmpty())
}
