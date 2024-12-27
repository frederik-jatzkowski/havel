package alu

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Operation struct {
	tool.Node[Operation]
	tool.NotImplemented[Operation]

	Name string                 `parser:"'alu' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (o Operation) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return nil
}
