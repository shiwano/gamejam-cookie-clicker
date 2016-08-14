package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type cookie struct {
	position  *sdl.Point
	imageRect *sdl.Rect
}

func (c *cookie) rect() *sdl.Rect {
	return &sdl.Rect{
		X: c.position.X,
		Y: c.position.Y,
		W: c.imageRect.W,
		H: c.imageRect.H,
	}
}
