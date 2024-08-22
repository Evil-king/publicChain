package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

func (block *Block) SetHash() {
	//1、将时间戳转化为字节数组
	timeString := strconv.FormatInt(block.TimeStamp, 2) // 2～36-->代表进制
	//fmt.Println(timeString)
	timeStampBytes := []byte(timeString)
	//2、拼接所有属性
	blockBytes := bytes.Join([][]byte{timeStampBytes, block.PrevBlockHash, block.Data}, []byte{})
	//3、生成hash
	hash := sha256.Sum256(blockBytes)
	//4、将hash赋给Hash属性字节
	block.Hash = hash[:]
}

// 工厂方法
func NewBlock(data string, prevBlockHash []byte) *Block {
	// 创建区块
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Hash:          []byte{},
	}
	// 设置当前区块的Hash值
	block.SetHash()
	return block
}

// 创建创世区块，并返回创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
