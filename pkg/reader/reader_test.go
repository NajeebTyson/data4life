package reader

import (
	"data4life/internal/repository"
	"data4life/pkg/generator"
	"os"
	"testing"
)

const (
	testDbName    = "data4life"
	testTokenFile = "test_token.txt"
	testTokenSize = 7
	seed          = 777      // seed for random number generator
	nTokens       = 10000000 // number of tokens to generate
)

func cleanupT(t *testing.T) {
	if err := os.Remove(testTokenFile); err != nil {
		t.Fatalf("Error cleaning up, error: %v", err)
	}
}

func TestReadTokensFile(t *testing.T) {
	defer cleanupT(t)

	if err := generator.GenerateTokensFile(testTokenFile, nTokens, testTokenSize, seed); err != nil {
		t.Fatal(err)
	}

	store, err := repository.NewTokenStoreMongodb(testDbName)
	if err != nil {
		t.Fatal(err)
	}

	count, err := ReadTokensFile(store, testTokenFile, testTokenSize)
	if err != nil {
		t.Fatal(err)
	}

	if count != nTokens {
		t.Fatalf("Expected tokens count: %v, got: %v", nTokens, count)
	}

}
