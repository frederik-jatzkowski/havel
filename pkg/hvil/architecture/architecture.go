package architecture

import "github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"

type Architecture interface {
	CalculateCallPlan(signature *types.FunctionType) CallPlan
	GeneralPurposeRegisters() []Register
	ScratchRegisters() []Register
}
