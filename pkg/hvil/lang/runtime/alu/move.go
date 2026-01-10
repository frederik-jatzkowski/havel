package alu

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Move struct {
	tool.Node[Move]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]

	Arg memory.Read `parser:"'move' '(' @@ ')'"`
}

func (node *Move) ResolveNames(ctx context.Context) error {
	return node.Arg.ResolveNames(ctx)
}

func (node *Move) ResolveTypes(target types.Type) error {
	arg := node.Arg.Type()

	if !target.CanBeAssigned(arg) {
		return node.Errorf("cannot assign %s result to %s", arg, target)
	}

	node.TypeCheckPass.Type = target

	return nil
}

func (node *Move) SetResultRegister(r architecture.Register) {
	//TODO implement me
	panic("implement me")
}

func (node *Move) GenerateVirtualMachineAssembly(p *assembly.P) error {
	//TODO implement me
	panic("implement me")
}

func (node *Move) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	sourceAddr := node.Arg.Addr(vm)
	length := node.TypeCheckPass.Type.Bytes()
	copy(unsafe.Slice((*byte)(result), length), unsafe.Slice((*byte)(sourceAddr), length))

	return nil
}
