package main

import (
	"fmt"
	"time"

	"github.com/pb-/tealib"
	server "github.com/pb-/tealib/modules/lineserver"
	"github.com/pb-/tealib/modules/os"
)

type received struct {
	id   int
	line string
}

func initialize([]string, map[string]string) (tealib.State, []tealib.Command) {
	return nil, []tealib.Command{server.StartServer{
		Port:       4711,
		Error:      func(err error) tealib.Message { return err },
		Connect:    func(_ int) tealib.Message { return nil },
		Receive:    func(id int, line string) tealib.Message { return received{id: id, line: line} },
		Disconnect: func(_ int) tealib.Message { return nil },
	}}
}

func update(s tealib.State, m tealib.Message, t time.Time) (tealib.State, []tealib.Command) {
	switch msg := m.(type) {
	case received:
		return s, []tealib.Command{server.Send{
			ID:    msg.id,
			Line:  fmt.Sprintf("at %s you said: %s", t.String(), msg.line),
			Error: func(err error) tealib.Message { return err },
		}}
	case error:
		return s, []tealib.Command{os.Exit{
			Status:         1,
			GoodbyeMessage: msg.Error(),
		}}
	}

	return s, tealib.None
}

func main() {
	tealib.Run(initialize, update, tealib.NoRender,
		tealib.Module(os.Main),
		tealib.Module(server.Main),
	)
}
