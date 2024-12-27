package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type Conditional struct {
	tool.Node[Conditional]
	tool.NotImplemented[Conditional]

	Condition memory.Read `parser:"'if':Keyword @@"`
	True      string      `parser:"'then':Keyword @Ident"`
	False     string      `parser:"'else':Keyword @Ident"`
}
