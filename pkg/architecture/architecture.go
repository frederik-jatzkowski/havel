package architecture

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
)

type Architecture interface {
	CalculateCallPlan(signature *types.Function) CallPlan
	GeneralPurposeRegisters() []Register
	ScratchRegisters() []Register
}
