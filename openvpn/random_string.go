package openvpn

import (
	"math/rand"
	"time"
)

var (
	randCharset = "abcdefghijklmnopqrstuvwxyABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func secureRandom() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomStringWithCharset(length int, charset string) string {
	seededRand := secureRandom()
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return RandomStringWithCharset(length, randCharset)
}
