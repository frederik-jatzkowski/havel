package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
)

type Operation interface {
	names.Resolver
	ResolveTypes(expected types.Type) error
	statistics.Calculator
	AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error)
	SetResultRegister(r architecture.Register)
	codegen.VirtualMachine
}
