package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	windowWidth     = 640
	windowHeight    = 480
	cookieImageName = "./cookie.png"
)

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

	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		return err
	}
	font, err := ttf.OpenFont("./font.ttf", 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		return err
	}
	defer font.Close()

	var cookies []*cookie
	var score int
	ticker := time.Tick(time.Second / 60)

loop:
	for {
		select {
		default:
		case <-ticker:
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
			renderer.SetDrawColor(255, 255, 255, 255)

			for _, c := range cookies {
				renderer.Copy(cookieTexture, nil, c.rect())
			}
			renderText(font, renderer, "Count: "+strconv.Itoa(score), &sdl.Point{X: 0, Y: 0})

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					break loop
				case *sdl.MouseButtonEvent:
					if t.State == 0 {
						cookies = append(cookies, newCookie(&sdl.Point{X: t.X, Y: t.Y}))
						score++
					}
				}
			}

			renderer.Present()

			aliveCookies := cookies[:0]
			for _, c := range cookies {
				if !c.IsDead() {
					aliveCookies = append(aliveCookies, c)
				}
			}
			cookies = aliveCookies
		}
	}
	return nil
}

func renderText(font *ttf.Font, renderer *sdl.Renderer, text string, point *sdl.Point) error {
	solid, err := font.RenderUTF8_Shaded(text, sdl.Color{R: 255, G: 255, B: 255, A: 255}, sdl.Color{R: 0, G: 0, B: 0, A: 0})
	if err != nil {
		return err
	}
	defer solid.Free()
	solidTexture, err := renderer.CreateTextureFromSurface(solid)
	if err != nil {
		return err
	}
	defer solidTexture.Destroy()
	renderer.Copy(solidTexture, nil, &sdl.Rect{X: point.X, Y: point.Y, W: solid.W, H: solid.H})
	return nil
}

func init() {
	runtime.LockOSThread()
}