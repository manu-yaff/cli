package api

import (
	fi "client-server/file"
	s "client-server/server"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClientsResponse struct {
	Name     string
	Ip       string
	Date     string
	Channels []string
}

type ChannelFiles struct {
	Channel     string
	NumberFiles int
	Files       map[string]*fi.File
}

type FilesResponse struct {
	TotalFiles     int
	TotalData      int64
	FilesByChannel []ChannelFiles
}

// define endpoints routes for the http server
func SetupRoutes(router *gin.Engine, server *s.Server) {
	router.GET("/clients", func(ctx *gin.Context) {
		GetClients(ctx, server)
	})
	router.GET("/files", func(ctx *gin.Context) {
		GetFiles(ctx, server)
	})
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

// retrieves all the files send to the tcp server, the number of files and total data
func GetFiles(c *gin.Context, serverInstance *s.Server) {
	var channelFilesArray []ChannelFiles

	var totalSize int64 = 0
	totalFiles := 0

	dir, err := ioutil.ReadDir("server-storage")
	if err != nil {
		fmt.Println("Error reading dir")
	}

	for _, f := range dir {
		if f.Name() != ".DS_Store" {
			totalFiles++
			totalSize += f.Size()
		}
	}

	// files by channel
	channels := serverInstance.Channels
	for _, channel := range channels {
		item := &ChannelFiles{
			Channel:     channel.Name,
			NumberFiles: len(channel.Files),
			Files:       channel.Files,
		}
		channelFilesArray = append(channelFilesArray, *item)
	}
	c.JSON(http.StatusOK, FilesResponse{
		TotalFiles:     totalFiles,
		TotalData:      totalSize,
		FilesByChannel: channelFilesArray,
	})
}
