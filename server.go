package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	IP            string
	Port          string
	OnlineMap     map[string]*User
	ServerChannel chan string
	mutex         sync.Mutex
}

// 创建服务器
func NewServer(ip string, port string) *Server {
	server := &Server{
		IP:            ip,
		Port:          port,
		OnlineMap:     make(map[string]*User),
		ServerChannel: make(chan string),
		mutex:         sync.Mutex{},
	}
	return server
}

// 连接处理
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("连接成功")
	user := NewUser(conn, s)  //创建用户
	user.Online()             //用户上线
	islive := make(chan bool) //用户活跃标志
	//消息处理线程
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) //读取用户消息
			if n == 0 {
				user.Offline() //n=0时用户下线
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("read error:", err)
				return
			}
			msg := strings.TrimSpace(string(buf[:n])) //转String并去除换行符
			user.DoMessage(msg)                       //处理用户消息
			islive <- true                            //用户活跃中
		}
	}()
	//超时处理
	for {
		select {
		case <-islive: //用户活跃，重置定时器

		case <-time.After(time.Second * 60): //用户60秒无消息则踢出
			user.SendMessage("你已被踢出服务器")
			close(user.UserChannel) //关闭用户消息通道
			conn.Close()            //关闭连接
			return
		}
	}
}

// 启动服务器
func (s *Server) Start() {
	//socket(listening)
	listener, eil := net.Listen("tcp", fmt.Sprintf("%s:%s", s.IP, s.Port))
	if eil != nil {
		fmt.Println("listen error:", eil)
	}
	//关闭连接
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("listener close error:", err)
		}
	}()
	go s.ListenGoroutine()
	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go s.Handler(conn)
	}
}

// 广播消息（写入ServerChannel）
func (s *Server) BroadCast(user *User, msg string) {
	Message := "[" + user.Useraddr + "]" + user.Username + ":" + msg + "\n"
	s.ServerChannel <- Message
}

// 监听ServerChannel并发送user
func (s *Server) ListenGoroutine() {
	for {
		msg := <-s.ServerChannel
		s.mutex.Lock()
		for _, user := range s.OnlineMap {
			user.UserChannel <- msg
		}
		s.mutex.Unlock()
	}
}
