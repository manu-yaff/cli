package clientHandler

import (
	req "client-server/request"
	"client-server/utils"
	"fmt"
	"net"
)

//
func HandleInputCommand(cmd string, args []string, conn *net.Conn) {
	switch cmd {
	case "/send":
		HandleSendCommand(cmd, args, conn)
	default:
		request := &req.Request{
			CommandName: cmd,
			Args:        args,
		}
		err := utils.WriteRequest(conn, request)
		if err != nil {
			fmt.Printf("%s Error sending request \n", utils.CurrentTime())
		}
	}
}

//
func HandleSendCommand(cmd string, args []string, conn *net.Conn) {
	// check arguments are in the correct format
	if len(args) != 2 {
		fmt.Printf("%s %s\n", utils.CurrentTime(), "Usage: /send [file] [channel]")
		return
	}

	filename := args[0]

	// check file exist
	if !utils.FileExists(filename) {
		fmt.Printf("'%s' does not exist\n", filename)
	} else {
		// file exists

		// read file content
		content, err := utils.ReadFile(filename)
		if err != nil {
			fmt.Printf("%s Error reading file: %s\n", utils.CurrentTime(), err.Error())
		}

		// create request object
		request := &req.Request{
			CommandName: cmd,
			Args:        args,
			Content:     content,
		}

		// send the request
		err = utils.WriteRequest(conn, request)
		if err != nil {
			fmt.Printf("%s Error sending request: %s\n", utils.CurrentTime(), err.Error())
		}
	}
}
