package block

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
)

type Terminator interface {
	names.Resolver
	statistics.Calculator
	codegen.VirtualMachine

	ResolveTypes() error
	AllocateRegisters(scope registeralloc.Scope) error
	Successors() []controlflow.Node
}
