package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 迭代器结构体
type BlockchainIterator struct {
	CurrentHash []byte   // 当前区块的Hash
	DB          *bolt.DB // 数据库
}

// 迭代器入口方法
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip, bc.DB}
}

// 迭代器执行方法
func (bi *BlockchainIterator) Next() *BlockchainIterator {
	var nextHash []byte
	if err := bi.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// 获取当前区块并且发序列化
		currentBlock := DeserializeBlock(b.Get(bi.CurrentHash))
		// 获取当前区块的下一个区块的Hash
		nextHash = currentBlock.PrevBlockHash
		return nil
	}); err != nil {
		log.Panic(err)
	}
	return &BlockchainIterator{nextHash, bi.DB}
}
