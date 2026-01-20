package stack

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Decl struct {
	tool.Node[Decl]
	address.Resolution[struct {
		RelAddr int
	}]
	registeralloc.RegisterAllocation[struct {
		Usages  int
		BoundTo architecture.Register
	}]

	Name         string     `parser:"@Ident"`
	DeclaredType types.Type `parser:"':' @@"`
}

func (node *Decl) Identifier() string {
	return node.Name
}

func (node *Decl) Type() types.Type {
	return node.DeclaredType
}

func (node *Decl) CalculateStatistics() {

}

func (node *Decl) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	stackAddr := vm.StackPointer + node.AddressResolutionPass.RelAddr
	return unsafe.Pointer(&vm.Stack[stackAddr])
}
