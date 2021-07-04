package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// create a client obj
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// connect with server
	serverAddr := fmt.Sprintf("%s:%d", serverIp, serverPort)
	con, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = con

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("Connection faile")
		return
	}

	fmt.Println("Connection successfully")

	select {}
}
