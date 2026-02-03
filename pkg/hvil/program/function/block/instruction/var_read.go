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

type VarRead struct {
	tool.Node[VarRead]
	names.NameResolution[struct {
		Decl VarDecl
	}]
	registeralloc.RegisterAllocation[struct {
		Register architecture.Register
	}]

	Ident string `parser:"@Ident"`
}

func (node *VarRead) Identifier() string {
	return node.NameResolutionPass.Decl.Identifier()
}

func (node *VarRead) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[VarDecl](ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

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
	decl := node.NameResolutionPass.Decl
	if reg := decl.BoundTo(); reg != nil && !decl.Volatile() {
		node.RegisterAllocationPass.Register = reg

		return nil, nil
	}

	reg, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate variable load register")
	}

	node.RegisterAllocationPass.Register = reg

	return []architecture.Register{reg}, nil
}

func (node *VarRead) GenerateVirtualMachineAssembly(p *assembly.P) error {
	decl := node.NameResolutionPass.Decl
	if reg := decl.BoundTo(); reg != nil && !decl.Volatile() {
		return nil
	}

	p.AddI1RLit(bytecode.OPStackPtr, node.Register().(bytecode.R), uint16(node.NameResolutionPass.Decl.RelAddr()), node.Position())

	op, err := bytecode.LoadForSize(node.NameResolutionPass.Decl.Type().Bytes())
	if err != nil {
		return node.Wrap(err)
	}

	p.AddI2R(op, node.Register().(bytecode.R), node.Register().(bytecode.R), node.Position())

	return nil
}

func (node *VarRead) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *VarRead) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}
