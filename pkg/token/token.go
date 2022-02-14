package token

import (
	"data4life/internal/utils"
	"math/rand"
)

// Token struct is the main struct to represent token
type Token struct {
	Data string
}

// Create a new Token instance with random string of length given by using given seed
func NewToken(length int, source *rand.Rand) Token {
	return Token{Data: utils.GenerateRandomString(length, source)}
}

// Create a new Token instance with random string of given length
func New(length int) Token {
	return Token{Data: utils.GenerateRandomStringSimple(length)}
}
