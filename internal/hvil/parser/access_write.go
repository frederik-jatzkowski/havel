package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Write interface {
	GenerateBackLinks(*BasicBlock)
	ResolveNames(*errors.Collector)
}

type WriteRegister struct {
	Pos        lexer.Position
	block      *BasicBlock
	Identifier string `"$" @Identifier`
	Type       Type   `":" @@`
}

func (write *WriteRegister) GenerateBackLinks(block *BasicBlock) {
	write.block = block
}

func (write *WriteRegister) ResolveNames(errorsCollector *errors.Collector) {
	_, exists := write.block.registerMap[write.Identifier]
	if exists {
		errorsCollector.Err(
			write.Pos,
			"NameError",
			fmt.Sprintf("the register %s is written twice in this basic block, this violates single static assignment form", write.Identifier),
		)
	}

	write.block.registerMap[write.Identifier] = write
}

type WriteVariable struct {
	Pos        lexer.Position
	block      *BasicBlock
	Identifier string `@Identifier`
}

func (write *WriteVariable) GenerateBackLinks(block *BasicBlock) {
	write.block = block
}

func (write *WriteVariable) ResolveNames(errorsCollector *errors.Collector) {
	_, exists := write.block.function.variableDeclarationMap[write.Identifier]
	if !exists {
		errorsCollector.Err(
			write.Pos,
			"NameError",
			fmt.Sprintf("variable %s is not found in the current scope", write.Identifier),
		)
	}
}
