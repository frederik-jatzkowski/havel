package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type RegRead struct {
	tool.Node[RegRead]
	names.NameResolution[struct {
		Decl *RegWrite
	}]

	Ident string `parser:"'$' @Ident"`
}

func (node *RegRead) Identifier() string {
	return node.NameResolutionPass.Decl.Identifier()
}

func (node *RegRead) ResolveNames(_ names.Scope[*stack.Decl], regs names.Scope[*RegWrite]) error {
	decl, err := regs.Find(node.Ident)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Decl = decl

	return nil
}

func (node *RegRead) Type() types.Type {
	return node.NameResolutionPass.Decl.RegType
}

func (node *RegRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return node.NameResolutionPass.Decl.Addr(vm)
}
