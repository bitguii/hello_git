package main

import (
	"blockchain/bolt"
	"log"
)

type BlockChainIterator struct {
	db *bolt.DB
	//游标，用于不断索引
	currentHashPointer []byte
}

//func NewIterator(bc *BlockChain)  {
//}

func (bc *BlockChain) NewIterator() *BlockChainIterator {

	return &BlockChainIterator{
		bc.db,
		//最初指向区块的最后一个，随着next()调用不断变化
		bc.tail,
	}
}

//迭代器属于区块链，next()方法属于迭代器
//1.返回当前区块
//2.指针前移
func (it *BlockChainIterator) Next() *Block {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时bucket不应为空，请检查！")
		}
		blockTmp := bucket.Get(it.currentHashPointer)

		//解码
		block = Deserialze(blockTmp)
		//游标哈希左移
		it.currentHashPointer = block.PrevHash
		return nil
	})
	return &block
}
