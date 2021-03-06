package main

/*
	在业务的高并发场景中, 会在短时间内产生大量消息; 在插入数据库之前,
	需要给这些消息、订单先打上一个ID, 然后再插入数据库; 对这个ID的要
	求是希望其中能带有一些时间信息, 这样即使后端的系统对消息进行了分库分表,
	也能够以时间顺序对这些消息进行排序;

	Twitter的snowflake算法是这种场景下的一个典型解法;(TODO)



										                    datacenter_id     sequence_id
                                                                  |                 |
                                                                  |                 |
		| <----------------41位----------------------------> |   \|/               \|/
	-------------------------------------------------------------------------------------------
	| 0 | 00000 00000 00000 00000 00000 00000 00000 00000 0  | 00000 | 00000 | 00000 00000 00 |
	-------------------------------------------------------------------------------------------
	  /|\                          /|\                                  /|\
       |                            |                                    |
       |                            |                                    |
       |                            |                                    |
	 未使用                   毫秒级时间戳                          worker_id


	- 数值为 int64类型, 被划分为4部分, 不含开头的第一位,因为这个位是符号位
	- 随后用41位来表示收到请求时的时间戳, 单位为毫秒
	- 用5位来表示数据中心的ID
	- 用5位来表示机器的实例ID
	- 12位的循环自增ID(到 达1111 1111 1111后会归零)

	这样的机制可以支持在同一台机器上在1毫秒内产生4096条消息, 1秒内共产生
	409.6万条消息; 从值域上来讲完全够用了; 数据中心加上实例ID共有10位,
	可以支持每个数据中心部署32台机器, 所有数据中心共1024台实例; 表示时间
	戳的41位, 可以支持使用69年; 这里的时间戳实际上只是相对于某个时间的增量;
	例如, 如果我们的系统上线时间是2018-08-01, 就可以把这个时间戳当作
	是从2018-08-01 00:00:00.000的偏移量;



*/

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake" // TODO: 熟悉应用
)

func main() {
	n, err := snowflake.NewNode(1)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	for i := 0; i < 3; i++ {
		id := n.Generate()
		fmt.Println("id", id)
		fmt.Println(
			"node: ", id.Node(),
			"step: ", id.Step(),
			"time: ", id.Time(),
		)
	}
}
