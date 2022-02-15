package repository

import (
	"data4life/pkg/token"
)

// TokenRepository is to handle the storage of tokens
type TokenRepository interface {
	// AddToken should be used to store single token
	AddToken(token.Token) error

	// AddTokenBatch should be used to store multiple tokens at a time
	AddTokenBatch([]token.Token) error

	// GetToken is to get a token from the store
	GetToken(token.Token) (*token.Token, error)

	// DeleteToken is to delete a token from the store
	DeleteToken(token.Token) error
}
