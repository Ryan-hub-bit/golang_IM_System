package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// create a user api to adduser
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListenMessage()

	return user
}

// 监听当前channle 的方法, 一旦又消息就发送给相应的client 因为要持续监听, 所以不应该是一个普通方法, 因为他要和 main 同时进行
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// make it to the client when it comes in channel
		this.conn.Write([]byte(msg + "\n"))

	}
}
