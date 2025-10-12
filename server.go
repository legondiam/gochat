package main

import (
	"fmt"
	"io"
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
	s.BroadCast(user, "已上线")
	go s.UserMessage(conn, user)
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

// 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	Message := "[" + user.Username + "]" + user.Useraddr + ":" + msg + "\n"
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
func (s *Server) UserMessage(conn net.Conn, user *User) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			s.BroadCast(user, "已下线")
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			return
		}
		msg := string(buf[:n-1])
		s.BroadCast(user, msg)
	}
}
