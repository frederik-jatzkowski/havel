package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type EQ struct {
	tool.Node[EQ]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Left  instruction.MemoryRead `parser:"'eq' '(' @@ ','"`
	Right instruction.MemoryRead `parser:"@@ ')'"`
}

func (node *EQ) ResolveNames(ctx context.Context) error {
	if err := node.Left.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *EQ) ResolveTypes(target types.Type) error {
	if !target.Equals(&types.Scalar{Size: 1}) {
		return node.Errorf("invalid target type: wants to assign 1 byte to %s", target)
	}

	node.TypeCheckPass.Type = target

	leftType := node.Left.Type()
	rightType := node.Right.Type()

	if leftType.Bytes() != rightType.Bytes() {
		return node.Errorf("unequally sized parameters %s and %s", leftType, rightType)
	}

	return nil
}

func (node *EQ) CalculateStatistics(ctx context.Context) {
	node.Left.CalculateStatistics(ctx)
	node.Right.CalculateStatistics(ctx)
}

func (node *EQ) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	leftRegs, err := node.Left.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	rightRegs, err := node.Right.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	return append(leftRegs, rightRegs...), nil
}

func (node *EQ) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *EQ) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Left.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Right.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI3R(
		bytecode.OPAluEq,
		node.RegisterAllocationPass.Result.(bytecode.R),
		node.Left.Register().(bytecode.R),
		node.Right.Register().(bytecode.R),
		node.Position(),
	)

	return nil
}
