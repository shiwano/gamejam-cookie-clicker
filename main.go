package main

import (
	"fmt"
	"os"
)

func main() {
	var host bool
	var serverURL string

	if len(os.Args) == 1 {
		host = true
		serverURL = "ws://localhost:5000"
		go runServer()
	} else if len(os.Args) == 2 {
		host = false
		serverURL = os.Args[1]
	} else {
		panic("Invalid args")
	}

	if err := gameLoop(host, serverURL); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
