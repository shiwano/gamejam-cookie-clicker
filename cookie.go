package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type cookie struct {
	position *sdl.Point
}

func newCookie(point *sdl.Point) *cookie {
	return &cookie{
		position: &sdl.Point{
			X: point.X - 30 + random(-30, 30),
			Y: point.Y - 27 + random(-30, 30),
		},
	}
}

func (c *cookie) rect() *sdl.Rect {
	return &sdl.Rect{
		X: c.position.X,
		Y: c.position.Y,
		W: 60,
		H: 55,
	}
}
