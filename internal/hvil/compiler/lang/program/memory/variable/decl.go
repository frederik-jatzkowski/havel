package variable

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/memory/types"
)

type Decl interface {
	Type() types.Type
}
