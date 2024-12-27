package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type VarWrite struct {
	tool.Node[Write]
	names.NameResolution[struct {
		Decl VarDecl
	}]
	tool.NotImplemented[Write]

	Ident string `parser:"@Ident"`
}

func (w VarWrite) ResolveNames(
	_ names.Scope[VarDecl],
	_ names.Scope[RegWrite],
) (errs []error) {
	return nil
}

func (w VarWrite) Type() types.Type {
	return w.NameResolutionPass.Decl.Type()
}
