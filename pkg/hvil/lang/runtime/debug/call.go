package debug

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Operation instruction.Operation

type Call struct {
	tool.Node[Call]

	Operation Operation `parser:"'debug' '.' @@"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
	return node.Operation.ResolveNames(ctx)
}

func (node *Call) ResolveTypes(target types.Type) error {
	return node.Operation.ResolveTypes(target)
}

func (node *Call) AllocateRegisters(arch architecture.Architecture) ([]architecture.Register, error) {
	return node.Operation.AllocateRegisters(arch)
}

func (node *Call) SetResultRegister(r architecture.Register) {
	node.Operation.SetResultRegister(r)
}

func (node *Call) GenerateVirtualMachineAssembly(p *assembly.P) error {
	return node.Operation.GenerateVirtualMachineAssembly(p)
}

func (node *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return node.Operation.Execute(vm, result)
}
