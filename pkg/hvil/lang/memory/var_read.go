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

func (read *VarRead) Identifier() string {
	return read.NameResolutionPass.Decl.Identifier()
}

func (read *VarRead) ResolveNames(vars names.Scope[*stack.Decl], _ names.Scope[*RegWrite]) (errs []error) {
	decl, err := vars.Find(read.Ident)
	if err != nil {
		return append(errs, read.Wrap(err))
	}

	read.NameResolutionPass.Decl = decl

	return nil
}

func (read *VarRead) Type() types.Type {
	return read.NameResolutionPass.Decl.Type()
}

func (read *VarRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return read.NameResolutionPass.Decl.Addr(vm)
}
