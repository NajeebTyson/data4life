package token

import (
	"testing"
)

func TestConvertToInterfaceSlice(t *testing.T) {
	testCases := []Token{
		Token("Hello"),
		Token("world"),
		Token("I love coding"),
	}

	res := ConvertToInterfaceSlice(testCases)
	for i := range testCases {
		if testCases[i] != res[i] {
			t.Errorf("Expected string %s, got %s", testCases[i], res[i])
		}
	}
}
