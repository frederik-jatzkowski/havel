package register

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Write struct {
	tool.Node[Write]

	Ident string     `parser:"'$' @Ident"`
	Type  types.Type `parser:"':' @@"`
}
