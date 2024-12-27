package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Write interface {
	ResolveNames(
		vars names.Scope[VarDecl],
		regs names.Scope[RegWrite],
	) (errs []error)
}
