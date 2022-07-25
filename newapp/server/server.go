package server

import (
	ch "client-server/channel"
	cl "client-server/client"
	er "client-server/constants/errors"
	notify "client-server/constants/notifications"
	req "client-server/request"
	res "client-server/response"
	"client-server/utils"
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	Listener       net.Listener
	Clients        map[net.Conn]*cl.Client
	Channels       map[string]*ch.Channel
	CurrentRequest chan req.Request
}

// creates a tcp server on localhost:1234 and returns a listener object
func (server *Server) StartServer(port string) {
	l, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		fmt.Printf("%s %s: %s\n", utils.CurrentTime(), er.ERROR_SERVER_START, err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s %s\n", utils.CurrentTime(), notify.SERVER_LISTENING)
	server.Listener = l
}

// listen for incoming client connections
func (server *Server) ListenForConnections() {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			fmt.Printf(er.ERROR_ACCEPT_CONN + err.Error())
			continue
		}
		fmt.Printf("%s %s %s\n", utils.CurrentTime(), notify.WELCOME, conn.RemoteAddr())
		go server.HandleClientConnection(&conn)
	}
}

// reads from commands from a channel
func (server *Server) ReadClientRequest() {
	for {
		request := <-server.CurrentRequest
		cmd := request.CommandName

		response := &res.Response{
			Message: "test",
		}
		utils.WriteResponse(&request.Client, response)

		switch cmd {
		case "/name":
			// server.HandleNameCommand(request)
		case "/list":
			// server.HandleListCommand(request)
		case "/create":
			// server.HandleCreateCommand(request)
		case "/join":
			// server.HandleJoinCommand(request)
		case "/send":
			// server.HandleSendFileCommand(request)
		default:
			// fmt.Printf("%s %s\n", utils.CurrentTime(), notify.INVALID_REQUEST)
			// response := &utils.Response{
			// Message: notify.INVALID_REQUEST,
			// }
			// utils.WriteToConn(request.Client, response)
		}
	}
}

func (server *Server) HandleClientConnection(conn *net.Conn) {
	// add client to server
	client := server.AddClient(conn)

	// read requests from client
	for {
		var clientInput req.Request
		err := utils.ReadRequest(conn, &clientInput)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s %s %s\n", utils.CurrentTime(), client.Conn.RemoteAddr(), notify.CLIENT_CONNECTION_CLOSED)
			} else {
				fmt.Printf("%s %s %s\n", utils.CurrentTime(), er.ERROR_READING_REQUEST, err.Error())
			}
			break
		}

		clientRequest := &req.Request{
			CommandName: clientInput.CommandName,
			Args:        clientInput.Args,
			Content:     clientInput.Content,
			Client:      *conn,
		}

		client.CurrentRequest <- *clientRequest
	}
}

func (server *Server) AddClient(conn *net.Conn) *cl.Client {
	newClient := &cl.Client{
		Conn:           *conn,
		Name:           "Anonymus",
		CurrentRequest: server.CurrentRequest,
		Date:           utils.CurrentTime(),
	}
	server.Clients[*conn] = newClient
	return server.Clients[*conn]
}
