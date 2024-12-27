package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type RegWrite struct {
	tool.Node[RegWrite]

	Ident string     `parser:"'$' @Ident"`
	Type  types.Type `parser:"':' @@"`
}

func (w RegWrite) Identifier() string {
	return w.Ident
}

func (w RegWrite) ResolveNames(
	_ names.Scope[VarDecl],
	regs names.Scope[RegWrite],
) (errs []error) {
	err := regs.Define(&w)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
