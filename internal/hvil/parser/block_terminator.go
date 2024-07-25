package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type BlockTerminator interface {
	GenerateBackLinks(*BasicBlock)
	ResolveNames(*errors.Collector)
}

type Return struct {
	Pos   lexer.Position
	block *BasicBlock
	Token string `@"return":Keyword`
}

func (terminator *Return) GenerateBackLinks(block *BasicBlock) {
	terminator.block = block
}

func (terminator *Return) ResolveNames(errorsCollector *errors.Collector) {}

type Jump struct {
	Pos    lexer.Position
	block  *BasicBlock
	Target string `@Identifier`
}

func (terminator *Jump) GenerateBackLinks(block *BasicBlock) {
	terminator.block = block
}

func (terminator *Jump) ResolveNames(errorsCollector *errors.Collector) {
	_, exists := terminator.block.function.blockMap[terminator.Target]
	if !exists {
		errorsCollector.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.Target),
		)
	}
}

type ConditionalJump struct {
	Pos       lexer.Position
	block     *BasicBlock
	Condition Read   `"if":Keyword @@`
	True      string `"then":Keyword @Identifier`
	False     string `"else":Keyword @Identifier`
}

func (terminator *ConditionalJump) GenerateBackLinks(block *BasicBlock) {
	terminator.Condition.GenerateBackLinks(block)

	terminator.block = block
}

func (terminator *ConditionalJump) ResolveNames(errorsCollector *errors.Collector) {
	terminator.Condition.ResolveNames(errorsCollector)

	_, exists := terminator.block.function.blockMap[terminator.True]
	if !exists {
		errorsCollector.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.True),
		)
	}

	_, exists = terminator.block.function.blockMap[terminator.False]
	if !exists {
		errorsCollector.Err(
			terminator.Pos,
			"NameError",
			fmt.Sprintf("the block %s does not exist in this function", terminator.False),
		)
	}
}
