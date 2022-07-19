package channel

import (
	"net"
	"tcp-server/client"
)

type Channel struct {
	Name    string
	Members map[net.Conn]client.Client
}
