package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type LocalCall struct {
	Name        string                   `"local" "." @Identifier`
	Args        CommaSeparatedList[Read] `"(" @@ ")"`
	Pos         lexer.Position
	Tokens      []lexer.Token
	block       *BasicBlock
	declaration *Function
}

func (op *LocalCall) GenerateBackLinks(block *BasicBlock) {
	op.block = block

	for _, arg := range op.Args.Items {
		arg.GenerateBackLinks(block)
	}
}

func (op *LocalCall) VisitLCR(visitor Visitor) {
	visitor.VisitLocalCall(op)

	for _, arg := range op.Args.Items {
		arg.VisitLCR(visitor)
	}
}
