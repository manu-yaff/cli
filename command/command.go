package command

import "net"

type Command struct {
	Name   string
	Client net.Conn
	Args   []string
}

func (c *Command) CheckArgs(numberArgs int) bool {
	if len(c.Args) < numberArgs {
		return false
	} else {
		return true
	}
}
