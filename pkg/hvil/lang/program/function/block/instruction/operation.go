package instruction

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Operation interface {
	names.Resolver
	ResolveTypes(expected types.Type) error
	SetResultRegister(r architecture.Register)
	codegen.VirtualMachine
	Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error
}
