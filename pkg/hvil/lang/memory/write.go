package memory

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"unsafe"
)

type Write interface {
	ResolveNames(
		vars names.Scope[VarDecl],
		regs names.Scope[RegWrite],
	) (errs []error)
	Type() types.Type
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
}
