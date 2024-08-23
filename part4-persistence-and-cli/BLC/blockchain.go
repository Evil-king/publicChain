package BLC

type Blockchain struct {
	Blocks []*Block // 存储有序的区块
}

// 新增区块
func (bc *Blockchain) AddBlock(data string) {
	// 1.获取上一个区块
	oldBlock := bc.Blocks[len(bc.Blocks)-1]
	// 2. 获取最新区块
	newBlock := NewBlock(data, oldBlock.Hash)
	// 3. 将新区块添加到区块链中
	bc.Blocks = append(bc.Blocks, newBlock)
}

// 创建一个带有创世区块的区块链
func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
