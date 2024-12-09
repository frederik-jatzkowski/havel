package terminator

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Jump struct {
	tool.Node[Jump]
	tool.NotImplemented[Jump]

	Target string `parser:"@Ident"`
}
