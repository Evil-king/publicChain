package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	// 时间戳 创建区块的时间
	TimeStamp int64
	// 上一个区块的hash
	PrevBlockHash []byte
	// Data，交易数据
	Data []byte
	// Hash 当前区块的hash
	Hash []byte
	// Nonce 随机数
	Nonce int
}

// 工厂方法
func NewBlock(data string, prevBlockHash []byte) *Block {
	// 创建区块
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
	}
	// 将block作为参数，创建一个pow对象
	pow := NewProofOfWork(block)

	// Run()执行一次工作量证明
	nonce, hash := pow.Run()

	// 设置区块Hash
	block.Hash = hash

	// 设置Nonce
	block.Nonce = nonce

	// 验证工作量证明
	//isValid := pow.Validate()
	//fmt.Println(isValid)

	// 返回block
	return block
}

// 将Block对象序列化成[]byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 将字节数组反序列化成Block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// 创建创世区块，并返回创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
