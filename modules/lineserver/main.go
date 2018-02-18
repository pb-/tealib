package lineserver

import (
	"errors"

	"github.com/pb-/tealib"
)

// StartServer starts a line-based server.
type StartServer struct {
	Port       int
	Error      func(err error) tealib.Message
	Connect    func(id int) tealib.Message
	Receive    func(id int, line string) tealib.Message
	Disconnect func(id int) tealib.Message
}

// Send sends a message.
type Send struct {
	ID    int
	Line  string
	Error func(err error) tealib.Message
}

// Disconnect disconnects a client.
type Disconnect struct {
	ID int
}

func Main(messages chan<- tealib.Message) chan<- tealib.Command {
	commands := make(chan tealib.Command)

	go loop(messages, commands)

	return commands
}

func loop(messages chan<- tealib.Message, commands <-chan tealib.Command) {
	clients := make(map[int]*connection)
	connInfo := make(chan *connectionInfo)
	ids := make(chan int)

	go func() {
		for id := 1; true; id++ {
			ids <- id
		}
	}()

	for {
		select {
		case command := <-commands:
			switch cmd := command.(type) {
			case StartServer:
				startServer(messages, cmd, connInfo, ids)
			case Send:
				send(messages, cmd, clients)
			case Disconnect:
				disconnect(cmd, clients)
			}
		case ci := <-connInfo:
			if ci.isConnected {
				clients[ci.id] = ci.connection
				messages <- ci.connection.server.startServer.Connect(ci.id)
			} else {
				delete(clients, ci.id)
				messages <- ci.connection.server.startServer.Disconnect(ci.id)
			}
		}
	}
}

func send(messages chan<- tealib.Message, cmd Send, clients map[int]*connection) {
	conn, ok := clients[cmd.ID]
	if !ok {
		if cmd.Error != nil {
			messages <- cmd.Error(errors.New("bad client id"))
		}
		return
	}

	select {
	case conn.send <- cmd:
		return
	default:
		if cmd.Error != nil {
			messages <- cmd.Error(errors.New("send queue is full"))
		}
		return
	}
}

func disconnect(cmd Disconnect, clients map[int]*connection) {
	c, ok := clients[cmd.ID]
	if !ok {
		return
	}

	c.conn.Close()
}
