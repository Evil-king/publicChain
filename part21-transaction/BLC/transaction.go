package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type Transaction struct {
	// 1. 交易ID
	ID []byte
	// 2. 交易输入
	Vin []TXInput
	// 3. 交易输出
	Vout []TXOutput
}

// 交易输入
type TXInput struct {
	// 1. 交易的ID
	Txid []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout int
	// 3. 用户名
	ScriptSig string
}

// 交易输出
type TXOutput struct {
	Value        int    // 分
	ScriptPubKey string // 用户名
}

// 创建一个新的 coinbase 交易(创世区块对应的交易)
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txInput := TXInput{[]byte{}, -1, data}
	txOutput := TXOutput{10, to}
	tx := &Transaction{nil, []TXInput{txInput}, []TXOutput{txOutput}}
	tx.SetID()
	return tx
}

// 设置交易hash(这里将Transaction序列化以后生成hash)
func (tx *Transaction) SetID() {

	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的字节数组生成256hash
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// 检查账号地址，解锁
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// 检查是否能够解锁账号
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
