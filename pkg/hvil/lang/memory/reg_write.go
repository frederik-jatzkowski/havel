package memory

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type RegWrite struct {
	tool.Node[RegWrite]
	address.Resolution[struct {
		RelAddr int
	}]
	registeralloc.RegisterAllocation[struct {
		Register architecture.Register
	}]

	Ident   string     `parser:"'$' @Ident"`
	RegType types.Type `parser:"':' @@"`
}

var _ Write = (*RegWrite)(nil)

func (node *RegWrite) Identifier() string {
	return node.Ident
}

func (node *RegWrite) ResolveNames(ctx context.Context) error {
	if err := DefineRegisterInCtx(ctx, node); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *RegWrite) GenerateVirtualMachineAssembly(p *assembly.P) error {
	return nil
}

func (node *RegWrite) Register() architecture.Register {
	return node.RegisterAllocationPass.Register
}

func (node *RegWrite) Type() types.Type {
	return node.RegType
}

func (node *RegWrite) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	stackAddr := vm.StackPointer + node.AddressResolutionPass.RelAddr
	return unsafe.Pointer(&vm.Stack[stackAddr])
}
