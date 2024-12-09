package register

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass"
)

type Read struct {
	tool.Node[Read]
	pass.NameResolution[struct {
		Decl *Write
	}]

	Ident string `parser:"'$' @Ident"`
}

func (read *Read) Type() types.Type {
	return read.NameResolutionPass.Decl.Type
}
