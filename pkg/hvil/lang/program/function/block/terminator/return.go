package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
)

type Return struct {
	tool.Node[Return]

	Token string `parser:"@'return':Keyword" json:"-"`
}
