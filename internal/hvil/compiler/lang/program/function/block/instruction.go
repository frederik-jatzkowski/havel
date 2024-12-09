package block

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Instruction struct {
	tool.Node[Instruction]

	Result    memory.Write `parser:"(@@ '=')?"`
	Operation Operation    `parser:"@@ ';'"`
}
