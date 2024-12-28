package debug

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type Dump struct {
	tool.Node[Dump]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Param memory.Read `parser:"'dump' '(' @@ ')'"`
}

func (d *Dump) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return d.Param.ResolveNames(vars, regs)
}

func (d *Dump) ResolveTypes(target types.Type) (errs []error) {
	if !target.CanBeAssigned(types.Void{}) {
		errs = append(errs, d.Errorf("cannot assign void to %s", target))
	}

	d.TypeCheckPass.Type = target

	return errs
}
