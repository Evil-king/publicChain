package main

import "publicChain/part17-persistence-cli-blockchain-object/BLC"

func main() {
	// 创建CLI对象
	cli := BLC.CLI{}

	// 调用cli的Run方法
	cli.Run()

}
