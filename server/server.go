// Package server contains the server struct and the functions related to it
package server

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"tcp-server/channel"
	"tcp-server/client"
	er "tcp-server/constants/errors"
	notify "tcp-server/constants/notifications"
	f "tcp-server/file"
	"tcp-server/utils"
)

type Server struct {
	Listener       net.Listener
	Clients        map[net.Conn]*client.Client
	Channels       map[string]*channel.Channel
	CurrentRequest chan utils.Request
	Files          map[string]*f.File
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
func (server *Server) HandleNameCommand(request utils.Request) {
	if !utils.CheckArgs(1, request.Args) {
		// args are not in the correct format
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_NAME)
		// utils.WriteToConn(request.Client, notify.USAGE_NAME)
		response := &utils.Response{
			Message: notify.USAGE_NAME,
		}
		utils.WriteToConn(request.Client, response)
	} else {
		clientName := request.Args[0]
		client := request.Client
		server.SetClientName(clientName, client)
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.CLIENT_CHANGED_NAME)
		response := &utils.Response{
			Message: fmt.Sprintf("You changed your name to '%s'", clientName),
		}
		utils.WriteToConn(request.Client, response)
	}
}

// handles 'list' command
func (server *Server) HandleListCommand(request utils.Request) {
	client := request.Client
	if len(server.Channels) == 0 {
		response := &utils.Response{
			Message: "There are no channels. You can create one with /create [channelName]",
		}
		utils.WriteToConn(client, response)
	} else {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.CLIENT_LIST_CHANNELS)
		channels := server.GetChannels()
		response := &utils.Response{
			Message: channels,
		}
		utils.WriteToConn(client, response)
	}
}

// handles 'create' command
func (server *Server) HandleCreateCommand(request utils.Request) {
	client := request.Client
	if !utils.CheckArgs(1, request.Args) {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_CREATE)
		response := &utils.Response{
			Message: notify.USAGE_CREATE,
		}
		utils.WriteToConn(client, response)
	} else {
		channelName := request.Args[0]
		result := server.CreateChannel(channelName)
		if result {
			fmt.Printf("%s %s\n", utils.CurrentTime(), notify.CLIENT_CREATED_CHANNEL)
			response := &utils.Response{
				Message: fmt.Sprintf("%s channel created", channelName),
			}
			utils.WriteToConn(client, response)
		} else {
			fmt.Printf("%s\n", "Client tried to create a channel with a name already in used")
			response := &utils.Response{
				Message: fmt.Sprintf("%s channel already exists!", channelName),
			}
			utils.WriteToConn(client, response)
		}
	}
}

// send message to all clients in a channel
func (server *Server) Broadcast(response *utils.Response, channel string, currentClient net.Conn) {
	members := server.Channels[channel].Members
	for _, member := range members {
		if member.Conn != currentClient {
			response.ClientIp = member.Conn.RemoteAddr().String()
			response.ClientName = member.Name
			utils.WriteToConn(member.Conn, response)
		}
	}
}

// handles 'join' command
func (server *Server) HandleJoinCommand(request utils.Request) {
	client := request.Client
	// if !utils.CheckArgs(1, request.Args) {
	if !utils.CheckArgs(1, request.Args) {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_JOIN)
		response := &utils.Response{
			Message: notify.USAGE_JOIN,
		}
		utils.WriteToConn(client, response)
	} else {
		channelName := request.Args[0]
		result := server.JoinChannel(client, channelName)
		if result {
			fmt.Printf("%s %s %s\n", utils.CurrentTime(), notify.CLIENT_JOIN_CHANNEL, channelName)
			response := &utils.Response{
				Message: fmt.Sprintf("You joined '%s'", channelName),
			}
			// add channel to the client
			server.Clients[client].Channels = append(server.Clients[client].Channels, channelName)

			utils.WriteToConn(client, response)
			// broadcast notification to members in the channel
			clientName := server.Clients[client].Name
			broadcastResponse := &utils.Response{
				Message: fmt.Sprintf("%s joined %s", clientName, channelName),
			}
			server.Broadcast(broadcastResponse, channelName, client)
		} else {
			fmt.Printf("%s %s\n", utils.CurrentTime(), "user tried to join a channel that does no exist")
			response := &utils.Response{
				Message: fmt.Sprintf("'%s' channel does not exist", channelName),
			}
			utils.WriteToConn(client, response)
		}
	}
}

