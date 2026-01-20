package memory

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
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

func (node *RegRead) CalculateStatistics() {

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
		var op bytecode.OP
		switch node.NameResolutionPass.Decl.Type().Bytes() {
		case 1:
			op = bytecode.OPLoadStack8
		case 2:
			op = bytecode.OPLoadStack16
		case 4:
			op = bytecode.OPLoadStack32
		case 8:
			op = bytecode.OPLoadStack64
		}

		p.AddI1RLit(op, node.Register().(bytecode.R), uint16(node.NameResolutionPass.Decl.AddressResolutionPass.RelAddr), node.Position())
	}

	return nil
}

func (node *RegRead) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *RegRead) Type() types.Type {
	return node.NameResolutionPass.Decl.RegType
}
