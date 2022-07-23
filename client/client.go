package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"tcp-server/utils"
)

type Client struct {
	Conn           net.Conn
	Name           string
	CurrentRequest chan<- utils.Request
	CurrentChannel string
}

func ConnectToServer(address string, port string) net.Conn {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println("Error connecting to the server")
		os.Exit(1)
	}
	return conn
}

func SendRequest(request *utils.Request, conn net.Conn) {
	err := gob.NewEncoder(conn).Encode(request)
	if err != nil {
		fmt.Println("Error sending file: ", err)
	}
}

func ReadServer(conn net.Conn) {
	for {
		var serverResponse utils.Response
		err := gob.NewDecoder(conn).Decode(&serverResponse)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed the connection")
				os.Exit(1)
			}
			fmt.Println(err)
		}
		if serverResponse.File.Content != nil {
			path := "storage-" + serverResponse.ClientName + "-" + serverResponse.ClientIp
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err := os.Mkdir(path, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}

			file, err := os.Create(path + "/" + serverResponse.File.Filename)
			if err != nil {
				fmt.Println(err)
			}

			b := bytes.NewReader(serverResponse.File.Content)
			_, err = io.Copy(file, b)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("> %s\n", serverResponse.Message)
		} else {
			fmt.Printf("> %s\n", serverResponse.Message)
		}
	}
}
