package debug

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Prefix struct {
	tool.Node[Prefix]

	Op Op `parser:"'debug' '.' @@"`
}
