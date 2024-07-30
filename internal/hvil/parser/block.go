package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type BasicBlock struct {
	Identifier   string          `parser:"'block':Keyword @Identifier '{'"`
	Instructions []*Instruction  `parser:"@@*"`
	Terminator   BlockTerminator `parser:"'}' '=>' @@ ';'"`
	Pos          lexer.Position
	Tokens       []lexer.Token
	RegisterMap  map[string]*WriteRegister
}

func (block *BasicBlock) VisitCLR(visitor Visitor) {
	visitor.SetCurrentBlock(block)

	visitor.VisitBlock(block)

	for _, instr := range block.Instructions {
		instr.VisitCLR(visitor)
	}

	block.Terminator.VisitCLR(visitor)
}
