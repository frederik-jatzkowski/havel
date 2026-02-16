package registeralloc

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture"
)

type Value interface {
	AllocateRegisters(allocator Scope) ([]architecture.Register, error)
	Register() architecture.Register
}
