# Collections - Go Collection Library | Go 集合库

[English](#english) | [中文](#chinese)

---

## English

Collections is a Go-based collection library that provides common collection types and tools found in Java SDK but missing from Go's standard library. This library follows Go language conventions and idioms, providing rich collection operation functionality for Go developers.

### Features

- Rich collection types: List, Set, Map, Queue, Stack, etc.
- Generic support (based on Go 1.18+)
- High-performance implementation
- Go-idiomatic API design
- Comprehensive documentation and testing
- Thread-safe implementations available
- Mathematical set operations (union, intersection, difference)
- Iterator pattern support

### Project Statistics

- **Total Go Files**: 28
- **Source Code Lines**: 4,731 (excluding tests)
- **Test Code Lines**: 3,915
- **Test Files**: 9
- **Average Test Coverage**: 70%+ (some modules reach 100%)

---

## Chinese

Collections 是一个基于Go语言的集合库，提供了Java SDK中常见但Go标准库中缺少的集合类型和工具。该库遵循Go语言的代码规范和惯例，为Go开发者提供丰富的集合操作功能。

### 特性

- 提供丰富的集合类型：List、Set、Map、Queue、Stack等
- 支持泛型（基于Go 1.18+）
- 高性能实现
- 符合Go语言习惯的API设计
- 完整的文档和测试
- 提供线程安全的实现
- 支持数学集合运算（并集、交集、差集）
- 迭代器模式支持

### 项目统计

- **Go文件总数**: 28个
- **源代码行数**: 4,731行（不含测试）
- **测试代码行数**: 3,915行
- **测试文件数**: 9个
- **平均测试覆盖率**: 70%+（部分模块达到100%）

### Installation

```bash
go get github.com/chenjianyu/collections
```

### Usage Examples

```go
package main

import (
	"fmt"
	"github.com/chenjianyu/collections/container/list"
	"github.com/chenjianyu/collections/container/set"
	"github.com/chenjianyu/collections/container/map"
)

func main() {
	// 创建ArrayList
	list := list.New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	fmt.Println("List size:", list.Size())
	fmt.Println("List contains 2:", list.Contains(2))

	// 创建HashSet
	set := set.New[string]()
	set.Add("apple")
	set.Add("banana")
	set.Add("apple") // 重复元素不会被添加

	fmt.Println("Set size:", set.Size())
	fmt.Println("Set elements:")
	set.ForEach(func(item string) {
		fmt.Println(item)
	})

	// 创建HashMap
	hashMap := maps.NewLinkedHashMap[string, int]()
	hashMap.Put("one", 1)
	hashMap.Put("two", 2)

	fmt.Println("Map size:", hashMap.Size())
	if val, ok := hashMap.Get("one"); ok {
		fmt.Println("Value for 'one':", val)
	}

	// 创建ConcurrentSkipListSet（线程安全的有序集合）
	skipSet := set.NewConcurrentSkipListSet[int]()
	skipSet.Add(5)
	skipSet.Add(2)
	skipSet.Add(8)
	skipSet.Add(1)

	fmt.Println("SkipSet (sorted):", skipSet) // 输出: [1, 2, 5, 8]
	if first, ok := skipSet.First(); ok {
		fmt.Println("First element:", first) // 输出: 1
	}
	if last, ok := skipSet.Last(); ok {
		fmt.Println("Last element:", last) // 输出: 8
	}

	// 创建TreeMap（有序映射）
	treeMap := maps.NewTreeMap[string, int]()
	treeMap.Put("one", 1)
	treeMap.Put("two", 2)
	treeMap.Put("three", 3)

	fmt.Println("TreeMap size:", treeMap.Size())
	if val, ok := treeMap.Get("two"); ok {
		fmt.Println("Value for 'two':", val)
	}
	
	// 创建LinkedHashMap（类似Java的HashMap，链表长度超过阈值时转换为红黑树）
	lhMap := maps.NewLinkedHashMap[string, int]()
	lhMap.Put("one", 1)
	lhMap.Put("two", 2)
	lhMap.Put("three", 3)

	fmt.Println("LinkedHashMap size:", lhMap.Size())
	if val, ok := lhMap.Get("one"); ok {
		fmt.Println("Value for 'one':", val)
	}
	
	// 创建一个小容量的LinkedHashMap，以便更容易触发冲突和树化
	collisionMap := maps.NewLinkedHashMapWithCapacity[int, string](4)
	// 添加足够多的元素以触发冲突和树化
	for i := 0; i < 20; i++ {
		collisionMap.Put(i, fmt.Sprintf("value%d", i))
	}
}
```

### Supported Collection Types

- **List**
  - ArrayList: Dynamic array-based List implementation
  - LinkedList: Doubly linked list-based List implementation

- **Set**
  - HashSet: Hash table-based Set implementation
  - TreeSet: Red-black tree-based ordered Set implementation
  - ConcurrentSkipListSet: Skip list-based thread-safe ordered Set implementation

- **Map**
  - HashMap: Hash table-based Map implementation
  - TreeMap: Red-black tree-based ordered Map implementation
  - RedisHashMap: Redis hash table design-based Map implementation with progressive resizing and chaining for collision resolution
  - LinkedHashMap: Chaining and red-black tree-based Map implementation, similar to Java's HashMap, converts chains to trees when threshold is exceeded

- **Queue**
  - LinkedQueue: Linked list-based Queue implementation
  - PriorityQueue: Priority queue implementation

- **Stack**
  - ArrayStack: Array-based Stack implementation
  - LinkedStack: Linked list-based Stack implementation

- **Others**
  - Pair: Key-value pair
  - MultiMap: Map with multiple values per key
  - BiMap: Bidirectional Map

### Architecture Design

The project follows a well-structured architecture with clear separation of concerns:

#### Core Interfaces
- **Container**: Basic container interface with `Size()`, `IsEmpty()`, `Clear()`, `Contains()`, `String()` methods
- **Iterable**: Provides iteration capability with `ForEach()` method
- **Iterator**: Standard iterator pattern with `HasNext()`, `Next()`, `Remove()` methods
- **Comparable**: Generic comparison interface for custom types

#### Module Structure
- `container/common`: Common interfaces and utilities
- `container/list`: List implementations (ArrayList, LinkedList)
- `container/set`: Set implementations (HashSet, TreeSet, ConcurrentSkipListSet)
- `container/map`: Map implementations (HashMap, TreeMap, LinkedHashMap)
- `container/queue`: Queue implementations
- `container/stack`: Stack implementations
- `examples`: Usage examples and demonstrations

### Implementation Highlights

#### Generic Support
- Full generic type support based on Go 1.18+
- Type-safe operations without runtime type assertions
- Clean and intuitive API design

#### Advanced Data Structures
- **ConcurrentSkipListSet**: Lock-free concurrent skip list with O(log n) operations
- **TreeMap/TreeSet**: Red-black tree implementation ensuring O(log n) complexity
- **LinkedHashMap**: Hybrid approach using both chaining and tree conversion for optimal performance

#### High-Performance Features
- Efficient memory management
- Optimized algorithms for common operations
- Benchmark tests ensuring performance standards

### Code Quality Analysis

#### Strengths
1. **Comprehensive Test Coverage**: 9 test files with 3,915 lines of test code
2. **Clean Architecture**: Well-defined interfaces and modular design
3. **Performance Optimization**: Includes benchmark tests for critical operations
4. **Thread Safety**: Provides concurrent implementations where needed
5. **Documentation**: Extensive examples and clear API documentation

#### Test Coverage Details
- `container/stack`: 100.0% coverage
- `container/queue`: High coverage for core operations
- `container/set`: 54.8% coverage (room for improvement)
- `container/map`: Good coverage for most implementations
- `container/list`: Comprehensive test suite

### 安装

```bash
go get github.com/chenjianyu/collections
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/chenjianyu/collections/container/list"
	"github.com/chenjianyu/collections/container/set"
	"github.com/chenjianyu/collections/container/map"
)

func main() {
	// 创建ArrayList
	list := list.New[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	fmt.Println("List size:", list.Size())
	fmt.Println("List contains 2:", list.Contains(2))

	// 创建HashSet
	set := set.New[string]()
	set.Add("apple")
	set.Add("banana")
	set.Add("apple") // 重复元素不会被添加

	fmt.Println("Set size:", set.Size())
	fmt.Println("Set elements:")
	set.ForEach(func(item string) {
		fmt.Println(item)
	})

	// 创建HashMap
	hashMap := maps.NewLinkedHashMap[string, int]()
	hashMap.Put("one", 1)
	hashMap.Put("two", 2)

	fmt.Println("Map size:", hashMap.Size())
	if val, ok := hashMap.Get("one"); ok {
		fmt.Println("Value for 'one':", val)
	}

	// 创建ConcurrentSkipListSet（线程安全的有序集合）
	skipSet := set.NewConcurrentSkipListSet[int]()
	skipSet.Add(5)
	skipSet.Add(2)
	skipSet.Add(8)
	skipSet.Add(1)

	fmt.Println("SkipSet (sorted):", skipSet) // 输出: [1, 2, 5, 8]
	if first, ok := skipSet.First(); ok {
		fmt.Println("First element:", first) // 输出: 1
	}
	if last, ok := skipSet.Last(); ok {
		fmt.Println("Last element:", last) // 输出: 8
	}

	// 创建TreeMap（有序映射）
	treeMap := maps.NewTreeMap[string, int]()
	treeMap.Put("one", 1)
	treeMap.Put("two", 2)
	treeMap.Put("three", 3)

	fmt.Println("TreeMap size:", treeMap.Size())
	if val, ok := treeMap.Get("two"); ok {
		fmt.Println("Value for 'two':", val)
	}
	
	// 创建LinkedHashMap（类似Java的HashMap，链表长度超过阈值时转换为红黑树）
	lhMap := maps.NewLinkedHashMap[string, int]()
	lhMap.Put("one", 1)
	lhMap.Put("two", 2)
	lhMap.Put("three", 3)

	fmt.Println("LinkedHashMap size:", lhMap.Size())
	if val, ok := lhMap.Get("one"); ok {
		fmt.Println("Value for 'one':", val)
	}
	
	// 创建一个小容量的LinkedHashMap，以便更容易触发冲突和树化
	collisionMap := maps.NewLinkedHashMapWithCapacity[int, string](4)
	// 添加足够多的元素以触发冲突和树化
	for i := 0; i < 20; i++ {
		collisionMap.Put(i, fmt.Sprintf("value%d", i))
	}
}
```

### 支持的集合类型

- **List**
  - ArrayList: 基于动态数组的List实现
  - LinkedList: 基于双向链表的List实现

- **Set**
  - HashSet: 基于哈希表的Set实现
  - TreeSet: 基于红黑树的有序Set实现
  - ConcurrentSkipListSet: 基于跳表的线程安全有序Set实现

- **Map**
  - HashMap: 基于哈希表的Map实现
  - TreeMap: 基于红黑树的有序Map实现
  - RedisHashMap: 基于Redis哈希表设计的Map实现，支持渐进式扩容和链表解决冲突
  - LinkedHashMap: 基于链地址法和红黑树的Map实现，类似Java的HashMap，当链表长度超过阈值时转换为红黑树

- **Queue**
  - LinkedQueue: 基于链表的Queue实现
  - PriorityQueue: 优先队列实现

- **Stack**
  - ArrayStack: 基于数组的Stack实现
  - LinkedStack: 基于链表的Stack实现

- **其他**
  - Pair: 键值对
  - MultiMap: 一个键对应多个值的Map
  - BiMap: 双向Map

### 架构设计

项目采用了良好的架构设计，具有清晰的关注点分离：

#### 核心接口
- **Container**: 基础容器接口，提供 `Size()`、`IsEmpty()`、`Clear()`、`Contains()`、`String()` 方法
- **Iterable**: 提供迭代能力，包含 `ForEach()` 方法
- **Iterator**: 标准迭代器模式，包含 `HasNext()`、`Next()`、`Remove()` 方法
- **Comparable**: 自定义类型的通用比较接口

#### 模块结构
- `container/common`: 通用接口和工具
- `container/list`: List实现（ArrayList、LinkedList）
- `container/set`: Set实现（HashSet、TreeSet、ConcurrentSkipListSet）
- `container/map`: Map实现（HashMap、TreeMap、LinkedHashMap）
- `container/queue`: Queue实现
- `container/stack`: Stack实现
- `examples`: 使用示例和演示

### 实现亮点

#### 泛型支持
- 基于Go 1.18+的完整泛型类型支持
- 类型安全操作，无需运行时类型断言
- 简洁直观的API设计

#### 高级数据结构
- **ConcurrentSkipListSet**: 无锁并发跳表，O(log n)操作复杂度
- **TreeMap/TreeSet**: 红黑树实现，确保O(log n)复杂度
- **LinkedHashMap**: 混合方法，使用链表和树转换以获得最佳性能

#### 高性能特性
- 高效的内存管理
- 常见操作的优化算法
- 基准测试确保性能标准

### 代码质量分析

#### 优势
1. **全面的测试覆盖**: 9个测试文件，3,915行测试代码
2. **清晰的架构**: 定义良好的接口和模块化设计
3. **性能优化**: 包含关键操作的基准测试
4. **线程安全**: 在需要的地方提供并发实现
5. **文档完善**: 丰富的示例和清晰的API文档

#### 测试覆盖率详情
- `container/stack`: 100.0% 覆盖率
- `container/queue`: 核心操作高覆盖率
- `container/set`: 54.8% 覆盖率（有改进空间）
- `container/map`: 大部分实现具有良好覆盖率
- `container/list`: 全面的测试套件

### 改进建议

#### 短期改进
1. **提高测试覆盖率**: 特别是 `container/set` 模块
2. **增强错误处理**: 添加更多边界条件检查
3. **性能优化**: 针对特定场景进行微调

#### 长期规划
1. **扩展集合类型**: 添加更多专用数据结构
2. **并发优化**: 提供更多线程安全的实现
3. **内存优化**: 实现更高效的内存管理策略

### Contributing | 贡献

Welcome to submit issues and Pull Requests! | 欢迎提交问题和Pull Request！

### License | 许可证

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. | 本项目采用MIT许可证 - 详见 [LICENSE](LICENSE) 文件。