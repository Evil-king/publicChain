package main

import (
	"fmt"
	"math/big"
	"publicChain/part13-persistence-and-cli/BLC"
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

	blockchainIterator := blockchain.Iterator()
	var hashInt big.Int

	for {
		fmt.Printf("%x\n", blockchainIterator.CurrentHash)
		// 获取下一个区块对象
		blockchainIterator = blockchainIterator.Next()
		// 将迭代器中的hash存储到hashInt
		hashInt.SetBytes(blockchainIterator.CurrentHash)

		/*
			// Cmp compares x and y and returns:
			//
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
		*/
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

}
