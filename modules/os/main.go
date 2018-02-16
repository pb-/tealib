package os

import (
	"fmt"
	"os"

	"github.com/pb-/tealib"
)

// Exit terminates the program with an optional message.
type Exit struct {
	Status         int
	GoodbyeMessage string
}

// Main initializes the module.
func Main(messages chan<- tealib.Message) chan<- tealib.Command {
	commands := make(chan tealib.Command)

	go func() {
		for {
			switch command := (<-commands).(type) {
			case Exit:
				exit(messages, command)
			}
		}
	}()

	return commands
}

func exit(messages chan<- tealib.Message, cmd Exit) {
	if cmd.GoodbyeMessage != "" {
		fmt.Println(cmd.GoodbyeMessage)
	}

	os.Exit(cmd.Status)
}
