package utils

import (
	req "client-server/request"
	res "client-server/response"
	"encoding/gob"
	"errors"
	"net"
	"os"
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

// writes the request to the given conn. Returns error if present
func WriteRequest(conn *net.Conn, request *req.Request) error {
	err := gob.NewEncoder(*conn).Encode(request)
	if err != nil {
		return err
	}
	return nil
}

// reads the request from the given conn. Returns error if present
func ReadRequest(conn *net.Conn, dest *req.Request) error {
	err := gob.NewDecoder(*conn).Decode(&dest)
	if err != nil {
		return err
	}
	return nil
}

// writes the response to the given conn. Returns error if present
func WriteResponse(conn *net.Conn, response *res.Response) error {
	err := gob.NewEncoder(*conn).Encode(response)
	if err != nil {
		return err
	}
	return nil
}

// reads the response from the given conn. Returns error if present
func ReadResponse(conn *net.Conn, dest *res.Response) error {
	err := gob.NewDecoder(*conn).Decode(&dest)
	if err != nil {
		return err
	}
	return nil
}

// checks if file exists, returns false if not
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// reads file and returns its content
func ReadFile(fileName string) ([]byte, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

// writes a file
func WriteFile(filepath string, content []byte) error {
	err := os.WriteFile("server-storage/"+filepath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
