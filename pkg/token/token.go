package token

import (
	"data4life/internal/utils"
	"math/rand"
)

// Token struct is the main struct to represent token
type Token string

// Create a new Token instance with random string of length given by using given seed
func NewToken(length int, source *rand.Rand) Token {
	return Token(utils.GenerateRandomString(length, source))
}

// Create a new Token instance with random string of given length
func New(length int) Token {
	return Token(utils.GenerateRandomStringSimple(length))
}

// Convert string slice into interface slice, because GO does not support passing
// string slice in variadic function which accpets variadic interface
func ConvertToInterfaceSlice(items []Token) []interface{} {
	itemsInterface := make([]interface{}, len(items))
	for i, item := range items {
		itemsInterface[i] = item
	}
	return itemsInterface
}
