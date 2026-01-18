package registeralloc

import "github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"

type Value interface {
	AllocateRegisters(allocator Scope) ([]architecture.Register, error)
	Register() architecture.Register
}
