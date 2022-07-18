package server

import (
	"fmt"
	"net"
	er "tcp-server/constants/errors"
)

func ListenForConnections(server net.Listener) {
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf(er.ERROR_ACCEPT_CONN)
			continue
		}
		fmt.Println("Welcome to the server: ", conn.RemoteAddr())
	}
}

func CreateServer() net.Listener {
	s, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf(er.ERROR_SERVER_START)
		return nil
	}

	fmt.Println("Server listening at localhost:1234")

	return s
}
