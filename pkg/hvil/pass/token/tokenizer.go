package token

import "github.com/alecthomas/participle/v2/lexer"

var Tokenizer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Comment", Pattern: `//.*\n`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
		{Name: "Size", Pattern: `1|2|4|8`},
		{Name: "Keyword", Pattern: `func|declare|block|return|if|then|else|goto`},
		{Name: "BitLiteral", Pattern: `0b[10]+`},
		{Name: "Ident", Pattern: `[a-z][_a-z0-9]*`},
		{Name: "Special", Pattern: `\(|\)|{|}|->|=>|:|\[|\]|,|;|\?|\$|=|\.`},
	},
})
