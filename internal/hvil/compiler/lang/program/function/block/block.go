package block

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Block struct {
	tool.Node[Block]

	Ident        string         `parser:"'block':Keyword @Ident '{'"`
	Instructions []*Instruction `parser:"@@*"`
	Terminator   Terminator     `parser:"'}' '=>' @@ ';'"`
}
