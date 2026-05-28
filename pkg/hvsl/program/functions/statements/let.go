package statements

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Let struct {
	tool.Node[Let]

	Name string         `parser:"'let':Keyword @Ident"`
	Typ  types.TypeName `parser:"@@"`
}

func (node *Let) Identifier() string {
	return node.Name
}

func (node *Let) ResolveNames(ctx context.Context) error {
	if err := node.Typ.ResolveNames(ctx); err != nil {
		return err
	}

	if err := contexttool.DefineInScope[scope.Object](ctx, node); err != nil {
		return err
	}

	return nil
}

func (node *Let) ResolveTypes(ctx context.Context) error {
	return node.Typ.ResolveTypes(ctx)
}

func (node *Let) Type() types.Type {
	return node.Typ.Type()
}
