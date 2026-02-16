package instruction

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Instruction struct {
	tool.Node[Instruction]
	statistics.Statistics[struct {
		InstructionID statistics.InstructionID
	}]

	ResultWrite MemoryWrite `parser:"(@@ '=')?"`
	Operation   Operation   `parser:"@@ ';'"`
}

func (node *Instruction) ResolveNames(ctx context.Context) error {
	ctx = contexttool.WithCurrent(ctx, node)

	if node.ResultWrite != nil {
		if err := node.ResultWrite.ResolveNames(ctx); err != nil {
			return err
		}
	}

	return node.Operation.ResolveNames(ctx)
}

func (node *Instruction) ResolveTypes() error {
	if node.ResultWrite != nil {
		target := node.ResultWrite.Type()

		if _, ok := target.(*types.Composite); ok {
			return node.Errorf("cannot have composite type %s in register", target)
		}

		return node.Operation.ResolveTypes(target)
	}

	return node.Operation.ResolveTypes(&types.Void{})
}

func (node *Instruction) CalculateStatistics(ctx context.Context) {
	id, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	node.StatisticsPass.InstructionID = id

	if node.ResultWrite != nil {
		node.ResultWrite.CalculateStatistics(ctx)
	}

	node.Operation.CalculateStatistics(ctx)
}

func (node *Instruction) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	scope.SetInstructionID(node.StatisticsPass.InstructionID)

	regs, err := node.Operation.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	defer scope.ReturnScratchRegisters(regs...)

	if node.ResultWrite == nil {
		return nil, nil
	}

	resultRegs, err := node.ResultWrite.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	node.Operation.SetResultRegister(node.ResultWrite.Register())

	return resultRegs, nil
}

func (node *Instruction) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Operation.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if node.ResultWrite != nil {
		if err := node.ResultWrite.GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return nil
}
