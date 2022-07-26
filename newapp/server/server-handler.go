package server

import (
	fi "client-server/file"
	req "client-server/request"
	res "client-server/response"
	"client-server/utils"
	"fmt"
	"os"
	"strings"
)

// handles logic for the join command
func (server *Server) HandleJoinCommand(request *req.Request) {
	channelName := request.Args[0]
	channel := server.Channels[channelName]
	client := request.Client
	// check channel exists
	if server.ChannelExists(channelName) {
		clientToAdd := server.Clients[client]

		// check if client is already in channel
		if channel.HasMember(clientToAdd) {
			fmt.Printf("%s Client is already in channel\n", utils.CurrentTime())
			utils.WriteResponse(&client,
				&res.Response{
					Message: fmt.Sprintf("You are already in '%s'", channelName),
				},
			)
		} else {
			// add client to channel
			server.AddToChannel(clientToAdd, channel)

			// add channel to client
			server.AddChannelToClient(clientToAdd, channel)

			// broadcast message for members
			channel.Broadcast(&res.Response{
				Message: fmt.Sprintf("%s joined %s", clientToAdd.Name, channelName),
			}, client)

			// print in server
			fmt.Printf("%s Client joined '%s' channel\n", utils.CurrentTime(), channelName)

			// response for sender
			utils.WriteResponse(&client,
				&res.Response{
					Message: fmt.Sprintf("You joined '%s'", channelName),
				},
			)
		}
	} else {
		utils.WriteResponse(&client,
			&res.Response{
				Message: fmt.Sprintf("%s channel does not exist", channelName),
			},
		)
	}
}

// handles logic for the send file command
func (server *Server) HandleSendFileCommand(request *req.Request) {
	filename := request.Args[0]
	channelName := request.Args[1]
	fileContent := request.Content
	conn := request.Client
	channel := server.Channels[channelName]

	// check channel exists
	if !server.ChannelExists(channelName) {
		response := &res.Response{
			Message: fmt.Sprintf("'%s' channel does not exist", channelName),
		}
		fmt.Printf("%s Channel does not exist\n", utils.CurrentTime())
		utils.WriteResponse(&conn, response)
		return
	}

	// check client is part of the channel
	if !server.Channels[channelName].HasMember(server.Clients[conn]) {
		fmt.Printf("%s User is not member of channel\n", utils.CurrentTime())
		utils.WriteResponse(&conn, &res.Response{
			Message: fmt.Sprintf("You are not member in '%s'", channelName),
		})
		return
	}

	// if client is part of channel, proceed

	// check if storage dir exists
	_, err := os.Stat("server-storage")
	if os.IsNotExist(err) {
		// create dir
		if err := os.Mkdir("server-storage", os.ModePerm); err != nil {
			fmt.Printf("%s Error creating folder for server storage", utils.CurrentTime())
			utils.WriteResponse(&conn, &res.Response{
				Message: "Server error",
			})
			return
		}
	}

	// save file in storage
	err = utils.WriteFile(filename, fileContent)
	if err != nil {
		fmt.Printf("%s Error writing file", utils.CurrentTime())
		utils.WriteResponse(&conn, &res.Response{
			Message: "Server error",
		})
		return

	} else {
		// check if that filename is in the channel
		delete(channel.Files, filename)

		// save file in channel
		newFile := &fi.File{
			Name: filename,
			Size: 0,
		}
		channel.Files[filename] = newFile

		// sender client
		sender := server.Clients[conn]

		// broadcast message
		response := &res.Response{
			Message: fmt.Sprintf("%s shared '%s' through '%s'", sender.Name, filename, channelName),
			File: &res.FileResponse{
				Filename: filename,
				Content:  fileContent,
			},
		}
		channel.Broadcast(response, conn)

		// write message for the client that sent file
		utils.WriteResponse(&conn, &res.Response{
			Message: fmt.Sprintf("You shared '%s' through '%s'", filename, channelName),
		})
		fmt.Printf("%s Client shared a file\n", utils.CurrentTime())
	}

}

// handles logic for name command
func (server *Server) HandleNameCommand(request *req.Request) {
	// change client name
	clientName := request.Args[0]
	conn := request.Client
	server.Clients[conn].Name = clientName

	// send response
	utils.WriteResponse(&conn, &res.Response{
		Message: fmt.Sprintf("You changed your name to '%s'", clientName),
	})

	// print in server
	fmt.Printf("%s Client change their name\n", utils.CurrentTime())
}

// list all commands in the server
func (server *Server) HandleListCommand(request *req.Request) {
	conn := request.Client
	channels := server.GetChannels()

	fmt.Printf("%s Client listed channels\n", utils.CurrentTime())
	formattedChannels := fmt.Sprintf("Channels\n%s", strings.Join(channels, "\n"))
	utils.WriteResponse(&conn, &res.Response{
		Message: formattedChannels,
	})
}

// creates a new channel
func (server *Server) HandleCreateCommand(request *req.Request) {
	channelName := request.Args[0]
	conn := request.Client

	// check if channel exists
	if server.ChannelExists(channelName) {
		utils.WriteResponse(&conn, &res.Response{
			Message: fmt.Sprintf("'%s' already exists", channelName),
		})
		fmt.Printf("%s Client tried to create channel that already exists\n", utils.CurrentTime())
		return
	}

	// client does not exist, then create it
	if server.CreateChannel(channelName) {
		utils.WriteResponse(&conn, &res.Response{
			Message: fmt.Sprintf("'%s' channel created", channelName),
		})
		fmt.Printf("%s Client created channel\n", utils.CurrentTime())
	} else {
		utils.WriteResponse(&conn, &res.Response{
			Message: fmt.Sprintf("Error creating '%s' channel", channelName),
		})
		fmt.Printf("%s Error creating '%s' channel \n", utils.CurrentTime(), channelName)
	}
}
