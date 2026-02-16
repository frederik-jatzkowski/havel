package alu

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Operation interface {
	tool.NodeLike
	instruction.Operation
}

type Call struct {
	tool.Node[Call]

	Operation Operation `parser:"'alu' '.' @@"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
	return node.Operation.ResolveNames(ctx)
}

func (node *Call) ResolveTypes(target types.Type) error {
	return node.Operation.ResolveTypes(target)
}

func (node *Call) CalculateStatistics(ctx context.Context) {
	node.Operation.CalculateStatistics(ctx)
}

func (node *Call) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Operation.AllocateRegisters(scope)
}

func (node *Call) SetResultRegister(r architecture.Register) {
	node.Operation.SetResultRegister(r)
}

func (node *Call) GenerateVirtualMachineAssembly(p *assembly.P) error {
	return node.Operation.GenerateVirtualMachineAssembly(p)
}
