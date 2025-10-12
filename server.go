package main

import (
	"fmt"
	"net"
	"sync"
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
	user := NewUser(conn) //创建用户
	s.mutex.Lock()
	s.OnlineMap[user.Username] = user
	s.mutex.Unlock()
	s.BroadCast(user)
	select {}
}

// 启动服务器
func (s *Server) Start() {
	//socket(listening)
	listener, eil := net.Listen("tcp", fmt.Sprintf("%s:%s", s.IP, s.Port))
	if eil != nil {
		fmt.Println("listrn error:", eil)
	}
	defer listener.Close()
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

// 广播用户上线消息
func (s *Server) BroadCast(user *User) {
	Message := "[" + user.Username + "]" + user.Useraddr + ": 已上线"
	s.ServerChannel <- Message
}

// 监听ServerChannel并发送广播给user
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
