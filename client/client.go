package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"tcp-server/command"
)

type Client struct {
	Conn           net.Conn
	Name           string
	CurrentCommand chan<- command.Command
	CurrentChannel string
}

func ConnectToServer(address string, port string) net.Conn {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println("Error connecting to the server")
		os.Exit(1)
	}
	return conn
}

func SendMessage(msg string, conn net.Conn) {
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("error while sending message: ", err)
	}
}

func ReadServer(conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		fmt.Print(msg)
	}
}
