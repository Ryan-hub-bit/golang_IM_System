package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// create a user api to adduser
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// broadcase message
	this.server.Broadcast(this, "已上线")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// broadcase message
	this.server.Broadcast(this, "已下线")
}

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "All users" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + user.Name + ":" + "online......\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if msg[:7] == "rename|" {

		newName := strings.Split(msg, "|")[1]

		// check if the newName has already been used
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("newName has already been used")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("Your name has been modified to:" + newName + "\n")
		}

	} else {
		this.server.Broadcast(this, msg)
	}
}

// 监听当前channle 的方法, 一旦又消息就发送给相应的client 因为要持续监听, 所以不应该是一个普通方法, 因为他要和 main 同时进行
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// make it to the client when it comes in channel
		this.conn.Write([]byte(msg + "\n"))

	}
}
