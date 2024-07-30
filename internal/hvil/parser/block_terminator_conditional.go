package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ConditionalJump struct {
	Condition Read   `parser:"'if':Keyword @@"`
	True      string `parser:"'then':Keyword @Identifier"`
	False     string `parser:"'else':Keyword @Identifier"`
	Pos       lexer.Position
	Tokens    []lexer.Token
}

func (terminator *ConditionalJump) VisitLCR(visitor Visitor) {
	visitor.VisitConditionalJump(terminator)
	terminator.Condition.VisitLCR(visitor)
}
