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
		Functions names.Scope[*function.Function]
	}]

	Functions []*function.Function `parser:"@@+"`
}

func (node *Program) ResolveNames() []error {
	node.NameResolutionPass.Functions = names.NewRootScope[*function.Function]("function")

	errs := node.NameResolutionPass.Functions.DefineAll(node.Functions)

	for i := 0; i < len(node.Functions); i++ {
		errs = append(errs, node.Functions[i].ResolveNames()...)
	}

	main, err := node.NameResolutionPass.Functions.Find("main")
	if err != nil {
		errs = append(errs, node.Errorf("no main function defined"))
	}

	node.NameResolutionPass.Main = main

	return errs
}

func (node *Program) ResolveTypes() (errs []error) {
	for i := 0; i < len(node.Functions); i++ {
		errs = append(errs, node.Functions[i].ResolveTypes()...)
	}

	return errs
}

func (node *Program) ResolveAddresses() (errs []error) {
	for i := 0; i < len(node.Functions); i++ {
		errs = append(errs, node.Functions[i].ResolveAddresses()...)
	}

	return errs
}

func (node *Program) Execute(vm *runtime.VirtualMachine) error {
	vm.CallStack = append(vm.CallStack, runtime.Call{
		Name: node.NameResolutionPass.Main.Name,
	})

	return node.NameResolutionPass.Main.Execute(vm)
}
