package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 数据库名字
const dbFile = "blockchain.db"

// 仓库
const blocksBucket = "blocks"

// 区块链结构体
type Blockchain struct {
	Tip []byte   // 区块链里面最新一个区块的Hash
	DB  *bolt.DB // 数据库
}

// 新增区块
//func (bc *Blockchain) AddBlock(data string) {
//	// 1.获取上一个区块
//	oldBlock := bc.Blocks[len(bc.Blocks)-1]
//	// 2. 获取最新区块
//	newBlock := NewBlock(data, oldBlock.Hash)
//	// 3. 将新区块添加到区块链中
//	bc.Blocks = append(bc.Blocks, newBlock)
//}

// 创建一个带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *Blockchain {
	var tip []byte // 获取最新一个区块的Hash
	// -----------数据库创建----------
	//  如果数据库存在，打开，如果不存在，创建一个数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// 查看表是否存在
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			// 创建创世区块
			genesis := NewGenesisBlock()
			// 创建表
			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			// 将创世区块序列化后存储到表中
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新区块的Hash
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesis.Hash

		} else { // 表存在
			// key: l
			// value : 最新一个区块的Hash
			tip = b.Get([]byte("l"))
		}
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	return &Blockchain{tip, db}
}
