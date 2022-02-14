package main

import (
	"data4life/internal/common"
	"data4life/pkg/generator"
	"log"
	"time"
)

const (
	nTokens = 10000000 // number of tokens to generate
	seed    = 777      // seed for random number generator
)

func main() {
	log.Println("Generating tokens")

	t1 := time.Now()
	if err := generator.GenerateTokensFile(common.TokenFilePath, nTokens, common.TokenLength, seed); err != nil {
		panic(err)
	}
	t2 := time.Now()

	log.Printf("Generated %v tokens in %v seconds", nTokens, t2.Sub(t1).Seconds())
}
