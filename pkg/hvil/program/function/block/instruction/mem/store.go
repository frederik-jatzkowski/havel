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

type Store struct {
	tool.Node[Store]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Ptr   instruction.MemoryRead `parser:"'store' '(' @@  ','"`
	Value instruction.MemoryRead `parser:"@@ ')'"`
}

func (node *Store) ResolveNames(ctx context.Context) error {
	if err := node.Ptr.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Value.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Store) ResolveTypes(target types.Type) error {
	if !target.CanBeAssigned(&types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	_, ok := node.Ptr.Type().(*types.Ref)
	if !ok {
		return node.Errorf("%s is not a ref type", node.Ptr.Type())
	}

	node.TypeCheckPass.Type = node.Value.Type()

	return nil
}

func (node *Store) CalculateStatistics(ctx context.Context) {
	node.Ptr.CalculateStatistics(ctx)
}

func (node *Store) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	ptrRegs, err := node.Ptr.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	valueRegs, err := node.Value.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	return append(ptrRegs, valueRegs...), nil
}

func (node *Store) SetResultRegister(r architecture.Register) {
	panic(fmt.Sprintf("target register assigned to %T, which returns void", node))
}

func (node *Store) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Ptr.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Value.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	op, err := bytecode.StoreForSize(node.TypeCheckPass.Type.Bytes())
	if err != nil {
		return node.Wrap(err)
	}

	p.AddI2R(op, node.Ptr.Register().(bytecode.R), node.Value.Register().(bytecode.R), node.Position())

	return nil
}
