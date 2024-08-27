package main

import "publicChain/part16-persistence-and-cli/BLC"

func main() {

	//blockchain := BLC.NewBlockChain()
	//
	//blockchain.AddBlock("Send 100 BTC To shaolin!!")
	//blockchain.AddBlock("Send 100 BTC To liming!!")
	//blockchain.AddBlock("Send 100 BTC To xiaoli!!")
	//blockchain.AddBlock("Send 100 BTC To fox!!")
	//blockchain.AddBlock("Send 100 BTC To hwq!!")
	//
	//fmt.Println("\n")
	//
	//blockchainIterator := blockchain.Iterator()
	//var hashInt big.Int
	//
	//for {
	//
	//	if err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
	//
	//		// 打开表
	//		b := tx.Bucket([]byte(blocksBucket))
	//		// 通过Hash获取区块字节数组
	//		blockBytes := b.Get(blockchainIterator.CurrentHash)
	//
	//		block := BLC.DeserializeBlock(blockBytes)
	//
	//		fmt.Printf("Data：%s \n", string(block.Data))
	//		fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
	//		fmt.Printf("Timestamp：%s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	//		fmt.Printf("Hash：%x \n", block.Hash)
	//		fmt.Printf("Nonce：%d \n", block.Nonce)
	//
	//		fmt.Println()
	//
	//		return nil
	//	}); err != nil {
	//		log.Panic(err)
	//	}
	//	// 获取下一个迭代器
	//	blockchainIterator = blockchainIterator.Next()
	//
	//	// 将迭代器中的hash存储到hashInt
	//	hashInt.SetBytes(blockchainIterator.CurrentHash)
	//
	//	/*
	//		// Cmp compares x and y and returns:
	//		//
	//		//   -1 if x <  y
	//		//    0 if x == y
	//		//   +1 if x >  y
	//	*/
	//	if hashInt.Cmp(big.NewInt(0)) == 0 { // 这一步的用途是 判断是否到达创世区块
	//		break
	//	}
	//}
	// 创建区块链
	blockchain := BLC.NewBlockChain()

	// 创建CLI对象
	cli := BLC.CLI{blockchain}

	// 调用cli的Run方法

	cli.Run()

}