// handles 'send file' command
func (server *Server) HandleSendFileCommand(request utils.Request) {
	client := request.Client
	if !utils.CheckArgs(2, request.Args) {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_SEND)
		response := &utils.Response{
			Message: notify.USAGE_SEND,
		}
		utils.WriteToConn(client, response)
	} else {
		// get request information
		clientName := server.Clients[client].Name
		fileName := request.Args[0]
		channelName := request.Args[1]
		channel := server.Channels[channelName]

		// check channel exist
		if !server.ChannelExists(channelName) || !channel.IsMember(client) {
			// if not channels doesn't exist
			response := &utils.Response{
				Message: fmt.Sprintf("'%s' channel does not exist or you are not part of it", channelName),
			}
			utils.WriteToConn(client, response)
			fmt.Printf("%s %s\n", utils.CurrentTime(), "channel doesn't exist or client is not a member")
		} else {

			// create empty file

			err := os.WriteFile("server-storage/"+fileName, request.Content, 0644)
			if err != nil {
				fmt.Println(err)
			}

			fi, err := os.Stat("server-storage/" + fileName)
			if err != nil {
				fmt.Println(err)
			}

			fileSize := fi.Size()

			// file, err := os.Create("server-storage/" + request.Args[0])
			// if err != nil {
			// 	fmt.Println(err)
			// }

			// // save file to server storage
			// b := bytes.NewReader(request.Content)
			// _, err = io.Copy(file, b)

			// if err != nil {
			// 	fmt.Println(err)
			// }

			fileResponse := &utils.FileResponse{
				Filename: fileName,
				Content:  request.Content,
			}
			// check that channel has more then 2 members
			if len(server.Channels[channelName].Members) < 2 {
				response := &utils.Response{
					Message: "There are no members in the specified channel",
				}
				fmt.Printf("%s %s\n", utils.CurrentTime(), "there are no members")
				utils.WriteToConn(client, response)
			} else {

				message := fmt.Sprintf("%s shared '%s' through '%s' channel", clientName, fileName, channelName)
				channelMembersResponse := &utils.Response{
					Message:  message,
					File:     *fileResponse,
					ClientIp: client.RemoteAddr().String(),
				}
				server.Broadcast(channelMembersResponse, request.Args[1], request.Client)

				// if file does not exist, then add it
				if _, ok := server.Files[fileName]; !ok {
					newFile := &f.File{
						Name: fileName,
						Size: fileSize,
					}
					server.Files[fileName] = newFile
					server.Files[fileName].Channels = append(server.Files[fileName].Channels, channelName)
				}

				clientMessage := fmt.Sprintf("You shared '%s' through '%s' channel", fileName, channelName)
				clientResponse := &utils.Response{
					Message: clientMessage,
				}
				utils.WriteToConn(client, clientResponse)
				fmt.Printf("%s %s\n", utils.CurrentTime(), "client shared file")
			}
		}
	}
}

// adds a client to the specified channel, returns true if ok
func (server *Server) JoinChannel(client net.Conn, channelName string) bool {
	channelExists := server.ChannelExists(channelName)
	clientToAdd := server.Clients[client]

	if channelExists {
		server.Channels[channelName].Members[clientToAdd.Conn] = clientToAdd
		clientToAdd.CurrentChannel = channelName
		return true
	} else {
		return false
	}
}

// checks if channel exists
func (server *Server) ChannelExists(channelName string) bool {
	if _, ok := server.Channels[channelName]; ok {
		return true
	} else {
		return false
	}
}

// creates a new channel and return true if created successfully
func (server *Server) CreateChannel(channelName string) bool {
	if server.ChannelExists(channelName) {
		return false
	} else {
		newChannel := &channel.Channel{
			Name:    channelName,
			Members: make(map[net.Conn]*client.Client),
		}
		server.Channels[channelName] = newChannel
		return true
	}
}

// returns the existing channels in an array
func (server *Server) GetChannels() string {
	var channels []string
	channels = append(channels, "Channels:")
	for _, channel := range server.Channels {
		channels = append(channels, "- "+channel.Name)
	}
	return strings.Join(channels, "\n")
}

// reads from commands from a channel
func (server *Server) ReadClientRequest() {
	for {
		request := <-server.CurrentRequest
		cmd := request.CommandName
		switch cmd {
		case "/name":
			server.HandleNameCommand(request)
		case "/list":
			server.HandleListCommand(request)
		case "/create":
			server.HandleCreateCommand(request)
		case "/join":
			server.HandleJoinCommand(request)
		case "/send":
			server.HandleSendFileCommand(request)
		default:
			fmt.Printf("%s %s\n", utils.CurrentTime(), notify.INVALID_REQUEST)
			response := &utils.Response{
				Message: notify.INVALID_REQUEST,
			}
			utils.WriteToConn(request.Client, response)
		}
	}
}

// sets the name for a given client
func (server *Server) SetClientName(clientName string, client net.Conn) {
	server.Clients[client].Name = clientName
}

// add a client to the loby
func (server *Server) AddClientToLoby(conn *net.Conn) *client.Client {
	newClient := &client.Client{
		Conn:           *conn,
		Name:           "Anonymus",
		CurrentRequest: server.CurrentRequest,
		Date:           utils.CurrentTime(),
	}
	server.Clients[*conn] = newClient
	return server.Clients[*conn]
}

// reads messages from the command from the client
func (server *Server) HandleClientConnection(conn *net.Conn) {
	client := server.AddClientToLoby(conn)
	for {
		var clientInput utils.Request
		err := gob.NewDecoder(*conn).Decode(&clientInput)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s %s %s\n", utils.CurrentTime(), client.Conn.RemoteAddr(), notify.CLIENT_CONNECTION_CLOSED)
			}
			break
		}

		clientRequest := &utils.Request{
			CommandName: clientInput.CommandName,
			Args:        clientInput.Args,
			Content:     clientInput.Content,
			Client:      *conn,
		}

		client.CurrentRequest <- *clientRequest
	}
}
