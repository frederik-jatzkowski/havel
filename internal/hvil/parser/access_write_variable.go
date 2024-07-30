package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type WriteVariable struct {
	Identifier  string `@Identifier`
	Declaration VariableDeclaration
	Pos         lexer.Position
	Tokens      []lexer.Token
	block       *BasicBlock
}

func (write *WriteVariable) GenerateBackLinks(block *BasicBlock) {
	write.block = block
}

func (write *WriteVariable) VisitLCR(visitor Visitor) {
	visitor.VisitWriteVariable(write)
}
