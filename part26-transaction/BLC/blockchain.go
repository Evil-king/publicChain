package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
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
	// 存储未花费输出的交易
	var unspentTXs []*Transaction
	// 存储交易输入对应的输出关系
	spentTXOs := make(map[string][]int)
	//1、遍历所有的区块链
	blockchainIterator := bc.Iterator()
	var hashInt big.Int

	for {
		if err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))

			blockByte := b.Get(blockchainIterator.CurrentHash)

			block := DeserializeBlock(blockByte)

			for _, transaction := range block.Transactions {
				// 将byte array类型转string
				txID := hex.EncodeToString(transaction.ID)

				// 查询交易数据的交易输出
			Outputs:
				for index, out := range transaction.Vout {
					// 是否已经被花费 ？
					if spentTXOs[txID] != nil {
						for _, value := range spentTXOs[txID] {
							if index == value {
								continue Outputs
							}
						}
					}
					if out.CanBeUnlockedWith(address) {
						unspentTXs = append(unspentTXs, transaction)
					}
				}

				// 判断是不是创世区块
				if transaction.IsCoinbase() == false {
					// 遍历交易输入
					for _, in := range transaction.Vin {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
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
	return unspentTXs
}

// 查找可用的未消费的输出信息
// 返回两个值，一个是未消费的金额，一个是对应输入中Vout的下标
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	// 查看未花费
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0
	// 遍历未花费的交易
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for index, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], index)
			}
			if accumulated >= amount {
				continue Work
			}
		}
	}
	return accumulated, unspentOutputs
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

// 根据交易的数组，打包新的区块
func (bc *Blockchain) MineBlock(txs []*Transaction) {

	if err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 新建block
		block := NewBlock(txs, bc.Tip)
		// 读表
		b := tx.Bucket([]byte(blocksBucket))
		// 将区块存到链上
		if err := b.Put(block.Hash, block.Serialize()); err != nil {
			log.Panic(err)
		}
		// 更新l对应的区块hash
		if err := b.Put([]byte("l"), block.Hash); err != nil {
			log.Panic(err)
		}
		// 更新本地对应的区块hash
		bc.Tip = block.Hash

		return nil
	}); err != nil {
		log.Panic(err)
	}
}
