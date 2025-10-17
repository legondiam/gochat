package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort string
	Conn       net.Conn
	ClientName string
}

func NewClient(serverIp string, serverPort string) *Client {
	client := &Client{
		ServerIP:   serverIp,
		ServerPort: serverPort,
	}
	conn, err := net.Dial("tcp", fmt.Sprint(serverIp, ":", serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.Conn = conn
	return client
}
func main() {
	client := NewClient("127.0.0.1", "8888")
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功")
}
