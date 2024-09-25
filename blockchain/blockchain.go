package main

import (
	"blockchain/bolt"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
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
func NewBlockChain(address string) *BlockChain {
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
			genesisBlock := GenesisBlock(address)

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
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "blockchain is future")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	for _,tx := range txs{
		if !bc.VerifyTransaction(tx) {
			fmt.Println("矿工发现无效交易！")
			return
		}
	}

	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		//完成数据添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket不应为空，请检查！")
		}
		//创建新的区块
		block := NewBlock(txs, lastHash)

		//添加到区块链db中
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("lastHashKey"), block.Hash)

		//更新内存中的区块链，指把最后的小尾巴tail更新一下
		bc.tail = block.Hash
		return nil
	})

}

func (bc *BlockChain) Printchain() {
	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		//从第一个key,value开始遍历，到最后一个固定key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}
			block := Deserialze(v)
			fmt.Printf("=========================== 区块高度:%d =======================\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号：%d\n", block.Version)
			fmt.Printf("前区块哈希：%x\n", block.PrevHash)
			fmt.Printf("默克尔根：%x\n", block.MerkelRoot)
			fmt.Printf("时间戳：%d\n", block.TimeStamp)
			fmt.Printf("难度值（随便写的）：%d\n", block.Difficulty)
			fmt.Printf("随机数：%d\n", block.Nonce)
			fmt.Printf("当前区块哈希：%x\n", block.Hash)
			fmt.Printf("区块数据：%s\n", block.Transactions[0].TXINPUT[0].PubKey)
			return nil
		})
		return nil
	})
}

// 找到指定地址的所有UTXO
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)

	for _, tx := range txs {
		for _, output := range tx.TXOUTPUT {
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

// 找到地址转账所需的UTXO
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxo集合
	utxos := make(map[string][]uint64)
	//找到的utxos里面包含的钱总数
	var calc float64

	txs := bc.FindUTXOTransactions(senderPubKeyHash)

	for _, tx := range txs {
		for i, output := range tx.TXOUTPUT {
			//直接比较是否相同返回true或false
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				if calc < amount {
					//1.把utxo加进来
					//array := utxos[string(tx.TXID)]
					//array = append(array,uint64(i))
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))

					//2.统计当前utxo总额
					calc += output.Value

					//3.比较是否满足转账需求
					//加完之后如果满足条件
					if calc >= amount {
						fmt.Printf("找到了满足的金额:%f\n", calc)
						return utxos, calc
					}
				} else {
					fmt.Printf("不满足转账金额，当前总额：%f,目标金额：%%f\n", calc, amount)
				}
			}
		}
	}
	return utxos, calc
}

// 找到所有相关UTXO
func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含utxo交易集合
	//定义一个map来保存消费过的output，key是这个output的交易id，value是这个交易中索引的数组
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	it := bc.NewIterator()
	//遍历区块
	for {
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {
			//遍历outpout，找到和地址相关的utxo(在添加output之前检查自己是否已经消耗过)
		OUTPUT:
			for i, output := range tx.TXOUTPUT {
				//做过滤，将所有消耗过的output和当前即将添加的output对比一下，如果相同则跳过
				//如果当前交易id存在于已经标识的map，那么说明这个交易中有消耗过的
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//当前准备添加的output已经消耗过，不需添加
							continue OUTPUT
						}
					}
				}
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					//UTXO = append(UTXO, output)
					//返回所有包含我的utxo集合
					txs = append(txs, tx)
				} else {

				}
			}
			//如果当前交易是挖矿交易，那么不做遍历，跳过
			if !tx.IsCoinbase() {
				//遍历input，找到自己花费过的utxo的集合(把自己消耗过的标识出来)
				for _, input := range tx.TXINPUT {
					//判断当前input和目标是否一致，一致说明这个是目标地址消耗过的output，就加入map
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			} else {
				//fmt.Println("this is coinbase")
			}
		}

		if len(block.PrevHash) == 0 {
			break
			fmt.Println("区块链遍历完成退出！")
		}
	}

	return txs
}

func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	//遍历区块链
	it := bc.NewIterator()
	for {
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {
			//比较交易，找到了直接退出
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历结束！")
			break
		}
	}
	//如果没找到，返回空的Transaction，同时返回错误状态
	return Transaction{}, errors.New("无效的交易id，请检查！")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//根据inputs来找，有多少input，就遍历多少次
	//找到目标交易（根据TXid来找）
	//添加到prevTXs
	for _, input := range tx.TXINPUT {
		//根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx
	}
	tx.Sign(privateKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//根据inputs来找，有多少input，就遍历多少次
	//找到目标交易（根据TXid来找）
	//添加到prevTXs
	for _, input := range tx.TXINPUT {
		//根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx
	}
	return tx.Verify(prevTXs)
}
