package main

import (
	"data4life/internal/common"
	"data4life/internal/repository"
	"data4life/pkg/reader"
	"log"
	"time"
)

const (
	dbName = "data4life"
)

func main() {
	log.Println("Reading token file: ", common.TokenFilePath)

	store, err := repository.NewTokenStoreMongodb(dbName)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	t1 := time.Now()
	if _, err := reader.ReadTokensFile(store, common.TokenFilePath, common.TokenLength); err != nil {
		panic(err)
	}
	t2 := time.Now()

	log.Printf("Successfully read file and insert tokens in DB, time: %.2f\n", t2.Sub(t1).Seconds())
}
