package generator

import (
	"data4life/internal/common"
	"data4life/pkg/token"
	"fmt"
	"os"
	"strings"
	"sync"
)

func GenerateTokensFile(filename string, n, tokenLength, seed int) (errResult error) {
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := fo.Close(); err != nil {
			errResult = err
		}
	}()

	var (
		nRoutines = common.NRoutines
		wgRead    = sync.WaitGroup{}
		wgWrite   = sync.WaitGroup{}
		tokenChan = make(chan []byte, nRoutines)
	)

	for i := 0; i < nRoutines; i++ {
		batch := n / nRoutines
		if i == (nRoutines - 1) { // to put the remaining size in last batch
			batch += (n % nRoutines)
		}
		wgWrite.Add(1)
		go generateTokens(&wgWrite, tokenChan, seed+i, batch, tokenLength)
	}

	wgRead.Add(1)
	go writeTokens(&wgRead, tokenChan, fo)

	func(tokenChan chan<- []byte) {
		wgWrite.Wait()
		close(tokenChan)
	}(tokenChan)

	wgRead.Wait()

	return nil
}

func writeTokens(wg *sync.WaitGroup, tokenChan <-chan []byte, fo *os.File) {
	defer wg.Done()

	for data := range tokenChan {
		if _, err := fo.Write(data); err != nil {
			panic(err)
		}
	}
}

func generateTokens(wg *sync.WaitGroup, tokenChan chan<- []byte, seed int, n int, tokenLength int) {
	defer wg.Done()

	tokenGenerator := token.NewTokenGenerator(seed, tokenLength)

	var builder strings.Builder
	for i := 0; i < n; i++ {
		token := tokenGenerator.NewToken()
		builder.WriteString(fmt.Sprintln(token))
	}

	tokenChan <- []byte(builder.String())
}
