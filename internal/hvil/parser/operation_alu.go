package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type AluOperation struct {
	Name   string                   `"alu" "." @Identifier`
	Args   CommaSeparatedList[Read] `"(" @@ ")"`
	Pos    lexer.Position
	Tokens []lexer.Token
	block  *BasicBlock
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
