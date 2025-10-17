package main

import (
	"net"
	"strings"
)

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
		message, ok := <-user.UserChannel
		if !ok {
			return
		}
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

// 单独用户发送消息
func (user *User) SendMessage(usermsg string) {
	user.conn.Write([]byte(usermsg))
}

// 用户消息处理
func (user *User) DoMessage(msg string) {
	//在线用户查询
	if msg == "who" {
		user.s.mutex.Lock()
		for _, OnlineUser := range user.s.OnlineMap {
			usermsg := "[" + OnlineUser.Useraddr + "]" + OnlineUser.Username + ":在线\n"
			user.SendMessage(usermsg)
		}
		user.s.mutex.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" { //用户名更改
		newname := msg[7:]
		_, ok := user.s.OnlineMap[newname]
		if ok {
			user.SendMessage("该用户名已被占用\n")
		} else {
			user.s.mutex.Lock()
			delete(user.s.OnlineMap, user.Username)
			user.s.OnlineMap[newname] = user
			user.s.mutex.Unlock()
			user.Username = newname
			user.SendMessage("用户名已更改为" + user.Username + "\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" { //私聊
		toUser := strings.Split(msg, "|")[1]
		if toUser == "" {
			user.SendMessage("请正确输入用户名，格式为to|用户名|消息内容\n")
			return
		} else if toUser == user.Username {
			user.SendMessage("不能私聊自己\n")
			return
		}

		_, ok := user.s.OnlineMap[toUser]
		if !ok {
			user.SendMessage("该用户不在线")
			return
		} else {
			usermsg := strings.Split(msg, "|")[2]
			if usermsg == "" {
				user.SendMessage("请正确输入消息内容，格式为to|用户名|消息内容\n")
			} else {
				user.s.OnlineMap[toUser].SendMessage(user.Username + "私聊你：" + usermsg)
			}
		}
	} else {
		user.s.BroadCast(user, msg)
	}
}
