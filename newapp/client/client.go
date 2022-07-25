package client

import (
	er "client-server/constants/errors"
	notify "client-server/constants/notifications"
	req "client-server/request"
	res "client-server/response"
	"client-server/utils"
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
		fmt.Printf("%s %s: %s\n", utils.CurrentTime(), er.ERROR_CONNECTING_SEVER, err.Error())
		os.Exit(1)
	}
	return conn
}

// reads requests from server
func ReadServer(conn *net.Conn) {
	for {
		var serverResponse res.Response
		err := utils.ReadResponse(conn, &serverResponse)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s %s\n", utils.CurrentTime(), notify.SERVER_CONNECTION_CLOSED)
				os.Exit(1)
			}
		}
		fmt.Printf("> %s\n", serverResponse.Message)
	}
}
