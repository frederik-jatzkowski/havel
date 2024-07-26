package parser

import "github.com/alecthomas/participle/v2/lexer"

type TupleType struct {
	Members []Type `"[" @@ ("," @@)* "]"`
	Tokens  []lexer.Token
}
