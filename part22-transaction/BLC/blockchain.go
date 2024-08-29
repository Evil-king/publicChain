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

// 创世区块里面的数据信息
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// 区块链结构体
type Blockchain struct {
	Tip []byte   // 区块链里面最新一个区块的Hash
	DB  *bolt.DB // 数据库
}

// 查询还未花费的交易集合
func (bc *Blockchain) FindUnspentTransactions(address string) []*Transaction {
	//1、遍历所有的区块链
	blockchainIterator := bc.Iterator()

	var hashInt big.Int

	for {
		if err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))

			blockByte := b.Get(blockchainIterator.CurrentHash)

			block := DeserializeBlock(blockByte)

			fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
			fmt.Printf("Hash：%x \n", block.Hash)
			fmt.Printf("Nonce：%d \n", block.Nonce)

			for _, transction := range block.Transactions {
				fmt.Printf("TransactionHash:%x\n", transction.ID)
			}

			fmt.Println()

			return nil
		}); err != nil {
			log.Panic(err)
		}
		// 获取下一个迭代器
		blockchainIterator = blockchainIterator.Next()
		// 将迭代器中的hash转为hashInt类型
		hashInt.SetBytes(blockchainIterator.CurrentHash)
		// 如果是创世区块 则退出循环
		if hashInt.Cmp(big.NewInt(0)) == 0 { // 如果等于0，则退出循环
			break
		}
	}
	//2、根据地址找到属于传入地址的未花销的交易
	return nil
}

// 创建一个带有创世区块的区块链
func NewBlockChain() *Blockchain {
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
			// 创建创世区块交易对象
			coinbaseTx := NewCoinbaseTX("hwq", genesisCoinbaseData)
			// 创建创世区块
			genesis := NewGenesisBlock(coinbaseTx)
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
