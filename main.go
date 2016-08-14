package main

import (
	"fmt"
	"os"
)

func main() {
	var serverURL string

	if len(os.Args) == 1 {
		serverURL = "ws://localhost:5000"
		go runServer(serverHub)
	} else if len(os.Args) == 2 {
		serverURL = os.Args[1]
	} else {
		panic("Invalid args")
	}

	if err := gameLoop(serverURL); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
