package lineserver

import (
	"fmt"
	"net"

	"github.com/pb-/tealib"
)

type connectionInfo struct {
	connection  *connection
	id          int
	isConnected bool
}

type server struct {
	startServer    *StartServer
	messages       chan<- tealib.Message
	connectionInfo chan<- *connectionInfo
}

func startServer(messages chan<- tealib.Message, ss StartServer, conns chan<- *connectionInfo, ids <-chan int) {
	server := &server{
		startServer:    &ss,
		messages:       messages,
		connectionInfo: conns,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ss.Port))
	if err != nil {
		if ss.Error != nil {
			messages <- ss.Error(err)
		}
		return
	}

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				messages <- ss.Error(err)
				return
			}

			id := <-ids
			conns <- &connectionInfo{
				connection:  newConnection(server, conn, id),
				id:          id,
				isConnected: true,
			}
		}
	}()
}
