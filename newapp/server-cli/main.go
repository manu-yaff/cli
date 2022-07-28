package main

import (
	api "client-server/api"
	ch "client-server/channel"
	cl "client-server/client"
	req "client-server/request"
	s "client-server/server"
	"fmt"
	"net"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// create tcp server
	server := &s.Server{
		Clients:        make(map[net.Conn]*cl.Client),
		Channels:       make(map[string]*ch.Channel),
		CurrentRequest: make(chan req.Request),
	}

	// create http server for browser client
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.Default())
	api.SetupRoutes(router, server)

	go router.Run("localhost:8888")
	fmt.Println("Http server listening on port: 8888")

	// listen for commands from the client
	go server.ReadClientRequest()

	// start server
	server.StartServer("1234")

	// defer closing server
	defer server.Listener.Close()

	// listen for connections
	server.ListenForConnections()

}
