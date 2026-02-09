package token

import "github.com/alecthomas/participle/v2/lexer"

var Tokenizer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Comment", Pattern: `//.*\n`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
		{Name: "Keyword", Pattern: `func|declare|block|return|if|then|else|goto`},
		{Name: "Number", Pattern: `(0b[10]+)|0x[0-9a-f]+|[0-9]+`},
		{Name: "Ident", Pattern: `[a-z][_a-z0-9]*`},
		{Name: "Special", Pattern: `\(|\)|{|}|:|,|;|\$|=|\.`},
	},
})
