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
func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

// 创建创世区块
func (cli *CLI) createGenesisBlockchain(address string) {

	blockchain := CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()
}

func (cli *CLI) printchain() {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.printChain()

}

// 先用它去查询余额
func (cli *CLI) getBalance(address string) {

	fmt.Println("地址：" + address)

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)

}

// 转账
func (cli *CLI) send(from []string, to []string, amount []string) {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)

}

func (cli *CLI) Run() {

	// 判断终端参数的个数，如果没有参数，直接打印Usage信息并且退出
	validateArgs()

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址......")
	flagTo := sendBlockCmd.String("to", "", "转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额......")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块的地址")
	getbalanceWithAddress := getBalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from, to, amount)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getBalanceCmd.Parsed() {
		if *getbalanceWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAddress)
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
