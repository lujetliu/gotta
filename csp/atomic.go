package main

/*
	面向CSP并发模型的channel原语和面向传统共享内存并发模型的sync包
	提供的原语已经足以满足Go语言应用并发设计中99.9%的并发同步需求了,
	而剩余那0.1%的需求, 可以使用Go标准库提供的atomic包来实现


	atomic包是Go语言提供的原子操作(atomic operation)原语的相关接口,
	原子操作是相对于普通指令操作而言的; 如下例:
		// var a int
		// a++
	a++ 这条语句需要以下3条普通机器指令来完成变量a的自增:
	- LOAD:将变量从内存加载到CPU寄存器
	- ADD:执行加法指令
	- STORE:将结果存储回原内存地址

	这3条普通指令在执行过程中是可中断的, 而原子操作的指令是不可中断的;
	好比一个事务, 要么不执行, 一旦执行就一次性全部执行完毕, 不可分割;
	正因如此, 原子操作可用于共享数据的并发同步;

	原子操作由底层硬件直接提供支持, 是一种硬件实现的指令级"事务",
	因此相比操作系统层面和Go运行时层面提供的同步技术而言, 它更为原始;
	atomic包封装了CPU实现的部分原子操作指令(TODO:封装实现), 为用户层
	提供体验良好的原子操作函数, 因此atomic包中提供的原语更接近硬件底层,
	也更为低级, 常被用于实现更为高级的并发同步技术(比如channel和sync
	包中的同步原语);

						    /|\
						     | 高级
				channel      |
						     |
						     |
				sync 包      |
						     |
						     |
				atomic 包    |
						     | 低级



	原子操作可以运用到以下场景中:
	- 对共享整型变量的无锁读写(针对整型变量, 包括有符号整型,
		无符号整型以及对应的指针类型) (./atomic1_test.go
	- 对共享自定义类型变量的无锁读写(针对自定义类型) (./atomic2_test.go)

	随着并发量提升, 使用atomic实现的共享变量的并发读写性能表现更为稳定,
	尤其是原子读操作, 这让atomic与sync包中的原语比起来表现出更好的伸缩
	性和更高的性能l; 所以atomic包更适合一些对性能十分敏感, 并发量较大
	且读多写少的场合;

	但atomic原子操作可用来同步的范围有较大限制, 仅支持整型变量或自定义
	类型变量, 如果要对一个复杂的临界区数据进行同步, 首选依旧是sync包中的原语;
*/
