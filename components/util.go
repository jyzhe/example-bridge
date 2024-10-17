package components

import (
	"math/rand"
	"time"
)

// some util functions for example purposes.
const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
