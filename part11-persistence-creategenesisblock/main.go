package main

import (
	"fmt"
	"publicChain/part11-persistence-creategenesisblock/BLC"
)

// 16 进制
// 64 个数字
// 88cc2fff6c2d5b12da3dfa060f0f7aa60ebb35370249113a01832150c00d73ed
// 10001000
// 32 字节
// 256 bit

func main() {

	blockchain := BLC.CreateBlockChainWithGenesisBlock()
	fmt.Println(blockchain)
	fmt.Printf("tip：%x\n", blockchain.Tip)
}
