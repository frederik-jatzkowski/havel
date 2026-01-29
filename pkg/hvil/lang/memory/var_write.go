package memory

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type VarWrite struct {
	tool.Node[instruction.MemoryWrite]
	names.NameResolution[struct {
		Decl *stack.Decl
	}]
	registeralloc.RegisterAllocation[struct {
		Register, Temp architecture.Register
	}]

	Ident string `parser:"@Ident"`
}

func (node *VarWrite) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[*stack.Decl](ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *VarWrite) CalculateStatistics(ctx context.Context) {
	instructionID, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	blockID, err := contexttool.CurrentFromContext[statistics.BlockID](ctx)
	if err != nil {
		panic(err)
	}

	if node.NameResolutionPass.Decl.StatisticsPass.Writes == nil {
		node.NameResolutionPass.Decl.StatisticsPass.Writes = make(map[statistics.BlockID][]statistics.InstructionID)
	}

	node.NameResolutionPass.Decl.StatisticsPass.Writes[blockID] = append(
		node.NameResolutionPass.Decl.StatisticsPass.Writes[blockID],
		instructionID,
	)
}

func (node *VarWrite) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	if reg := node.NameResolutionPass.Decl.RegisterAllocationPass.BoundTo; reg != nil {
		node.RegisterAllocationPass.Register = reg

		return nil, nil
	}

	reg, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate variable store register")
	}

	temp, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate variable store tmp register")
	}

	node.RegisterAllocationPass.Register = reg
	node.RegisterAllocationPass.Temp = temp

	scope.ReturnScratchRegisters(reg)

	return nil, nil
}

func (node *VarWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if reg := node.NameResolutionPass.Decl.RegisterAllocationPass.BoundTo; reg != nil {
		return nil
	}

	p.AddI1RLit(bytecode.OPStackPtr, node.RegisterAllocationPass.Temp.(bytecode.R), uint16(node.NameResolutionPass.Decl.AddressResolutionPass.RelAddr), node.Position())

	op, err := bytecode.StoreForSize(node.NameResolutionPass.Decl.Type().Bytes())
	if err != nil {
		return node.Wrap(err)
	}

	p.AddI2R(op, node.RegisterAllocationPass.Temp.(bytecode.R), node.Register().(bytecode.R), node.Position())

	return nil
}

func (node *VarWrite) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *VarWrite) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}
