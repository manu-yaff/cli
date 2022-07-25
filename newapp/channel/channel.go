package channel

import (
	cl "client-server/client"
	fi "client-server/file"
	"net"
)

type Channel struct {
	Name    string
	Members map[net.Conn]*cl.Client
	Files   []fi.File
}
