package utils

import (
	"fmt"
	"testing"
)

func TestNewRandomGenerator(t *testing.T) {
	var (
		n       = 10
		seed    = 100
		maxN    = 1000
		numbers = make([]int, 0)
	)

	source := NewRandomGenerator(seed)
	for i := 0; i < n; i++ {
		numbers = append(numbers, source.Intn(maxN))
	}

	for tc := 0; tc < 5; tc++ {
		t.Run("", func(t *testing.T) {
			source := NewRandomGenerator(seed)
			for i := 0; i < n; i++ {
				num := source.Intn(maxN)
				if num != numbers[i] {
					t.Errorf("Expected number %v, got %v", numbers[i], num)
				}
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	testCases := []struct {
		length int
	}{
		{length: 4},
		{length: 7},
		{length: 100},
		{length: 9999},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Length: %v", tc.length), func(t *testing.T) {
			source := NewRandomGenerator(tc.length) // giving length as seed
			str := GenerateRandomString(tc.length, source)
			if len(str) != tc.length {
				t.Fatalf("Expected string length %v, got %v", tc.length, len(str))
			}
		})
	}
}

func TestGenerateRandomStringSimple(t *testing.T) {
	testCases := []struct {
		length int
	}{
		{length: 4},
		{length: 7},
		{length: 100},
		{length: 9999},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Length: %v", tc.length), func(t *testing.T) {
			str := GenerateRandomStringSimple(tc.length)
			if len(str) != tc.length {
				t.Fatalf("Expected string length %v, got %v", tc.length, len(str))
			}
		})
	}
}
