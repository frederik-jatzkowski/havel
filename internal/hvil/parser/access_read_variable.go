package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ReadVariable struct {
	Identifier  string `parser:"@Identifier"`
	Declaration VariableDeclaration
	Pos         lexer.Position
	Tokens      []lexer.Token
}

func (read *ReadVariable) VisitLCR(visitor Visitor) {
	visitor.VisitReadVariable(read)
}
