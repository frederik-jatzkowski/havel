package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type VarDecl interface {
	names.ScopedObject
}
