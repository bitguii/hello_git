package main

import "fmt"

// 打印区块链
func (cli *CLI) printBlockChain() {
	//创建迭代器
	it := cli.bc.NewIterator()

	//调用迭代器，返回每一个区块数据
	for {
		//返回区块，然后左移
		block := it.Next()

		fmt.Printf("=============================\n")
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前区块哈希：%x\n", block.PrevHash)
		fmt.Printf("默克尔根：%x\n", block.MerkelRoot)
		fmt.Printf("时间戳：%d\n", block.TimeStamp)
		fmt.Printf("难度值（随便写的）：%d\n", block.Difficulty)
		fmt.Printf("随机数：%d\n", block.Nonce)
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("区块数据：%s\n", block.Data)

		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}

// 添加区块
func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
}
