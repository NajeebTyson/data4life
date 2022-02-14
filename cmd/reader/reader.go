package main

import (
	"data4life/internal/common"
	"data4life/internal/repository"
	"data4life/pkg/reader"
	"log"
	"time"
)

func main() {
	log.Println("Reading token file: ", common.TokenFilePath)

	store, err := repository.NewTokenStore()
	if err != nil {
		panic(err)
	}

	t1 := time.Now()
	if _, err := reader.ReadTokensFile(store, common.TokenFilePath, common.TokenLength); err != nil {
		panic(err)
	}
	t2 := time.Now()

	log.Println("Successfully read file and insert tokens in DB, time:", t2.Sub(t1).Seconds())
}
