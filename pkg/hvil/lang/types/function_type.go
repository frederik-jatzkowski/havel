package types

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type FunctionType struct {
	tool.Node[FunctionType]
	tool.NotImplemented[FunctionType]

	Parameters  tool.List[Type] `parser:"'func' '(' @@ ')'"`
	ReturnValue Type            `parser:"( '=>' @@ )?"`
}

func (t FunctionType) CanBeAssigned(other Type) bool {
	return false
}

func (t FunctionType) BitSize() int {
	return 64
}
