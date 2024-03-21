package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var currency = []string{"USD", "IND", "EUR", "JPY"}

// RandomInt  returns a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString  generates a random alphanumeric string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		t := alphabet[rand.Intn(k)]
		sb.WriteByte(t)
	}

	return sb.String()
}

// RandomOwner creates an owner to use in tests
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney  creates money value to use in tests
func RandomMoney() int64 {
	return RandomInt(10, 1000)
}

// RandomCurrency  selects a random currency from the list of ["USD", "INR", "EUR", "JPY"]
func RandomCurrency() string {
	return currency[rand.Intn(len(currency))]
}
