package token

import "github.com/alecthomas/participle/v2/lexer"

var Tokenizer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Comment", Pattern: `//.*\n`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
		{Name: "BitSize", Pattern: `8|16|32|64`},
		{Name: "Keyword", Pattern: `func|declare|block|return|if|then|else`},
		{Name: "BitLiteral", Pattern: `0b[10]+`},
		{Name: "Identifier", Pattern: `[a-z][_a-z0-9]*`},
		{Name: "Special", Pattern: `\(|\)|{|}|=>|:|\[|\]|,|;|\?|\$|=|\.`},
	},
})
