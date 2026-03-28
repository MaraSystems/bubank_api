package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

var alphabets = "abcdecfjijklmnopqrxtuvwxyz"

func RandomAmount(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	l := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(4))
}

func RandomCurrency() string {
	keys := make([]string, len(supportedCurrencies))
	for k := range supportedCurrencies {
		keys = append(keys, k)
	}

	return keys[rand.Intn(len(keys))]
}
