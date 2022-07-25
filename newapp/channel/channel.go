package channel

import (
	cl "client-server/client"
	fi "client-server/file"
	res "client-server/response"
	utils "client-server/utils"
	"net"
)

type Channel struct {
	Name    string
	Members map[net.Conn]*cl.Client
	Files   map[string]*fi.File
}

func (channel *Channel) Broadcast(response *res.Response, currentClient net.Conn) {
	members := channel.Members
	for _, member := range members {
		if member.Conn != currentClient {
			response.ClientIp = member.Conn.RemoteAddr().String()
			response.ClientName = member.Name
			utils.WriteResponse(&member.Conn, response)
		}
	}
}
