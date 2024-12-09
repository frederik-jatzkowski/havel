package function

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Type struct {
	tool.Node[Type]
	tool.NotImplemented[Type]

	Parameters  tool.List[types.Type] `parser:"'func' '(' @@ ')'"`
	ReturnValue types.Type            `parser:"( '=>' @@ )?"`
}

func (t Type) Equals(other types.Type) bool {
	return t.String() == other.String()
}
