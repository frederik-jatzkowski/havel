package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type PrimitiveLiteral struct {
	Pos   lexer.Position
	block *BasicBlock
	Value uint64 `@BitLiteral`
}

func (op *PrimitiveLiteral) GenerateBackLinks(block *BasicBlock) {
	op.block = block
}

func (op *PrimitiveLiteral) ResolveNames(errorsCollector *errors.Collector) {}
