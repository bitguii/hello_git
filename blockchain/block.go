package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
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
	//Data []byte
	//b.真实交易数据
	Transactions []*Transaction
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
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		00,
		prevBlockHash,
		[]byte{},
		uint64(time.Now().Unix()),
		3,
		0,
		[]byte{},
		//[]byte(data),
		txs,
	}

	block.MerkelRoot = block.MakeMerkelRoot()

	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	//不停进行哈希运算
	hash, nonce := pow.Run()
	//根据挖矿结果对区块进行更新
	block.Hash = hash
	block.Nonce = nonce

	return &block
}

// 序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer

	//使用gob进行序列化得到字节流
	//1.定义一个编码器
	//2.使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic("序列化失败！")
	}

	return buffer.Bytes()
}

// 反序列化
func Deserialze(data []byte) Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("反序列化失败！")
	}
	return block
}

// 模拟默克尔根生成，只对交易数据做简单的拼接，不做二叉树处理
func (block *Block) MakeMerkelRoot() []byte {
	var info []byte
	//将交易的哈希值拼接起来，整体再次哈希
	for _, tx := range block.Transactions {
		info = append(info, tx.TXID...)
	}
	hash := sha256.Sum256(info)
	return hash[:]
}
