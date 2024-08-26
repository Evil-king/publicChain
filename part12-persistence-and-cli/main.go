package main

import (
	"fmt"
	"publicChain/part12-persistence-and-cli/BLC"
)

// 16 进制
// 64 个数字
// 88cc2fff6c2d5b12da3dfa060f0f7aa60ebb35370249113a01832150c00d73ed
// 10001000
// 32 字节
// 256 bit

func main() {

	blockchain := BLC.NewBlockChain()
	fmt.Println(blockchain)
	fmt.Printf("tip：%x\n", blockchain.Tip)
	blockchain.AddBlock("Send 100 BTC To shaolin!!")
	fmt.Printf("tip：%x\n", blockchain.Tip)
}
