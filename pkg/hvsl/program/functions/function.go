package functions

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Function struct {
	tool.Node[Function]
	names.NameResolution[struct {
		Scope scope.Scope[scope.Object]
	}]

	Name    string               `parser:"'func':Keyword @Ident"`
	Params  tool.List[Parameter] `parser:"'(' @@ ')'"`
	Result  *types.TypeName      `parser:"( '->' @@ )?"`
	Members []Member             `parser:"'{' (@@)* '}'"`
}

func (node *Function) Identifier() string {
	return node.Name
}

func (node *Function) ResolveNames(ctx context.Context) error {
	ctx = contexttool.WithCurrent(ctx, node)

	node.NameResolutionPass.Scope, ctx = contexttool.WithScope[scope.Object](ctx, names.KindIdentifier)

	for i := range node.Params.Items {
		if err := node.Params.Items[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	if node.Result != nil {
		if err := node.Result.ResolveNames(ctx); err != nil {
			return err
		}
	}

	for i := range node.Members {
		if err := node.Members[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Function) ResolveTypes(ctx context.Context) error {
	for i := range node.Members {
		if err := node.Members[i].ResolveTypes(ctx); err != nil {
			return err
		}
	}

	return nil
}
