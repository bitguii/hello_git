package main

import (
	"blockchain/bolt"
	"log"
)

// 区块链结构
type BlockChain struct {
	//区块链数组
	//blocks []*Block
	db *bolt.DB

	//存储最后一个区块哈希
	tail []byte
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

// 创建区块链
func NewBlockChain() *BlockChain {
	//return &BlockChain{
	//	blocks: []*Block{genesisBlock},
	//}

	//最后一个区块哈希
	var lastHash []byte
	//1.打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	//defer db.Close()

	//操作数据库(改写)
	db.Update(func(tx *bolt.Tx) error {
		//2.找到bucket(如果没有就创建)
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉，创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建blockBucket失败")
			}
			//创建一个创世区块并添加至区块链中
			genesisBlock := GenesisBlock()

			//写数据
			//hash作为key，block的字节流作为value，尚未实现
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("lastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte("lastHashKey"))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

// 创世区块
func GenesisBlock() *Block {
	return NewBlock("blockchain is the future", []byte{})
}

// 添加区块
func (bc *BlockChain) AddBlock(data string) {
	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		//完成数据添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket不应为空，请检查！")
		}
		//创建新的区块
		block := NewBlock(data, lastHash)

		//添加到区块链db中
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("lastHashKey"), block.Hash)

		//更新内存中的区块链，指把最后的小尾巴tail更新一下
		bc.tail = block.Hash
		return nil
	})

}
