package stack

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
)

type Decl struct {
	tool.Node[Decl]
	tool.NotImplemented[Decl]

	Name         string     `parser:"@Ident"`
	DeclaredType types.Type `parser:"':' @@"`
}

func (d Decl) Identifier() string {
	return d.Name
}

func (d Decl) Type() types.Type {
	return d.DeclaredType
}
