package utils

import (
	"net"
	"strings"
	"time"
)

// returns the current date in the format: yyyy/mm/dd hh:mm:ss
func CurrentTime() string {
	t := time.Now()
	date := t.Format("2006/01/02 15:04:05")
	return date
}

// takes the user input and returns two strings, the command and the args
func FormatUserInput(userInput string) (string, []string) {
	userInput = strings.Trim(userInput, "\n")
	slice := strings.Split(userInput, " ")
	cmd := slice[0]
	args := slice[1:]
	return cmd, args
}

// wirtes message to the given conn
func WriteToConn(conn net.Conn, message string) {
	conn.Write([]byte("> " + message + "\n"))
}
