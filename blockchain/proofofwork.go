package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 定义工作量证明结构
type ProofOfWork struct {
	//a.block
	block *Block
	//b.目标值(需要非常大数)
	target *big.Int
}

// 创建pow函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//制定难度值
	targetStr := "000100000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}
	//将难度值转化为big.int
	tmpInt.SetString(targetStr, 16)

	pow.target = &tmpInt
	return &pow
}

// 提供不断计算hash的函数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	block := pow.block
	var hash [32]byte

	fmt.Println("开始挖矿...")
	for {
		//拼数据
		tmp := [][]byte{
			uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			uint64ToByte(block.TimeStamp),
			uint64ToByte(block.Difficulty),
			uint64ToByte(nonce),
			//只对区块头做哈希，区块体通过默克尔根产生影响
			//block.Data,
		}
		blockInfo := bytes.Join(tmp, []byte{})

		//哈希运算
		hash = sha256.Sum256(blockInfo)

		//与pow中的target进行比较
		tmpInt := big.Int{}
		//将计算的哈希转化为bigInt
		tmpInt.SetBytes(hash[:])

		//bigInt比较
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash:%x,nonce:%d\n", hash, nonce)
			//break
			return hash[:], nonce
		} else {
			nonce++
		}
	}

}
