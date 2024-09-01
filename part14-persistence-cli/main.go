package main

import "publicChain/part14-persistence-cli/BLC"

func main() {
	// 创建区块链
	blockchain := BLC.CreateBlockchainWithGenesisBlock()

	// 创建CLI对象
	cli := BLC.CLI{blockchain}

	// 调用cli的Run方法
	cli.Run()

}
