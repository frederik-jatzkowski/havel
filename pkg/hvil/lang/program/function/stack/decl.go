package stack

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"unsafe"
)

type Decl struct {
	tool.Node[Decl]
	tool.NotImplemented[Decl]

	Name         string     `parser:"@Ident"`
	DeclaredType types.Type `parser:"':' @@"`
}

func (d Decl) Identifier() string {
	return d.Name
}

func (d Decl) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return nil
}

func (d Decl) Type() types.Type {
	return d.DeclaredType
}

func (d Decl) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	return unsafe.Pointer(uintptr(1))
}
