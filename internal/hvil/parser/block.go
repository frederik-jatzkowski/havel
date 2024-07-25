package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type BasicBlock struct {
	Pos          lexer.Position
	function     *Function
	Identifier   string          `"block":Keyword @Identifier "{"`
	Instructions []*Instruction  `@@*`
	Terminator   BlockTerminator `"}" "=>" @@ ";"`
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

func (block *BasicBlock) ResolveNames(errorsCollector *errors.Collector) {
	block.registerMap = make(map[string]*WriteRegister, len(block.Instructions))

	for _, instr := range block.Instructions {
		if instr.Result != nil {
			(*instr.Result).ResolveNames(errorsCollector)
		}

		instr.Operation.ResolveNames(errorsCollector)
	}

	block.Terminator.ResolveNames(errorsCollector)
}
