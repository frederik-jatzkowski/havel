package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type LocalCall struct {
	Name        string                   `parser:"'local' '.' @Identifier"`
	Args        CommaSeparatedList[Read] `parser:"'(' @@ ')'"`
	Pos         lexer.Position
	Tokens      []lexer.Token
	Declaration *Function `parser:"" json:"-"`
}

func (op *LocalCall) VisitCLR(visitor Visitor) {
	visitor.VisitLocalCall(op)

	for _, arg := range op.Args.Items {
		arg.VisitCLR(visitor)
	}
}
