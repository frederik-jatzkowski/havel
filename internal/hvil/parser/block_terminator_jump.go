package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Jump struct {
	Target string `parser:"@Identifier"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (terminator *Jump) VisitLCR(visitor Visitor) {
	visitor.VisitJump(terminator)
}
