package random

import (
	"golang.org/x/exp/rand"
	"time"
)

func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
