package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"unsafe"
)

type VarRead struct {
	tool.Node[VarRead]
	names.NameResolution[struct {
		Decl VarDecl
	}]
	tool.NotImplemented[VarRead]

	Ident string `parser:"@Ident"`
}

func (read *VarRead) Identifier() string {
	return read.NameResolutionPass.Decl.Identifier()
}

func (read *VarRead) ResolveNames(vars names.Scope[VarDecl], _ names.Scope[RegWrite]) (errs []error) {
	decl, exists := vars.Find(read.Ident)
	if !exists {
		return append(errs, read.Errorf("register '%s' is not defined", read.Ident))
	}

	read.NameResolutionPass.Decl = *decl

	return nil
}

func (read *VarRead) Type() types.Type {
	return read.NameResolutionPass.Decl.Type()
}

func (read *VarRead) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return read.NameResolutionPass.Decl.Addr(vm)
}
