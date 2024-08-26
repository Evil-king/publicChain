package main

import (
	"fmt"
	"publicChain/part11-persistence-and-cli/BLC"
	"time"
)

// 16 进制
// 64 个数字
// 88cc2fff6c2d5b12da3dfa060f0f7aa60ebb35370249113a01832150c00d73ed
// 10001000
// 32 字节
// 256 bit

func main() {

	blockchain := BLC.NewBlockChain()

	blockchain.AddBlock("Send 20 BTC To HaoLin From Liyuechun")

	blockchain.AddBlock("Send 10 BTC To SaoLin From Liyuechun")

	blockchain.AddBlock("Send 30 BTC To HaoTian From Liyuechun")

	for _, block := range blockchain.Blocks {

		fmt.Printf("Data：%s \n", string(block.Data))
		fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x \n", block.Hash)
		fmt.Printf("Nonce：%d \n", block.Nonce)

		fmt.Println()
	}
}
