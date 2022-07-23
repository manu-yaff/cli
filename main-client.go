package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	c "tcp-server/client"
	"tcp-server/utils"
)

func main() {
	conn := c.ConnectToServer("localhost", "1234")
	go c.ReadServer(conn)
	defer conn.Close()

	for {
		// read from console
		reader := bufio.NewReader(os.Stdin)
		clientInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		// format input
		cmd, args := utils.FormatUserInput(clientInput)

		// create request object
		request := &utils.Request{
			CommandName: cmd,
			Args:        args,
		}

		if cmd == "/send" {
			if _, err := os.Stat("sample-files/" + args[0]); errors.Is(err, os.ErrNotExist) {
				fmt.Println("> Specified file does not exist")
			} else {
				content := utils.ReadFile("sample-files/" + args[0])
				request.Content = content
				c.SendRequest(request, conn)
			}
		} else {
			// send request
			c.SendRequest(request, conn)
		}

	}
}
