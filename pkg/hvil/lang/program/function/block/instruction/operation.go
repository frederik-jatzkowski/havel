package instruction

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Operation interface {
	names.Resolver
	ResolveTypes(expected types.Type) error
	statistics.Calculator
	AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error)
	SetResultRegister(r architecture.Register)
	codegen.VirtualMachine
}
