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

type VarRead struct {
	tool.Node[VarRead]
	names.NameResolution[struct {
		Decl *stack.Decl
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

func (node *VarRead) GenerateVirtualMachineAssembly(p *assembly.P) error {
	//TODO implement me
	panic("implement me")
}

func (node *VarRead) Register() architecture.Register {
	//TODO implement me
	panic("implement me")
}

func (node *VarRead) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}

func (node *VarRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
