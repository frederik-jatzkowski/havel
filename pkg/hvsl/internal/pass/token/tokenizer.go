package token

import "github.com/alecthomas/participle/v2/lexer"

var Tokenizer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Comment", Pattern: `//.*\n`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
		{Name: "Keyword", Pattern: `func|let|return|if|else|for|type|struct|debug|literal|alu|mem|call`},
		{Name: "Ident", Pattern: `[_a-zA-Z][_a-zA-Z0-9]*`},
		{Name: "Number", Pattern: `(0b[10]+)|0x[0-9a-f]+|[0-9]+`},
		{Name: "Special", Pattern: `\(|\)|{|}|,|;|->|\.|\*`},
	},
})
