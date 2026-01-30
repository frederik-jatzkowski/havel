package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
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

	Left  instruction.MemoryRead `parser:"'sub_u' '(' @@ ','"`
	Right instruction.MemoryRead `parser:"@@ ')'"`
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
	node.TypeCheckPass.Type = target

	return resolveBinOpTypesWithTarget(node, node.Left, node.Right, target)
}

func (node *SubU) CalculateStatistics(ctx context.Context) {
	node.Left.CalculateStatistics(ctx)
	node.Right.CalculateStatistics(ctx)
}

func (node *SubU) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
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
		op = bytecode.OPAluSubU8
	case 2:
		op = bytecode.OPAluSubU16
	case 4:
		op = bytecode.OPAluSubU32
	case 8:
		op = bytecode.OPAluSubU64
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
