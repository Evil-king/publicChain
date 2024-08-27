package main

import "publicChain/part18-persistence-and-cli/BLC"

func main() {
	// 创建区块链
	blockchain := BLC.NewBlockChain()

	// 创建CLI对象
	cli := BLC.CLI{blockchain}

	// 调用cli的Run方法

	cli.Run()

}
