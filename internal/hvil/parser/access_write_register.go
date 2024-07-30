package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type WriteRegister struct {
	Identifier string `"$" @Identifier`
	Type       Type   `":" @@`
	Pos        lexer.Position
	Tokens     []lexer.Token
	block      *BasicBlock
}

func (write *WriteRegister) GenerateBackLinks(block *BasicBlock) {
	write.block = block
}

func (write *WriteRegister) VisitLCR(visitor Visitor) {
	visitor.VisitWriteRegister(write)
}
