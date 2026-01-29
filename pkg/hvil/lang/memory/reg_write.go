package memory

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type RegWrite struct {
	tool.Node[RegWrite]
	statistics.Statistics[struct {
		Reads []statistics.InstructionID
		Decl  statistics.InstructionID
	}]
	address.Resolution[struct {
		RelAddr int
	}]
	registeralloc.RegisterAllocation[struct {
		Register, Temp architecture.Register
		Spilled        bool
	}]

	Ident   string     `parser:"'$' @Ident"`
	RegType types.Type `parser:"':' @@"`
}

var _ instruction.MemoryWrite = (*RegWrite)(nil)

func (node *RegWrite) Identifier() string {
	return node.Ident
}

func (node *RegWrite) ResolveNames(ctx context.Context) error {
	if err := contexttool.DefineInScope(ctx, node); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *RegWrite) CalculateStatistics(ctx context.Context) {
	id, err := contexttool.CurrentFromContext[statistics.InstructionID](ctx)
	if err != nil {
		panic(err)
	}

	node.StatisticsPass.Decl = id
}

func (node *RegWrite) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	reg, ok := scope.GetGeneralPurposeRegister()
	if !ok {
		node.RegisterAllocationPass.Spilled = true
		reg, ok := scope.GetScratchRegister()
		if !ok {
			return nil, node.Errorf("cannot allocate spill register")
		}

		temp, ok := scope.GetScratchRegister()
		if !ok {
			return nil, node.Errorf("cannot allocate spill register")
		}

		node.RegisterAllocationPass.Register = reg
		node.RegisterAllocationPass.Temp = temp

		scope.ReturnScratchRegisters(reg, temp)

		return nil, nil
	}

	node.RegisterAllocationPass.Register = reg

	return []architecture.Register{reg}, nil
}

func (node *RegWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if node.RegisterAllocationPass.Spilled {
		temp := node.RegisterAllocationPass.Temp.(bytecode.R)
		p.AddI1RLit(bytecode.OPStackPtr, temp, uint16(node.AddressResolutionPass.RelAddr), node.Position())

		op, err := bytecode.StoreForSize(node.RegType.Bytes())
		if err != nil {
			return node.Wrap(err)
		}

		p.AddI2R(op, temp, node.Register().(bytecode.R), node.Position())
	}

	return nil
}

func (node *RegWrite) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *RegWrite) Type() types.Type {
	return node.RegType
}
