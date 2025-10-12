package main

import "net"

type User struct {
	Username    string
	Useraddr    string
	UserChannel chan string
	conn        net.Conn
}

// 创建新用户
func NewUser(conn net.Conn) *User {
	user := &User{
		Username:    conn.RemoteAddr().String(),
		Useraddr:    conn.RemoteAddr().String(),
		UserChannel: make(chan string),
		conn:        conn,
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
