package utils

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Request struct {
	CommandName string
	Args        []string
	Content     []byte
	Client      net.Conn
}

type FileResponse struct {
	Filename string
	Content  []byte
}

type Response struct {
	Message    string
	File       FileResponse
	ClientName string
	ClientIp   string
}

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
func WriteToConn(conn net.Conn, response *Response) {
	err := gob.NewEncoder(conn).Encode(response)
	if err != nil {
		fmt.Println("Error sending response: ", err)
	}
}

// checks that the arguments commands are the required
func CheckArgs(expectedArgs int, actualArgs []string) bool {
	if len(actualArgs) < expectedArgs {
		return false
	} else {
		return true
	}
}

// reads file and returns its content and extension
func ReadFile(fileName string) []byte {
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	return content
}

// gets the extension file
func GetFileExtension(fileName string) string {
	fileExtension := strings.Split(fileName, ".")[1]
	return "." + fileExtension
}
