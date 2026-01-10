package block

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Terminator interface {
	ResolveNames(ctx context.Context) error
	ResolveTypes() error
	GenerateVirtualMachineAssembly(p *assembly.P, isMain bool) error
	Execute(vm *runtime.VirtualMachine) (*Block, error)
}
