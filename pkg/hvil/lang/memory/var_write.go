package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type VarWrite struct {
	tool.Node[Write]
	names.NameResolution[struct {
		Decl *stack.Decl
	}]

	Ident string `parser:"@Ident"`
}

func (node *VarWrite) ResolveNames(
	vars names.Scope[*stack.Decl],
	_ names.Scope[*RegWrite],
) (errs []error) {
	decl, err := vars.Find(node.Ident)
	if err != nil {
		return append(errs, node.Wrap(err))
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *VarWrite) Type() types.Type {
	return node.NameResolutionPass.Decl.Type()
}

func (node *VarWrite) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
