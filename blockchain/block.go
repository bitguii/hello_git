package main

import (
	"bytes"
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
		3,
		0,
		[]byte{},
		[]byte(data),
	}

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

//序列化
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

//反序列化
func Deserialze(data []byte) Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("反序列化失败！")
	}
	return block
}

// 生成哈希
/*func (block *Block) SetHash() {
	var blockInfo []byte
	//拼装数据
	//blockInfo = append(blockInfo, uint64ToByte(block.Version)...)
	//blockInfo = append(blockInfo, block.PrevHash...)
	//blockInfo = append(blockInfo, block.MerkelRoot...)
	//blockInfo = append(blockInfo, uint64ToByte(block.TimeStamp)...)
	//blockInfo = append(blockInfo, uint64ToByte(block.Difficulty)...)
	//blockInfo = append(blockInfo, uint64ToByte(block.Nonce)...)
	//blockInfo = append(blockInfo, block.Data...)
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
}*/
