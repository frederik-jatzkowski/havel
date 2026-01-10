package instruction

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Instruction struct {
	tool.Node[Instruction]

	ResultWrite memory.Write `parser:"(@@ '=')?"`
	Operation   Operation    `parser:"@@ ';'"`
}

func (node *Instruction) ResolveNames(ctx context.Context) error {
	if node.ResultWrite != nil {
		if err := node.ResultWrite.ResolveNames(ctx); err != nil {
			return err
		}
	}

	return node.Operation.ResolveNames(ctx)
}

func (node *Instruction) ResolveTypes() error {
	if node.ResultWrite != nil {
		return node.Operation.ResolveTypes(node.ResultWrite.Type())
	}

	return node.Operation.ResolveTypes(&types.Void{})
}

func (node *Instruction) Execute(vm *runtime.VirtualMachine) error {
	var result unsafe.Pointer
	if node.ResultWrite != nil {
		result = node.ResultWrite.Addr(vm)
	}

	return node.Operation.Execute(vm, result)
}

func (node *Instruction) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Operation.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if node.ResultWrite != nil {
		if err := node.Result().GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return nil
}

func (node *Instruction) Result() memory.Write {
	return node.ResultWrite
}
