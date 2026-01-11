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

type DivU struct {
	tool.Node[DivU]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Left  memory.Read `parser:"'div_u' '(' @@ ','"`
	Right memory.Read `parser:"@@ ')'"`
}

func (node *DivU) ResolveNames(ctx context.Context) error {
	if err := node.Left.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *DivU) ResolveTypes(target types.Type) error {
	left := node.Left.Type()
	right := node.Right.Type()

	_, ok := left.(*types.ScalarType)
	if !ok {
		return node.Errorf("operands must be a scalar type but was %s", left)
	}

	if !left.Equals(right) {
		return node.Errorf("cannot divide %s by %s", left, right)
	}

	if !target.CanBeAssigned(left) {
		return node.Errorf("cannot assign %s result to %s", left, target)
	}

	node.TypeCheckPass.Type = left

	return nil
}

func (node *DivU) AllocateRegisters(arch architecture.Architecture) ([]architecture.Register, error) {
	leftRegs, err := node.Left.AllocateRegisters(arch)
	if err != nil {
		return nil, err
	}

	rightRegs, err := node.Right.AllocateRegisters(arch)
	if err != nil {
		return nil, err
	}

	return append(leftRegs, rightRegs...), nil
}

func (node *DivU) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *DivU) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Left.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Right.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	var op bytecode.OP
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		op = bytecode.OPAluDivU8
	case 2:
		op = bytecode.OPAluDivU16
	case 4:
		op = bytecode.OPAluDivU32
	case 8:
		op = bytecode.OPAluDivU64
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

func (node *DivU) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		if *(*byte)(node.Right.Addr(vm)) == 0 {
			return node.Errorf("division by zero")
		}
		*(*byte)(result) = *(*byte)(node.Left.Addr(vm)) / *(*byte)(node.Right.Addr(vm))
	case 2:
		if *(*uint16)(node.Right.Addr(vm)) == 0 {
			return node.Errorf("division by zero")
		}
		*(*uint16)(result) = *(*uint16)(node.Left.Addr(vm)) / *(*uint16)(node.Right.Addr(vm))
	case 4:
		if *(*uint32)(node.Right.Addr(vm)) == 0 {
			return node.Errorf("division by zero")
		}
		*(*uint32)(result) = *(*uint32)(node.Left.Addr(vm)) / *(*uint32)(node.Right.Addr(vm))
	case 8:
		if *(*uint64)(node.Right.Addr(vm)) == 0 {
			return node.Errorf("division by zero")
		}
		*(*uint64)(result) = *(*uint64)(node.Left.Addr(vm)) / *(*uint64)(node.Right.Addr(vm))
	}

	return nil
}
