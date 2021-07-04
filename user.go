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
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + user.Name + ":" + "online......\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {

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
	} else if len(msg) > 3 && msg[:3] == "to|" {

		// find username
		receiver := strings.Split(msg, "|")[1]

		// find the user object
		_, ok := this.server.OnlineMap[receiver]
		if !ok {

			this.SendMessage("The user is not online or it does not exist")
			return
		}
		// send content
		sendMsg := strings.Split(msg, "|")[2]
		if sendMsg == "" {
			this.SendMessage("please input something to send\n")
		}
		receiverObj := this.server.OnlineMap[receiver]
		receiverObj.SendMessage(this.Name + ":" + sendMsg + "\n")

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
