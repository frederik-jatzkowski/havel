package instruction

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type VarRead struct {
	tool.Node[VarRead]
	names.NameResolution[struct {
		Decl   VarDecl
		Type   types.Type
		Offset uint
	}]
	registeralloc.RegisterAllocation[struct {
		Register architecture.Register
	}]

	Ident        string `parser:"@Ident"`
	Dereferences []uint `parser:"( '[' @Number ']' )*"`
}

func (node *VarRead) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[VarDecl](ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	dereferenced, offset, err := node.NameResolutionPass.Decl.Type().Dereference(node.Dereferences)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Type = dereferenced
	node.NameResolutionPass.Offset = offset

	return nil
}

func (node *VarRead) CalculateStatistics(ctx context.Context) {
	instructionID, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	blockID, err := contexttool.CurrentFromContext[statistics.BlockID](ctx)
	if err != nil {
		panic(err)
	}

	node.NameResolutionPass.Decl.AddReadToStatistic(blockID, instructionID)
}

func (node *VarRead) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	reg, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate variable load register")
	}

	node.RegisterAllocationPass.Register = reg

	return []architecture.Register{reg}, nil
}

func (node *VarRead) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.NameResolutionPass.Decl.AddBytecodeVirtualmachinePtrInstruction(
		p,
		node.Register().(bytecode.R),
		node.Dereferences,
	); err != nil {
		return node.Wrap(err)
	}

	if node.NameResolutionPass.Type.Bytes() <= 8 {
		op, err := bytecode.LoadForSize(node.NameResolutionPass.Type.Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, node.Register().(bytecode.R), node.Register().(bytecode.R), node.Position())
	}

	return nil
}

func (node *VarRead) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *VarRead) Type() types.Type {
	return node.NameResolutionPass.Type
}
