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
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type VarWrite struct {
	tool.Node[Write]
	names.NameResolution[struct {
		Decl *stack.Decl
	}]

	Ident string `parser:"@Ident"`
}

func (node *VarWrite) ResolveNames(ctx context.Context) error {
	decl, err := stack.FromCtx(ctx, node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *VarWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	//TODO implement me
	panic("implement me")
}

func (node *VarWrite) Register() architecture.Register {
	//TODO implement me
	panic("implement me")
}

func (node *VarWrite) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}

func (node *VarWrite) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
