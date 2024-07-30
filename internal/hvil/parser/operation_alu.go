package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type AluOperation struct {
	Name   string                   `"alu" "." @Identifier`
	Args   CommaSeparatedList[Read] `"(" @@ ")"`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
}

func (op *AluOperation) GenerateBackLinks(block *BasicBlock) {
	op.block = block

	for _, arg := range op.Args.Items {
		arg.GenerateBackLinks(block)
	}
}

func (op *AluOperation) VisitLCR(visitor Visitor) {
	visitor.VisitAluOperation(op)

	for _, arg := range op.Args.Items {
		arg.VisitLCR(visitor)
	}
}
