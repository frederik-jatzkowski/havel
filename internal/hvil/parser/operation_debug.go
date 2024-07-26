package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type DebugOperation struct {
	Name  string                   `"debug" "." @Identifier`
	Args  CommaSeparatedList[Read] `"(" @@ ")"`
	Pos   lexer.Position
	block *BasicBlock
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
