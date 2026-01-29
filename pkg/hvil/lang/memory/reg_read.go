package memory

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type RegRead struct {
	tool.Node[RegRead]
	names.NameResolution[struct {
		Decl *RegWrite
	}]
	registeralloc.RegisterAllocation[struct {
		Register architecture.Register
	}]

	Ident string `parser:"'$' @Ident"`
}

func (node *RegRead) Identifier() string {
	return node.NameResolutionPass.Decl.Identifier()
}

func (node *RegRead) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[*RegWrite](ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *RegRead) CalculateStatistics(ctx context.Context) {
	id, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	node.NameResolutionPass.Decl.StatisticsPass.Reads = append(node.NameResolutionPass.Decl.StatisticsPass.Reads, id)
}

func (node *RegRead) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	if node.NameResolutionPass.Decl.RegisterAllocationPass.Spilled {
		reg, ok := scope.GetScratchRegister()
		if !ok {
			return nil, node.Errorf("cannot allocate spill register")
		}

		node.RegisterAllocationPass.Register = reg

		return []architecture.Register{reg}, nil
	}

	node.RegisterAllocationPass.Register = node.NameResolutionPass.Decl.RegisterAllocationPass.Register

	scope.UseRegisters(node.RegisterAllocationPass.Register)

	return nil, nil
}

func (node *RegRead) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if node.NameResolutionPass.Decl.RegisterAllocationPass.Spilled {
		p.AddI1RLit(bytecode.OPStackPtr, node.Register().(bytecode.R), uint16(node.NameResolutionPass.Decl.AddressResolutionPass.RelAddr), node.Position())

		op, err := bytecode.LoadForSize(node.NameResolutionPass.Decl.Type().Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, node.Register().(bytecode.R), node.Register().(bytecode.R), node.Position())
	}

	return nil
}

func (node *RegRead) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *RegRead) Type() types.Type {
	return node.NameResolutionPass.Decl.RegType
}
