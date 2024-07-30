package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type BasicBlock struct {
	Identifier   string          `"block":Keyword @Identifier "{"`
	Instructions []*Instruction  `@@*`
	Terminator   BlockTerminator `"}" "=>" @@ ";"`
	Pos          lexer.Position
	Tokens       []lexer.Token
	RegisterMap  map[string]*WriteRegister
}

func (block *BasicBlock) VisitLCR(visitor Visitor) {
	visitor.SetCurrentBlock(block)

	visitor.VisitBlock(block)

	for _, instr := range block.Instructions {
		instr.VisitLCR(visitor)
	}

	block.Terminator.VisitLCR(visitor)
}
