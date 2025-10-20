package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

	go client.DealResponse()
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

// 修改用户名
func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&client.ClientName)
	newname := "rename|" + client.ClientName
	_, err := client.Conn.Write([]byte(newname))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

// 公聊模式
func (client *Client) PublicChat() {
	var chatMsg string
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendmsg := chatMsg
			_, err := client.Conn.Write([]byte(sendmsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		fmt.Println("请输入聊天内容（exit退出）")
		chatMsg = ""
		fmt.Scanln(&chatMsg)
	}
}

// 查询用户
func (client *Client) OnlineUsers() {
	sendmsg := "who"
	_, err := client.Conn.Write([]byte(sendmsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
	}
}

// 私聊模式
func (client *Client) PrivateChat() {
	var chatMsg string
	client.OnlineUsers()
	var remoteName string
	fmt.Println("请输入聊天对象用户名(exit退出)")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		for chatMsg != "exit" {
			if chatMsg != "" {
				sendmsg := "to|" + remoteName + "|" + chatMsg //发送消息时 server才知道remotename
				_, err := client.Conn.Write([]byte(sendmsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					break
				}
			}
			fmt.Println("请输入聊天内容(exit退出)")
			chatMsg = ""
			fmt.Scanln(&chatMsg)
		}
		fmt.Println("请输入聊天对象用户名(exit退出)")
		remoteName = ""
		fmt.Scanln(&remoteName)
	}
}

// 监听消息
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.Conn)
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}
		switch client.flag {
		case 1:
			//公聊模式
			client.PublicChat()
			break
		case 2:
			//私聊模式
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			client.UpdateName()
			break
		}
	}
}
