package functions

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Parameter struct {
	tool.Node[Parameter]

	Name string         `parser:"@Ident"`
	Typ  types.TypeName `parser:"@@"`
}

func (node *Parameter) Identifier() string {
	return node.Name
}

func (node *Parameter) ResolveNames(ctx context.Context) error {
	if err := contexttool.DefineInScope[scope.Object](ctx, node); err != nil {
		return err
	}

	return node.Typ.ResolveNames(ctx)
}

func (node *Parameter) Type() types.Type {
	return node.Typ.Type()
}
