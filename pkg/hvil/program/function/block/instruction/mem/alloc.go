package mem

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

type Alloc struct {
	tool.Node[Alloc]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Size instruction.MemoryRead `parser:"'alloc' '(' @@ ')'"`
}

func (node *Alloc) ResolveNames(ctx context.Context) error {
	return node.Size.ResolveNames(ctx)
}

func (node *Alloc) ResolveTypes(target types.Type) error {
	_, ok := target.(*types.Ref)
	if !ok {
		return node.Errorf("%s is not a ref type", target)
	}

	_, ok = node.Size.Type().(*types.Scalar)
	if !ok {
		return node.Size.Errorf("%s is not a scalar type", node.Size.Type())
	}

	node.TypeCheckPass.Type = target

	return nil
}

func (node *Alloc) CalculateStatistics(ctx context.Context) {
	node.Size.CalculateStatistics(ctx)
}

func (node *Alloc) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Size.AllocateRegisters(scope)
}

func (node *Alloc) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Alloc) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Size.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI2R(bytecode.OPAlloc, node.RegisterAllocationPass.Result.(bytecode.R), node.Size.Register().(bytecode.R), node.Position())

	return nil
}
