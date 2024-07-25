package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Operation interface {
	GenerateBackLinks(*BasicBlock)
	ResolveNames(*errors.Collector)
}

type PrimitiveLiteral struct {
	Pos   lexer.Position
	block *BasicBlock
	Value uint64 `@BitLiteral`
}

func (op *PrimitiveLiteral) GenerateBackLinks(block *BasicBlock) {
	op.block = block
}

func (op *PrimitiveLiteral) ResolveNames(errorsCollector *errors.Collector) {}

type AluOperation struct {
	Pos   lexer.Position
	block *BasicBlock
	Name  string                   `"alu" "." @Identifier`
	Args  CommaSeparatedList[Read] `"(" @@ ")"`
}

func (op *AluOperation) GenerateBackLinks(block *BasicBlock) {
	op.block = block

	for _, arg := range op.Args.Items {
		arg.GenerateBackLinks(block)
	}
}

func (op *AluOperation) ResolveNames(errorsCollector *errors.Collector) {
	switch op.Name {
	case "add_u_32", "sub_u_32", "lt_u_32":
	default:
		errorsCollector.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("there is no ALU operation called '%s'", op.Name),
		)
	}

	for _, arg := range op.Args.Items {
		arg.ResolveNames(errorsCollector)
	}
}

type LocalCall struct {
	Pos   lexer.Position
	block *BasicBlock
	Name  string                   `"local" "." @Identifier`
	Args  CommaSeparatedList[Read] `"(" @@ ")"`
}

func (op *LocalCall) GenerateBackLinks(block *BasicBlock) {
	op.block = block

	for _, arg := range op.Args.Items {
		arg.GenerateBackLinks(block)
	}
}

func (op *LocalCall) ResolveNames(errorsCollector *errors.Collector) {
	_, exists := op.block.function.pkg.functionsMap[op.Name]
	if !exists {
		errorsCollector.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("the function %s is not defined locally in this package", op.Name),
		)
	}

	for _, arg := range op.Args.Items {
		arg.ResolveNames(errorsCollector)
	}
}

type DebugOperation struct {
	Pos   lexer.Position
	block *BasicBlock
	Name  string                   `"debug" "." @Identifier`
	Args  CommaSeparatedList[Read] `"(" @@ ")"`
}

func (op *DebugOperation) GenerateBackLinks(block *BasicBlock) {
	op.block = block

	for _, arg := range op.Args.Items {
		arg.GenerateBackLinks(block)
	}
}

func (op *DebugOperation) ResolveNames(errorsCollector *errors.Collector) {
	switch op.Name {
	case "print_u_32":
	default:
		errorsCollector.Err(
			op.Pos,
			"NameError",
			fmt.Sprintf("there is no debug operation called '%s'", op.Name),
		)
	}

	for _, arg := range op.Args.Items {
		arg.ResolveNames(errorsCollector)
	}
}
