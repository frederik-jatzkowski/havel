package call

import (
	"context"
	"fmt"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/scope"
)

type Local struct {
	tool.Node[Local]
	names.NameResolution[struct {
		Decl *functions.Function
	}]

	Name string                           `parser:"'local' '.' @Ident"`
	Args tool.List[statements.Expression] `parser:"'(' @@ ')'"`
}

func (node *Local) ResolveNames(ctx context.Context) error {
	object, err := contexttool.FromCtx[scope.Object](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	fn, ok := object.(*functions.Function)
	if !ok {
		return node.Wrap(fmt.Errorf("%s object is not a function declaration", node.Name))
	}

	node.NameResolutionPass.Decl = fn

	for i := range node.Args.Items {
		if err := node.Args.Items[i].ResolveNames(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (node *Local) ResolveTypes(ctx context.Context) error {
	for i := range node.Args.Items {
		if err := node.Args.Items[i].ResolveTypes(ctx); err != nil {
			return err
		}

		paramType := node.NameResolutionPass.Decl.Params.Items[i].Type()
		argType := node.Args.Items[i].Type()
		if !argType.Equals(paramType) {
			return node.Errorf("argument %d must be of type %s, got %s", i, paramType, argType)
		}
	}

	return nil
}

func (node *Local) Type() types.Type {
	return node.NameResolutionPass.Decl.Result.Type()
}
