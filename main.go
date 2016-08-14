package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		go runServer()
	} else if len(os.Args) == 2 {
		url := os.Args[1]
		fmt.Println(url)
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
