package token_test

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass/token"
	"strings"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "whitespace and comment",
			input: `
// hi
			`,
			expected: []string{
				"",
				"// hi",
				"",
			},
		},
		{
			name: "identifier",
			input: `
// hi asd
asd
			`,
			expected: []string{
				"",
				"// hi asd",
				"asd",
				"",
			},
		},
		{
			name:  "special",
			input: `(){}[]:=>,;`,
			expected: []string{
				"(",
				")",
				"{",
				"}",
				"[",
				"]",
				":",
				"=>",
				",",
				";",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer, err := token.Tokenizer.LexString(test.name, test.input)
			if err != nil {
				t.Error(err)
			}

			for _, expectedString := range test.expected {
				actual, err := lexer.Next()
				if err != nil {
					t.Error(err)
				}

				actualString := strings.TrimSpace(actual.Value)
				if actualString != expectedString {
					t.Errorf("expected '%s' but got '%s'", expectedString, actualString)
				}
			}
		})
	}
}
