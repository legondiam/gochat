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
	flag       int
}

func NewClient(serverIp string, serverPort string) *Client {
	client := &Client{
		ServerIP:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

	//客户端业务
	client.Run()
	select {}
}
func (client *Client) Menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的数字")
		return false
	}
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式")
			break
		case 2:
			fmt.Println("私聊模式")
			break
		case 3:
			fmt.Println("更新用户名")
			break
		}
	}
}
