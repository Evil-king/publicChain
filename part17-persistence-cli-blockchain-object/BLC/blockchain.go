package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)

// 数据库名字
const dbName = "blockchain.db"

// 仓库
const blockTableName = "blocks"

// 区块链结构体
type Blockchain struct {
	Tip []byte   // 区块链里面最新一个区块的Hash
	DB  *bolt.DB // 数据库
}

// 判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 迭代器
func (blc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blc.Tip, blc.DB}
}

// 新增区块
func (blc *Blockchain) AddBlockToBlockchain(data string) {
	// 存储到数据库中
	if err := blc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 从表中获取当前最新的区块hash,并且序列化为block对象
			block := DeserializeBlock(b.Get(blc.Tip))
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
			blc.Tip = newBlock.Hash
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

// 创建一个带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(data string) {
	// 判断数据库是否存在
	if dbExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")
	// 创建一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// 创建表
		var b *bolt.Bucket
		b, err = tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			// 创建创世区块
			genesis := CreateGenesisBlock(data)
			// 将创世区块序列化后存储到表中
			if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
				log.Panic(err)
			}
			// 存储最新区块的Hash
			if err = b.Put([]byte("l"), genesis.Hash); err != nil {
				log.Panic(err)
			}
		}
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("创建创世区块完成.......")
}

// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {

	fmt.Println("开始打印输出...........")
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

// 返回Blockchain对象
func BlockchainObject() *Blockchain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var tip []byte

	if err = db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(blockTableName)); b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
	return &Blockchain{tip, db}
}
