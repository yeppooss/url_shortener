package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	var result string
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxuz1234567890"

	for i := 0; i < length; i++ {
		char := rnd.Int() % len(chars)
		result += string(chars[char])
	}

	return result
}
