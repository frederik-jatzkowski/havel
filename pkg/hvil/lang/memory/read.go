package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Read interface {
	tool.NodeLike
	ResolveNames(vars names.Scope[VarDecl], regs names.Scope[RegWrite]) (errs []error)
	Type() types.Type
}
