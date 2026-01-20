package block

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/codegen"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
)

type Terminator interface {
	names.Resolver
	codegen.VirtualMachine

	ResolveTypes() error
	AllocateRegisters(scope registeralloc.Scope) error
}
