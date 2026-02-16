package mem

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Copy struct {
	tool.Node[Copy]

	Dest instruction.MemoryRead `parser:"'copy' '(' @@  ','"`
	Src  instruction.MemoryRead `parser:"@@ ','"`
	Size instruction.MemoryRead `parser:"@@ ')'"`
}

func (node *Copy) ResolveNames(ctx context.Context) error {
	if err := node.Dest.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Src.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Size.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Copy) ResolveTypes(target types.Type) error {
	if !target.Equals(&types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	if _, ok := node.Dest.Type().(*types.Ref); !ok {
		return node.Errorf("%s is not a ref type", node.Dest.Type())
	}

	if _, ok := node.Src.Type().(*types.Ref); !ok {
		return node.Errorf("%s is not a ref type", node.Src.Type())
	}

	if err := node.Size.Type().EqualsDetailed(&types.Scalar{Size: 8}); err != nil {
		return node.Size.Wrap(err)
	}

	return nil
}

func (node *Copy) CalculateStatistics(ctx context.Context) {
	node.Dest.CalculateStatistics(ctx)
	node.Src.CalculateStatistics(ctx)
	node.Size.CalculateStatistics(ctx)
}

func (node *Copy) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	destRegs, err := node.Dest.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	srcRegs, err := node.Src.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	sizeRegs, err := node.Size.AllocateRegisters(scope)
	if err != nil {
		return nil, err
	}

	return append(append(destRegs, srcRegs...), sizeRegs...), nil
}

func (node *Copy) SetResultRegister(r architecture.Register) {
	panic(fmt.Sprintf("target register assigned to %T, which returns void", node))
}

func (node *Copy) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Dest.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Src.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	if err := node.Size.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI3R(
		bytecode.OPCopy,
		node.Dest.Register().(bytecode.R),
		node.Src.Register().(bytecode.R),
		node.Size.Register().(bytecode.R),
		node.Position(),
	)

	return nil
}
