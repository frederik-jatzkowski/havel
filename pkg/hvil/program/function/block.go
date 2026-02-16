package function

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Block interface {
	scope.Object
	names.Resolver
	codegen.VirtualMachine
	controlflow.Node

	FullyQualifiedIdentifier() string
	RegisterScope() scope.Scope[*instruction.RegWrite]
	ResolveTypes() error
	CalculateStatistics(ctx context.Context, blockID statistics.BlockID, current statistics.InstructionID) (next statistics.InstructionID)
	ResolveAddresses(offset int) int
	AllocateRegisters(scope registeralloc.Scope) error
}
