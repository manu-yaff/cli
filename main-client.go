package main

import (
	"bufio"
	"os"
	"strings"
	c "tcp-server/client"
)

func main() {
	conn := c.ConnectToServer("localhost", "1234")
	go c.ReadServer(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		c.SendMessage(cmd, conn)
	}
}
