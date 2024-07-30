package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type WriteRegister struct {
	Identifier string `parser:"'$' @Identifier"`
	Type       Type   `parser:"':' @@"`
	Pos        lexer.Position
	Tokens     []lexer.Token
}

func (write *WriteRegister) VisitLCR(visitor Visitor) {
	visitor.VisitWriteRegister(write)
}
