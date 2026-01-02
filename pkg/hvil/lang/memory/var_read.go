package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
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

func (node *VarRead) ResolveNames(vars names.Scope[*stack.Decl], _ names.Scope[*RegWrite]) error {
	decl, err := vars.Find(node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *VarRead) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}

func (node *VarRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
