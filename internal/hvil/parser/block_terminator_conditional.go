package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ConditionalJump struct {
	Condition Read   `"if":Keyword @@`
	True      string `"then":Keyword @Identifier`
	False     string `"else":Keyword @Identifier`
	Pos       lexer.Position
	Tokens    []lexer.Token
	block     *BasicBlock
}

func (terminator *ConditionalJump) GenerateBackLinks(block *BasicBlock) {
	terminator.Condition.GenerateBackLinks(block)

	terminator.block = block
}

func (terminator *ConditionalJump) VisitLCR(visitor Visitor) {
	visitor.VisitConditionalJump(terminator)
	terminator.Condition.VisitLCR(visitor)
}
