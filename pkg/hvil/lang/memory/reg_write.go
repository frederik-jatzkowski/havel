package memory

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type RegWrite struct {
	tool.Node[RegWrite]
	address.Resolution[struct {
		RelAddr int
	}]
	registeralloc.RegisterAllocation[struct {
		Register architecture.Register
		Spilled  bool
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

func (node *RegWrite) CalculateStatistics() {

}

func (node *RegWrite) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	reg, ok := scope.GetGeneralPurposeRegister()
	if !ok {
		node.RegisterAllocationPass.Spilled = true
		reg, ok := scope.GetScratchRegister()
		if !ok {
			return nil, node.Errorf("cannot allocate spill register")
		}

		node.RegisterAllocationPass.Register = reg
		scope.ReturnScratchRegisters(reg)

		return nil, nil
	}

	node.RegisterAllocationPass.Register = reg

	return []architecture.Register{reg}, nil
}

func (node *RegWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if node.RegisterAllocationPass.Spilled {
		var op bytecode.OP
		switch node.RegType.Bytes() {
		case 1:
			op = bytecode.OPStoreStack8
		case 2:
			op = bytecode.OPStoreStack16
		case 4:
			op = bytecode.OPStoreStack32
		case 8:
			op = bytecode.OPStoreStack64
		}

		p.AddI1RLit(op, node.Register().(bytecode.R), uint16(node.AddressResolutionPass.RelAddr), node.Position())
	}

	return nil
}

func (node *RegWrite) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *RegWrite) Type() types.Type {
	return node.RegType
}
