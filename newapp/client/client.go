package client

import (
	"bytes"
	req "client-server/request"
	res "client-server/response"
	"client-server/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	Conn           net.Conn
	Name           string
	CurrentRequest chan<- req.Request
	CurrentChannel string
	Date           string
	Channels       []string
}

// connects to the server and returns the conn object
func ConnectToServer(address string, port string) net.Conn {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("%s Error connecting to server: %ss\n", utils.CurrentTime(), err.Error())
		os.Exit(1)
	}
	return conn
}

// reads requests from server
func ReadServer(conn *net.Conn, c chan string) {
	for {
		var serverResponse res.Response
		err := utils.ReadResponse(conn, &serverResponse)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s %s\n", utils.CurrentTime(), "Server connection closed")
				os.Exit(1)
			}
		}

		// check is response has file
		if serverResponse.File != nil {
			// create dir
			if _, err := os.Stat(serverResponse.ClientIp); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir(serverResponse.ClientIp, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}

			// create file
			file, err := os.Create(serverResponse.ClientIp + "/" + serverResponse.File.Filename)
			if err != nil {
				fmt.Println(err)
			}

			// send bytes to file
			filesBytes := bytes.NewReader(serverResponse.File.Content)
			_, err = io.Copy(file, filesBytes)
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Printf("> %s\n", serverResponse.Message)
		c <- "@" + serverResponse.ClientName
	}
}
