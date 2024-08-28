package BLC

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
