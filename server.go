package main

import (
	"fmt"
	"net"
)

// we want Sever to be pulbic used, so it is cap
type Server struct {
	Ip   string
	Port int
}

// because we want to create a new server, and we also can modify it in real address, so we need to return the address of the server
// it is just a common interface
func NewServer(ip string, port int) *Server {
	// get the address of the server so we can modify it and change it
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(con net.Conn) {
	// current business
	fmt.Println("do something....")
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
