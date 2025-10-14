package main

import "net"

type User struct {
	Username    string
	Useraddr    string
	UserChannel chan string
	conn        net.Conn
	s           *Server
}

// 创建新用户
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Username:    conn.RemoteAddr().String(),
		Useraddr:    conn.RemoteAddr().String(),
		UserChannel: make(chan string),
		conn:        conn,
		s:           server,
	}
	go user.ListenUser()
	return user
}

// 监听用户方法
func (user *User) ListenUser() {
	for {
		message := <-user.UserChannel
		user.conn.Write([]byte(message))
	}
}

// 用户上线
func (user *User) Online() {
	user.s.mutex.Lock()
	user.s.OnlineMap[user.Username] = user
	user.s.mutex.Unlock()
	user.s.BroadCast(user, "已上线")
}

// 用户下线
func (user *User) Offline() {
	user.s.mutex.Lock()
	delete(user.s.OnlineMap, user.Username)
	user.s.mutex.Unlock()
	user.s.BroadCast(user, "已下线")
}

// 发送在线列表
func (user *User) SendOnlineUsers(usermsg string) {
	user.conn.Write([]byte(usermsg))
}

// 用户消息处理
func (user *User) DoMessage(msg string) {
	//在线用户查询
	if msg == "who" {
		user.s.mutex.Lock()
		for _, OnlineUser := range user.s.OnlineMap {
			usermsg := "[" + OnlineUser.Username + "]" + OnlineUser.Useraddr + ":在线\n"
			user.SendOnlineUsers(usermsg)
		}
		user.s.mutex.Unlock()
	} else {
		user.s.BroadCast(user, msg)
	}
}
