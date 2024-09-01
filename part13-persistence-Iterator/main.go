package main

import "publicChain/part13-persistence-Iterator/BLC"

func main() {
	// 创建区块链
	blockchain := BLC.CreateBlockchainWithGenesisBlock()

	//新区块
	blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang")

	blockchain.AddBlockToBlockchain("Send 200RMB To changjingkong")

	blockchain.AddBlockToBlockchain("Send 300RMB To juncheng")

	blockchain.AddBlockToBlockchain("Send 50RMB To haolin")

	blockchain.Printchain()

}
