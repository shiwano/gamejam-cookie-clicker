package main

import (
	"math/rand"
	"time"
)

func random(min, max int) int32 {
	rand.Seed(time.Now().Unix())
	return int32(rand.Intn(max-min) + min)
}
