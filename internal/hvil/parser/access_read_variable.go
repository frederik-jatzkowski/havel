package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ReadVariable struct {
	Identifier  string `@Identifier`
	Declaration VariableDeclaration
	Pos         lexer.Position
	Tokens      []lexer.Token
	block       *BasicBlock
}

func (read *ReadVariable) GenerateBackLinks(block *BasicBlock) {
	read.block = block
}

func (read *ReadVariable) VisitLCR(visitor Visitor) {
	visitor.VisitReadVariable(read)
}
