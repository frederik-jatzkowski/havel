package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type InitializerCallPtr struct {
	tool.Node[InitializerCallPtr]
	names.NameResolution[struct {
		Function *function.Function
	}]

	Name string `parser:"'call' '.' 'ptr' '(' @Ident ')'"`
}

func (node *InitializerCallPtr) ResolveNames(ctx context.Context) error {
	fn, err := contexttool.FromCtx[*function.Function](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Function = fn

	return nil
}

func (node *InitializerCallPtr) ResolveTypes(expected types.Type) error {
	if err := expected.EqualsDetailed(node.NameResolutionPass.Function.Signature()); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *InitializerCallPtr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddSLabel(node.NameResolutionPass.Function.NameResolutionPass.Entry.FullyQualifiedIdentifier())

	return nil
}
