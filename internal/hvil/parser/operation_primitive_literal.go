package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type PrimitiveLiteral struct {
	Value  uint64 `parser:"@BitLiteral"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (op *PrimitiveLiteral) VisitLCR(visitor Visitor) {
	visitor.VisitPrimitiveLiteral(op)
}
