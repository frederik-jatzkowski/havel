package scalar

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Literal struct {
	tool.Node[Literal]

	Value uint64 `parser:"@BitLiteral"`
}

func (l *Literal) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return nil
}
