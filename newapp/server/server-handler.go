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

func (server *Server) HandleSendFileCommand(request *req.Request) {
	filename := request.Args[0]
	channel := request.Args[1]
	content := request.Content
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
		err = utils.WriteFile(filename, content)
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

			// broadcast message
			response := &res.Response{
				Message: "You will receive a file",
			}
			channelObj.Broadcast(response, conn)
		}
	}
}
