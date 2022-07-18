// Package server contains the server struct and the functions related to it
package server

import (
	"fmt"
	"net"
	er "tcp-server/constants/errors"
)

// Creates a tcp server on localhost:1234 and returns a listener object
func CreateServer() net.Listener {
	s, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf(er.ERROR_SERVER_START + err.Error())
		return nil
	}

	fmt.Println("Server listening at localhost:1234")

	return s
}

// listen for incoming client connections
func ListenForConnections(server net.Listener) {
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf(er.ERROR_ACCEPT_CONN + err.Error())
			continue
		}
		fmt.Println("Welcome to the server: ", conn.RemoteAddr())
	}
}
