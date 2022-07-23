// main package for the server
package main

import (
	"fmt"
	"net"
	"tcp-server/api"
	"tcp-server/channel"
	"tcp-server/client"
	s "tcp-server/server"
	"tcp-server/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// create tcp server

func main() {
	server := &s.Server{
		Clients:        make(map[net.Conn]*client.Client),
		Channels:       make(map[string]*channel.Channel),
		CurrentRequest: make(chan utils.Request),
	}

	// create http server for browser client
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.Default())
	api.SetupRoutes(router, server)

	go router.Run("localhost:8888")
	fmt.Println("Http server listening on port: 8888")

	//----------------------------------------------\\

	// listen for commands from client in channel
	go server.ReadClientRequest()

	// start server
	server.StartServer("1234")

	// defer closing server
	defer server.Listener.Close()

	// listen for connections
	server.ListenForConnections()

}
