package program

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Program struct {
	tool.Node[Program]
	names.NameResolution[struct {
		Main      *function.Function
		Functions names.Scope[*function.Function]
	}]

	Functions []*function.Function `parser:"@@+"`
}

func (node *Program) ResolveNames(ctx context.Context) error {
	node.NameResolutionPass.Functions = names.NewRootScope[*function.Function](names.KindFunction)
	ctx = contexttool.WithScope(ctx, node.NameResolutionPass.Functions)

	for _, f := range node.Functions {
		if err := node.NameResolutionPass.Functions.Define(f); err != nil {
			return f.Wrap(err)
		}
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
