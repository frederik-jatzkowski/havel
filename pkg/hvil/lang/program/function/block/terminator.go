package block

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
)

type Terminator interface {
	ResolveNames(ctx context.Context) error
	ResolveTypes() error
	Execute(vm *runtime.VirtualMachine) (*Block, error)
}
