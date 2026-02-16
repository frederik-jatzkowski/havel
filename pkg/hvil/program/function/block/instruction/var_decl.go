package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type VarDecl interface {
	scope.Object

	Type() types.Type
	AddReadToStatistic(blockID statistics.BlockID, instructionID statistics.InstructionID)
	AddWriteToStatistic(blockID statistics.BlockID, instructionID statistics.InstructionID)
	SetPtrTaken()
	AddBytecodeVirtualmachinePtrInstruction(p *assembly.P, target bytecode.R, dereferences []uint) (err error)
}
