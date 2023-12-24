// main.go
package main

import "time"

// main 函数启动PBFT节点
func main() {
	node1 := NewPBFTNode(1)
	node2 := NewPBFTNode(2)
	node3 := NewPBFTNode(3)

	node1.Peers[2] = node2
	node1.Peers[3] = node3

	node2.Peers[1] = node1
	node2.Peers[3] = node3

	node3.Peers[1] = node1
	node3.Peers[2] = node2

	go StartRPCServer(node1)
	go StartRPCServer(node2)
	go StartRPCServer(node3)

	// 模拟交易
	time.Sleep(time.Second)
	transaction := Transaction{
		ID:     "1",
		From:   "Alice",
		To:     "Bob",
		Amount: 10,
	}

	node1.RequestTransaction(transaction, nil)

	time.Sleep(time.Second * 10)
}
