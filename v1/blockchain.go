package main

// 区块链结构
type BlockChain struct {
	//区块链数组
	blocks []*Block
}

// 创建区块链
func NewBlockChain() *BlockChain {
	//创建一个创世区块并添加至区块链中
	genesisBlock := GenesisBlock()
	return &BlockChain{
		[]*Block{genesisBlock},
	}
}

// 创世区块
func GenesisBlock() *Block {
	return NewBlock("blockchain is the future", []byte{})
}

// 添加区块
func (bc *BlockChain) AddBlock(data string) {
	//1.创建新区块
	block := NewBlock(data, bc.blocks[len(bc.blocks)-1].Hash)
	//2.添加至区块链中
	bc.blocks = append(bc.blocks, block)
}
