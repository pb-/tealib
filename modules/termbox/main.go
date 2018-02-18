package termbox

import (
	"github.com/nsf/termbox-go"

	"github.com/pb-/tealib"
)

// Init calls termbox.Init().
type Init struct {
	Error func(error) tealib.Message
}

// SetInputMode calls termbox.SetInputMode().
type SetInputMode struct {
	Mode      InputMode
	InputMode func(termbox.InputMode) tealib.Message
}

// Shutdown calls termbox.Close().
type Shutdown struct {
}

func Main(messages chan<- tealib.Message) chan<- tealib.Command {
	commands := make(chan tealib.Command)
	events := make(chan *termbox.Event)

	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	go func() {
		select {
		case command := <-commands:
		case event := <-events:
			messages <- event
		}
	}()

	return commands
}
