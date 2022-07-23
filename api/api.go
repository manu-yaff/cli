package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	f "tcp-server/file"
	s "tcp-server/server"

	"github.com/gin-gonic/gin"
)

type ClientsResponse struct {
	Name     string
	Ip       string
	Date     string
	Channels []string
}

type FilesResponse struct {
	TotalFiles int
	TotalData  int64
	Files      []f.File
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
	var files []f.File
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

	for _, file := range serverInstance.Files {
		item := &f.File{
			Name:     file.Name,
			Size:     file.Size,
			Channels: file.Channels,
		}
		files = append(files, *item)
	}

	response := &FilesResponse{
		TotalFiles: totalFiles,
		TotalData:  totalSize,
		Files:      files,
	}

	c.JSON(http.StatusOK, response)
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
