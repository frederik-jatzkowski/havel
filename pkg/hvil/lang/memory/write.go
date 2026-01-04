package memory

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Write interface {
	names.Resolver
	Type() types.Type
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
}
