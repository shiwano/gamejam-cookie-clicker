package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"os"
	"runtime"
	"time"
)

const (
	windowWidth     = 640
	windowHeight    = 480
	cookieImageName = "./cookie.png"
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

	cookieImage, err := img.Load(cookieImageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		return err
	}
	defer cookieImage.Free()
	cookieTexture, err := renderer.CreateTextureFromSurface(cookieImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return err
	}
	defer cookieTexture.Destroy()
	cookieTextureRect := &sdl.Rect{X: 0, Y: 0, W: 60, H: 55}

	ticker := time.Tick(time.Second / 60)

loop:
	for {
		select {
		default:
		case <-ticker:
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
			renderer.SetDrawColor(255, 255, 255, 255)

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					break loop
				case *sdl.MouseButtonEvent:
					if t.State == 0 {
						r := &sdl.Rect{X: t.X, Y: t.Y, W: cookieTextureRect.W, H: cookieTextureRect.H}
						renderer.Copy(cookieTexture, cookieTextureRect, r)
					}
				}
			}

			renderer.Present()
		}
	}
	return nil
}

func init() {
	runtime.LockOSThread()
}
