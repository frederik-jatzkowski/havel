package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type Jump struct {
	tool.Node[Jump]
	tool.NotImplemented[Jump]

	Target string `parser:"@Ident"`
}
