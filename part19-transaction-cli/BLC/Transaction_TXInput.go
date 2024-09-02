package BLC

// 交易输入
type TXInput struct {
	// 1. 交易的hash
	TxHash []byte
	// 2. 存储TXOutput在Vout里面的索引---->通俗理解就是输入和输出的对应关系
	Vout int
	// 3. 用户名
	ScriptSig string
}

// 检查账号地址，解锁
func (in *TXInput) UnLockWithAddress(address string) bool {
	return in.ScriptSig == address
}
