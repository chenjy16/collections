package main

import (
	"fmt"
	"github.com/chenjianyu/collections/container/set"
)

func LinkedHashSetExample() {
	fmt.Println("=== LinkedHashSet 示例 ===")
	
	// 创建一个新的 LinkedHashSet
	linkedSet := set.NewLinkedHashSet[string]()
	
	// 添加元素 - 保持插入顺序
	fmt.Println("\n1. 添加元素（保持插入顺序）:")
	elements := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, element := range elements {
		added := linkedSet.Add(element)
		fmt.Printf("添加 '%s': %t\n", element, added)
	}
	
	// 尝试添加重复元素
	fmt.Println("\n2. 尝试添加重复元素:")
	duplicate := linkedSet.Add("apple")
	fmt.Printf("添加重复元素 'apple': %t\n", duplicate)
	
	// 显示集合内容（保持插入顺序）
	fmt.Printf("\n3. 集合内容（插入顺序）: %s\n", linkedSet.String())
	fmt.Printf("集合大小: %d\n", linkedSet.Size())
	
	// 检查元素是否存在
	fmt.Println("\n4. 检查元素是否存在:")
	testElements := []string{"banana", "grape", "cherry"}
	for _, element := range testElements {
		contains := linkedSet.Contains(element)
		fmt.Printf("包含 '%s': %t\n", element, contains)
	}
	
	// 使用 ForEach 遍历（保持插入顺序）
	fmt.Println("\n5. 使用 ForEach 遍历:")
	linkedSet.ForEach(func(element string) {
		fmt.Printf("- %s\n", element)
	})
	
	// 使用迭代器遍历（保持插入顺序）
	fmt.Println("\n6. 使用迭代器遍历:")
	iterator := linkedSet.Iterator()
	for iterator.HasNext() {
		if element, ok := iterator.Next(); ok {
			fmt.Printf("- %s\n", element)
		}
	}
	
	// 删除元素
	fmt.Println("\n7. 删除元素:")
	removed := linkedSet.Remove("cherry")
	fmt.Printf("删除 'cherry': %t\n", removed)
	fmt.Printf("删除后的集合: %s\n", linkedSet.String())
	
	// 添加新元素到末尾
	fmt.Println("\n8. 添加新元素:")
	linkedSet.Add("fig")
	fmt.Printf("添加 'fig' 后的集合: %s\n", linkedSet.String())
	
	// 转换为切片（保持插入顺序）
	fmt.Println("\n9. 转换为切片:")
	slice := linkedSet.ToSlice()
	fmt.Printf("切片: %v\n", slice)
	
	// 集合操作
	fmt.Println("\n10. 集合操作:")
	
	// 创建另一个 LinkedHashSet
	otherSet := set.LinkedHashSetFromSlice([]string{"cherry", "date", "fig", "grape"})
	fmt.Printf("另一个集合: %s\n", otherSet.String())
	
	// 并集（保持第一个集合的顺序，然后是第二个集合的新元素）
	union := linkedSet.Union(otherSet)
	fmt.Printf("并集: %s\n", union.String())
	
	// 交集（保持第一个集合的顺序）
	intersection := linkedSet.Intersection(otherSet)
	fmt.Printf("交集: %s\n", intersection.String())
	
	// 差集
	difference := linkedSet.Difference(otherSet)
	fmt.Printf("差集: %s\n", difference.String())
	
	// 子集检查
	subset := set.LinkedHashSetFromSlice([]string{"apple", "banana"})
	fmt.Printf("子集 %s 是否为原集合的子集: %t\n", subset.String(), subset.IsSubsetOf(linkedSet))
	
	// 与普通 HashSet 的对比
	fmt.Println("\n11. 与普通 HashSet 的对比:")
	
	// 创建普通 HashSet
	hashSet := set.New[string]()
	for _, element := range elements {
		hashSet.Add(element)
	}
	
	fmt.Printf("LinkedHashSet（保持插入顺序）: %s\n", linkedSet.String())
	fmt.Printf("HashSet（无序）: %s\n", hashSet.String())
	
	// 清空集合
	fmt.Println("\n12. 清空集合:")
	linkedSet.Clear()
	fmt.Printf("清空后的集合: %s\n", linkedSet.String())
	fmt.Printf("是否为空: %t\n", linkedSet.IsEmpty())
	
	// 演示大数据集的插入顺序维护
	fmt.Println("\n13. 大数据集插入顺序维护:")
	largeSet := set.NewLinkedHashSet[int]()
	
	// 添加一些数字
	numbers := []int{42, 17, 89, 3, 56, 21, 78, 9, 34, 67}
	for _, num := range numbers {
		largeSet.Add(num)
	}
	
	fmt.Printf("原始顺序: %v\n", numbers)
	fmt.Printf("LinkedHashSet: %s\n", largeSet.String())
	fmt.Println("可以看到 LinkedHashSet 完全保持了插入顺序！")
}