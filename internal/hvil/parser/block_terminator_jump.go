package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Jump struct {
	Target string `@Identifier`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
}

func (terminator *Jump) GenerateBackLinks(block *BasicBlock) {
	terminator.block = block
}

func (terminator *Jump) VisitLCR(visitor Visitor) {
	visitor.VisitJump(terminator)
}
