package alu

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Move struct {
	tool.Node[Move]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Arg instruction.MemoryRead `parser:"'move' '(' @@ ')'"`
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

func (node *Move) CalculateStatistics() {
	node.Arg.CalculateStatistics()
}

func (node *Move) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Arg.AllocateRegisters(scope)
}

func (node *Move) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Move) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Arg.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI2R(
		bytecode.OPAluMove,
		node.RegisterAllocationPass.Result.(bytecode.R),
		node.Arg.Register().(bytecode.R),
		node.Position(),
	)

	return nil
}

func (node *Move) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	sourceAddr := node.Arg.Addr(vm)
	length := node.TypeCheckPass.Type.Bytes()
	copy(unsafe.Slice((*byte)(result), length), unsafe.Slice((*byte)(sourceAddr), length))

	return nil
}
