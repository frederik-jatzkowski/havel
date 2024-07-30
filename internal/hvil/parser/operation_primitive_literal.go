package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type PrimitiveLiteral struct {
	Value  uint64 `@BitLiteral`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
}

func (op *PrimitiveLiteral) GenerateBackLinks(block *BasicBlock) {
	op.block = block
}

func (op *PrimitiveLiteral) VisitLCR(visitor Visitor) {
	visitor.VisitPrimitiveLiteral(op)
}
