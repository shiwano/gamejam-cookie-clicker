package main

import (
	"fmt"
	"github.com/shiwano/websocket-conn"
	"net/http"
)

func runServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := conn.New()
		c.TextMessageHandler = func(text string) {
		}
		if err := c.UpgradeFromHTTP(w, r); err != nil {
			fmt.Printf("Failed to connect: %v", err)
			return
		}
		fmt.Println("Client connected")
	})
	fmt.Println("Start server")
	http.ListenAndServe(":5000", nil)
}
