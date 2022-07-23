package channel

import (
	"net"
	"tcp-server/client"
)

type Channel struct {
	Name    string
	Members map[net.Conn]*client.Client
}

func (channel *Channel) IsMember(client net.Conn) bool {
	if _, ok := channel.Members[client]; ok {
		return true
	} else {
		return false
	}
}
