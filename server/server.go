// Package server contains the server struct and the functions related to it
package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"tcp-server/channel"
	"tcp-server/client"
	c "tcp-server/client"
	"tcp-server/command"
	er "tcp-server/constants/errors"
	notify "tcp-server/constants/notifications"
	"tcp-server/utils"
)

type Server struct {
	Listener       net.Listener
	Clients        map[net.Conn]*c.Client
	Channels       map[string]*channel.Channel
	CurrentCommand chan command.Command
	CurrentClient  *c.Client
}

// creates a tcp server on localhost:1234 and returns a listener object
func (server *Server) StartServer(port string) {
	l, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		fmt.Printf(er.ERROR_SERVER_START + err.Error())
		os.Exit(1)
	}

	fmt.Println("Server listening at localhost:1234")
	server.Listener = l
}

// listen for incoming client connections
func (server *Server) ListenForConnections() {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			fmt.Printf(er.ERROR_ACCEPT_CONN + err.Error())
			continue
		}
		fmt.Printf("%s Welcome to the server: %s\n", utils.CurrentTime(), conn.RemoteAddr())
		go server.HandleClientConnection(&conn)
	}
}

// handles the 'name' command
func (server *Server) HandleNameCommand(cmd command.Command) {
	if !cmd.CheckArgs(1) {
		// args are not in the correct format
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_NAME)
		utils.WriteToConn(cmd.Client, notify.USAGE_NAME)
	} else {
		server.SetClientName(cmd.Args[0], cmd.Client)
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.CLIENT_CHANGED_NAME)
		utils.WriteToConn(cmd.Client, "You changed your name to "+"'"+cmd.Args[0]+"'")
	}
}

// handles 'list' command
func (server *Server) HandleListCommand(cmd command.Command) {
	if len(server.Channels) == 0 {
		utils.WriteToConn(cmd.Client, "There are no channels. You can create one with /create [channelName]")
	} else {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.CLIENT_LIST_CHANNELS)
		channels := server.GetChannels()
		utils.WriteToConn(cmd.Client, channels)
	}
}

// handles 'create' command
func (server *Server) HandleCreateCommand(cmd command.Command) {
	if !cmd.CheckArgs(1) {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_CREATE)
		utils.WriteToConn(cmd.Client, notify.USAGE_CREATE)
	} else {
		channelName := cmd.Args[0]
		result := server.CreateChannel(channelName)
		if result {
			fmt.Printf("%s %s\n", utils.CurrentTime(), notify.CLIENT_CREATED_CHANNEL)
			utils.WriteToConn(cmd.Client, channelName+" channel created")
		} else {
			fmt.Printf("%s\n", "Client tried to create a channel with a name already in used")
			utils.WriteToConn(cmd.Client, channelName+" channel already exists!")
		}
	}
}

// send message to all clients in a channel
func (server *Server) Broadcast(message string, channel string, currentClient net.Conn) {
	members := server.Channels[channel].Members
	clientObj := server.Clients[currentClient]
	for _, member := range members {
		if member.Conn != clientObj.Conn {
			// fmt.Println(member.Name)
			member.Conn.Write([]byte("> " + clientObj.Name + " joined " + channel + "\n"))
		}
	}
}

// handles 'join' command
func (server *Server) HandleJoinCommand(cmd command.Command) {
	if !cmd.CheckArgs(1) {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_JOIN)
		utils.WriteToConn(cmd.Client, notify.USAGE_JOIN)
	} else {
		channelName := cmd.Args[0]
		result := server.JoinChannel(cmd.Client, channelName)
		if result {
			fmt.Printf("%s %s %s\n", utils.CurrentTime(), notify.CLIENT_JOIN_CHANNEL, channelName)
			utils.WriteToConn(cmd.Client, fmt.Sprintf("You joined '%s'", channelName))
			// broadcast notification to members in the channel
			server.Broadcast("", channelName, cmd.Client)
		} else {
			fmt.Printf("%s %s", utils.CurrentTime(), "User tried to join a channel that does no exist")
			utils.WriteToConn(cmd.Client, fmt.Sprintf("'%s' channel does not exist", channelName))
		}
	}
}

// adds a client to the specified channel, returns true if ok
func (server *Server) JoinChannel(client net.Conn, channelName string) bool {
	channelExists := server.channelExists(channelName)
	clientToAdd := server.Clients[client]

	if channelExists {
		server.Channels[channelName].Members[clientToAdd.Conn] = *clientToAdd
		clientToAdd.CurrentChannel = channelName
		return true
	} else {
		return false
	}
}

// checks if channel exists
func (server *Server) channelExists(channelName string) bool {
	if _, ok := server.Channels[channelName]; ok {
		return true
	} else {
		return false
	}
}

// creates a new channel and return true if created successfully
func (server *Server) CreateChannel(channelName string) bool {
	if server.channelExists(channelName) {
		return false
	} else {
		newChannel := &channel.Channel{
			Name:    channelName,
			Members: make(map[net.Conn]client.Client),
		}
		server.Channels[channelName] = newChannel
		return true
	}
}

// returns the existing channels in an array
func (server *Server) GetChannels() string {
	var channels []string
	for _, channel := range server.Channels {
		fmt.Println(channel.Name)
		for _, member := range channel.Members {
			fmt.Println("\t", member.Name)
		}
	}
	for _, channel := range server.Channels {
		channels = append(channels, channel.Name)
	}
	return strings.Join(channels, ", ")
}

// reads from commands from a channel
func (server *Server) ReadCommandsFromClient() {
	for {
		cmd := <-server.CurrentCommand
		switch cmd.Name {
		case "/name":
			server.HandleNameCommand(cmd)
		case "/list":
			server.HandleListCommand(cmd)
		case "/create":
			server.HandleCreateCommand(cmd)
		case "/join":
			server.HandleJoinCommand(cmd)
		default:
			fmt.Printf("%s %s\n", utils.CurrentTime(), notify.INVALID_REQUEST)
			utils.WriteToConn(cmd.Client, notify.INVALID_REQUEST)
		}
	}
}

// sets the name for a given client
func (server *Server) SetClientName(clientName string, client net.Conn) {
	server.Clients[client].Name = clientName
}

// add a client to the loby
func (server *Server) AddClientToLoby(conn *net.Conn) *c.Client {
	newClient := &c.Client{
		Conn:           *conn,
		Name:           "Anonymus",
		CurrentCommand: server.CurrentCommand,
	}
	server.Clients[*conn] = newClient
	return server.Clients[*conn]
}

// reads messages from the command from the client
func (server *Server) HandleClientConnection(conn *net.Conn) {
	client := server.AddClientToLoby(conn)
	for {
		clientRequest, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s %s %s\n", utils.CurrentTime(), client.Conn.RemoteAddr(), notify.CLIENT_CONNECTION_CLOSED)
			}
			break
		}

		cmdName, args := utils.FormatUserInput(clientRequest)

		cmd := &command.Command{
			Name:   cmdName,
			Client: *conn,
			Args:   args,
		}
		client.CurrentCommand <- *cmd
	}
}
