package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
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
	// 2. 存储TXOutput在Vout里面的索引---->通俗理解就是输入和输出的对应关系
	Vout int
	// 3. 用户名
	ScriptSig string
}

// 交易输出
type TXOutput struct {
	Value        int    // 分
	ScriptPubKey string // 用户名
}

// 判断当前交易是否是CoinbaseTX(创世区块的交易)
func (tx *Transaction) IsCoinbase() bool {

	return len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1 && len(tx.Vin) == 1
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

// 建立转账交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	// 输入
	var inputs []TXInput
	// 输出
	var outputs []TXOutput

	// 获取该地址对应的可用余额
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	fmt.Printf("%d", acc)
	g.Dump(validOutputs)
	// 可用余额
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	// 建立输入
	for index, out := range validOutputs {
		txID, err := hex.DecodeString(index)
		if err != nil {
			log.Panic(err)
		}
		for _, value := range out {
			inputs = append(inputs, TXInput{
				Txid:      txID,
				Vout:      value,
				ScriptSig: from,
			})
		}
	}
	// 建立输出，转账
	outputs = append(outputs, TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	})
	// 建立输出，找零
	outputs = append(outputs, TXOutput{
		Value:        acc - amount,
		ScriptPubKey: from,
	})
	// 创建交易
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
