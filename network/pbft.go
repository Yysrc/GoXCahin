package network

import (
	"MYchain/blockchain"
	"MYchain/crypto"
	"crypto/ecdsa"
	"fmt"
	"net/rpc"
	"time"
)

const f = 1

// PBFT 消息
type PBFTMessage struct {
	Type       string           // 消息类型：PrePrepare/Prepare/Commit
	NodeID     int              // 发送消息的节点 ID
	NodePubkey *ecdsa.PublicKey // 发送消息的节点的公钥
	Timestamp  string           // 消息生成的时间戳
	Block      blockchain.Block // 相关区块
	Signature  string           // 签名
}

// 将 PBFT 消息转化为 string 便于验证
func messageToString(message PBFTMessage) string {
	return fmt.Sprintf("%s%d%s%s",
		message.Type,
		message.NodeID,
		message.Timestamp,
		message.Block.Index,
		message.Block.Hash,
		message.Block.PrevHash)
}

// RotateLeader 轮换主节点
func (node *PBFTNode) RotateLeader() {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	// 当前主节点ID加1，如果超过节点总数，则回到第一个节点
	node.LeaderID = (node.LeaderID + 1) % len(node.Peers)

	fmt.Printf("节点 %d 成为新的主节点\n", node.LeaderID)
}

func (node *PBFTNode) LeaderRoutine() {
	// 定期轮换主节点
	for {
		time.Sleep(10 * time.Second)
		node.RotateLeader()
	}
}

// 将 PBFT 消息广播给其他节点
func (node *PBFTNode) BroadcastPBFTMessage(message PBFTMessage) {
	for _, peer := range node.Peers {
		if peer.ID != node.ID {
			go func(peer *PBFTNode) {
				client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", peer.ID+8000))
				if err != nil {
					fmt.Println("连接到节点时发生错误:", err)
					return
				}
				defer client.Close()

				var response bool
				err = client.Call("RPCService.ReceivePBFTMessage", message, &response)
				if err != nil {
					fmt.Println("广播 PBFT 消息时发生错误:", err)
				}
			}(peer)
		}
	}
}

// 处理接收到的 PBFT 消息
func (r *RPCService) ReceivePBFTMessage(message PBFTMessage, response *bool) error {
	r.Node.Mutex.Lock()
	defer r.Node.Mutex.Unlock()

	FinishPrepare := false

	switch message.Type {
	case "PrePrepare":
		r.Node.PreparePhase(message)
	case "Prepare":
		if FinishPrepare {
			r.Node.CommitPhase(message)
		}
		if r.Node.ID == r.Node.LeaderID {
			Exist := false
			for PBFTMessage, num := range r.Node.PrepareMsgs {
				if message.Block == PBFTMessage.Block {
					Exist = true
					num++
					if num >= 2*f {
						FinishPrepare = true
						r.Node.CommitPhase(message)
						break
					}
				}
			}
			if !Exist {
				r.Node.PrepareMsgs[&message] = 1
			}
		}
	case "Commit":
		Exist := false
		for PBFTMessage, num := range r.Node.CommitMsgs {
			if message.Block == PBFTMessage.Block {
				Exist = true
				num++
				if num >= 2*f {
					r.Node.Blockchain = append(r.Node.Blockchain, message.Block)
					break
				}
			}
		}
		if !Exist {
			r.Node.PrepareMsgs[&message] = 1
		}
	}

	return nil
}

// 主节点执行 PBFT 中的预准备阶段
func (node *PBFTNode) PrePreparePhase(block blockchain.Block) {
	// 是否是主节点
	if node.ID != node.LeaderID {
		fmt.Println("PrePrepare: 当前节点不是主节点")
		return
	}
	// 生成 PrePrepare 消息
	preprepareMessage := PBFTMessage{
		Type:       "PrePrepare",
		Block:      block,
		NodeID:     node.ID,
		NodePubkey: node.PublicKey,
		Timestamp:  time.Now().String(),
	}
	// 对消息进行签名
	signature, err := crypto.ECDSASign(node.PrivateKey, messageToString(preprepareMessage))
	if err != nil {
		fmt.Println("生成 PrePrepare 消息签名时发生错误:", err)
		return
	}
	preprepareMessage.Signature = signature
	// 将消息广播给其他节点
	node.BroadcastPBFTMessage(preprepareMessage)
}

// 执行 PBFT 中的准备阶段
func (node *PBFTNode) PreparePhase(message PBFTMessage) {
	// 验证消息签名
	if !crypto.ECDSAVerify(message.NodePubkey, messageToString(message), message.Signature) {
		fmt.Println("PreparePhase: 消息签名验证失败")
		return
	}

	// 生成 Prepare 消息
	prepareMessage := PBFTMessage{
		Type:       "Prepare",
		Block:      message.Block,
		NodeID:     node.ID,
		NodePubkey: node.PublicKey,
		Timestamp:  time.Now().String(),
	}
	// 对消息进行签名
	signature, err := crypto.ECDSASign(node.PrivateKey, messageToString(prepareMessage))
	if err != nil {
		fmt.Println("生成 Prepare 消息签名时发生错误:", err)
		return
	}

	prepareMessage.Signature = signature

	// 将 Prepare 消息广播给其他节点
	node.BroadcastPBFTMessage(prepareMessage)
}

// 执行 PBFT 中的提交阶段
func (node *PBFTNode) CommitPhase(message PBFTMessage) {
	// 验证消息签名
	if !crypto.ECDSAVerify(message.NodePubkey, messageToString(message), message.Signature) {
		fmt.Println("消息签名验证失败")
		return
	}
	// 生成 Commit 消息
	commitMessage := PBFTMessage{
		Type:       "Commit",
		Block:      message.Block,
		NodeID:     node.ID,
		NodePubkey: node.PublicKey,
		Timestamp:  time.Now().String(),
	}

	// 对消息进行签名
	signature, err := crypto.ECDSASign(node.PrivateKey, messageToString(commitMessage))
	if err != nil {
		fmt.Println("生成 Commit 消息签名时发生错误:", err)
		return
	}

	commitMessage.Signature = signature

	// 将 Commit 消息广播给其他节点
	node.BroadcastPBFTMessage(commitMessage)
}
