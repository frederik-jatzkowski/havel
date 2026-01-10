package virtualmachine

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Architecture struct {
	argRegisters            []bytecode.R
	resultRegisters         []bytecode.R
	generalPurposeRegisters []bytecode.R
	scratchRegisters        []bytecode.R
}

func NewArchitecture() *Architecture {
	arch := &Architecture{}
	for i := range 256 {
		r := bytecode.R(i)
		switch {
		case i < 8:
			arch.argRegisters = append(arch.argRegisters, r)
		case i <= 16:
			arch.resultRegisters = append(arch.resultRegisters, r)
		case i <= 248:
			arch.generalPurposeRegisters = append(arch.generalPurposeRegisters, r)
		default:
			arch.scratchRegisters = append(arch.scratchRegisters, r)
		}
	}

	return arch
}

var _ architecture.Architecture = &Architecture{}

func (a *Architecture) GetArgRegister() (architecture.Register, bool) {
	return a.pop(&a.argRegisters)
}

func (a *Architecture) GetResultRegister() (architecture.Register, bool) {
	return a.pop(&a.resultRegisters)
}

func (a *Architecture) GetGeneralPurposeRegister() (architecture.Register, bool) {
	return a.pop(&a.generalPurposeRegisters)
}

func (a *Architecture) GetScratchRegister() (architecture.Register, bool) {
	return a.pop(&a.scratchRegisters)
}

func (a *Architecture) pop(registers *[]bytecode.R) (architecture.Register, bool) {
	if len(*registers) == 0 {
		return nil, false
	}

	r := (*registers)[len(*registers)-1]
	*registers = (*registers)[:len(*registers)-1]

	return r, true
}

func (a *Architecture) push(registers *[]bytecode.R, archR architecture.Register) {
	r, ok := archR.(bytecode.R)
	if !ok {
		panic(fmt.Sprintf("invalid architecture register type: %T", archR))
	}

	*registers = append(*registers, r)
}
