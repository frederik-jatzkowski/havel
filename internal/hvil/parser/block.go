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
	function     *Function
	registerMap  map[string]*WriteRegister
}

func (block *BasicBlock) GenerateBackLinks(function *Function) {
	block.function = function

	for _, instr := range block.Instructions {
		instr.GenerateBackLinks(block)
		if instr.Result != nil {
			(*instr.Result).GenerateBackLinks(block)
		}

		instr.Operation.GenerateBackLinks(block)
	}

	block.Terminator.GenerateBackLinks(block)
}

func (block *BasicBlock) VisitLCR(visitor Visitor) {
	visitor.SetCurrentBlock(block)

	visitor.VisitBlock(block)

	for _, instr := range block.Instructions {
		instr.VisitLCR(visitor)
	}

	block.Terminator.VisitLCR(visitor)
}
