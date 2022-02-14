package utils

import (
	"math/rand"
	"time"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyz"
)

// Factory function to create instance of random number generator with seed
func NewRandomGenerator(seed int) *rand.Rand {
	return rand.New(rand.NewSource(int64(seed)))
}

// Generate random string of n characters
func GenerateRandomString(n int, randSource *rand.Rand) string {
	ll := len(letters)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[randSource.Intn(ll)]

	}
	return string(b)
}

// Generate random string of n characters without needing any random number generator
func GenerateRandomStringSimple(n int) string {
	randSource := NewRandomGenerator(int(time.Now().UnixNano()))
	return GenerateRandomString(n, randSource)
}

// Convert string slice into interface slice, because GO does not support passing
// string slice in variadic function which accpets variadic interface
func ConvertToInterfaceSlice(items []string) []interface{} {
	itemsInterface := make([]interface{}, len(items))
	for i, item := range items {
		itemsInterface[i] = item
	}
	return itemsInterface
}
