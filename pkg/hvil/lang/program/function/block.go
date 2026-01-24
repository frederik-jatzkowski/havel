package function

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Block interface {
	names.ScopedObject
	names.Resolver
	codegen.VirtualMachine
	controlflow.Node

	FullyQualifiedIdentifier() string
	RegisterScope() names.Scope[*memory.RegWrite]
	ResolveTypes() error
	CalculateStatistics(ctx context.Context, blockID statistics.BlockID, current statistics.InstructionID) (next statistics.InstructionID)
	ResolveAddresses(offset int) int
	AllocateRegisters(scope registeralloc.Scope) error
}
