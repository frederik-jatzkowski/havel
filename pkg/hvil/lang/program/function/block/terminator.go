package block

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Terminator interface {
	names.Resolver
	codegen.VirtualMachine

	ResolveTypes() error
	AllocateRegisters(arch architecture.Architecture) error
	CalculateLiveRanges(ctx context.Context) error
	Execute(vm *runtime.VirtualMachine) (function.Block, error)
}
