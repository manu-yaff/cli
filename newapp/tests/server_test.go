package main

import (
	ch "client-server/channel"
	cl "client-server/client"
	req "client-server/request"
	s "client-server/server"
	"fmt"
	"net"
	"testing"
)

var testServer = &s.Server{
	Clients:        make(map[net.Conn]*cl.Client),
	Channels:       make(map[string]*ch.Channel),
	CurrentRequest: make(chan req.Request),
}

func clearMap(m map[net.Conn]*cl.Client) {
	for k := range m {
		delete(m, k)
	}
}

// test client is added to the server
func TestAddClient(t *testing.T) {
	var conn net.Conn
	expected := &cl.Client{
		Conn:     conn,
		Name:     "Anonymus",
		Channels: []string{},
	}
	tt := []struct {
		test    string
		payload net.Conn
		want    *cl.Client
	}{
		{
			"Adding a new client to the server",
			conn,
			expected,
		},
	}

	for _, tc := range tt {
		actualResult := testServer.AddClient(&tc.payload)
		if actualResult.Name != tc.want.Name {
			t.Errorf("got %s, expected %s", actualResult.Name, tc.want.Name)
		}
		if actualResult.Conn != tc.want.Conn {
			t.Errorf("got %s, expected %s", actualResult.Conn, tc.want.Conn)
		}
	}

	expectedClients := 1
	actualClients := len(testServer.Clients)

	if expectedClients != actualClients {
		t.Errorf("got %d, expected %d", actualClients, expectedClients)
	}
	clearMap(testServer.Clients)
}

// test when a client joins a channel, then it is added to the channel, the added in the client
func TestJoinExistingChannel(t *testing.T) {
	fmt.Println("Should add client to channel")

	// payload
	var conn net.Conn

	payloadChannel := &ch.Channel{
		Name:    "dev",
		Members: make(map[net.Conn]*cl.Client),
	}

	client := testServer.AddClient(&conn)

	testServer.Channels[payloadChannel.Name] = payloadChannel

	// call function
	testServer.AddToChannel(client, payloadChannel)
	testServer.AddChannelToClient(client, payloadChannel)

	// actual result
	actualMembers := len(testServer.Channels[payloadChannel.Name].Members)
	actualChannels := len(testServer.Clients[client.Conn].Channels)

	// expected result
	expectedMembers := 1
	expectedChannels := 1

	if actualMembers != expectedMembers {
		t.Errorf("got %d, expected %d", actualMembers, expectedMembers)
	}

	if actualChannels != expectedChannels {
		t.Errorf("got %d, expected %d", actualChannels, expectedChannels)
	}
	clearMap(testServer.Clients)
}

// test when a client joins a channel, then it is added to the channel, the added in the client
func TestJoinNonExistingChannel(t *testing.T) {
	fmt.Println("Should not add client")

	// payload
	var conn net.Conn

	payloadChannel := &ch.Channel{
		Name:    "frontend",
		Members: make(map[net.Conn]*cl.Client),
	}

	client := testServer.AddClient(&conn)

	// call function
	testServer.AddToChannel(client, payloadChannel)

	// check that channel doesnt exist
	if _, ok := testServer.Channels[payloadChannel.Name]; ok {
		t.Errorf("%s channel exists", payloadChannel.Name)
	}

	clearMap(testServer.Clients)
}

// test the client's name is changed
func TestChangeNameExistingClient(t *testing.T) {
	fmt.Println("Should change name of the client")

	var payloadClient net.Conn

	client := testServer.AddClient(&payloadClient)

	// change name
	expectedName := "jonh"
	testServer.SetClientName(expectedName, &client.Conn)
	actualName := testServer.Clients[client.Conn].Name

	if actualName != expectedName {
		t.Errorf("got %s, expected %s", actualName, expectedName)
	}
	clearMap(testServer.Clients)
}

func TestChangeNameNonExistingClient(t *testing.T) {
	fmt.Println("Should not add client")

	// payload
	conn, _ := net.Dial("tcp", "golang.org:80")
	client := cl.Client{
		Conn: conn,
		Name: "test",
	}

	// call function
	testServer.SetClientName("rick", &client.Conn)

	// check that channel doesnt exist
	if len(testServer.Clients) != 0 {
		t.Errorf("got %d, expected %d", len(testServer.Clients), 0)
	}
}
