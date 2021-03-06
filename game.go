package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/shiwano/websocket-conn"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	windowWidth     = 640
	windowHeight    = 480
	cookieImageName = "./cookie.png"
)

func gameLoop(host bool, serverURL string) error {
	userID := uuid.NewV4().String()

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	defer renderer.Destroy()

	cookieImage, err := img.Load(cookieImageName)
	if err != nil {
		return err
	}
	defer cookieImage.Free()
	cookieTexture, err := renderer.CreateTextureFromSurface(cookieImage)
	if err != nil {
		return err
	}
	defer cookieTexture.Destroy()

	if err := ttf.Init(); err != nil {
		return err
	}
	font, err := ttf.OpenFont("./font.ttf", 32)
	if err != nil {
		return err
	}
	defer font.Close()

	messageCh := make(chan string, 100)
	c := conn.New()
	c.TextMessageHandler = func(m string) { messageCh <- m }
	if _, err := c.Connect(serverURL, nil); err != nil {
		return err
	}

	var cookies []*cookie
	var lastCookieUserID string
	var lastCookieAddedAt int
	var score int
	ticker := time.Tick(time.Second / 60)

	c.WriteTextMessage(fmt.Sprintf("%v sync", userID))

loop:
	for {
		select {
		case message := <-messageCh:
			fmt.Println("Server say: " + message)
			params := strings.Split(message, " ")
			messageUserID := params[0]

			switch params[1] {
			case "cookie":
				x, _ := strconv.Atoi(params[2])
				y, _ := strconv.Atoi(params[3])
				cookies = append(cookies, newCookie(&sdl.Point{X: int32(x), Y: int32(y)}))
				score++
				addedAt, _ := strconv.Atoi(params[4])
				if lastCookieUserID != messageUserID && addedAt-lastCookieAddedAt < 2 {
					score += 10
				}
				lastCookieUserID = messageUserID
				lastCookieAddedAt = addedAt
			case "score":
				if !host {
					score, _ = strconv.Atoi(params[2])
				}
			case "sync":
				if host {
					c.WriteTextMessage(fmt.Sprintf("%v score %v", userID, score))
				}
			}
		case <-ticker:
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
			renderer.SetDrawColor(255, 255, 255, 255)

			for _, c := range cookies {
				renderer.Copy(cookieTexture, nil, c.rect())
			}
			renderText(font, renderer, "Count: "+strconv.Itoa(score), &sdl.Point{X: 0, Y: 0})
			renderer.Present()

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					break loop
				case *sdl.MouseButtonEvent:
					if t.State == 0 {
						c.WriteTextMessage(fmt.Sprintf("%v cookie %v %v %v", userID, t.X, t.Y, time.Now().Unix()))
					}
				}
			}

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
