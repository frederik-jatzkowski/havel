package instruction

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type MemoryRead interface {
	tool.NodeLike
	names.ScopedObject
	names.Resolver
	registeralloc.Value
	codegen.VirtualMachine
	Type() types.Type
	CalculateLiveRanges(ctx context.Context) error
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
}

type MemoryWrite interface {
	names.Resolver
	registeralloc.Value
	codegen.VirtualMachine
	Type() types.Type
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
	CalculateLiveRanges(ctx context.Context) error
}
