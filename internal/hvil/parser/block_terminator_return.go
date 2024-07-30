package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Return struct {
	Token  string `@"return":Keyword`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
}

func (terminator *Return) GenerateBackLinks(block *BasicBlock) {
	terminator.block = block
}

func (terminator *Return) VisitLCR(visitor Visitor) {
	visitor.VisitReturn(terminator)
}
