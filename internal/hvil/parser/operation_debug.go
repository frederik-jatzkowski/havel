package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type DebugOperation struct {
	Name   string                   `"debug" "." @Identifier`
	Args   CommaSeparatedList[Read] `"(" @@ ")"`
	Pos    lexer.Position
	Tokens []lexer.Token
}

func (op *DebugOperation) VisitLCR(visitor Visitor) {
	visitor.VisitDebugOperation(op)
}
