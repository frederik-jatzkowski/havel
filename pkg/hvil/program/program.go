package program

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/global"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Program struct {
	tool.Node[Program]
	names.NameResolution[struct {
		Main      *function.Function
		Functions scope.Scope[*function.Function]
	}]

	Declaration global.Block         `parser:"( @@ )?"`
	Functions   []*function.Function `parser:"@@+"`
}

func (node *Program) ResolveNames(ctx context.Context) error {
	node.NameResolutionPass.Functions, ctx = contexttool.WithScope[*function.Function](ctx, names.KindFunction)

	for _, f := range node.Functions {
		if err := node.NameResolutionPass.Functions.Define(f); err != nil {
			return f.Wrap(err)
		}
	}

	node.Declaration.NameResolutionPass.Declarations, ctx = contexttool.WithScope[instruction.VarDecl](ctx, names.KindVariable)
	if err := node.Declaration.ResolveNames(ctx); err != nil {
		return err
	}

	for i := 0; i < len(node.Functions); i++ {
		if err := node.Functions[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	main, err := node.NameResolutionPass.Functions.Find(names.SpecialMain)
	if err != nil {
		return node.Errorf("no main function defined")
	}

	node.NameResolutionPass.Main = main

	return err
}

func (node *Program) ResolveTypes() error {
	if err := node.Declaration.ResolveTypes(); err != nil {
		return err
	}

	for i := 0; i < len(node.Functions); i++ {
		if err := node.Functions[i].ResolveTypes(); err != nil {
			return err
		}
	}

	return nil
}

func (node *Program) CalculateStatistics(ctx context.Context) {
	for _, f := range node.Functions {
		f.CalculateStatistics(ctx)
	}
}

func (node *Program) ResolveAddresses(arch architecture.Architecture) error {
	if err := node.Declaration.ResolveAddresses(arch); err != nil {
		return err
	}

	for i := 0; i < len(node.Functions); i++ {
		if err := node.Functions[i].ResolveAddresses(arch); err != nil {
			return err
		}
	}

	return nil
}

func (node *Program) GenerateVirtualMachineAssembly() (*assembly.P, error) {
	if err := node.allocateRegisters(registeralloc.NewAllocator(virtualmachine.NewArchitecture())); err != nil {
		return nil, err
	}

	p := assembly.NewP()

	if err := node.Declaration.GenerateVirtualMachineAssembly(p); err != nil {
		return nil, err
	}

	p.AddJumpToLabel("main.entry", node.Position())

	return p, node.generateVirtualMachineAssembly(p)
}

func (node *Program) allocateRegisters(allocator registeralloc.Allocator) error {
	for i := 0; i < len(node.Functions); i++ {
		if err := node.Functions[i].AllocateRegisters(allocator); err != nil {
			return err
		}
	}

	return nil
}

func (node *Program) generateVirtualMachineAssembly(p *assembly.P) error {
	for i := 0; i < len(node.Functions); i++ {
		if err := node.Functions[i].GenerateVirtualMachineAssembly(p); err != nil {
			return err
		}
	}

	return nil
}
