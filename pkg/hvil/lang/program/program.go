package program

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Program struct {
	tool.Node[Program]
	names.NameResolution[struct {
		Main      *function.Function
		Functions names.Scope[function.Function]
	}]

	Functions []function.Function `parser:"@@+"`
}

func (p *Program) ResolveNames() []error {
	p.NameResolutionPass.Functions = names.NewRootScope[function.Function]("function")

	errs := p.NameResolutionPass.Functions.DefineAll(p.Functions)

	for i := 0; i < len(p.Functions); i++ {
		errs = append(errs, p.Functions[i].ResolveNames()...)
	}

	main, exists := p.NameResolutionPass.Functions.Find("main")
	if !exists {
		errs = append(errs, p.Errorf("no main function defined"))
	}

	p.NameResolutionPass.Main = main

	return errs
}

func (p *Program) ResolveTypes() (errs []error) {
	for i := 0; i < len(p.Functions); i++ {
		errs = append(errs, p.Functions[i].ResolveTypes()...)
	}

	return errs
}

func (p *Program) ResolveAddresses() (errs []error) {
	for i := 0; i < len(p.Functions); i++ {
		errs = append(errs, p.Functions[i].ResolveAddresses()...)
	}

	return errs
}

func (p *Program) Execute(vm *runtime.VirtualMachine) error {
	vm.CallStack = append(vm.CallStack, runtime.Call{
		Name: p.NameResolutionPass.Main.Name,
	})

	return p.NameResolutionPass.Main.Execute(vm)
}
