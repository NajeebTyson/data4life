package repository

import (
	"data4life/pkg/token"
	"testing"
)

func TestNewTokenStoreMongodb(t *testing.T) {
	store, err := NewTokenStoreMongodb(dbnameMongodb)
	if err != nil {
		t.Fatal(err)
	}
	store.Close()
}

func TestAddTokenMongodb(t *testing.T) {
	store, err := NewTokenStoreMongodb(dbnameMongodb)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	token := token.New(7)

	if err := store.AddToken(token); err != nil {
		t.Fatal(err)
	}
	queryToken, err := store.GetToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if queryToken == nil {
		t.Fatal("token not found")
	}

	if err := store.DeleteToken(token); err != nil {
		t.Fatal(err)
	}
}

func TestAddTokenBatchMongodb(t *testing.T) {
	store, err := NewTokenStoreMongodb(dbnameMongodb)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	tokens := []token.Token{
		token.New(testTokenSize),
		token.New(testTokenSize),
		token.New(testTokenSize),
	}

	if err := store.AddTokenBatch(tokens); err != nil {
		t.Fatal(err)
	}

	for _, token := range tokens {
		queryToken, err := store.GetToken(token)
		if err != nil {
			t.Error(err)
		}
		if queryToken == nil {
			t.Error("token not found")
		}
	}

	for _, token := range tokens {
		if err := store.DeleteToken(token); err != nil {
			t.Error(err)
		}
	}
}
