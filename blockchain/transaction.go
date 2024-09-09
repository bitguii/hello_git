package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 50

// 1.定义交易结构
type Transaction struct {
	TXID     []byte     //交易ID
	TXINPUT  []TXInput  //交易输入
	TXOUTPUT []TXOutput //交易输出
}

// 定义交易输入
type TXInput struct {
	TXid  []byte //交易ID
	Index int64  //索引
	Sig   string //解锁脚本，先用地址来模拟
}

// 定义交易输出
type TXOutput struct {
	Value      float64 //转账金额
	PutKeyHash string  //锁定脚本，用地址模拟
}

// 设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 实现一个函数判断当前函数是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXINPUT) == 1 {
		//1.交易input只有一个
		//2.交易input的id为空
		//3.交易的index为-1
		if len(tx.TXINPUT) == 1 && len(tx.TXINPUT[0].TXid) == 0 && tx.TXINPUT[0].Index == -1 {
			return true
		}
	}
	return false
}

// 2.提供创建交易方法(挖矿交易)
func NewCoinbaseTX(address string, data string) *Transaction {
	//挖矿交易只有一个input，无需引用交易id，无需引用index
	//矿工由于挖矿时无需指定签名，所以sig字段可以由矿工自由填写，一般填写矿池的名字
	input := TXInput{
		[]byte{},
		-1,
		data,
	}

	output := TXOutput{
		reward,
		address,
	}

	tx := &Transaction{
		[]byte{},
		[]TXInput{input},
		[]TXOutput{output},
	}
	tx.SetHash()

	return tx
}

// 创建普通转账交易

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// a.找到最合理的utxo集合，map[string][]uint64
	utxos, resValue := bc.FindNeedUTXOs(from, amount)

	if resValue < amount {
		fmt.Println("余额不足，交易失败！")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	// b.创建交易输入，将这些utxo逐一转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{
				[]byte(id),
				int64(i),
				from,
			}
			inputs = append(inputs, input)
		}
	}

	// c.创建outputs
	output := TXOutput{
		amount,
		to,
	}
	outputs = append(outputs, output)

	// d.如果有零钱，要找零
	if resValue > amount {
		//找零
		outputs = append(outputs, TXOutput{
			resValue - amount,
			from,
		})
	}
	tx := Transaction{
		[]byte{},
		inputs,
		outputs,
	}
	tx.SetHash()
	return &tx
}
