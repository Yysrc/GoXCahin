package network

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// RPCService 表示节点通信的RPC服务
type RPCService struct {
	Node *PBFTNode // 指向 PBFTNode 类型的指针
}

// StartRPCServer 启动节点通信的RPC服务器
func StartRPCServer(node *PBFTNode) {
	rpcService := &RPCService{Node: node}
	rpc.Register(rpcService)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", node.ID+8000))
	if err != nil {
		fmt.Println("启动RPC服务器时发生错误:", err)
		return
	}
	fmt.Printf("节点 %d 的RPC服务器正在监听端口 %d\n", node.ID, node.ID+8000)
	http.Serve(listener, nil)
}
