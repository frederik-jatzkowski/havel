package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Store struct {
	tool.Node[Store]

	Ptr  statements.Expression `parser:"'store' '(' @@"`
	Expr statements.Expression `parser:"',' @@ ')'"`
}

func (node *Store) ResolveNames(ctx context.Context) error {
	if err := node.Ptr.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Expr.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Store) ResolveTypes(ctx context.Context) error {
	if err := node.Ptr.ResolveTypes(ctx); err != nil {
		return err
	}

	if err := node.Expr.ResolveTypes(ctx); err != nil {
		return err
	}

	ptr := &types.Ptr{Underlying: node.Expr.Type()}

	if !node.Ptr.Type().Equals(ptr) {
		return node.Errorf("cannot store %s to %s, requires %s", node.Expr.Type(), node.Ptr.Type(), ptr)
	}

	return nil
}

func (node *Store) Type() types.Type {
	return &types.Void{}
}
