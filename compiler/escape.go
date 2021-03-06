package main

/*
* TODO: 逃逸分析, 堆, 栈
 * 逃逸分析是 go 中重要的优化阶段, 用于标识变量内存应该被分配在栈区还是堆区.
 * 在 c 或 c++ 语言中, 开发者经常会犯的错误是函数返回了一个栈上的对象指针, 在
 * 函数执行完成, 栈被销毁后, 继续访问被销毁栈上的指针对象, 导致出现问题;
 * go 语言能够通过编译时的逃逸分析识别这种问题, 自动将该变量放置到堆区, 并借
 * 助于 go 运行时的垃圾回收机制自动释放内存; 编译器会尽可能的将变量放置到栈
 * 中, 因为栈中的对象随着函数调用结束会自动销毁, 减轻运行时分配和垃圾回收的
 * 负担.
 *
 *
 * 在 go 中, 开发者模糊了栈区和堆区的差别, 不管是字符串, 数组字面量, 还是通过
 * new, make 表示符创建的对象, 都既可能被分配到栈中, 也可能被分配到堆中, 遵循:
 * 原则1: 指向栈上对象的指针不能被存储到堆中
 * 原则2: 指向栈上对象的指针不能超过该栈对象的生命周期
 *
*/

// TODO: 图
// go 通过对抽象语法树的静态数据流分析(static data-flow analysis) 实现逃逸
// 分析, 这种方式构建了带权重的有向图.
var z *int

func escape() {
	a := 1
	z = &a
}

// 变量 z 引用了变量 a 的地址, 如果变量 a 被分配到栈中, 那么程序将违背原则2, 即
// 变量 z 超过了变量 a 的生命周期, 因此变量 a 被分配到堆中; 可以通过在编译时
// 加入 -m=2(?) 标志打印出编译时的逃逸分析信息:
// go tool compile -m=2 memory.go
// memory.go:26:6: can inline escape with cost 9 as: func() { a := 1; z = &a }
// memory.go:27:2: a escapes to heap:
// memory.go:27:2:   flow: {heap} = &a:
// memory.go:27:2:     from &a (address-of) at memory.go:28:6
// memory.go:27:2:     from z = &a (assign) at memory.go:28:4
// memory.go:27:2: moved to heap(堆): a

/*
 * go 语言在编译时构建了带权重的有向图, 其中权重可以表明当前变量引用与解引用
 * 的数量; 下例中为 p 引用 q 时的权重, 当权重大于 0 时, 代表存在解引用(*)操作,
 * 当权重为 -1 时, 代表存在引用(&)操作
 * p = &q        // -1
 * p = q         // 0
 * p = *q        // 1
 * p = **q       // 2
 * p = **&**&q   // 2
 *
 */

//并不是权重为 -1 就一定要逃逸, 如下, 虽然 z 引用了变量 a 的地址, 但是由于变量
//z 并没有超过变量 a 的生命周期, 因此变量 a 和变量 z 都不需要逃逸.
func f() int {
	a := 1
	z := &a
	return *z
}

// 更复杂的例子:
var o *int

func main() {
	l := new(int)
	*l = 42
	m := &l
	n := &m
	o = **n
}

// TODO:
// 最终编译器在逃逸分析中的数据流分析, 会解析为带权重的有向图:
//             权重:2          权重:1          权重:0           权重:-1
//               /|\            /|\             /|\              /|\
//                |              |               |                |
// o  <----2---- n <---- -1 ---- m <---- -1 ---- l <---- -1 ---- new(int)
// 节点代表变量, 边代表变量之间的赋值, 箭头代表赋值的方向, 边上的数字代表当前
// 赋值的引用或解引用的个数. 节点的权重=前一个节点的权重+箭头上的数字

// 遍历和计算有向权重图的目的是找到权重为 -1 的节点, 它的节点变量地址会被传递
// 到根节点 o 中, 这时还需要考虑逃逸分析的分配原则, o 节点为全局变量, 不能被分
// 配在栈中, 因此, new(int) 节点创建的变量会被分配到堆中(TODO)
// go tool compile -m escape.go
// escape.go:66:10: new(int) escapes to heap

/*
 * TODO: escape.go 核心逻辑
 * 实际的情况更加复杂, 因为一个节点可能拥有多个多条边(如结构体), 而节点之间
 * 可能出现环, go 语言采用 Bellman Ford(TODO) 算法遍历有向图中权重小于 0 的
 * 节点, 核心逻辑位于 $GOROOT/src/cmd/compile/internal/gc/escape.go 中.
 */
