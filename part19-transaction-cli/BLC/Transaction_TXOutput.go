package BLC

// 交易输出
type TXOutput struct {
	Value        int64  // 分
	ScriptPubKey string // 用户名
}

// 检查是否能够解锁账号
func (out *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return out.ScriptPubKey == address
}
