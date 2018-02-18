package lineserver

import (
	"bufio"
	"net"
)

type connection struct {
	server *server
	id     int
	conn   net.Conn
	send   chan Send
	closed chan bool
}

func newConnection(s *server, conn net.Conn, id int) *connection {
	c := &connection{
		server: s,
		id:     id,
		conn:   conn,
		send:   make(chan Send, 5),

		closed: make(chan bool),
	}

	go c.sender()
	go c.receiver()

	return c
}

func (c *connection) sender() {
	for {
		select {
		case send := <-c.send:
			_, err := c.conn.Write(append([]byte(send.Line), '\n'))
			if err != nil {
				c.conn.Close()
			}
		case <-c.closed:
			return
		}
	}
}

func (c *connection) receiver() {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.server.connectionInfo <- &connectionInfo{
				connection:  c,
				id:          c.id,
				isConnected: false,
			}

			c.closed <- true
			return
		}

		c.server.messages <- c.server.startServer.Receive(c.id, line[:len(line)-1])
	}
}
