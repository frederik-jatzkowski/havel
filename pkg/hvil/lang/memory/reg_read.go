package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type RegRead struct {
	tool.Node[RegRead]
	names.NameResolution[struct {
		Decl *RegWrite
	}]

	Ident string `parser:"'$' @Ident"`
}

func (read *RegRead) ResolveNames(_ names.Scope[VarDecl], regs names.Scope[RegWrite]) (errs []error) {
	decl, exists := regs.Find(read.Ident)
	if !exists {
		return append(errs, read.Errorf("register '%s' is not defined", read.Ident))
	}

	read.NameResolutionPass.Decl = decl

	return nil
}

func (read *RegRead) Type() types.Type {
	return read.NameResolutionPass.Decl.Type
}
