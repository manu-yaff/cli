// main package for the server
package main

import (
	"net"
	"tcp-server/channel"
	"tcp-server/client"
	s "tcp-server/server"
	"tcp-server/utils"
)

func main() {

	// create server
	server := &s.Server{
		Clients:        make(map[net.Conn]*client.Client),
		Channels:       make(map[string]*channel.Channel),
		CurrentRequest: make(chan utils.Request),
	}

	// listen for commands from client in channel
	go server.ReadClientRequest()

	// start server
	server.StartServer("1234")

	// defer closing server
	defer server.Listener.Close()

	// listen for connections
	server.ListenForConnections()

}
