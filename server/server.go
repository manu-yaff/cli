package server

import (
	ch "client-server/channel"
	cl "client-server/client"
	fi "client-server/file"
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
			fmt.Printf("Error accepting connection: %s\n", err.Error())
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
			server.HandleNameCommand(&request)
		case "/list":
			server.HandleListCommand(&request)
		case "/create":
			server.HandleCreateCommand(&request)
		case "/join":
			server.HandleJoinCommand(&request)
		case "/send":
			server.HandleSendFileCommand(&request)
		case "/leave":
			server.HandleLeaveCommand(&request)
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

// -------------------------------- Client functions -------------------------------- \\

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

// changes the name of a given client
func (server *Server) SetClientName(clientName string, client *net.Conn) {
	if _, ok := server.Clients[*client]; ok {
		server.Clients[*client].Name = clientName
	}
}

// removes channel from client's array
func (server *Server) RemoveChannelFromClient(client *cl.Client, channel string) {
	channelsArray := server.Clients[client.Conn].Channels
	var index = -1

	for _, val := range channelsArray {
		index++
		if val == channel {
			break
		}
	}

	if len(channelsArray) <= 1 {
		server.Clients[client.Conn].Channels = make([]string, 0)
		return
	}

	result := append(channelsArray[0:index], channelsArray[index+1:]...)
	server.Clients[client.Conn].Channels = result
}

// -------------------------------- Channel functions -------------------------------- \\

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

// returns the existing channels in the server
func (server *Server) GetChannels() []string {
	var channels []string
	for _, channel := range server.Channels {
		channels = append(channels, "- "+channel.Name)
	}
	return channels
}

// creates channel only when the channel does not exist already
func (server *Server) CreateChannel(channelName string) bool {
	if server.ChannelExists(channelName) {
		return false
	}

	newChannel := &ch.Channel{
		Name:    channelName,
		Members: make(map[net.Conn]*cl.Client),
		Files:   make(map[string]*fi.File),
	}

	server.Channels[channelName] = newChannel
	return true
}
