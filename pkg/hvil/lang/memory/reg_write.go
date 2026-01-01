package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type RegWrite struct {
	tool.Node[RegWrite]
	address.Resolution[struct {
		RelAddr int
	}]

	Ident   string     `parser:"'$' @Ident"`
	RegType types.Type `parser:"':' @@"`
}

var _ Write = (*RegWrite)(nil)

func (w *RegWrite) Identifier() string {
	return w.Ident
}

func (w *RegWrite) ResolveNames(
	_ names.Scope[*stack.Decl],
	regs names.Scope[*RegWrite],
) (errs []error) {
	err := regs.Define(w)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (w *RegWrite) Type() types.Type {
	return w.RegType
}

func (w *RegWrite) Addr(vm *runtime.VirtualMachine) unsafe.Pointer {
	stackAddr := vm.StackPointer + w.AddressResolutionPass.RelAddr
	return unsafe.Pointer(&vm.Stack[stackAddr])
}
