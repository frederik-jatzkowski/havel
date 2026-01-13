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
	for i := range 32 {
		r := bytecode.R(i)
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

func (a *Architecture) GetArgRegister() (architecture.Register, bool) {
	return a.pop(&a.argRegisters)
}

func (a *Architecture) GetResultRegister() (architecture.Register, bool) {
	return a.pop(&a.resultRegisters)
}

func (a *Architecture) GetGeneralPurposeRegister() (architecture.Register, bool) {
	return a.pop(&a.generalPurposeRegisters)
}

func (a *Architecture) ReturnGeneralPurposeRegisters(r ...architecture.Register) {
	a.push(&a.generalPurposeRegisters, r)
}

func (a *Architecture) GetScratchRegister() (architecture.Register, bool) {
	return a.pop(&a.scratchRegisters)
}

func (a *Architecture) ReturnScratchRegisters(r ...architecture.Register) {
	a.push(&a.scratchRegisters, r)
}

func (a *Architecture) pop(registers *[]bytecode.R) (architecture.Register, bool) {
	if len(*registers) == 0 {
		return nil, false
	}

	r := (*registers)[len(*registers)-1]
	*registers = (*registers)[:len(*registers)-1]

	return r, true
}

func (a *Architecture) push(registers *[]bytecode.R, archRs []architecture.Register) {
	for _, archR := range archRs {
		r, ok := archR.(bytecode.R)
		if !ok {
			panic(fmt.Sprintf("invalid architecture register type: %T", archR))
		}

		*registers = append(*registers, r)
	}
}
