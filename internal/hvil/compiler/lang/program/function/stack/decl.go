package stack

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/variable"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

var _ variable.Decl = (*Decl)(nil)

type Decl struct {
	tool.Node[Decl]
	tool.NotImplemented[Decl]

	Name         string     `parser:"@Ident"`
	DeclaredType types.Type `parser:"':' @@"`
}

func (Decl *Decl) Type() types.Type {
	return Decl.DeclaredType
}
