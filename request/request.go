package request

import "net"

type Request struct {
	CommandName string
	Args        []string
	Content     []byte
	Client      net.Conn
}
