// main package for the server
package main

import (
	"fmt"
	"net"
	"net/http"
	"tcp-server/channel"
	"tcp-server/client"
	s "tcp-server/server"
	"tcp-server/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ClientResponse struct {
	Name     string
	Ip       string
	Date     string
	Channels []string
}

var server = &s.Server{
	Clients:        make(map[net.Conn]*client.Client),
	Channels:       make(map[string]*channel.Channel),
	CurrentRequest: make(chan utils.Request),
}

func getClients(c *gin.Context) {
	clients := server.Clients
	var serverResponseObject []ClientResponse
	for _, c := range clients {
		item := &ClientResponse{
			Name:     c.Name,
			Ip:       c.Conn.RemoteAddr().String(),
			Date:     c.Date,
			Channels: c.Channels,
		}

		serverResponseObject = append(serverResponseObject, *item)
	}
	c.JSON(http.StatusOK, serverResponseObject)
}

type Grahph struct {
	Channel      string
	FilesNumbers int
}

func getFiles(c *gin.Context) {
	a := []Grahph{
		{
			Channel:      "dev",
			FilesNumbers: 9,
		},
		{
			Channel:      "frontend",
			FilesNumbers: 4,
		},
		{
			Channel:      "general",
			FilesNumbers: 7,
		},
	}
	// s := &Grahph{
	// }
	c.JSON(http.StatusOK, a)
}

func main() {
	// http server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.Default())
	router.GET("/clients", getClients)
	router.GET("/files", getFiles)
	go router.Run("localhost:8888")
	fmt.Println("http server running on port :8888")

	// create server

	// listen for commands from client in channel
	go server.ReadClientRequest()

	// start server
	server.StartServer("1234")

	// defer closing server
	defer server.Listener.Close()

	// listen for connections
	server.ListenForConnections()

}
