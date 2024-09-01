package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"time"
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

// 迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip, bc.DB}
}

// 新增区块
func (bc *Blockchain) AddBlockToBlockchain(data string) {
	// 存储到数据库中
	if err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			// 从表中获取当前最新的区块hash,并且序列化为block对象
			block := DeserializeBlock(b.Get(bc.Tip))
			// 新增一个区块
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			// 存储新区块的数据
			if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
				log.Panic(err)
			}
			// 存储最新区块的Hash
			if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
				log.Panic(err)
			}
			// 将最新的区块的Hash存储到blockchain的Tip中
			bc.Tip = newBlock.Hash
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

// 创建一个带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {
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
			genesis := CreateGenesisBlock("Genesis Block...")
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

// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {

	fmt.Println("PrintchainPrintchainPrintchainPrintchain")
	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("PrevBlockHash：%x\n", block.PrevBlockHash)
		fmt.Printf("Data：%s\n", block.Data)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)

		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}
