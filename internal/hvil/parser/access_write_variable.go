package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type WriteVariable struct {
	Identifier  string `parser:"@Identifier"`
	Declaration VariableDeclaration
	Pos         lexer.Position
	Tokens      []lexer.Token
}

func (write *WriteVariable) VisitLCR(visitor Visitor) {
	visitor.VisitWriteVariable(write)
}
