package main

import (
	"bufio"
	cl "client-server/client"
	hd "client-server/client-handler"
	"client-server/utils"
	"fmt"
	"os"
)

func main() {
	conn := cl.ConnectToServer("localhost", "1234")
	go cl.ReadServer(&conn)
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

		// handle input command
		hd.HandleInputCommand(cmd, args, &conn)
	}
}
