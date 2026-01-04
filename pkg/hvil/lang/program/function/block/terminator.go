package block

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Terminator interface {
	ResolveNames(
		vars names.Scope[*stack.Decl],
		regs names.Scope[*memory.RegWrite],
		blocks names.Scope[*Block],
	) error
	ResolveTypes() error
	Execute(vm *runtime.VirtualMachine) (*Block, error)
}
