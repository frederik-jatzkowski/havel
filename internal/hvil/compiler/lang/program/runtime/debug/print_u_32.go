package debug

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type PrintU32 struct {
	tool.Node[PrintU32]

	Param memory.Read `parser:"'print_u_32' '(' @@ ')'"`
}
