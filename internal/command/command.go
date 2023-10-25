package command

import (
	"strings"
	"time"
)

// Declaring time layout constant
const timelayout = "2006-Jan-02"

type command string

const (
	HelpCommand    command = "help"
	EchoCommand    command = "echo"
	ReserveCommand command = "reserv"
	GetCommand     command = "get"
	CancelCommand  command = "cancl"
)

type Command interface {
	Command() command
}

func (c command) Command() command {
	return c
}

type resource string

const (
	DeskOption    resource = "desk"
	ParkingOption resource = "parking"
)

type Resource interface {
	Resource() resource
}

func (r resource) Resource() resource {
	return r
}

type CommandCall struct {
	Command  Command
	Resource Resource
	Date     time.Time
}

func Parse(message string) CommandCall {
	returnValue := CommandCall{}
	words := strings.Split(strings.ToLower(message), " ")
	if len(words) > 0 {
		returnValue.Command = command(words[0])
		if len(words) > 1 {
			returnValue.Resource = resource(words[1])
			if len(words) > 2 {
				returnValue.Date, _ = time.Parse(timelayout, words[2])
			}
		}
	}
	return returnValue
}

func (c *CommandCall) IsValid() bool {
	returnValue := true

	return returnValue
}

/*
	help

	echo

	reserve desk date
	reserve parking date

	cancel desk date
	cancel parking date

	get desk
	get parking
	get desk date
	get parking date
*/
