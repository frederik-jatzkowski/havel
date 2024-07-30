package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type DebugOperation struct {
	Name   string                   `"debug" "." @Identifier`
	Args   CommaSeparatedList[Read] `"(" @@ ")"`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
}

func (op *DebugOperation) GenerateBackLinks(block *BasicBlock) {
	op.block = block

	for _, arg := range op.Args.Items {
		arg.GenerateBackLinks(block)
	}
}

func (op *DebugOperation) VisitLCR(visitor Visitor) {
	visitor.VisitDebugOperation(op)

	for _, arg := range op.Args.Items {
		arg.VisitLCR(visitor)
	}
}
