package debug

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type PrintU32 struct {
	tool.Node[PrintU32]

	Param memory.Read `parser:"'print_u_32' '(' @@ ')'"`
}

func (p *PrintU32) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return p.Param.ResolveNames(vars, regs)
}
