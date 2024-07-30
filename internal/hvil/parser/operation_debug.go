package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type DebugOperation struct {
	Name   string                   `parser:"'debug' '.' @Identifier"`
	Args   CommaSeparatedList[Read] `parser:"'(' @@ ')'"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (op *DebugOperation) VisitCLR(visitor Visitor) {
	visitor.VisitDebugOperation(op)

	for _, arg := range op.Args.Items {
		arg.VisitCLR(visitor)
	}
}
