package instruction

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Operation interface {
	ResolveNames(
		vars names.Scope[*stack.Decl],
		regs names.Scope[*memory.RegWrite],
	) error
	ResolveTypes(expected types.Type) error
	Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error
}
