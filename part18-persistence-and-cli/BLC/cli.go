package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)

type CLI struct {
	Blockchain *Blockchain
}

// 打印参数信息
func (cli *CLI) printUsage() {

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("\tprintchain - print all the blocks of the blockchain")

}

// 判断终端参数的个数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
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
		fmt.Println("Data：" + *addBlockData)
	}

	if printChainCmd.Parsed() {
		// 通过迭代器遍历区块链中的区块信息
		cli.printChain()
	}
}

func (cli *CLI) printChain() {
	blockchainIterator := cli.Blockchain.Iterator()
	var hashInt big.Int

	for {

		if err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {

			// 打开表
			b := tx.Bucket([]byte(blocksBucket))
			// 通过Hash获取区块字节数组
			blockBytes := b.Get(blockchainIterator.CurrentHash)

			block := DeserializeBlock(blockBytes)

			fmt.Printf("Data：%s \n", string(block.Data))
			fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
			fmt.Printf("Hash：%x \n", block.Hash)
			fmt.Printf("Nonce：%d \n", block.Nonce)

			fmt.Println()

			return nil
		}); err != nil {
			log.Panic(err)
		}
		// 获取下一个迭代器
		blockchainIterator = blockchainIterator.Next()

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

func (cli *CLI) addBlock(data string) {
	cli.Blockchain.AddBlock(data)
}
