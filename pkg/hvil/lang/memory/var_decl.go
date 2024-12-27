package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type VarDecl interface {
	names.ScopedObject
	Type() types.Type
}
