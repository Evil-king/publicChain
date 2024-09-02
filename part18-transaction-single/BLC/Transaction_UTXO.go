package BLC

// 未花费结构体
type UTXO struct {
	TxHash []byte
	Index  int
	Output *TXOutput
}
