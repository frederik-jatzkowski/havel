package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type AluOperation struct {
	Name   string                   `parser:"'alu' '.' @Identifier"`
	Args   CommaSeparatedList[Read] `parser:"'(' @@ ')'"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (op *AluOperation) VisitCLR(visitor Visitor) {
	visitor.VisitAluOperation(op)

	for _, arg := range op.Args.Items {
		arg.VisitCLR(visitor)
	}
}
