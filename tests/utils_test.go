package main

import (
	"tcp-server/utils"
	"testing"
)

func TestCurrenTime(t *testing.T) {
	now := utils.CurrentTime()
	if len(now) != 19 {
		t.Errorf("Date should be in the format yyyy/mm/dd hh:mm:ss")
	}
}

func TestFormatUserInput(t *testing.T) {
	var test = struct {
		name string
		args []string
	}{
		"/name",
		[]string{"jony"},
	}
	userInput := "/name jony"
	cmd, args := utils.FormatUserInput(userInput)

	if cmd != test.name {
		t.Errorf("got %s, expected %s", cmd, test.name)
	}

	if len(args) != len(test.args) {
		t.Errorf("got %d, expected %d", len(args), len(test.args))
	}

}
