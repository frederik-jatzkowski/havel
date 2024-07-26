package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type LocalCall struct {
	Name  string                   `"local" "." @Identifier`
	Args  CommaSeparatedList[Read] `"(" @@ ")"`
	Pos   lexer.Position
	block *BasicBlock
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
