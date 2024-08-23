package main

import (
	"fmt"
	"publicChain/part3-Serialize-DeserializeBlock/BLC"
)

func main() {
	block := &BLC.Block{[]byte("Send 3 BTC to ZhangBozhi"), 1000}
	fmt.Printf("%s\n", block.Data)
	fmt.Printf("%d\n", block.Nonce)
	fmt.Printf("\n\n")

	bytes := block.Serialize()

	fmt.Println(bytes)
	fmt.Printf("\n\n")

	blc := BLC.DeserializeBlock(bytes)

	fmt.Println(blc)
	fmt.Printf("%s\n", blc.Data)
	fmt.Printf("%d\n", blc.Nonce)
	fmt.Printf("\n\n")
}
