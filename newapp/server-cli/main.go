package main

import (
	ch "client-server/channel"
	cl "client-server/client"
	req "client-server/request"
	s "client-server/server"
	"net"
)

func main() {
	// create tcp server
	server := &s.Server{
		Clients:        make(map[net.Conn]*cl.Client),
		Channels:       make(map[string]*ch.Channel),
		CurrentRequest: make(chan req.Request),
	}

	server.Channels["dev"] = &ch.Channel{
		Name: "dev",
	}

	// listen for commands from the client
	go server.ReadClientRequest()

	// start server
	server.StartServer("1234")

	// defer closing server
	defer server.Listener.Close()

	// listen for connections
	server.ListenForConnections()

}
