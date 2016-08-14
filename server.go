package main

import (
	"fmt"
	"github.com/shiwano/websocket-conn"
	"net/http"
)

type serverHub struct {
	connections     []*conn.Conn
	receivedMessage chan string
	sendMessageCh   chan string
	addCh           chan *conn.Conn
	removeCh        chan *conn.Conn
}

func newServerHub() *serverHub {
	return &serverHub{
		connections:     make([]*conn.Conn, 0),
		receivedMessage: make(chan string, 100),
		sendMessageCh:   make(chan string, 100),
		addCh:           make(chan *conn.Conn),
		removeCh:        make(chan *conn.Conn),
	}
}

func (c *serverHub) run() {
	for {
		select {
		case message := <-c.sendMessageCh:
			for _, c := range c.connections {
				c.WriteTextMessage(message)
			}
		case connection := <-c.addCh:
			c.connections = append(c.connections, connection)
		case connection := <-c.removeCh:
			connections := c.connections[:0]
			for _, c := range c.connections {
				if c != connection {
					connections = append(connections, c)
				}
			}
			c.connections = connections
		}
	}
}

func runServer(hub *serverHub) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := conn.New()
		c.TextMessageHandler = func(text string) {
			fmt.Println("Received: " + text)
			hub.receivedMessage <- text
		}
		c.DisconnectHandler = func() {
			hub.removeCh <- c
		}
		if err := c.UpgradeFromHTTP(w, r); err != nil {
			fmt.Printf("Failed to connect: %v", err)
			return
		}
		hub.addCh <- c
		fmt.Println("Client connected")
	})

	fmt.Println("Start server")
	http.ListenAndServe(":5000", nil)
}
