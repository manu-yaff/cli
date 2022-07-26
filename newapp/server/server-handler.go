package server

import (
	fi "client-server/file"
	req "client-server/request"
	res "client-server/response"
	"client-server/utils"
	"fmt"
	"log"
	"os"
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
				Message: fmt.Sprintf("%s Channel does not exist", channelName),
			},
		)
	}
}

// handles logic for the send file command
func (server *Server) HandleSendFileCommand(request *req.Request) {
	filename := request.Args[0]
	channel := request.Args[1]
	fileContent := request.Content
	conn := request.Client

	// check channel exists
	if !server.ChannelExists(channel) {
		response := &res.Response{
			Message: fmt.Sprintf("'%s' channel does not exist", channel),
		}
		fmt.Printf("%s Channel does not exist\n", utils.CurrentTime())
		utils.WriteResponse(&conn, response)
	} else {
		// check if storage dir exists
		_, err := os.Stat("server-storage")
		if os.IsNotExist(err) {
			// create dir
			if err := os.Mkdir("server-storage", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
		// save file in storage
		err = utils.WriteFile(filename, fileContent)
		if err != nil {
			fmt.Println("Error writting file")
		} else {
			// check if that filename is in the channel
			delete(server.Channels[channel].Files, filename)

			// save file in channel
			newFile := &fi.File{
				Name: filename,
				Size: 0,
			}
			channelObj := server.Channels[channel]
			channelObj.Files = make(map[string]*fi.File)
			channelObj.Files[filename] = newFile

			// sender client
			sender := server.Clients[conn]

			// broadcast message
			response := &res.Response{
				Message: fmt.Sprintf("%s shared '%s' through '%s'", sender.Name, filename, channelObj.Name),
				File: &res.FileResponse{
					Filename: filename,
					Content:  fileContent,
				},
			}
			channelObj.Broadcast(response, conn)

			// write message for the client that sent file
			utils.WriteResponse(&conn, &res.Response{
				Message: fmt.Sprintf("You shared '%s' through '%s'", filename, channelObj.Name),
			})
		}
	}
}
