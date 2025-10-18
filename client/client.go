package main

import (
	"flag"
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

var serverIP string
var serverPort string

func init() { //初始化命令行参数
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器IP地址（默认127.0.0.1）")
	flag.StringVar(&serverPort, "port", "8888", "设置服务器端口号（默认8888）")
}
func main() {
	//解析命令行参数
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功")
	select {}
}
