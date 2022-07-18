// this is package
package main

import s "tcp-server/server"

func main() {
	// create server
	server := s.CreateServer()
	defer server.Close()

	// // listen for connections
	s.ListenForConnections(server)

}
