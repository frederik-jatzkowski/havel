package codegen

import (
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
)

type VirtualMachine interface {
	GenerateVirtualMachineAssembly(p *assembly.P) error
}
