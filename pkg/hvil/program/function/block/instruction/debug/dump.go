package debug

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Dump struct {
	tool.Node[Dump]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Param instruction.MemoryRead `parser:"'dump' '(' @@ ')'"`
}

func (node *Dump) ResolveNames(ctx context.Context) error {
	return node.Param.ResolveNames(ctx)
}

func (node *Dump) ResolveTypes(target types.Type) error {
	if !target.Equals(&types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	node.TypeCheckPass.Type = node.Param.Type()

	return nil
}

func (node *Dump) CalculateStatistics(ctx context.Context) {
	node.Param.CalculateStatistics(ctx)
}

func (node *Dump) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Param.AllocateRegisters(scope)
}

func (node *Dump) SetResultRegister(r architecture.Register) {
	panic(fmt.Sprintf("target register assigned to %T, which returns void", node))
}

func (node *Dump) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Param.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI1R(bytecode.OPDebugDump, node.Param.Register().(bytecode.R), node.Position())

	return nil
}
