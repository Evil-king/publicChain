package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gogf/gf/v2/frame/g"
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
	fmt.Println("\tgetbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("\tcreateblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("\tprintchain - Print all the blocks of the blockchain:")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")

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

	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	genesisAddress := createBlockchainCmd.String("address", "", "创建创世区块，并且将数据打包到数据库.")

	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	balanceAddress := getBalanceCmd.String("address", "", "余额查询....")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "源地址...")
	sendTo := sendCmd.String("to", "", "目标地址...")
	sendAmount := sendCmd.Int("amount", 0, "转账的额度....")

	switch os.Args[1] {
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *genesisAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		// 添加区块
		fmt.Println("创建创世区块并且存储到数据库....")
	}

	if getBalanceCmd.Parsed() {
		if *balanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		// 查询余额
		fmt.Printf("查询 %s 的余额....\n", *balanceAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			cli.printUsage()
			os.Exit(1)
		}
		// 转账
		fmt.Printf("from：%s to：%s amount：%d\n", *sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		// 通过迭代器遍历区块链中的区块信息
		cli.printChain()
	}
}

func (cli *CLI) printChain() {

	// 判断数据库是否存在
	if !DbExists() {
		cli.printUsage()
		return
	}

	blockchainIterator := cli.Blockchain.Iterator()
	var hashInt big.Int

	for {

		if err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {

			// 打开表
			b := tx.Bucket([]byte(blocksBucket))
			// 通过Hash获取区块字节数组
			blockBytes := b.Get(blockchainIterator.CurrentHash)

			block := DeserializeBlock(blockBytes)

			g.Dump(block.Transactions)
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

func (cli *CLI) sendToken() {

	// 1. 10 -> liyuechun
	// 2. 3 -> 转给 chenhaolinls
	// (1) 新建一个交易
	tx1 := NewUTXOTransaction("hwq", "xiaoming", 3, cli.Blockchain)
	tx2 := NewUTXOTransaction("hwq", "saolin", 2, cli.Blockchain)
	cli.Blockchain.MineBlock([]*Transaction{tx1, tx2})
}

func (cli *CLI) addBlock(data string) {
	cli.sendToken()
}
