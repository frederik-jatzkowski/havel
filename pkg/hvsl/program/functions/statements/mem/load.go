package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/functions/statements"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type Load struct {
	tool.Node[Load]
	typecheck.TypeCheck[struct {
		Underlying types.Type
	}]

	Select statements.Expression `parser:"'load' '(' @@ ')'"`
}

func (node *Load) ResolveNames(ctx context.Context) error {
	if err := node.Select.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *Load) ResolveTypes(ctx context.Context) error {
	if err := node.Select.ResolveTypes(ctx); err != nil {
		return err
	}

	ptr, ok := node.Select.Type().(*types.Ptr)
	if !ok {
		return node.Errorf("non pointer type in load statement: %s", node.Select.Type())
	}

	node.TypeCheckPass.Underlying = ptr.Underlying

	return nil
}

func (node *Load) Type() types.Type {
	return node.TypeCheckPass.Underlying
}
