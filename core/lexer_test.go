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

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		tokens []Token
		err    error
	}{
		{
			input: "select a",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(selectKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: "a",
					kind:  identifierKind,
				},
			},
		},
		{
			input: "select true",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(selectKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: "true",
					kind:  boolKind,
				},
			},
		},
		{
			input: "select 1",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(selectKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: "1",
					kind:  numericKind,
				},
			},
			err: nil,
		},
		{
			input: "select 'foo' || 'bar';",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(selectKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: "foo",
					kind:  stringKind,
				},
				{
					loc:   location{col: 13, line: 0},
					value: string(concatSymbol),
					kind:  symbolKind,
				},
				{
					loc:   location{col: 16, line: 0},
					value: "bar",
					kind:  stringKind,
				},
				{
					loc:   location{col: 21, line: 0},
					value: string(semicolonSymbol),
					kind:  symbolKind,
				},
			},
			err: nil,
		},
		{
			input: "CREATE TABLE u (id INT, name TEXT)",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(createKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: string(tableKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 13, line: 0},
					value: "u",
					kind:  identifierKind,
				},
				{
					loc:   location{col: 15, line: 0},
					value: "(",
					kind:  symbolKind,
				},
				{
					loc:   location{col: 16, line: 0},
					value: "id",
					kind:  identifierKind,
				},
				{
					loc:   location{col: 19, line: 0},
					value: "int",
					kind:  keywordKind,
				},
				{
					loc:   location{col: 22, line: 0},
					value: ",",
					kind:  symbolKind,
				},
				{
					loc:   location{col: 24, line: 0},
					value: "name",
					kind:  identifierKind,
				},
				{
					loc:   location{col: 29, line: 0},
					value: "text",
					kind:  keywordKind,
				},
				{
					loc:   location{col: 33, line: 0},
					value: ")",
					kind:  symbolKind,
				},
			},
		},
		{
			input: "insert into users values (105, 233)",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(insertKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: string(intoKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 12, line: 0},
					value: "users",
					kind:  identifierKind,
				},
				{
					loc:   location{col: 18, line: 0},
					value: string(valuesKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 25, line: 0},
					value: "(",
					kind:  symbolKind,
				},
				{
					loc:   location{col: 26, line: 0},
					value: "105",
					kind:  numericKind,
				},
				{
					loc:   location{col: 30, line: 0},
					value: ",",
					kind:  symbolKind,
				},
				{
					loc:   location{col: 32, line: 0},
					value: "233",
					kind:  numericKind,
				},
				{
					loc:   location{col: 36, line: 0},
					value: ")",
					kind:  symbolKind,
				},
			},
			err: nil,
		},
		{
			input: "SELECT id FROM users;",
			tokens: []Token{
				{
					loc:   location{col: 0, line: 0},
					value: string(selectKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 7, line: 0},
					value: "id",
					kind:  identifierKind,
				},
				{
					loc:   location{col: 10, line: 0},
					value: string(fromKeyword),
					kind:  keywordKind,
				},
				{
					loc:   location{col: 15, line: 0},
					value: "users",
					kind:  identifierKind,
				},
				{
					loc:   location{col: 20, line: 0},
					value: ";",
					kind:  symbolKind,
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		tokens, err := LexParse(test.input)
		assert.Equal(t, test.err, err, test.input)
		assert.Equal(t, len(test.tokens), len(tokens), test.input)

		for i, tok := range tokens {
			assert.Equal(t, &test.tokens[i], tok, test.input)
		}
	}
}
