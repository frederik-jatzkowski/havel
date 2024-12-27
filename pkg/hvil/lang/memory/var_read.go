package memory

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type VarRead struct {
	tool.Node[VarRead]
	names.NameResolution[struct {
		Decl VarDecl
	}]
	tool.NotImplemented[VarRead]

	Ident string `parser:"@Ident"`
}

func (read *VarRead) ResolveNames(vars names.Scope[VarDecl], _ names.Scope[RegWrite]) (errs []error) {
	decl, exists := vars.Find(read.Ident)
	if !exists {
		return append(errs, read.Errorf("register '%s' is not defined", read.Ident))
	}

	read.NameResolutionPass.Decl = *decl

	return nil
}

func (read *VarRead) Position() lexer.Position {
	return read.Pos
}
