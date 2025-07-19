package utils

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	r = rand.New(src)
}

func RandomInt(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + int64(r.Int63n(max-min+1))
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[r.Intn(len(letters))]
	}
	return string(s)
}
