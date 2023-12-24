package blockchain

import (
	"MYchain/transaction"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block 结构定义了区块链中的一个区块
type Block struct {
	Index        int                       // 区块在区块链中的索引
	Timestamp    string                    // 区块生成的时间戳
	Transactions []transaction.Transaction // 区块包含的交易列表
	Hash         string                    // 区块的哈希值，用于唯一标识区块
	PrevHash     string                    // 上一个区块的哈希值，建立区块链中的链接
	Validator    int                       // 区块的验证者（节点）ID
}

// createGenesisBlock 创建区块链的创世区块
func createGenesisBlock() Block {
	return Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []transaction.Transaction{},
		Hash:         calculateHash(0, "", time.Now().String(), ""),
		PrevHash:     "",
		Validator:    0,
	}
}

// calculateHash 计算区块的哈希值
func calculateHash(index int, prevHash, timestamp, data string) string {
	payload := fmt.Sprintf("%d%s%s%s", index, prevHash, timestamp, data)
	hash := sha256.New()
	hash.Write([]byte(payload))
	return hex.EncodeToString(hash.Sum(nil))
}
