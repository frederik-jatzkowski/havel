package mem

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

type Load struct {
	tool.Node[Load]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Param instruction.MemoryRead `parser:"'load' '(' @@ ')'"`
}

func (node *Load) ResolveNames(ctx context.Context) error {
	return node.Param.ResolveNames(ctx)
}

func (node *Load) ResolveTypes(target types.Type) error {
	_, ok := node.Param.Type().(*types.RefType)
	if !ok {
		return node.Errorf("%s is not a ref type", node.Param.Type())
	}

	node.TypeCheckPass.Type = target

	return nil
}

func (node *Load) CalculateStatistics(ctx context.Context) {
	node.Param.CalculateStatistics(ctx)
}

func (node *Load) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Param.AllocateRegisters(scope)
}

func (node *Load) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Load) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if node.TypeCheckPass.Type.Bytes() == 0 {
		return nil // noop
	}

	if err := node.Param.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	op, err := bytecode.LoadForSize(node.TypeCheckPass.Type.Bytes())
	if err != nil {
		return node.Wrap(err)
	}

	p.AddI2R(op, node.RegisterAllocationPass.Result.(bytecode.R), node.Param.Register().(bytecode.R), node.Position())

	return nil
}
