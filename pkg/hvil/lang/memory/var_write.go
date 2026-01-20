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
		Register architecture.Register
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

func (node *VarWrite) CalculateStatistics() {

}

func (node *VarWrite) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	node.NameResolutionPass.Decl.RegisterAllocationPass.Usages++

	if reg := node.NameResolutionPass.Decl.RegisterAllocationPass.BoundTo; reg != nil {
		node.RegisterAllocationPass.Register = reg

		return nil, nil
	}

	reg, ok := scope.GetScratchRegister()
	if !ok {
		return nil, node.Errorf("cannot allocate variable store register")
	}

	node.RegisterAllocationPass.Register = reg

	scope.ReturnScratchRegisters(reg)

	return nil, nil
}

func (node *VarWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if reg := node.NameResolutionPass.Decl.RegisterAllocationPass.BoundTo; reg != nil {
		return nil
	}

	var op bytecode.OP
	switch node.NameResolutionPass.Decl.Type().Bytes() {
	case 1:
		op = bytecode.OPStoreStack8
	case 2:
		op = bytecode.OPStoreStack16
	case 4:
		op = bytecode.OPStoreStack32
	case 8:
		op = bytecode.OPStoreStack64
	}

	p.AddI1RLit(op, node.Register().(bytecode.R), uint16(node.NameResolutionPass.Decl.AddressResolutionPass.RelAddr), node.Position())

	return nil
}

func (node *VarWrite) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *VarWrite) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}
