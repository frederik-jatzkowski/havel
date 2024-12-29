package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/address"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"unsafe"
)

type RegWrite struct {
	tool.Node[RegWrite]
	address.Resolution[struct {
		RelAddr int
	}]

	Ident   string     `parser:"'$' @Ident"`
	RegType types.Type `parser:"':' @@"`
}

func (w RegWrite) Identifier() string {
	return w.Ident
}

func (w *RegWrite) ResolveNames(
	_ names.Scope[VarDecl],
	regs names.Scope[RegWrite],
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
