package variable

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass"
)

type Write struct {
	tool.Node[Write]
	pass.NameResolution[struct {
		Decl Decl
	}]
	tool.NotImplemented[Write]

	Ident string `parser:"@Ident"`
}
