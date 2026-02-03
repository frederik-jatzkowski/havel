package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
)

type VarDecl interface {
	names.ScopedObject

	Type() types.Type
	AddReadToStatistic(blockID statistics.BlockID, instructionID statistics.InstructionID)
	AddWriteToStatistic(blockID statistics.BlockID, instructionID statistics.InstructionID)
	SetPtrTaken()
	BoundTo() architecture.Register
	Volatile() bool
	RelAddr() int
}
