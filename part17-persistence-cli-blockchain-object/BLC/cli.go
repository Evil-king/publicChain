package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

// 打印参数信息
func printUsage() {

	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchainwithgenesis -data -- 交易数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")
}

// 判断终端参数的个数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	if !dbExists() {
		fmt.Println("数据不存在......")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.AddBlockToBlockchain(data)
}

func (cli *CLI) printchain() {
	if !dbExists() {
		fmt.Println("数据不存在......")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}

func (cli *CLI) createGenesisBlockchain(data string) {
	CreateBlockchainWithGenesisBlock(data)
}

func (cli *CLI) Run() {

	// 判断终端参数的个数，如果没有参数，直接打印Usage信息并且退出
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "", "交易数据......")
	flagCreateBlockchainWithData := createBlockchainCmd.String("data", "Genesis block data......", "创世区块交易数据......")

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
	case "crateblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockData)
	}

	if createBlockchainCmd.Parsed() {
		if *flagCreateBlockchainWithData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithData)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}
}
