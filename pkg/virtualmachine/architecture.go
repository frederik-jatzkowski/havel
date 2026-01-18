package virtualmachine

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Architecture struct {
	argRegisters            []architecture.Register
	resultRegisters         []architecture.Register
	generalPurposeRegisters []architecture.Register
	scratchRegisters        []architecture.Register
}

func NewArchitecture() *Architecture {
	arch := &Architecture{}
	for i := range 32 {
		r := architecture.Register(bytecode.R(i))
		switch {
		case i < 2:
			// reserved for pc and sp
		case i < 8:
			arch.argRegisters = append(arch.argRegisters, r)
		case i < 16:
			arch.resultRegisters = append(arch.resultRegisters, r)
		case i < 24:
			arch.generalPurposeRegisters = append(arch.generalPurposeRegisters, r)
		case i < 32:
			arch.scratchRegisters = append(arch.scratchRegisters, r)
		}
	}

	return arch
}

var _ architecture.Architecture = &Architecture{}

func (a *Architecture) GeneralPurposeRegisters() []architecture.Register {
	return append([]architecture.Register(nil), a.generalPurposeRegisters...)
}

func (a *Architecture) ScratchRegisters() []architecture.Register {
	return append([]architecture.Register(nil), a.scratchRegisters...)
}
