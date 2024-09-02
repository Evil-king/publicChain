package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gogf/gf/v2/frame/g"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// 数据库名字
const dbName = "blockchain.db"

// 表明
const blockTableName = "blocks"

// 创世区块里面的数据信息
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// 区块链结构体
type Blockchain struct {
	Tip []byte   // 区块链里面最新一个区块的Hash
	DB  *bolt.DB // 数据库
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip, bc.DB}
}

// 查询还未花费的交易集合
func (bc *Blockchain) UnUTXOs(address string) []*UTXO {
	// 存储未花费输出的交易
	var unUTXOs []*UTXO
	var hashInt big.Int
	// 存储交易输入对应的输出关系
	spentTXOs := make(map[string][]int)
	// 获取下一个迭代器
	blockchainIterator := bc.Iterator()

	for {
		block := blockchainIterator.Next()

		for _, tx := range block.Txs {
			// 是不是创世区块产生的交易
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vin {
					// 能够解锁
					if in.UnLockWithAddress(address) {
						key := hex.EncodeToString(tx.TxHash)
						spentTXOs[key] = append(spentTXOs[key], in.Vout)
					}
				}
			}
			// Vouts==输出
		Work:
			for index, out := range tx.Vouts {
				// 交易输出是否能解锁
				if out.UnLockScriptPubKeyWithAddress(address) {
					if len(spentTXOs) != 0 {
						var isSpentUTXO bool
						for txHash, indexArray := range spentTXOs {
							for _, i := range indexArray {
								if index == i && txHash == hex.EncodeToString(tx.TxHash) {
									// 说明是已经消费了
									isSpentUTXO = true
									continue Work
								}
							}
						}
						if isSpentUTXO == false {
							unUTXOs = append(unUTXOs, &UTXO{tx.TxHash, index, out})
						}
					} else {
						unUTXOs = append(unUTXOs, &UTXO{tx.TxHash, index, out})
					}
				}
			}
		}

		// 将迭代器中的hash转为hashInt类型
		hashInt.SetBytes(block.Hash)
		// 如果是创世区块 则退出循环
		if hashInt.Cmp(big.NewInt(0)) == 0 { // 如果等于0，则退出循环
			break
		}
	}
	return unUTXOs
}

// 查找可用的未消费的输出信息
// 返回两个值，一个是未消费的金额，一个是对应输入中Vout的下标
func (bc *Blockchain) FindSpendableUTXOS(form string, amount int) (int64, map[string][]int) {
	spendableUTXO := make(map[string][]int)
	// 查看未花费
	utxos := bc.UnUTXOs(form)
	var value int64
	// 遍历未花费的交易
	for _, utxo := range utxos {
		value += utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}
	if value < int64(amount) {
		fmt.Printf("%s's fund is 不足\n", form)
		os.Exit(1)
	}
	return value, spendableUTXO
}

// 创建一个带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {
	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")
	//  如果数据库存在，打开，如果不存在，创建一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var tip []byte
	err = db.Update(func(tx *bolt.Tx) error {
		// 查看表是否存在
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			// 创建创世区块交易对象
			coinbaseTx := NewCoinbaseTransaction(address, genesisCoinbaseData)
			// 创建创世区块
			genesis := CreateGenesisBlock([]*Transaction{coinbaseTx})
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
		}
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	return &Blockchain{tip, db}
}

// 根据交易的数组，打包新的区块
func (bc *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	// 建立一笔交易
	value, _ := strconv.Atoi(amount[0])
	NewSimpleTransaction(from[0], to[0], value, bc)

	var txs []*Transaction
	var block *Block
	if err := bc.DB.View(func(tx *bolt.Tx) error {
		// 读表
		if b := tx.Bucket([]byte(blockTableName)); b != nil {
			blockByte := b.Get([]byte("l"))
			block = DeserializeBlock(blockByte)
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}

	// 建立新的区块
	block = NewBlock(txs, block.Height+1, block.PrevBlockHash)

	// 将新区块存储到数据库
	if err := bc.DB.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(blockTableName)); b != nil {
			//序列化存入数据库
			if err := b.Put(block.Hash, block.Serialize()); err != nil {
				log.Panic(err)
			}
			//更新l对应的区块hash
			if err := b.Put([]byte("l"), block.Hash); err != nil {
				log.Panic(err)
			}
			// 更新本地bc中的Tip
			bc.Tip = block.Hash
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

// 查询余额
func (bc *Blockchain) GetBalance(address string) int64 {
	// 查询未消费
	utxos := bc.UnUTXOs(address)
	var money int64
	for _, utxo := range utxos {
		money = money + utxo.Output.Value
	}
	return money
}

func (bc *Blockchain) printChain() {
	blockchainIterator := bc.Iterator()
	var hashInt big.Int

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("PrevBlockHash：%x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)
		g.Dump(block.Txs)

		fmt.Println()

		// 将迭代器中的hash存储到hashInt
		hashInt.SetBytes(blockchainIterator.CurrentHash)

		/*
			// Cmp compares x and y and returns:
			//
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
		*/
		if hashInt.Cmp(big.NewInt(0)) == 0 { // 这一步的用途是 判断是否到达创世区块
			break
		}
	}
}
