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
	Token  string `@"return":Keyword`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
}

func (terminator *Return) GenerateBackLinks(block *BasicBlock) {
	terminator.block = block
}

func (terminator *Return) ResolveNames(errorsCollector *errors.Collector) {}

type Jump struct {
	Target string `@Identifier`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
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
	Condition Read   `"if":Keyword @@`
	True      string `"then":Keyword @Identifier`
	False     string `"else":Keyword @Identifier`
	Pos       lexer.Position
	Tokens    []lexer.Token
	block     *BasicBlock
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
