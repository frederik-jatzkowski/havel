package registeralloc

import "github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"

type Value interface {
	AllocateRegisters(arch architecture.Architecture) ([]architecture.Register, error)
	Register() architecture.Register
}
