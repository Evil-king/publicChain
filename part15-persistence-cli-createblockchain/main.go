package main

import "publicChain/part15-persistence-cli-createblockchain/BLC"

func main() {
	// 创建CLI对象
	cli := BLC.CLI{}

	// 调用cli的Run方法
	cli.Run()

}
