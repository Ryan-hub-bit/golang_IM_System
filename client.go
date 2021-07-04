package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// create a client obj
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999, // default value
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.Public chat mode")
	fmt.Println("2.Private chat mode")
	fmt.Println("3.Change the name mode")
	fmt.Println("0.Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>>>>Please input a valid number<<<<<<<<<<<")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("Please input your newName:\n")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Update name error:", err)
		return false
	}
	return true
}

// need a goroutine to recieve the message from server concurrently
func (client *Client) DealResponse() {
	// read data from client.conn, and copy it to os.Stdout and output, it supports multi-time execute
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PublicChat() {
	// suggest user input message
	fmt.Println("Please input the message you want to enter into pulbic chatting room: (exit end)")
	var sendMsg string
	fmt.Scanln(&sendMsg)
	// fmt.Println(sendMsg)

	for sendMsg != "exit" {
		if len(sendMsg) > 0 {
			chatMsg := sendMsg + "\n"
			_, err := client.conn.Write([]byte(chatMsg))
			if err != nil {
				fmt.Println("public chat error:", err)
				break
			}
		}
		fmt.Println("Please input the message you want to enter into pulbic chatting room: (exit end)")
		sendMsg = ""
		fmt.Scanln(&sendMsg)
		fmt.Println(&sendMsg)

	}
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"

	// send the msg to server, it will handle it and broadcast
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Select User Error:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteUser string
	var sendMsg string

	client.SelectUser()
	fmt.Println("Please input the username you want to send private msg: (exit end)")
	fmt.Scanln(&remoteUser)

	for remoteUser != "exit" {
		fmt.Println("Please input the chat msg you want to send to this user:(exit end)")
		fmt.Scanln(&sendMsg)

		for sendMsg != "exit" {
			if len(remoteUser) != 0 {
				chatMsg := "to|" + remoteUser + "|" + sendMsg + "\n"
				_, err := client.conn.Write([]byte(chatMsg))
				if err != nil {
					fmt.Println("private chat error:", err)
					break
				}

			}
			// sendMsg = ""
			// fmt.Println("Please input the chat msg you want to send to this user:(exit end)")
			// fmt.Scanln(&sendMsg)
		}
		sendMsg = ""
		fmt.Println("Please input the chat msg you want to send to this user:(exit end)")
		fmt.Scanln(&sendMsg)

		// client.SelectUser()
		// fmt.Println("Please input the username you want to send private msg: (exit end)")
		// fmt.Scanln(&remoteUser)

	}
	client.SelectUser()
	fmt.Println("Please input the username you want to send private msg: (exit end)")
	fmt.Scanln(&remoteUser)
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println("You are in public chat mode")
			client.PublicChat()
			break
		case 2:
			fmt.Println("You are in privatte chat mode")
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

var (
	serverIp   string
	serverPort int
)

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set default ip address to localhost")
	flag.IntVar(&serverPort, "port", 8888, "set default port  to 8888")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Connection failed")
		return
	}

	go client.DealResponse()

	fmt.Println("Connection successfully")

	client.Run()
}
