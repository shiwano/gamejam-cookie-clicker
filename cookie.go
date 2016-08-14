package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type cookie struct {
	createdAt time.Time
	position  *sdl.Point
}

func newCookie(point *sdl.Point) *cookie {
	return &cookie{
		createdAt: time.Now(),
		position: &sdl.Point{
			X: point.X - 30,
			Y: point.Y - 27,
		},
	}
}

func (c *cookie) IsDead() bool {
	return time.Now().Sub(c.createdAt).Seconds() > 5
}

func (c *cookie) rect() *sdl.Rect {
	return &sdl.Rect{
		X: c.position.X,
		Y: c.position.Y,
		W: 60,
		H: 55,
	}
}
