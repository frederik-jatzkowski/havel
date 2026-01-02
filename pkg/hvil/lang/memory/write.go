package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Write interface {
	ResolveNames(
		vars names.Scope[*stack.Decl],
		regs names.Scope[*RegWrite],
	) error
	Type() types.Type
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
}
