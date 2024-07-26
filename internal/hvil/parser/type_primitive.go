package parser

import "github.com/alecthomas/participle/v2/lexer"

type PrimitiveType struct {
	BitSize uint8 `@BitSize`
	Tokens  []lexer.Token
}
