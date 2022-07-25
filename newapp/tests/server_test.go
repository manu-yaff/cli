package main

import (
	ch "client-server/channel"
	cl "client-server/client"
	req "client-server/request"
	s "client-server/server"
	"net"
	"testing"
)

var testServer = &s.Server{
	Clients:        make(map[net.Conn]*cl.Client),
	Channels:       make(map[string]*ch.Channel),
	CurrentRequest: make(chan req.Request),
}

var conn net.Conn

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
