package main

import (
	"fmt"
	"os"
)

func main() {
	var url string
	var connectionContainer *connectionContainer

	if len(os.Args) == 1 {
		url := "ws://localhost:5000"
		connectionContainer = newConnectionContainer()
		connectionContainer.run()
		go runServer(connectionContainer)
	} else if len(os.Args) == 2 {
		url = os.Args[1]
	} else {
		panic("Invalid args")
	}

	if err := gameLoop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
