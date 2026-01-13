package function

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Block interface {
	names.ScopedObject
	names.Resolver
	codegen.VirtualMachine

	FullyQualifiedIdentifier() string
	RegisterScope() names.Scope[*memory.RegWrite]
	ResolveTypes() error
	ResolveAddresses(offset int) int
	AllocateRegisters(arch architecture.Architecture) error
	Execute(vm *runtime.VirtualMachine) (Block, error)
}
