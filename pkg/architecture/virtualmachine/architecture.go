package virtualmachine

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
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

func (a *Architecture) InitialStackOffset() int {
	return 16 // 8 bytes reserved for return address, 8 byte for return stack pointer
}

func (a *Architecture) ArgRegisters() []architecture.Register {
	return append([]architecture.Register(nil), a.argRegisters...)
}

func (a *Architecture) GeneralPurposeRegisters() []architecture.Register {
	return append([]architecture.Register(nil), a.generalPurposeRegisters...)
}

func (a *Architecture) ScratchRegisters() []architecture.Register {
	return append([]architecture.Register(nil), a.scratchRegisters...)
}

func (a *Architecture) CalculateCallPlan(signature *types.Function) architecture.CallPlan {
	call := architecture.CallPlan{}
	argRegs := a.ArgRegisters()
	offset := a.InitialStackOffset()
	for i, param := range signature.Parameters.Items {
		if i > len(argRegs)-1 {
			call.Params = append(call.Params, architecture.MemoryAllocation{
				RelAddr: offset,
				Bytes:   param.Bytes(),
			})
		} else {
			r := argRegs[i]
			call.Params = append(call.Params, architecture.MemoryAllocation{
				BoundTo: r,
				RelAddr: offset,
				Bytes:   param.Bytes(),
			})
		}

		offset += param.Bytes()
	}

	call.Result = architecture.MemoryAllocation{
		BoundTo: a.resultRegisters[0],
		RelAddr: offset,
		Bytes:   signature.ReturnType().Bytes(),
	}

	offset += call.Result.Bytes

	call.Offset = offset

	return call
}
