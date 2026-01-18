package registeralloc

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Scope interface {
	GetGeneralPurposeRegister() (architecture.Register, bool)
	ReturnGeneralPurposeRegisters(r ...architecture.Register)

	GetScratchRegister() (architecture.Register, bool)
	ReturnScratchRegisters(r ...architecture.Register)
}

type scope struct {
	generalPurposeRegisters []architecture.Register
	scratchRegisters        []architecture.Register
}

var _ Scope = (*scope)(nil)

func (s *scope) GetGeneralPurposeRegister() (architecture.Register, bool) {
	return s.pop(&s.generalPurposeRegisters)
}

func (s *scope) ReturnGeneralPurposeRegisters(r ...architecture.Register) {
	s.push(&s.generalPurposeRegisters, r)
}

func (s *scope) GetScratchRegister() (architecture.Register, bool) {
	return s.pop(&s.scratchRegisters)
}

func (s *scope) ReturnScratchRegisters(r ...architecture.Register) {
	s.push(&s.scratchRegisters, r)
}

func (s *scope) pop(registers *[]architecture.Register) (architecture.Register, bool) {
	if len(*registers) == 0 {
		return nil, false
	}

	r := (*registers)[len(*registers)-1]
	*registers = (*registers)[:len(*registers)-1]

	return r, true
}

func (s *scope) push(registers *[]architecture.Register, archRs []architecture.Register) {
	for _, archR := range archRs {
		r, ok := archR.(bytecode.R)
		if !ok {
			panic(fmt.Sprintf("invalid architecture register type: %T", archR))
		}

		*registers = append(*registers, r)
	}
}
