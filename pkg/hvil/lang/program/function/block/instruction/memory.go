package instruction

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type MemoryRead interface {
	tool.NodeLike
	names.ScopedObject
	names.Resolver
	statistics.Calculator
	registeralloc.Value
	codegen.VirtualMachine
	Type() types.Type
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
}

type MemoryWrite interface {
	names.Resolver
	statistics.Calculator
	registeralloc.Value
	codegen.VirtualMachine
	Type() types.Type
	Addr(vm *runtime.VirtualMachine) unsafe.Pointer
}
