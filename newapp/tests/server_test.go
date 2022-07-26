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

var conn net.Conn

// test client is added to the server
func TestAddClient(t *testing.T) {
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

}
