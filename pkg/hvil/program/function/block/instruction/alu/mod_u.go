package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type ModU struct {
	tool.Node[ModU]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Left  instruction.MemoryRead `parser:"'mod_u' '(' @@ ','"`
	Right instruction.MemoryRead `parser:"@@ ')'"`
}

func (node *ModU) ResolveNames(ctx context.Context) error {
	if err := node.Left.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *ModU) ResolveTypes(target types.Type) error {
	node.TypeCheckPass.Type = target

	return resolveBinOpTypesWithTarget(node, node.Left, node.Right, target)
}

func (node *ModU) CalculateStatistics(ctx context.Context) {
	node.Left.CalculateStatistics(ctx)
	node.Right.CalculateStatistics(ctx)
}

func (node *ModU) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
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

func (node *ModU) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *ModU) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Left.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Right.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	var op bytecode.OP
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		op = bytecode.OPAluModU8
	case 2:
		op = bytecode.OPAluModU16
	case 4:
		op = bytecode.OPAluModU32
	case 8:
		op = bytecode.OPAluModU64
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
