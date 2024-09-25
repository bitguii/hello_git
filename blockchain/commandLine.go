package main

import (
	"fmt"
	"time"
)

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
		fmt.Printf("区块数据：%s\n", block.Transactions[0].TXINPUT[0].PubKey)

		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}

// 反向打印
func (cli *CLI) PrintBlockChainReverse() {
	bc := cli.bc

	//创建迭代器
	it := bc.NewIterator()

	//调用迭代器，返回每一个区块数据
	for {
		//返回区块，然后左移
		block := it.Next()
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")

		fmt.Printf("=============================\n")
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前区块哈希：%x\n", block.PrevHash)
		fmt.Printf("默克尔根：%x\n", block.MerkelRoot)
		fmt.Printf("时间戳：%s\n", timeFormat)
		fmt.Printf("难度值（随便写的）：%d\n", block.Difficulty)
		fmt.Printf("随机数：%d\n", block.Nonce)
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("区块数据：%s\n", block.Transactions[0].TXINPUT[0].PubKey)

		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}

// 获取指定地址钱包余额
func (cli *CLI) GetBalance(address string) {
	//校验地址
	if !IsValidAddress(address) {
		fmt.Printf("地址无效:%s\n", address)
		return
	}

	//生成公钥哈希
	pubKeyHash := GetPubKeyFromAddress(address)

	utxos := cli.bc.FindUTXOs(pubKeyHash)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("|%s|的余额为：%f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	if IsValidAddress(from) {
		fmt.Printf("地址无效 from：%s\n", from)
		return
	}
	if IsValidAddress(to) {
		fmt.Printf("地址无效 to：%s\n", to)
		return
	}
	if IsValidAddress(miner) {
		fmt.Printf("地址无效 miner：%s\n", miner)
		return
	}
	//1.创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	//2.创建普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		return
	}
	//3.将交易添加到区块中
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Println("转账成功！")
}

func (cli *CLI) NewWallet() {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址：%s\n", address)
	//for address := range ws.WalletsMap {
	//	fmt.Printf("地址：%s\n", address)
	//}
}

func (cli *CLI) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddress()
	for _, address := range addresses {
		fmt.Printf("地址：%s\n", address)
	}
}
