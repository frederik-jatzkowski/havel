package block

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Terminator interface {
	names.Resolver
	statistics.Calculator
	codegen.VirtualMachine

	ResolveTypes() error
	AllocateRegisters(scope registeralloc.Scope) error
}
