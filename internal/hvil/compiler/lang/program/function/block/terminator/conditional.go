package terminator

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Conditional struct {
	tool.Node[Conditional]
	tool.NotImplemented[Conditional]

	Condition memory.Read `parser:"'if':Keyword @@"`
	True      string      `parser:"'then':Keyword @Ident"`
	False     string      `parser:"'else':Keyword @Ident"`
}
