package reader

import (
	"data4life/internal/repository"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
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

	tokens := make(map[string]int)
	var (
		wg  = &sync.WaitGroup{}
		mtx = &sync.RWMutex{}
	)

	count := 0
	routines := 0
readLoop:
	for {

		buf := make([]byte, (tokenLength+1)*1024*128)
		c, err := fi.Read(buf)
		buf = buf[:c]

		switch {
		case err == io.EOF:
			break readLoop
		case err != nil:
			return 0, err
		}

		wg.Add(1)
		routines += 1
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			n := processChunk(string(buf), tokens, mtx)
			count += n
		}(wg)
	}
	fmt.Printf("Total routines %d\n", routines)
	wg.Wait()

	fmt.Printf("Processed %d tokens and %d unique tokens, now storing in Database\n", count, len(tokens))
	storeTokens(store, tokens)
	return count, nil
}

func processChunk(data string, tokensMap map[string]int, mtx *sync.RWMutex) int {
	tokens := strings.Split(data, "\n")
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	for _, token := range tokens {
		mtx.RLock()
		count, ok := tokensMap[token]
		mtx.RUnlock()

		mtx.Lock()
		if !ok {
			tokensMap[token] = 0
		}
		tokensMap[token] = count + 1
		mtx.Unlock()
	}

	return len(tokens)
}

func storeTokens(store repository.TokenRepository, tokens map[string]int) error {
	var (
		// poolCount       = runtime.NumCPU()
		poolCount       = 20
		insertBatchSize = 60000
		tokensPool      = sync.Pool{New: func() interface{} { return make([]string, insertBatchSize) }}
		iter            = 0
		wg              = &sync.WaitGroup{}
		ch              = make(chan []string, poolCount)
		// uniqueTokens    = make([]string, insertBatchSize)
	)

	fmt.Printf("pool size: %d, map size: %d\n", poolCount, len(tokens))

	for i := 0; i < poolCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tokens := range ch {
				insertTokens(store, tokens, &tokensPool)
			}
		}()
	}

	batchTokens := tokensPool.Get().([]string)
	for token := range tokens {
		batchTokens[iter] = token
		iter++
		if iter == insertBatchSize {
			ch <- batchTokens
			iter = 0
			// uniqueTokens = make([]string, insertBatchSize)
			batchTokens = tokensPool.Get().([]string)
		}
	}

	if iter < insertBatchSize {
		ch <- batchTokens[:iter]
	}

	close(ch)
	wg.Wait()

	return nil
}

func insertTokens(store repository.TokenRepository, tokens []string, tokensPool *sync.Pool) {
	err := store.AddTokenBatch(tokens)
	if err != nil {
		log.Println("Error in inserting batch tokens, err: ", err)
	}
	tokensPool.Put(tokens)
}
