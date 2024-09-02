package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

type CLI struct {
}

// 打印参数信息
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -address -- 输出区块信息.")

}

// 判断终端参数的个数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {

	// 判断终端参数的个数，如果没有参数，直接打印Usage信息并且退出
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		// 添加区块
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		// 通过迭代器遍历区块链中的区块信息
		cli.printChain()
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
