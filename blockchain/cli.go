package main

import (
	"fmt"
	"os"
	"strconv"
)

// 此文件用来接收命令参数并控制区块链运作
type CLI struct {
	bc *BlockChain
}

const Usage = `
    printChain              		'print all blockchain data'
	printChainR              		'print all blockchain data Reverse'
	getBalance -address ADDRESS		'get the balance all UXTO'
	send FROM TO AMOUNT MINER DATA  "transfer amount from 'from' to 'to',mine by "miner",and write data at the same time"
`

// 接受参数的动作
func (cli *CLI) Run() {
	//得到命令
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}

	//分析命令
	cmd := args[1]
	switch cmd {
	case "printChain":
		fmt.Printf("正向打印区块链\n")
		cli.printBlockChain()
	case "printChainR":
		fmt.Printf("反向打印区块链\n")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Println("获取余额")
		if len(args) == 4 && args[2] == "-address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Println("转账开始...")
		if len(args) != 7 {
			fmt.Println("参数错误，请检查！")
			fmt.Println(Usage)
		}
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	default:
		fmt.Println("无效的命令，请检查")
		fmt.Println(Usage)
	}
}
