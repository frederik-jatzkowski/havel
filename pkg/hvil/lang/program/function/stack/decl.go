package stack

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
)

type Decl struct {
	tool.Node[Decl]
	address.Resolution[struct {
		RelAddr int
	}]

	Name         string     `parser:"@Ident"`
	DeclaredType types.Type `parser:"':' @@"`
}

func (d *Decl) Identifier() string {
	return d.Name
}

func (d *Decl) Type() types.Type {
	return d.DeclaredType
}

func (d *Decl) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	stackAddr := vm.StackPointer + d.AddressResolutionPass.RelAddr
	return unsafe.Pointer(&vm.Stack[stackAddr])
}
