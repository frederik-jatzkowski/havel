package alu

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type SubU struct {
	tool.Node[SubU]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Left  memory.Read `parser:"'sub_u' '(' @@ ','"`
	Right memory.Read `parser:"@@ ')'"`
}

func (node *SubU) ResolveNames(ctx context.Context) error {
	if err := node.Left.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *SubU) ResolveTypes(target types.Type) error {
	left := node.Left.Type()
	right := node.Right.Type()

	_, ok := left.(*types.ScalarType)
	if !ok {
		return node.Errorf("operands must be a scalar type but was %s", left)
	}

	if !left.Equals(right) {
		return node.Errorf("cannot add %s and %s", left, right)
	}

	if !target.CanBeAssigned(left) {
		return node.Errorf("cannot assign %s result to %s", left, target)
	}

	node.TypeCheckPass.Type = left

	return nil
}

func (node *SubU) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *SubU) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Left.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Right.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	var op bytecode.OP
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		op = bytecode.OPAluSubU1
	case 2:
		op = bytecode.OPAluSubU2
	case 4:
		op = bytecode.OPAluSubU4
	case 8:
		op = bytecode.OPAluSubU8
	}

	p.AddI3R(
		op,
		node.RegisterAllocationPass.Result.(bytecode.R),
		node.Left.Register().(bytecode.R),
		node.Right.Register().(bytecode.R),
		node.Position(),
	)

	return nil
}

func (node *SubU) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		*(*byte)(result) = *(*byte)(node.Left.Addr(vm)) - *(*byte)(node.Right.Addr(vm))
	case 2:
		*(*uint16)(result) = *(*uint16)(node.Left.Addr(vm)) - *(*uint16)(node.Right.Addr(vm))
	case 4:
		*(*uint32)(result) = *(*uint32)(node.Left.Addr(vm)) - *(*uint32)(node.Right.Addr(vm))
	case 8:
		*(*uint64)(result) = *(*uint64)(node.Left.Addr(vm)) - *(*uint64)(node.Right.Addr(vm))
	}

	return nil
}
