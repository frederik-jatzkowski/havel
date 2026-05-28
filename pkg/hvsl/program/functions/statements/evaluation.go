package statements

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Evaluation struct {
	tool.Node[Evaluation]

	Expr Expression `parser:"@@"`
}

func (node *Evaluation) ResolveNames(ctx context.Context) error {
	if err := node.Expr.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Evaluation) ResolveTypes(ctx context.Context) error {
	return node.Expr.ResolveTypes(ctx)
}
