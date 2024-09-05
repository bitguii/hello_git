package main

func main() {
	bc := NewBlockChain()
	cli := CLI{bc}
	cli.Run()

	//bc.AddBlock("1111111111111")
	//bc.AddBlock("2222222222222")
	//
	////创建迭代器
	//it := bc.NewIterator()
	//
	////调用迭代器，返回每一个区块数据
	//for {
	//	//返回区块，然后左移
	//	block := it.Next()
	//
	//	fmt.Printf("=============================\n")
	//	fmt.Printf("前区块哈希：%x\n", block.PrevHash)
	//	fmt.Printf("当前区块哈希：%x\n", block.Hash)
	//	fmt.Printf("区块数据：%s\n", block.Data)
	//
	//	if len(block.PrevHash) == 0 {
	//		fmt.Println("区块链遍历结束！")
	//		break
	//	}
	//}

}
