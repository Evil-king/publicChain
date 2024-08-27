package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
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
		fmt.Println("printchain,printchain,printchain!")
	}
}
