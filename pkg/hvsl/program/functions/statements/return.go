package statements

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Return struct {
	tool.Node[Return]
	names.NameResolution[struct {
		Current *functions.Function
	}]

	Expr Expression `parser:"'return':Keyword @@?"`
}

func (node *Return) ResolveNames(ctx context.Context) error {
	if node.Expr != nil {
		if err := node.Expr.ResolveNames(ctx); err != nil {
			return err
		}
	}

	currentFunction, err := contexttool.CurrentFromContext[*functions.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Current = currentFunction

	return nil
}

func (node *Return) ResolveTypes(ctx context.Context) error {
	var exprType types.Type = &types.Void{}

	if node.Expr != nil {
		if err := node.Expr.ResolveTypes(ctx); err != nil {
			return err
		}

		exprType = node.Expr.Type()
	}

	resultType := node.NameResolutionPass.Current.Result.Type()

	if !resultType.Equals(exprType) {
		return node.Errorf("return type %s does not match function return type %s", exprType, resultType)
	}

	return nil
}
