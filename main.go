package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"runtime"
	"time"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

func main() {
	if err := gameLoop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func gameLoop() error {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return err
	}
	defer renderer.Destroy()

	ticker := time.Tick(time.Second / 60)

loop:
	for {
		select {
		default:
		case <-ticker:
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.Present()

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					break loop
				case *sdl.MouseMotionEvent:
					fmt.Println("Move ", t)
				case *sdl.MouseButtonEvent:
					if t.State == 0 {
						fmt.Println("Click")
					}
				}
			}
		}
	}
	return nil
}

func init() {
	runtime.LockOSThread()
}
