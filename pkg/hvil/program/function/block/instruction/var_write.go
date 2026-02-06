package instruction

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type VarWrite struct {
	tool.Node[MemoryWrite]
	names.NameResolution[struct {
		Decl VarDecl
	}]
	registeralloc.RegisterAllocation[struct {
		Register, Temp architecture.Register
	}]

	Ident string `parser:"@Ident"`
}

func (node *VarWrite) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[VarDecl](ctx, node.Ident)
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

	node.NameResolutionPass.Decl.AddWriteToStatistic(blockID, instructionID)
}

func (node *VarWrite) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	decl := node.NameResolutionPass.Decl
	if reg := decl.BoundTo(); reg != nil && !decl.Volatile() {
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

	scope.ReturnScratchRegisters(reg, temp)

	return nil, nil
}

func (node *VarWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	decl := node.NameResolutionPass.Decl
	if reg := decl.BoundTo(); reg != nil && !decl.Volatile() {
		return nil
	}

	node.NameResolutionPass.Decl.AddBytecodeVirtualmachinePtrInstruction(p, node.RegisterAllocationPass.Temp.(bytecode.R))

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
