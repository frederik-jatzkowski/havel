package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"unsafe"
)

type Op interface {
	ResolveNames(
		vars names.Scope[memory.VarDecl],
		regs names.Scope[memory.RegWrite],
	) (errs []error)
	ResolveTypes(expected types.Type) (errs []error)
	Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error
}
