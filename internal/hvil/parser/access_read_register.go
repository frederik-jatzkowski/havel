package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ReadRegister struct {
	Identifier  string `parser:"'$' @Identifier"`
	Declaration *WriteRegister
	Pos         lexer.Position
	Tokens      []lexer.Token
}

func (read *ReadRegister) VisitCLR(visitor Visitor) {
	visitor.VisitReadRegister(read)
}
