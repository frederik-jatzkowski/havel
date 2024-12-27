package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Instruction struct {
	tool.Node[Instruction]

	Result    memory.Write `parser:"(@@ '=')?"`
	Operation Op           `parser:"@@ ';'"`
}

func (i Instruction) ResolveNames(
	vars names.Scope[memory.VarDecl],
	regs names.Scope[memory.RegWrite],
) (errs []error) {
	if i.Result != nil {
		errs = append(errs, i.Result.ResolveNames(vars, regs)...)
	}

	errs = append(errs, i.Operation.ResolveNames(vars, regs)...)

	return errs
}
