package instruction

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Operation interface {
	names.Resolver
	ResolveTypes(expected types.Type) error
	CalculateLiveRanges(ctx context.Context) error
	AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error)
	SetResultRegister(r architecture.Register)
	codegen.VirtualMachine
	Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error
}
