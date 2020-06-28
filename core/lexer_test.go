package core

import (
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
