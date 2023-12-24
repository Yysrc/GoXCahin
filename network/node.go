package network

import (
	"MYchain/blockchain"
	"MYchain/crypto"
	"MYchain/transaction"
	"crypto/ecdsa"
	"fmt"
	"sync"
)

// PBFTNode 表示 PBFT 网络中的一个节点
type PBFTNode struct {
	ID          int                       // 节点 ID
	LeaderID    int                       // 主节点ID
	Blockchain  []blockchain.Block        // 节点维护的区块链
	PendingTx   []transaction.Transaction // 待处理的交易池
	Peers       map[int]*PBFTNode         // 与该节点相连的其他节点列表
	Mutex       sync.Mutex                // 用于保护节点数据的互斥锁
	PublicKey   *ecdsa.PublicKey          // 公钥
	PrivateKey  *ecdsa.PrivateKey         // 私钥
	PrepareMsgs map[*PBFTMessage]int      //
	CommitMsgs  map[*PBFTMessage]int      //
}

// NewPBFTNode 创建具有指定ID的新PBFT节点
func NewPBFTNode(id int) *PBFTNode {
	// 生成 ECDSA 密钥对
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		fmt.Println("生成密钥对时发生错误:", err)
		return nil
	}

	return &PBFTNode{
		ID:         id,
		Blockchain: []blockchain.Block{blockchain.createGenesisBlock()},
		PendingTx:  []transaction.Transaction{},
		Peers:      make(map[int]*PBFTNode),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// GetPublicKey 返回节点的公钥
func (node *PBFTNode) GetPublicKey() *ecdsa.PublicKey {
	return node.PublicKey
}
