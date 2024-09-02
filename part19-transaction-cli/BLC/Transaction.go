package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	//1. 交易hash
	TxHash []byte
	// 2. 交易输入
	Vin []*TXInput
	// 3. 交易输出
	Vouts []*TXOutput
}

// 判断当前交易是否是CoinbaseTX(创世区块的交易)
func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.Vin[0].TxHash) == 0 && tx.Vin[0].Vout == -1 && len(tx.Vin) == 1
}

// 1. Transaction 创建分两种情况
// 创世区块创建时的Transaction
func NewCoinbaseTransaction(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txInput := &TXInput{[]byte{}, -1, data}
	txOutput := &TXOutput{10, to}
	tx := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	tx.HashTransaction()
	return tx
}

// 2. 转账时产生的Transaction
func NewSimpleTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	// 输入
	var txInputs []*TXInput
	// 输出
	var txOutputs []*TXOutput

	// 获取该地址对应的可用余额
	money, spendableUTXODic := bc.FindSpendableUTXOS(from, amount)

	// 建立输入
	for txHash, indexArray := range spendableUTXODic {
		txID, _ := hex.DecodeString(txHash)
		for _, value := range indexArray {
			txInputs = append(txInputs, &TXInput{
				TxHash:    txID,
				Vout:      value,
				ScriptSig: from,
			})
		}
	}
	// 建立输出，转账
	txOutputs = append(txOutputs, &TXOutput{
		Value:        int64(amount),
		ScriptPubKey: to,
	})
	// 建立输出，找零
	txOutputs = append(txOutputs, &TXOutput{
		Value:        money - int64(amount),
		ScriptPubKey: from,
	})
	// 创建交易
	tx := Transaction{[]byte{}, txInputs, txOutputs}
	// 设置hash值
	tx.HashTransaction()
	return &tx
}

func (tx *Transaction) HashTransaction() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}
