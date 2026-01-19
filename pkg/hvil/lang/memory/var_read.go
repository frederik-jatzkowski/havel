package memory

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type VarRead struct {
	tool.Node[VarRead]
	names.NameResolution[struct {
		Decl *stack.Decl
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
	decl, err := stack.FromCtx(ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *VarRead) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	node.NameResolutionPass.Decl.RegisterAllocationPass.Usages++

	if reg := node.NameResolutionPass.Decl.RegisterAllocationPass.BoundTo; reg != nil {
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
	if reg := node.NameResolutionPass.Decl.RegisterAllocationPass.BoundTo; reg != nil {
		return nil
	}

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

	return nil
}

func (node *VarRead) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *VarRead) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}

func (node *VarRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
