package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Return struct {
	Token  string `parser:"@'return':Keyword"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (terminator *Return) VisitLCR(visitor Visitor) {
	visitor.VisitReturn(terminator)
}
