package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Jump struct {
	Target string `parser:"@Identifier"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (terminator *Jump) VisitCLR(visitor Visitor) {
	visitor.VisitJump(terminator)
}
