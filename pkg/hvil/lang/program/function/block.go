package function

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Block interface {
	names.ScopedObject
	names.Resolver
	codegen.VirtualMachine

	FullyQualifiedIdentifier() string
	RegisterScope() names.Scope[*memory.RegWrite]
	ResolveTypes() error
	ResolveAddresses(offset int) int
	AllocateRegisters(scope registeralloc.Scope) error
	CalculateLiveRanges(ctx context.Context) error
	Execute(vm *runtime.VirtualMachine) (Block, error)
}
