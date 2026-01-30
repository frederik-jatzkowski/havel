package registeralloc

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
)

type Scope interface {
	Architecture() architecture.Architecture

	SetInstructionID(id statistics.InstructionID)
	GetInstructionID() statistics.InstructionID

	GetGeneralPurposeRegister() (architecture.Register, bool)
	ReturnGeneralPurposeRegisters(r ...architecture.Register)

	GetScratchRegister() (architecture.Register, bool)
	ReturnScratchRegisters(r ...architecture.Register)

	UseRegisters(r ...architecture.Register)
	IsLiveAt(r architecture.Register, id statistics.InstructionID) bool
}

type scope struct {
	arch                    architecture.Architecture
	instructionID           statistics.InstructionID
	generalPurposeRegisters []architecture.Register
	liveRanges              map[architecture.Register][]Range
	scratchRegisters        []architecture.Register
}

func newScope(arch architecture.Architecture) Scope {
	s := &scope{
		arch:                    arch,
		generalPurposeRegisters: arch.GeneralPurposeRegisters(),
		liveRanges:              make(map[architecture.Register][]Range),
		scratchRegisters:        arch.ScratchRegisters(),
	}

	for _, register := range s.generalPurposeRegisters {
		s.liveRanges[register] = []Range{}
	}

	return s
}

var _ Scope = (*scope)(nil)

func (s *scope) Architecture() architecture.Architecture {
	return s.arch
}

func (s *scope) SetInstructionID(id statistics.InstructionID) {
	s.instructionID = id
}

func (s *scope) GetInstructionID() statistics.InstructionID {
	return s.instructionID
}

func (s *scope) GetGeneralPurposeRegister() (architecture.Register, bool) {
	r, ok := s.pop(&s.generalPurposeRegisters)
	if !ok {
		return nil, false
	}

	s.liveRanges[r] = append(s.liveRanges[r], Range{
		Start: s.instructionID,
		End:   s.instructionID,
	})

	return r, true
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

func (s *scope) UseRegisters(rs ...architecture.Register) {
	for _, r := range rs {
		ranges := s.liveRanges[r]
		ranges[len(ranges)-1].Start = min(ranges[len(ranges)-1].Start, s.instructionID)
		ranges[len(ranges)-1].End = max(ranges[len(ranges)-1].End, s.instructionID)
	}
}

func (s *scope) IsLiveAt(r architecture.Register, id statistics.InstructionID) bool {
	for _, lr := range s.liveRanges[r] {
		if lr.Start < id && lr.End > id {
			return true
		}
	}

	return false
}
