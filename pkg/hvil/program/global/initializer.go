package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
)

type Initializer interface {
	ResolveNames(ctx context.Context) error
	ResolveTypes(expected types.Type) error
	GenerateVirtualMachineAssembly(p *assembly.P) error
}
