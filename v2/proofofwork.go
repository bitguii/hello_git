package main

import "math/big"

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

	targetStr := "000010000000000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr, 16)

	pow.target = &tmpInt
	return &pow
}

// 提供不断计算hash的函数
func (pow *ProofOfWork) Run() ([]byte, uint) {
	//todo
	return []byte{}, 0
}
