package clientHandler

import (
	req "client-server/request"
	"client-server/utils"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
)

// handles the logic according to the input command
func HandleInputCommand(cmd string, args []string, conn *net.Conn) {
	switch cmd {
	case "/send":
		HandleSendCommand(cmd, args, conn)
	case "/join":
		HandleJoinCommand(cmd, args, conn)
	case "/name":
		HandleNameCommand(cmd, args, conn)
	case "/list":
		HandleListChannels(cmd, args, conn)
	case "/create":
		HandleCreateCommand(cmd, args, conn)
	case "/leave":
		HandleLeaveCommand(cmd, args, conn)
	case "/open":
		HadnleOpenCommand(cmd, args, conn)
	case "/files":
		HandleFilesCommand()
	case "/help":
		HandleHelpCommand()

	default:
		fmt.Println("Error: command not supported. Run /help to see the available commands")
	}
}

// checks the arguments are correct. Checks if file exists and reads it to send it in the request
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

// checks the arguments are correct for the join command
func HandleJoinCommand(cmd string, args []string, conn *net.Conn) {
	if len(args) != 1 {
		fmt.Printf("%s\n", "Usage: /join [channel]")
		return
	}
	request := &req.Request{
		CommandName: cmd,
		Args:        args,
	}

	err := utils.WriteRequest(conn, request)
	if err != nil {
		fmt.Printf("%s Error sending request \n", utils.CurrentTime())
	}
}

// checks the arguments are correct for the name command
func HandleNameCommand(cmd string, args []string, conn *net.Conn) {
	if len(args) != 1 {
		fmt.Printf("%s \n", "Usage: /name [name]")
		return
	}
	request := &req.Request{
		CommandName: cmd,
		Args:        args,
	}

	err := utils.WriteRequest(conn, request)
	if err != nil {
		fmt.Printf("%s Error sending request \n", utils.CurrentTime())
	}
}

// sends the request since not arguments are required in list command
func HandleListChannels(cmd string, args []string, conn *net.Conn) {
	err := utils.WriteRequest(conn, &req.Request{
		CommandName: cmd,
		Args:        args,
	})

	if err != nil {
		fmt.Printf("%s Error sending request \n", utils.CurrentTime())
	}
}

// checks arguments are correct for the create command
func HandleCreateCommand(cmd string, args []string, conn *net.Conn) {
	if len(args) != 1 {
		fmt.Printf("%s \n", "Usage: /create [channel]")
		return
	}

	err := utils.WriteRequest(conn, &req.Request{
		CommandName: cmd,
		Args:        args,
	})

	if err != nil {
		fmt.Printf("%s Error sending request \n", utils.CurrentTime())
	}
}

// checks arguments are correct for the leave command
func HandleLeaveCommand(cmd string, args []string, conn *net.Conn) {
	if len(args) != 1 {
		fmt.Printf("%s \n", "Usage: /leave [channel]")
		return
	}

	err := utils.WriteRequest(conn, &req.Request{
		CommandName: cmd,
		Args:        args,
	})

	if err != nil {
		fmt.Printf("%s Error sending request \n", utils.CurrentTime())
	}
}

// opens the given file
func HadnleOpenCommand(cmd string, args []string, conn *net.Conn) {
	if len(args) != 1 {
		fmt.Printf("%s \n", "Usage: /leave [channel]")
		return
	}
	file := args[0]
	go func() {
		cmd := exec.Command("open", file)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
		}
	}()
}

// list the files in the current dir
func HandleFilesCommand() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Printf("Error reading dir: %s\n", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() && file.Name() != ".DS_Store" && file.Name() != "main.go" {
			fmt.Println("-", file.Name())
		}
	}
}

// displays information about the supported commands
func HandleHelpCommand() {
	fmt.Println(" /name [name] -> sets the name of your user, by default is Anonymus")
	fmt.Println(" /create [channel] -> creates a channel")
	fmt.Println(" /join [channel] -> joins a channel")
	fmt.Println(" /list -> displays all the channels available")
	fmt.Println(" /leave [channel] -> leaves a channel. Stops receiving information from there")
	fmt.Println(" /create [channel] -> creates a channel")
	fmt.Println(" /send [file] [channel] -> sends file through channel")
	fmt.Println(" /files -> list all files in the current user's directory")
	fmt.Println(" /open [file] -> opens file")
}
