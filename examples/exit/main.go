package main

import (
	"time"

	"github.com/pb-/tealib"
	"github.com/pb-/tealib/modules/os"
)

func initialize([]string, map[string]string) (tealib.State, []tealib.Command) {
	return nil, []tealib.Command{
		os.Exit{
			Status:         1,
			GoodbyeMessage: "Bye!",
		},
	}
}

func update(s tealib.State, m tealib.Message, t time.Time) (tealib.State, []tealib.Command) {
	return s, tealib.None
}

func main() {
	tealib.Run(initialize, update, tealib.NoRender,
		tealib.Module(os.Main))
}
