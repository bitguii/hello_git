package main

import (
	"fmt"
	"os"
)

// 此文件用来接收命令参数并控制区块链运作
type CLI struct {
	bc *BlockChain
}

const Usage = `
    addBlock --data DATA    'add data to blockchain'
    printChain              'print all blockchain data'
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
	case "addBlock":
		//添加区块
		if len(args) == 4 && args[2] == "--data" {
			//获取命令行数据
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Println("添加区块参数使用不当，请检查")
			fmt.Println(Usage)
		}
	case "printChain":
		//打印区块链
		cli.printBlockChain()
	default:
		fmt.Println("无效的命令，请检查")
		fmt.Println(Usage)
	}
}
