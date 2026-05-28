package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Ptr struct {
	tool.Node[Ptr]
	typecheck.TypeCheck[struct {
		Result types.Type
	}]

	Ident Ident `parser:"'ptr' '(' @@ ')'"`
}

func (node *Ptr) ResolveNames(ctx context.Context) error {
	if err := node.Ident.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Ptr) ResolveTypes(ctx context.Context) error {
	if err := node.Ident.ResolveTypes(ctx); err != nil {
		return err
	}

	node.TypeCheckPass.Result = &types.Ptr{
		Underlying: node.Ident.Type(),
	}

	return nil
}

func (node *Ptr) Type() types.Type {
	return node.TypeCheckPass.Result
}
