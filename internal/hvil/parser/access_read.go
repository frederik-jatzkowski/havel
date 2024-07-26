package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Read interface {
	GenerateBackLinks(*BasicBlock)
	ResolveNames(*errors.Collector)
}

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

func (read *ReadRegister) ResolveNames(errorsCollector *errors.Collector) {
	declaration, exists := read.block.registerMap[read.Identifier]
	if !exists {
		errorsCollector.Err(
			read.Pos,
			"NameError",
			fmt.Sprintf("the register %s is not yet defined in the current scope", read.Identifier),
		)
	} else {
		read.Declaration = declaration
	}
}

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

func (read *ReadVariable) ResolveNames(errorsCollector *errors.Collector) {
	localDeclaration, exists := read.block.function.variableDeclarationMap[read.Identifier]
	if !exists {
		errorsCollector.Err(
			read.Pos,
			"NameError",
			fmt.Sprintf("variable %s is not found in the current scope", read.Identifier),
		)
	} else {
		read.Declaration = localDeclaration
	}
}
