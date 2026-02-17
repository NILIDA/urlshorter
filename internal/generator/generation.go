package generator

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const length = 10

func Generate() string{
	b := make([]byte, length)
	for i := range b{
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
