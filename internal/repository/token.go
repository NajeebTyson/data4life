package repository

import (
	"data4life/pkg/token"
)

type TokenRepository interface {
	AddToken(*token.Token) error
	AddTokenBatch([]string) error
	GetToken(string) (*token.Token, error)
	DeleteToken(string) error
}
