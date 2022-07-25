package clientHandler

import (
	er "client-server/constants/errors"
	notify "client-server/constants/notifications"
	req "client-server/request"
	"client-server/utils"
	"fmt"
	"net"
)

func HandleInputCommand(cmd string, args []string, conn *net.Conn) {
	switch cmd {
	case "/send":
		HandleSendCommand(cmd, args, conn)
	}
}

func HandleSendCommand(cmd string, args []string, conn *net.Conn) {
	// check arguments are in the correct format
	if len(args) != 2 {
		fmt.Printf("%s %s\n", utils.CurrentTime(), notify.USAGE_SEND)
		return
	}

	filename := args[0]

	// check file exist
	if !utils.FileExists(filename) {
		fmt.Printf("%s '%s %s\n", utils.CurrentTime(), filename, er.ERROR_FILE_NOT_EXISTS)
	} else {
		// file exists

		// read file content
		content, err := utils.ReadFile(filename)
		if err != nil {
			fmt.Printf("%s %s %s\n", utils.CurrentTime(), er.ERROR_READING_FILE, err.Error())
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
			fmt.Printf("%s %s %s\n", utils.CurrentTime(), er.ERROR_SENDING_REQUEST, err.Error())
		}
	}
}
