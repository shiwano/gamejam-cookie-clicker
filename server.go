package main

import (
	"fmt"
	"github.com/shiwano/websocket-conn"
	"net/http"
)

type connectionContainer struct {
	connections     []*conn.Conn
	receivedMessage chan string
	sendMessageCh   chan string
	addCh           chan *conn.Conn
	removeCh        chan *conn.Conn
}

func newConnectionContainer() *connectionContainer {
	return &connectionContainer{
		connections:     make([]*conn.Conn, 0),
		receivedMessage: make(chan string, 100),
		sendMessageCh:   make(chan string, 100),
		addCh:           make(chan *conn.Conn),
		removeCh:        make(chan *conn.Conn),
	}
}

func (c *connectionContainer) run() {
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

func runServer(container *connectionContainer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := conn.New()
		c.TextMessageHandler = func(text string) {
			fmt.Println("Received: " + text)
			container.receivedMessage <- text
		}
		c.DisconnectHandler = func() {
			container.removeCh <- c
		}
		if err := c.UpgradeFromHTTP(w, r); err != nil {
			fmt.Printf("Failed to connect: %v", err)
			return
		}
		container.addCh <- c
		fmt.Println("Client connected")
	})

	fmt.Println("Start server")
	http.ListenAndServe(":5000", nil)
}
