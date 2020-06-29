package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberLexer(t *testing.T) {
	tests := []struct {
		inputString      string
		expectedIsNumber bool
		expectedValue    string
	}{
		{
			inputString:      "105",
			expectedIsNumber: true,
			expectedValue:    "105",
		},
		{
			inputString:      "105 ",
			expectedIsNumber: true,
			expectedValue:    "105",
		},
		{

			inputString:      "123.",
			expectedIsNumber: true,
			expectedValue:    "123",
		},
		{
			inputString:      "123.145",
			expectedIsNumber: true,
			expectedValue:    "123",
		},
		{
			inputString:      "1e5",
			expectedIsNumber: true,
			expectedValue:    "1",
		},
		{
			inputString:      "1.1e-2",
			expectedIsNumber: true,
			expectedValue:    "1",
		},
		// false tests
		{
			inputString:      "e4",
			expectedIsNumber: false,
			expectedValue:    "",
		},
		{
			inputString:      "..5",
			expectedIsNumber: false,
			expectedValue:    "",
		},
		{
			inputString:      " 5",
			expectedIsNumber: false,
			expectedValue:    "",
		},
	}

	for _, test := range tests {
		gotToken, _, ok := numberLexer(test.inputString, cursor{})
		assert.Equal(t, test.expectedIsNumber, ok)

		if ok {
			assert.Equal(t, test.expectedValue, gotToken.value)
		}
	}
}

func TestSymbolLexer(t *testing.T) {
	tests := []struct {
		symbol bool
		value  string
	}{
		{
			symbol: true,
			value:  "= ",
		},
		{
			symbol: true,
			value:  "||",
		},
	}

	for _, test := range tests {
		tok, _, ok := symbolLexer(test.value, cursor{})
		assert.Equal(t, test.symbol, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, test.value, tok.value, test.value)
		}
	}
}

func TestKeywordLexer(t *testing.T) {
	tests := []struct {
		expectedIsKeyword bool
		input             string
	}{
		{
			expectedIsKeyword: true,
			input:             "as",
		},
		{
			expectedIsKeyword: true,
			input:             "select ",
		},
		{
			expectedIsKeyword: true,
			input:             "from",
		},
		{
			expectedIsKeyword: true,
			input:             "SELECT",
		},
		{
			expectedIsKeyword: true,
			input:             "into",
		},
		// false tests
		{
			expectedIsKeyword: false,
			input:             " into",
		},
		{
			expectedIsKeyword: false,
			input:             "flubbrety",
		},
	}

	for _, test := range tests {
		tok, _, ok := keywordLexer(test.input, cursor{})
		assert.Equal(t, test.expectedIsKeyword, ok, test.input)
		if ok {
			test.input = strings.TrimSpace(test.input)
			assert.Equal(t, strings.ToLower(test.input), tok.value, test.input)
		}
	}
}
