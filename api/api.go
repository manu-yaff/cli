package api

import (
	"net/http"

	s "tcp-server/server"

	"github.com/gin-gonic/gin"
)

type ClientsResponse struct {
	Name     string
	Ip       string
	Date     string
	Channels []string
}

// retrieves all the clients currently connected to the tcp server
func GetClients(c *gin.Context, serverInstance *s.Server) {
	clients := serverInstance.Clients
	var serverResponseObject []ClientsResponse
	for _, c := range clients {
		item := &ClientsResponse{
			Name:     c.Name,
			Ip:       c.Conn.RemoteAddr().String(),
			Date:     c.Date,
			Channels: c.Channels,
		}
		serverResponseObject = append(serverResponseObject, *item)
	}
	c.JSON(http.StatusOK, serverResponseObject)
}

// define endpoints routes for the http server
func SetupRoutes(router *gin.Engine, server *s.Server) {
	router.GET("/clients", func(ctx *gin.Context) {
		GetClients(ctx, server)
	})
}
