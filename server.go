package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// we want Sever to be pulbic used, so it is cap
type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// because we want to create a new server, and we also can modify it in real address, so we need to return the address of the server
// it is just a common interface
func NewServer(ip string, port int) *Server {
	// get the address of the server so we can modify it and change it
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// listen broadcase
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(con net.Conn) {
	// current business
	user := NewUser(con, this)

	user.Online()

	// 监控当前用户是否活跃的channel
	isLive := make(chan bool)

	// receive message from client
	go func() {
		buf := make([]byte, 4096)
		for {
			n, error := con.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if error != nil && error != io.EOF {
				fmt.Println("Conn Read Error:", error)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			// do nothing, it will activate the following case to reset the timer
		case <-time.After(time.Second * 80):
			user.SendMessage("You are offline because of timeout")
			close(user.C)
			con.Close()
			return
		}
	}
}

// create a method of current class
func (this *Server) Start() {
	// socket listen
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if error != nil {
		fmt.Println("net.listen error", error)
		return
	}
	// close listen socket
	defer listener.Close()
	go this.ListenMessager()
	// TCP is based on bit transmit, so we have to use a for loop to receive them one by one
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}
