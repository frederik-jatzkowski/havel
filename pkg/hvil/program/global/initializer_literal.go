package global

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type InitializerLiteral struct {
	tool.Node[InitializerLiteral]
	typecheck.TypeCheck[struct {
		Size int
	}]

	Value uint64 `parser:"@Number"`
}

func (node *InitializerLiteral) ResolveNames(ctx context.Context) error {
	return nil
}

func (node *InitializerLiteral) ResolveTypes(expected types.Type) error {
	scalar, ok := expected.(*types.Scalar)
	if !ok {
		return node.Errorf("expected scalar but got %s", expected)
	}

	node.TypeCheckPass.Size = scalar.Size

	return nil
}

func (node *InitializerLiteral) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddSLit(node.TypeCheckPass.Size, node.Value)

	return nil
}
