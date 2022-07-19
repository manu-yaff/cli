// main package for the server
package main

import (
	"net"
	"tcp-server/client"
	"tcp-server/command"
	s "tcp-server/server"
)

func main() {

	// create server
	server := &s.Server{
		Clients:        make(map[net.Conn]*client.Client),
		CurrentCommand: make(chan command.Command),
	}

	// listen for commands from client in channel
	go server.ReadCommandsFromClient()

	// start server
	server.StartServer("1234")

	// defer closing server
	defer server.Listener.Close()

	// listen for connections
	server.ListenForConnections()

}
