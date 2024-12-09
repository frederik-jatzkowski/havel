package scalar

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
)

type Literal struct {
	tool.Node[Literal]

	Value uint64 `parser:"@BitLiteral"`
}
