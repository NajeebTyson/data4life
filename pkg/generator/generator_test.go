package generator

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

const (
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

func cleanupB(b *testing.B) {
	if err := os.Remove(testTokenFile); err != nil {
		b.Fatalf("Error cleaning up, error: %v", err)
	}
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func TestGenerateTokensFile(t *testing.T) {
	defer cleanupT(t)

	n := 5
	if err := GenerateTokensFile(testTokenFile, n, testTokenSize, seed); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(testTokenFile)
	if err != nil {
		t.Fatal(err)
	}

	lines, err := lineCounter(f)
	if err != nil {
		t.Fatal(err)
	}

	if n != lines {
		t.Fatalf("Excpected %v tokens, got %v tokens", n, lines)
	}
}

func benchmarkGenerateTokensFile(nt int, b *testing.B) {
	defer cleanupB(b)

	for n := 0; n < b.N; n++ {
		t1 := time.Now()
		GenerateTokensFile(testTokenFile, nt, testTokenSize, seed)
		t2 := time.Now()
		b.Logf("N tokens: %v, n: %v, time: %v", nt, n, t2.Sub(t1).Seconds())
	}
}

func BenchmarkGTF1000(b *testing.B)       { benchmarkGenerateTokensFile(1000, b) }
func BenchmarkGTF10Million(b *testing.B)  { benchmarkGenerateTokensFile(nTokens, b) }
func BenchmarkGTF100Million(b *testing.B) { benchmarkGenerateTokensFile(nTokens*10, b) }
