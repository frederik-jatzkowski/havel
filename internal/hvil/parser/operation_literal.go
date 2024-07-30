package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ScalarLiteral struct {
	Value  uint64 `parser:"@BitLiteral"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (op *ScalarLiteral) VisitCLR(visitor Visitor) {
	visitor.VisitScalarLiteral(op)
}
