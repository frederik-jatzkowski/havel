package alu

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Operation struct {
	tool.Node[Operation]
	tool.NotImplemented[Operation]

	Name string                 `parser:"'alu' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}
