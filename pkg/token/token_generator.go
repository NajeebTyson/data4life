package token

import (
	"math/rand"
)

// TokenGenerator class generates tokens of the given length
type TokenGenerator struct {
	generator *rand.Rand
	length    int
}

// NewTokenGenerator create a token generator with given seed
func NewTokenGenerator(seed int, length int) *TokenGenerator {
	return &TokenGenerator{
		generator: rand.New(rand.NewSource(int64(seed))),
		length:    length,
	}
}

// Generate new random token
func (tg *TokenGenerator) NewToken() Token {
	return NewToken(tg.length, tg.generator)
}
