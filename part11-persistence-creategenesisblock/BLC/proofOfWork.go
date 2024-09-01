package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 20

type ProofOfWork struct {
	// 区块
	block *Block
	// 目标难度值
	target *big.Int
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.TimeStamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// ProofOfWork对象的方法
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		// pow.target hash
		// hashInt < pow.target -1
		// hashInt = pow.target 0
		// hashInt > pow.target 1
		// 如果hashInt < pow.target ，break
		// 000010000000
		// 000001000000
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

// 创建工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) // 左移动
	pow := &ProofOfWork{block, target}
	return pow
}
