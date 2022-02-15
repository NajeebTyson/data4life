package reader

import (
	"data4life/internal/repository"
	"data4life/pkg/token"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func ReadTokensFile(store repository.TokenRepository, filename string, tokenLength int) (tokenCount int, errResult error) {
	fi, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			errResult = err
			tokenCount = 0
		}
	}()

	var (
		tokens = make(map[token.Token]int)
		wg     = &sync.WaitGroup{}
		mtx    = &sync.RWMutex{}
		count  = 0
		t1     = time.Now()
	)

readLoop:
	for {
		buf := make([]byte, (tokenLength+1)*1024*64)
		c, err := fi.Read(buf)
		buf = buf[:c]

		switch {
		case err == io.EOF:
			break readLoop
		case err != nil:
			return 0, err
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			n := processChunk(string(buf), tokens, mtx)
			count += n
		}(wg)
	}
	wg.Wait()
	t2 := time.Now()

	log.Printf("Processed %d tokens and %d unique tokens in %.2f seconds\n", count, len(tokens), t2.Sub(t1).Seconds())
	dt1 := time.Now()

	storeTokens(store, tokens)

	dt2 := time.Now()
	log.Printf("Stored tokens in %.2f seconds\n", dt2.Sub(dt1).Seconds())

	return count, nil
}

func processChunk(data string, tokensMap map[token.Token]int, mtx *sync.RWMutex) int {
	tokens := strings.Split(data, "\n")
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	for _, t := range tokens {
		tok := token.Token(t)
		mtx.RLock()
		count, ok := tokensMap[tok]
		mtx.RUnlock()

		mtx.Lock()
		if !ok {
			tokensMap[tok] = 0
		}
		tokensMap[tok] = count + 1
		mtx.Unlock()
	}

	return len(tokens)
}

func storeTokens(store repository.TokenRepository, tokens map[token.Token]int) error {
	var (
		poolCount       = 40
		insertBatchSize = 500
		iter            = 0
		wg              = &sync.WaitGroup{}
		ch              = make(chan []token.Token, poolCount)
		uniqueTokens    = make([]token.Token, insertBatchSize)
	)

	log.Printf("Storing tokens, pool size: %d, map size: %d\n", poolCount, len(tokens))

	// creating go routines pool which stores tokens into DB
	for i := 0; i < poolCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tokens := range ch {
				insertTokens(store, tokens)
			}
		}()
	}

	// passing batch of tokens into channel which consumed in go routines
	for t := range tokens {
		uniqueTokens[iter] = t
		iter++
		if iter == insertBatchSize {
			ch <- uniqueTokens
			iter = 0
			uniqueTokens = make([]token.Token, insertBatchSize)
		}
	}
	if iter < insertBatchSize {
		ch <- uniqueTokens[:iter]
	}

	close(ch) // closign the channel
	wg.Wait() // waiting for go-routines to complete

	return nil
}

func insertTokens(store repository.TokenRepository, tokens []token.Token) {
	err := store.AddTokenBatch(tokens)
	if err != nil {
		log.Println("Error in inserting batch tokens, err: ", err)
	}
}
