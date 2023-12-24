package transaction

import (
	"MYchain/crypto"
	"crypto/ecdsa"
	"fmt"
	"net/rpc"
)

// Transaction 表示区块链中的交易
type Transaction struct {
	ID     string
	From   string
	To     string
	Amount int
}

// 使用ECDSA签名验证交易
func validateTransaction(transaction Transaction, publicKey *ecdsa.PublicKey) bool {
	message := fmt.Sprintf("%s%s%s%d", transaction.ID, transaction.From, transaction.To, transaction.Amount)
	return crypto.ECDSAVerify(publicKey, message, transaction.ID)
}

// RequestTransaction 处理请求新交易的RPC请求
func (r *RPCService) RequestTransaction(request Transaction, response *bool) error {
	r.Node.Mutex.Lock()
	defer r.Node.Mutex.Unlock()

	// 验证交易
	if !validateTransaction(request, &r.Node.GetPublicKey()) {
		fmt.Println("无效的交易签名")
		return nil
	}

	// 将交易添加到待处理交易池
	r.Node.PendingTx = append(r.Node.PendingTx, request)

	// 将交易广播给其他节点
	for _, peer := range r.Node.Peers {
		if peer.ID != r.Node.ID {
			go func(peer *PBFTNode) {
				client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", peer.ID+8000))
				if err != nil {
					fmt.Println("连接到节点时发生错误:", err)
					return
				}
				defer client.Close()

				var response bool
				err = client.Call("RPCService.BroadcastTransaction", request, &response)
				if err != nil {
					fmt.Println("广播交易时发生错误:", err)
				}
			}(peer)
		}
	}

	return nil
}

// BroadcastTransaction 处理广播交易给其他节点的RPC请求
func (r *RPCService) BroadcastTransaction(request Transaction, response *bool) error {
	r.Node.Mutex.Lock()
	defer r.Node.Mutex.Unlock()

	// 验证交易
	if !validateTransaction(request, &r.Node.GetPublicKey()) {
		fmt.Println("无效的交易签名")
		return nil
	}

	// 将交易添加到待处理交易池
	r.Node.PendingTx = append(r.Node.PendingTx, request)

	return nil
}
