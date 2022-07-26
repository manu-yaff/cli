package server

import (
	ch "client-server/channel"
	cl "client-server/client"
	req "client-server/request"
	"client-server/utils"
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	Listener       net.Listener
	Clients        map[net.Conn]*cl.Client
	Channels       map[string]*ch.Channel
	CurrentRequest chan req.Request
}

// creates a tcp server on localhost:1234 and returns a listener object
func (server *Server) StartServer(port string) {
	l, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		fmt.Printf("%s Error starting server: %s\n", utils.CurrentTime(), err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s %s\n", utils.CurrentTime(), "Server listening at localhost:1234")
	server.Listener = l
}

// listen for incoming client connections
func (server *Server) ListenForConnections() {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			fmt.Printf("%s Error accepting connection %s\n", err.Error())
			continue
		}
		fmt.Printf("%s Welcome to the server %s\n", utils.CurrentTime(), conn.RemoteAddr())
		go server.HandleClientConnection(&conn)
	}
}

// reads from commands from a channel
func (server *Server) ReadClientRequest() {
	for {
		request := <-server.CurrentRequest
		cmd := request.CommandName

		switch cmd {
		case "/name":
			// server.HandleNameCommand(request)
		case "/list":
			// server.HandleListCommand(request)
		case "/create":
			// server.HandleCreateCommand(request)
		case "/join":
			server.HandleJoinCommand(&request)
		case "/send":
			server.HandleSendFileCommand(&request)
		default:
			// fmt.Printf("%s %s\n", utils.CurrentTime(), notify.INVALID_REQUEST)
			// response := &utils.Response{
			// Message: notify.INVALID_REQUEST,
			// }
			// utils.WriteToConn(request.Client, response)
		}
	}
}

// listens for the incoming request and sends the message to the channel
func (server *Server) HandleClientConnection(conn *net.Conn) {
	// add client to server
	client := server.AddClient(conn)

	// read requests from client
	for {
		var clientInput req.Request
		err := utils.ReadRequest(conn, &clientInput)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s %s Client closed connection\n", utils.CurrentTime(), client.Conn.RemoteAddr())
				(*conn).Close()
				server.RemoveClient(conn)
			} else {
				fmt.Printf("%s Error reading request: %s\n", utils.CurrentTime(), err.Error())
			}
			break
		}

		clientRequest := &req.Request{
			CommandName: clientInput.CommandName,
			Args:        clientInput.Args,
			Content:     clientInput.Content,
			Client:      *conn,
		}

		client.CurrentRequest <- *clientRequest
	}
}

// adds a client to the server instance
func (server *Server) AddClient(conn *net.Conn) *cl.Client {
	newClient := &cl.Client{
		Conn:           *conn,
		Name:           "Anonymus",
		CurrentRequest: server.CurrentRequest,
		Date:           utils.CurrentTime(),
	}
	server.Clients[*conn] = newClient
	return server.Clients[*conn]
}

// removes a client from a server instance
func (server *Server) RemoveClient(conn *net.Conn) {
	delete(server.Clients, *conn)
}

// checks if channel exists
func (server *Server) ChannelExists(channelName string) bool {
	if _, ok := server.Channels[channelName]; ok {
		return true
	}
	return false
}

// adds a client to the specified channel, returns true if ok
func (server *Server) AddToChannel(client *cl.Client, channel *ch.Channel) {
	channel.Members[client.Conn] = client
}

// adds a channel to client
func (server *Server) AddChannelToClient(client *cl.Client, channel *ch.Channel) {
	client.Channels = append(client.Channels, channel.Name)
}
