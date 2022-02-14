package token

import (
	"fmt"
	"testing"
)

func TestTokenGenerator(t *testing.T) {
	testCases := []struct {
		token     Token
		tokenSize int
		seed      int
	}{
		{
			token:     NewTokenGenerator(100, 7).NewToken(),
			tokenSize: 7,
			seed:      100,
		},
		{
			token:     NewTokenGenerator(200, 77).NewToken(),
			tokenSize: 77,
			seed:      200,
		},
		{
			token:     NewTokenGenerator(300, 777).NewToken(),
			tokenSize: 777,
			seed:      300,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("size: %v, seed: %v", tc.tokenSize, tc.seed), func(t *testing.T) {
			gen := NewTokenGenerator(tc.seed, tc.tokenSize)
			token := gen.NewToken()

			if tc.token.Data != token.Data {
				t.Fatalf("Expected token: %v, got: %v", tc.token, token)
			}
		})
	}
}
