package debug

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type PrintU32 struct {
	tool.Node[PrintU32]

	Param memory.Read `parser:"'print_u_32' '(' @@ ')'"`
}

func (p *PrintU32) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return p.Param.ResolveNames(vars, regs)
}

func (p *PrintU32) ResolveTypes(target types.Type) (errs []error) {
	{
		expected := types.ScalarType{Size: 32}
		actual := p.Param.Type()

		assignable := expected.CanBeAssigned(actual)
		if !assignable {
			errs = append(errs, p.Errorf("cannot assign %s to %s", actual, expected))
		}
	}

	{
		if !target.CanBeAssigned(types.Void{}) {
			errs = append(errs, p.Errorf("cannot assign void to %s", target))
		}
	}

	return errs
}
