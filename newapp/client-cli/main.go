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
	c := make(chan string)
	go cl.ReadServer(&conn, c)
	defer conn.Close()

	for {
		// read from console
		fmt.Print("$ ")
		reader := bufio.NewReader(os.Stdin)
		clientInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		// format input
		cmd, args := utils.FormatUserInput(clientInput)

		// handle input command
		hd.HandleInputCommand(cmd, args, &conn)
		<-c
	}
}
