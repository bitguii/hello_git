package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

// 定义区块结构
type Block struct {
	//1.版本号
	Version uint64
	//2.前区块hash
	PrevHash []byte
	//3.Merkel根
	MerkelRoot []byte
	//4.时间戳
	TimeStamp uint64
	//5.难度值
	Difficulty uint64
	//6.随机数
	Nonce uint64

	//a.当前hash(正常比特币区块中没有当前哈希，此处简化了)
	Hash []byte
	//b.数据
	Data []byte
}

// 辅助函数，将uint64转换为[]byte
func uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		00,
		prevBlockHash,
		[]byte{},
		uint64(time.Now().Unix()),
		0,
		0,
		[]byte{},
		[]byte(data),
	}

	block.SetHash()

	return &block
}

// 生成哈希
func (block *Block) SetHash() {
	var blockInfo []byte
	//拼装数据
	/*blockInfo = append(blockInfo, uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, block.PrevHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, uint64ToByte(block.Nonce)...)
	blockInfo = append(blockInfo, block.Data...)*/
	tmp := [][]byte{
		uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		uint64ToByte(block.TimeStamp),
		uint64ToByte(block.Difficulty),
		uint64ToByte(block.Nonce),
		block.Data,
	}

	//将二维数组切片连接起来，返回一个一维切片
	blockInfo = bytes.Join(tmp, []byte{})

	//sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
