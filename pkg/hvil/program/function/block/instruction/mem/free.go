package mem

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

type Free struct {
	tool.Node[Free]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Ptr instruction.MemoryRead `parser:"'free' '(' @@  ')'"`
}

func (node *Free) ResolveNames(ctx context.Context) error {
	if err := node.Ptr.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Free) ResolveTypes(target types.Type) error {
	if !target.CanBeAssigned(&types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	refType, ok := node.Ptr.Type().(*types.Ref)
	if !ok {
		return node.Errorf("%s is not a ref type", node.Ptr.Type())
	}

	node.TypeCheckPass.Type = refType

	return nil
}

func (node *Free) CalculateStatistics(ctx context.Context) {
	node.Ptr.CalculateStatistics(ctx)
}

func (node *Free) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Ptr.AllocateRegisters(scope)
}

func (node *Free) SetResultRegister(r architecture.Register) {
	panic(fmt.Sprintf("target register assigned to %T, which returns void", node))
}

func (node *Free) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Ptr.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI1R(bytecode.OPFree, node.Ptr.Register().(bytecode.R), node.Position())

	return nil
}
