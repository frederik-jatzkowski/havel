package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ReadRegister struct {
	Identifier  string `"$" @Identifier`
	Declaration *WriteRegister
	Pos         lexer.Position
	Tokens      []lexer.Token
	block       *BasicBlock
}

func (read *ReadRegister) GenerateBackLinks(block *BasicBlock) {
	read.block = block
}

func (read *ReadRegister) VisitLCR(visitor Visitor) {
	visitor.VisitReadRegister(read)
}
